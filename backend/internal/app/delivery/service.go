// Package delivery is the worker side of the pipeline: it takes a queued
// message, walks the project's providers in priority order with fallback,
// records the outcome and emits webhooks.
package delivery

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/hibiken/asynq"
	"github.com/jackc/pgx/v5"
	"github.com/rs/zerolog"

	"github.com/Esca6585dev/habarchy/backend/internal/adapters/postgres"
	"github.com/Esca6585dev/habarchy/backend/internal/adapters/postgres/sqlcgen"
	rds "github.com/Esca6585dev/habarchy/backend/internal/adapters/redis"
	"github.com/Esca6585dev/habarchy/backend/internal/app/contacts"
	"github.com/Esca6585dev/habarchy/backend/internal/app/events"
	"github.com/Esca6585dev/habarchy/backend/internal/app/providers"
	"github.com/Esca6585dev/habarchy/backend/internal/app/webhooks"
	"github.com/Esca6585dev/habarchy/backend/internal/domain"
	"github.com/Esca6585dev/habarchy/backend/internal/ports"
	"github.com/Esca6585dev/habarchy/backend/internal/queue"
	"github.com/Esca6585dev/habarchy/backend/pkg/metrics"
)

// Service executes deliveries.
type Service struct {
	db        *postgres.DB
	redis     *rds.Client
	providers *providers.Service
	webhooks  *webhooks.Service
	contacts  *contacts.Service
	log       zerolog.Logger
	maxRetry  int
	// Events publishes live status updates for the dashboard (optional).
	Events *events.Publisher
}

// New creates the service. maxRetry is the number of asynq retries after
// the first attempt (spec: 3).
func New(db *postgres.DB, redis *rds.Client, p *providers.Service, w *webhooks.Service, c *contacts.Service, log zerolog.Logger, maxRetry int) *Service {
	if maxRetry < 0 {
		maxRetry = 3
	}
	return &Service{db: db, redis: redis, providers: p, webhooks: w, contacts: c, log: log, maxRetry: maxRetry}
}

// ErrRetry tells asynq to retry the task later.
var ErrRetry = errors.New("delivery: retryable failure")

// HandleTask is the asynq handler for every send:* task.
func (s *Service) HandleTask(ctx context.Context, t *asynq.Task) error {
	var p queue.SendPayload
	if err := json.Unmarshal(t.Payload(), &p); err != nil {
		return fmt.Errorf("bad payload: %w", err) // never retried: SkipRetry via asynq.SkipRetry
	}
	retried, _ := asynq.GetRetryCount(ctx)
	maxRetry, _ := asynq.GetMaxRetry(ctx)
	return s.Deliver(ctx, p.MessageID, retried, maxRetry)
}

// Deliver runs one attempt for the message. attempt is the number of
// previous attempts (asynq retry count).
func (s *Service) Deliver(ctx context.Context, id uuid.UUID, attempt, maxRetry int) error {
	msg, err := s.db.Queries.MarkMessageProcessing(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			// Cancelled, already sent, or deleted: nothing to do.
			return nil
		}
		return err
	}
	log := s.log.With().Str("message_id", id.String()).Str("channel", string(msg.Channel)).Int("attempt", attempt+1).Logger()
	ctx = log.WithContext(ctx)

	project, err := s.db.Queries.GetProject(ctx, msg.ProjectID)
	if err != nil {
		return err
	}
	if project.Status != sqlcgen.ProjectStatusActive {
		return s.fail(ctx, &msg, nil, "project_suspended", "project is "+string(project.Status), nil)
	}

	var candidates []sqlcgen.Provider
	if msg.IsTest {
		candidates = []sqlcgen.Provider{{ID: uuid.Nil, Name: "sandbox", Type: sqlcgen.ProviderType("sandbox"), Channel: msg.Channel}}
	} else {
		candidates, err = s.providers.ActiveForChannel(ctx, msg.ProjectID, domain.Channel(msg.Channel))
		if err != nil {
			return err
		}
		if len(candidates) == 0 {
			return s.fail(ctx, &msg, nil, string(domain.CodeProviderUnavailable), "no active provider for channel "+string(msg.Channel), nil)
		}
	}

	var lastErr *ports.ProviderError
	var lastProvider *sqlcgen.Provider
	for i := range candidates {
		prov := &candidates[i]
		res, perr := s.attempt(ctx, &msg, prov)
		if perr == nil {
			return s.succeed(ctx, &msg, prov, res)
		}
		lastErr, lastProvider = perr, prov
		s.event(ctx, msg.ID, sqlcgen.EventTypeAttempt, providerID(prov), map[string]any{
			"provider": prov.Name, "type": prov.Type, "error_code": perr.Code, "error": perr.Message, "retryable": perr.Retryable, "raw": perr.Raw,
		})
		if !perr.Retryable {
			// Permanent (invalid number, blocked user...): no fallback helps.
			return s.fail(ctx, &msg, prov, perr.Code, perr.Message, perr.Raw)
		}
		log.Warn().Str("provider", prov.Name).Str("code", perr.Code).Msg("provider failed, trying next")
	}

	// Every provider failed with a retryable error.
	if attempt < maxRetry {
		if err := s.db.Queries.RequeueMessage(ctx, sqlcgen.RequeueMessageParams{ID: msg.ID, ErrorCode: lastErr.Code, ErrorMessage: lastErr.Message}); err != nil {
			return err
		}
		return fmt.Errorf("%w: %s", ErrRetry, lastErr.Message)
	}
	return s.fail(ctx, &msg, lastProvider, lastErr.Code, lastErr.Message+" (retries exhausted)", lastErr.Raw)
}

func (s *Service) attempt(ctx context.Context, msg *sqlcgen.Message, prov *sqlcgen.Provider) (*ports.SendResult, *ports.ProviderError) {
	var adapter any
	if prov.ID == uuid.Nil {
		adapter = providers.Sandbox(domain.Channel(msg.Channel))
	} else {
		if prov.RateLimitPerSec > 0 {
			if err := s.redis.Wait(ctx, "provider:"+prov.ID.String(), float64(prov.RateLimitPerSec), int(prov.RateLimitPerSec)); err != nil {
				return nil, &ports.ProviderError{Code: "rate_limit_wait", Message: err.Error(), Retryable: true}
			}
		}
		var err error
		if adapter, err = s.providers.Build(ctx, prov); err != nil {
			return nil, &ports.ProviderError{Code: "provider_config", Message: err.Error(), Retryable: false}
		}
	}

	sendCtx, cancel := context.WithTimeout(ctx, 60*time.Second)
	defer cancel()
	started := time.Now()
	defer func() {
		metrics.ProviderDuration.WithLabelValues(string(prov.Type)).Observe(time.Since(started).Seconds())
	}()
	var res *ports.SendResult
	var err error
	switch domain.Channel(msg.Channel) {
	case domain.ChannelSMS:
		res, err = adapter.(ports.SMSProvider).Send(sendCtx, ports.SMSMessage{To: msg.ToAddress, Text: msg.RenderedBody, Ref: msg.ID.String()})
	case domain.ChannelEmail:
		res, err = adapter.(ports.EmailProvider).Send(sendCtx, emailFrom(msg))
	case domain.ChannelTelegram:
		res, err = adapter.(ports.TelegramProvider).Send(sendCtx, ports.TelegramMessage{ChatID: msg.ToAddress, Text: msg.RenderedBody, ParseMode: metaString(msg, "parse_mode")})
	case domain.ChannelWhatsApp, domain.ChannelSlack:
		res, err = adapter.(ports.ChatProvider).Send(sendCtx, ports.ChatMessage{To: msg.ToAddress, Text: msg.RenderedBody, Subject: msg.RenderedSubject, Extra: metaMap(msg)})
	case domain.ChannelPush:
		var pr *ports.PushResult
		pr, err = adapter.(ports.PushProvider).Send(sendCtx, ports.PushMessage{
			Tokens: []string{msg.ToAddress}, Title: msg.RenderedSubject, Body: msg.RenderedBody, Data: metaStringMap(msg, "data"),
		})
		if pr != nil {
			res = &pr.SendResult
			for _, tok := range pr.InvalidTokens {
				_ = s.contacts.DisableDevice(ctx, tok)
			}
		}
	default:
		return nil, &ports.ProviderError{Code: "unsupported_channel", Message: string(msg.Channel)}
	}
	if err != nil {
		metrics.ProviderCalls.WithLabelValues(string(prov.Type), "error").Inc()
		var pe *ports.ProviderError
		if errors.As(err, &pe) {
			return nil, pe
		}
		if errors.Is(err, context.DeadlineExceeded) {
			return nil, &ports.ProviderError{Code: "timeout", Message: "provider call timed out", Retryable: true}
		}
		return nil, &ports.ProviderError{Code: "provider_error", Message: err.Error(), Retryable: true}
	}
	metrics.ProviderCalls.WithLabelValues(string(prov.Type), "ok").Inc()
	return res, nil
}

func (s *Service) succeed(ctx context.Context, msg *sqlcgen.Message, prov *sqlcgen.Provider, res *ports.SendResult) error {
	currency := res.Currency
	if currency == "" {
		currency = "TMT"
	}
	if err := s.db.Queries.MarkMessageSent(ctx, sqlcgen.MarkMessageSentParams{
		ID: msg.ID, ProviderID: providerID(prov), ProviderMessageID: res.ProviderMessageID, CostMicros: res.CostMicros, Currency: currency,
	}); err != nil {
		return err
	}
	s.event(ctx, msg.ID, sqlcgen.EventTypeSent, providerID(prov), map[string]any{"provider": prov.Name, "provider_message_id": res.ProviderMessageID, "raw": res.Raw})
	msg.Status, msg.ProviderMessageID = sqlcgen.MessageStatusSent, res.ProviderMessageID
	msg.ErrorCode, msg.ErrorMessage = "", ""
	s.live(ctx, msg, prov.Name)
	s.emit(ctx, msg, domain.WebhookMessageSent)
	// Sandbox and channels without receipts are final at "sent"; push and
	// telegram have no DLR concept, so count them delivered.
	if msg.IsTest || msg.Channel == sqlcgen.ChannelPush || msg.Channel == sqlcgen.ChannelTelegram || msg.Channel == sqlcgen.ChannelSlack {
		return s.MarkDelivered(ctx, msg.ID, map[string]any{"implicit": true})
	}
	s.afterTerminal(ctx, msg)
	return nil
}

func (s *Service) fail(ctx context.Context, msg *sqlcgen.Message, prov *sqlcgen.Provider, code, message string, raw map[string]any) error {
	if err := s.db.Queries.MarkMessageFailed(ctx, sqlcgen.MarkMessageFailedParams{ID: msg.ID, ProviderID: providerID(prov), ErrorCode: code, ErrorMessage: truncate(message, 500)}); err != nil {
		return err
	}
	s.event(ctx, msg.ID, sqlcgen.EventTypeFailed, providerID(prov), map[string]any{"error_code": code, "error": message, "raw": raw})
	msg.Status, msg.ErrorCode, msg.ErrorMessage = sqlcgen.MessageStatusFailed, code, message
	provName := ""
	if prov != nil {
		provName = prov.Name
	}
	s.live(ctx, msg, provName)
	s.emit(ctx, msg, domain.WebhookMessageFailed)
	s.afterTerminal(ctx, msg)
	return nil
}

// MarkDelivered records a delivery receipt (DLR) for a sent message.
func (s *Service) MarkDelivered(ctx context.Context, id uuid.UUID, raw map[string]any) error {
	n, err := s.db.Queries.MarkMessageDelivered(ctx, id)
	if err != nil {
		return err
	}
	if n == 0 {
		return nil // not in "sent" state (duplicate DLR, or failed meanwhile)
	}
	msg, err := s.db.Queries.GetMessageByID(ctx, id)
	if err != nil {
		return err
	}
	s.event(ctx, id, sqlcgen.EventTypeDelivered, msg.ProviderID, raw)
	s.live(ctx, &msg, "")
	s.emit(ctx, &msg, domain.WebhookMessageDelivered)
	s.afterTerminal(ctx, &msg)
	return nil
}

// MarkFailedByReceipt records a negative DLR for a sent message.
func (s *Service) MarkFailedByReceipt(ctx context.Context, id uuid.UUID, code, message string, raw map[string]any) error {
	msg, err := s.db.Queries.GetMessageByID(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil
		}
		return err
	}
	if msg.Status != sqlcgen.MessageStatusSent {
		return nil
	}
	var prov *sqlcgen.Provider
	if msg.ProviderID != nil {
		prov = &sqlcgen.Provider{ID: *msg.ProviderID}
	}
	return s.fail(ctx, &msg, prov, code, message, raw)
}

// HandleReceipt resolves a provider message id to a message and applies
// the receipt. Used by HTTP DLR callbacks and the SMPP deliver_sm handler.
func (s *Service) HandleReceipt(ctx context.Context, providerID uuid.UUID, providerMessageID string, delivered bool, status string, raw map[string]any) error {
	msg, err := s.db.Queries.GetMessageByProviderMessageID(ctx, sqlcgen.GetMessageByProviderMessageIDParams{ProviderID: &providerID, ProviderMessageID: providerMessageID})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			s.log.Debug().Str("provider_message_id", providerMessageID).Msg("receipt for unknown message")
			return nil
		}
		return err
	}
	if raw == nil {
		raw = map[string]any{}
	}
	raw["status"] = status
	if delivered {
		return s.MarkDelivered(ctx, msg.ID, raw)
	}
	return s.MarkFailedByReceipt(ctx, msg.ID, "undelivered", "delivery report: "+status, raw)
}

// afterTerminal refreshes batch counters and fires batch.completed once.
func (s *Service) afterTerminal(ctx context.Context, msg *sqlcgen.Message) {
	if msg.BatchID == nil {
		return
	}
	b, err := s.db.Queries.RefreshBatchCounters(ctx, *msg.BatchID)
	if err != nil {
		s.log.Error().Err(err).Msg("refresh batch counters")
		return
	}
	done, err := s.db.Queries.CompleteBatchIfDone(ctx, b.ID)
	if err != nil {
		if !errors.Is(err, pgx.ErrNoRows) {
			s.log.Error().Err(err).Msg("complete batch")
		}
		return
	}
	bid := done.ID
	_ = s.webhooks.Emit(ctx, done.ProjectID, domain.WebhookBatchCompleted, nil, &bid, map[string]any{
		"batch_id": done.ID, "total": done.Total, "sent": done.Sent, "delivered": done.Delivered, "failed": done.Failed, "completed_at": done.CompletedAt,
	})
}

// live publishes the status change to dashboards and bumps metrics.
func (s *Service) live(ctx context.Context, msg *sqlcgen.Message, provider string) {
	metrics.MessagesTotal.WithLabelValues(string(msg.Channel), string(msg.Status)).Inc()
	if s.Events == nil {
		return
	}
	id := msg.ID
	s.Events.Publish(ctx, events.Event{
		Type: "message." + string(msg.Status), ProjectID: msg.ProjectID, MessageID: &id, Channel: string(msg.Channel),
		Status: string(msg.Status), To: msg.ToAddress, Provider: provider, ErrorCode: msg.ErrorCode,
	})
}

func (s *Service) emit(ctx context.Context, msg *sqlcgen.Message, event domain.WebhookEvent) {
	mid := msg.ID
	var meta map[string]any
	_ = json.Unmarshal(msg.Metadata, &meta)
	data := map[string]any{
		"message_id": msg.ID, "status": msg.Status, "channel": msg.Channel, "to": msg.ToAddress, "template": msg.TemplateKey,
		"provider_message_id": msg.ProviderMessageID, "error_code": msg.ErrorCode, "error_message": msg.ErrorMessage,
		"metadata": meta, "batch_id": msg.BatchID, "contact_id": msg.ContactID, "is_test": msg.IsTest,
	}
	if err := s.webhooks.Emit(ctx, msg.ProjectID, event, &mid, msg.BatchID, data); err != nil {
		s.log.Error().Err(err).Str("event", string(event)).Msg("emit webhook")
	}
}

func (s *Service) event(ctx context.Context, id uuid.UUID, typ sqlcgen.EventType, prov *uuid.UUID, payload map[string]any) {
	b, err := json.Marshal(payload)
	if err != nil {
		b = []byte("{}")
	}
	if _, err := s.db.Queries.CreateMessageEvent(ctx, sqlcgen.CreateMessageEventParams{MessageID: id, Type: typ, ProviderID: prov, Payload: b}); err != nil {
		s.log.Error().Err(err).Msg("write message event")
	}
}

func providerID(p *sqlcgen.Provider) *uuid.UUID {
	if p == nil || p.ID == uuid.Nil {
		return nil
	}
	id := p.ID
	return &id
}

func emailFrom(msg *sqlcgen.Message) ports.EmailMessage {
	em := ports.EmailMessage{To: msg.ToAddress, Subject: msg.RenderedSubject, ReplyTo: metaString(msg, "reply_to")}
	body := msg.RenderedBody
	if looksLikeHTML(body) {
		em.HTML = body
		if txt := metaString(msg, "text"); txt != "" {
			em.Text = txt
		}
	} else {
		em.Text = body
	}
	var meta struct {
		Attachments []struct {
			Filename    string `json:"filename"`
			ContentType string `json:"content_type"`
			Content     []byte `json:"content"` // base64 in JSON
			URL         string `json:"url"`
		} `json:"attachments"`
	}
	_ = json.Unmarshal(msg.Metadata, &meta)
	for _, a := range meta.Attachments {
		em.Attachments = append(em.Attachments, ports.EmailAttachment{Filename: a.Filename, ContentType: a.ContentType, Content: a.Content, URL: a.URL})
	}
	return em
}

func looksLikeHTML(s string) bool {
	t := strings.TrimSpace(strings.ToLower(s))
	return strings.HasPrefix(t, "<!doctype") || strings.HasPrefix(t, "<html") || strings.Contains(t, "<p") || strings.Contains(t, "<div") || strings.Contains(t, "<table") || strings.Contains(t, "<br")
}

func metaString(msg *sqlcgen.Message, key string) string {
	var meta map[string]any
	if json.Unmarshal(msg.Metadata, &meta) != nil {
		return ""
	}
	if v, ok := meta[key].(string); ok {
		return v
	}
	return ""
}

// metaMap returns the whole metadata object (provider-specific options).
func metaMap(msg *sqlcgen.Message) map[string]any {
	var m map[string]any
	_ = json.Unmarshal(msg.Metadata, &m)
	return m
}

func metaStringMap(msg *sqlcgen.Message, key string) map[string]string {
	var meta map[string]any
	if json.Unmarshal(msg.Metadata, &meta) != nil {
		return nil
	}
	raw, ok := meta[key].(map[string]any)
	if !ok {
		return nil
	}
	out := make(map[string]string, len(raw))
	for k, v := range raw {
		out[k] = fmt.Sprint(v)
	}
	return out
}

func truncate(s string, n int) string {
	if len(s) <= n {
		return s
	}
	return s[:n]
}
