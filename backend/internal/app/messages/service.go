// Package messages implements the send pipeline entry points: validate,
// resolve recipient and template, enforce quota and idempotency, persist
// and enqueue.
package messages

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/mail"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"

	"github.com/Esca6585dev/habarchy/backend/internal/adapters/postgres"
	"github.com/Esca6585dev/habarchy/backend/internal/adapters/postgres/sqlcgen"
	rds "github.com/Esca6585dev/habarchy/backend/internal/adapters/redis"
	"github.com/Esca6585dev/habarchy/backend/internal/app/contacts"
	"github.com/Esca6585dev/habarchy/backend/internal/app/events"
	"github.com/Esca6585dev/habarchy/backend/internal/app/templates"
	"github.com/Esca6585dev/habarchy/backend/internal/domain"
	"github.com/Esca6585dev/habarchy/backend/internal/ports"
	"github.com/Esca6585dev/habarchy/backend/pkg/ids"
	"github.com/Esca6585dev/habarchy/backend/pkg/phone"
	"github.com/Esca6585dev/habarchy/backend/pkg/render"
)

// Limits are the tunables the service needs from config.
type Limits struct {
	IdempotencyTTL     time.Duration
	BatchMaxRecipients int
	// MaxScheduleAhead bounds scheduled_at (default 30 days).
	MaxScheduleAhead time.Duration
}

// Service holds the message use cases.
type Service struct {
	db        *postgres.DB
	redis     *rds.Client
	queue     ports.Queue
	contacts  *contacts.Service
	templates *templates.Service
	limits    Limits
	now       func() time.Time
	// Events publishes live "queued" updates (optional).
	Events *events.Publisher
}

// New creates the service.
func New(db *postgres.DB, redis *rds.Client, queue ports.Queue, c *contacts.Service, t *templates.Service, limits Limits) *Service {
	if limits.IdempotencyTTL <= 0 {
		limits.IdempotencyTTL = 24 * time.Hour
	}
	if limits.BatchMaxRecipients <= 0 {
		limits.BatchMaxRecipients = 1000
	}
	if limits.MaxScheduleAhead <= 0 {
		limits.MaxScheduleAhead = 30 * 24 * time.Hour
	}
	return &Service{db: db, redis: redis, queue: queue, contacts: c, templates: t, limits: limits, now: time.Now}
}

// Recipient is either a raw address or a contact reference.
type Recipient struct {
	Address    string
	ContactID  *uuid.UUID
	ExternalID string
}

// SendInput is the payload of POST /messages.
type SendInput struct {
	Channel        domain.Channel
	To             Recipient
	Template       string
	Data           map[string]any
	Subject        string
	Body           string
	Locale         domain.Locale
	ScheduledAt    *time.Time
	Priority       domain.Priority
	IdempotencyKey string
	Metadata       map[string]any
	// Push extras
	Title string
}

// Caller identifies the project and key type making the request.
type Caller struct {
	Project sqlcgen.Project
	IsTest  bool
}

// Result is the API response of a send.
type Result struct {
	Message   sqlcgen.Message
	Duplicate bool // true when an idempotency key replayed an earlier message
}

// Send validates, persists and enqueues one message.
func (s *Service) Send(ctx context.Context, caller Caller, in SendInput) (*Result, error) {
	if err := s.validate(&in); err != nil {
		return nil, err
	}
	if in.IdempotencyKey != "" {
		if existing, err := s.replay(ctx, caller.Project.ID, in.IdempotencyKey); err != nil || existing != nil {
			return existing, err
		}
	}
	prepared, err := s.prepare(ctx, caller, in)
	if err != nil {
		return nil, err
	}
	// A contact with several devices yields several push messages; the
	// single-send API returns the first and queues all of them.
	if err := s.checkQuota(ctx, caller, len(prepared)); err != nil {
		return nil, err
	}
	if in.IdempotencyKey != "" {
		if err := s.redis.ReserveIdempotency(ctx, idemKey(caller.Project.ID, in.IdempotencyKey), prepared[0].ID.String(), s.limits.IdempotencyTTL); err != nil {
			var coll *rds.ErrIdempotencyCollision
			if errors.As(err, &coll) {
				return s.replayByID(ctx, caller.Project.ID, coll.Existing)
			}
			return nil, err
		}
	}
	msgs, err := s.persistAndEnqueue(ctx, caller, nil, prepared, in.IdempotencyKey)
	if err != nil {
		if in.IdempotencyKey != "" {
			_ = s.redis.ReleaseIdempotency(ctx, idemKey(caller.Project.ID, in.IdempotencyKey))
		}
		return nil, err
	}
	return &Result{Message: msgs[0]}, nil
}

// BatchInput is the payload of POST /messages/batch.
type BatchInput struct {
	Channel        domain.Channel
	Template       string
	Subject        string
	Body           string
	Locale         domain.Locale
	ScheduledAt    *time.Time
	Priority       domain.Priority
	IdempotencyKey string
	Metadata       map[string]any
	Title          string
	Recipients     []BatchRecipient
	// GroupIDs expands to every member of the contact groups (deduplicated
	// with Recipients that reference the same contact).
	GroupIDs []uuid.UUID
	// Data is shared template data; a recipient's own Data overlays it.
	Data map[string]any
}

// BatchRecipient is one entry of a batch.
type BatchRecipient struct {
	Recipient
	Data map[string]any
}

// BatchResult is the API response of a batch send.
type BatchResult struct {
	Batch     sqlcgen.Batch
	Accepted  int
	Rejected  []BatchRejection
	Duplicate bool
}

// BatchRejection explains why one recipient was skipped.
type BatchRejection struct {
	Index  int    `json:"index"`
	To     string `json:"to"`
	Reason string `json:"reason"`
}

// SendBatch creates one batch and up to BatchMaxRecipients messages.
// Invalid recipients are reported, not fatal, unless all are invalid.
func (s *Service) SendBatch(ctx context.Context, caller Caller, in BatchInput) (*BatchResult, error) {
	if len(in.GroupIDs) > 0 {
		ids, err := s.db.Queries.ListGroupContactIDs(ctx, sqlcgen.ListGroupContactIDsParams{GroupIds: in.GroupIDs, ProjectID: caller.Project.ID})
		if err != nil {
			return nil, err
		}
		seen := map[uuid.UUID]bool{}
		for _, r := range in.Recipients {
			if r.ContactID != nil {
				seen[*r.ContactID] = true
			}
		}
		for _, id := range ids {
			if seen[id] {
				continue
			}
			seen[id] = true
			cid := id
			in.Recipients = append(in.Recipients, BatchRecipient{Recipient: Recipient{ContactID: &cid}})
		}
		if len(in.Recipients) == 0 {
			return nil, domain.ErrInvalidRecipient.WithMessage("the selected groups have no members")
		}
	}
	if len(in.Recipients) == 0 {
		return nil, domain.ErrValidation.WithDetails(map[string]any{"recipients": "required"})
	}
	if len(in.Recipients) > s.limits.BatchMaxRecipients {
		return nil, domain.ErrValidation.WithDetails(map[string]any{"recipients": fmt.Sprintf("at most %d", s.limits.BatchMaxRecipients)})
	}
	base := SendInput{Channel: in.Channel, Template: in.Template, Subject: in.Subject, Body: in.Body, Locale: in.Locale,
		ScheduledAt: in.ScheduledAt, Priority: in.Priority, Metadata: in.Metadata, Title: in.Title, To: Recipient{Address: "placeholder"}}
	if err := s.validate(&base); err != nil {
		return nil, err
	}
	if in.IdempotencyKey != "" {
		b, err := s.db.Queries.GetBatchByIdempotencyKey(ctx, sqlcgen.GetBatchByIdempotencyKeyParams{ProjectID: caller.Project.ID, IdempotencyKey: &in.IdempotencyKey})
		if err == nil {
			return &BatchResult{Batch: b, Accepted: int(b.Total), Duplicate: true}, nil
		}
		if !errors.Is(err, pgx.ErrNoRows) {
			return nil, err
		}
	}

	var prepared []prepared
	var rejected []BatchRejection
	for i, r := range in.Recipients {
		one := base
		one.To, one.Data = r.Recipient, overlay(in.Data, r.Data)
		p, err := s.prepare(ctx, caller, one)
		if err != nil {
			rejected = append(rejected, BatchRejection{Index: i, To: describe(r.Recipient), Reason: reason(err)})
			continue
		}
		prepared = append(prepared, p...)
	}
	if len(prepared) == 0 {
		return nil, domain.ErrInvalidRecipient.WithMessage("no valid recipients").WithDetails(map[string]any{"rejected": rejected})
	}
	if err := s.checkQuota(ctx, caller, len(prepared)); err != nil {
		return nil, err
	}
	batchID := ids.New()
	var idem *string
	if in.IdempotencyKey != "" {
		idem = &in.IdempotencyKey
	}
	var batch sqlcgen.Batch
	err := s.db.WithTx(ctx, func(q *sqlcgen.Queries) error {
		var err error
		batch, err = q.CreateBatch(ctx, sqlcgen.CreateBatchParams{
			ID: batchID, ProjectID: caller.Project.ID, TemplateKey: in.Template, Channel: sqlcgen.Channel(prepared[0].Channel),
			Total: int32(len(prepared)), IdempotencyKey: idem, //nolint:gosec // bounded by BatchMaxRecipients
		})
		if err != nil {
			if postgres.IsUniqueViolation(err) {
				return domain.ErrDuplicateRequest.WithMessage("batch with this idempotency key already exists")
			}
		}
		return err
	})
	if err != nil {
		return nil, err
	}
	if _, err := s.persistAndEnqueue(ctx, caller, &batchID, prepared, ""); err != nil {
		return nil, err
	}
	if rejected == nil {
		rejected = []BatchRejection{}
	}
	return &BatchResult{Batch: batch, Accepted: len(prepared), Rejected: rejected}, nil
}

// prepared is a fully resolved message ready to be stored.
type prepared struct {
	ID              uuid.UUID
	Channel         domain.Channel
	To              string
	ContactID       *uuid.UUID
	TemplateKey     string
	TemplateVersion *int32
	Subject         string
	Body            string
	Priority        domain.Priority
	ScheduledAt     *time.Time
	Metadata        json.RawMessage
}

func (s *Service) validate(in *SendInput) error {
	details := map[string]any{}
	if in.Channel == "" {
		in.Channel = domain.ChannelAuto
	}
	if in.Channel != domain.ChannelAuto && !in.Channel.Valid() {
		details["channel"] = "sms, email, push, telegram or auto"
	}
	if in.To.Address == "" && in.To.ContactID == nil && in.To.ExternalID == "" {
		details["to"] = "required"
	}
	if in.Template == "" && strings.TrimSpace(in.Body) == "" {
		details["body"] = "either template or body is required"
	}
	if in.Template != "" && in.Body != "" {
		details["body"] = "use template or body, not both"
	}
	if in.Priority == "" {
		in.Priority = domain.PriorityNormal
	}
	if !in.Priority.Valid() {
		details["priority"] = "high, normal or low"
	}
	if in.Locale != "" && !in.Locale.Valid() {
		details["locale"] = "tk, ru or en"
	}
	if in.ScheduledAt != nil {
		if in.ScheduledAt.After(s.now().Add(s.limits.MaxScheduleAhead)) {
			details["scheduled_at"] = "too far in the future"
		}
		if in.ScheduledAt.Before(s.now().Add(-time.Minute)) {
			in.ScheduledAt = nil // already due: send now
		}
	}
	if len(in.IdempotencyKey) > 128 {
		details["idempotency_key"] = "max 128 characters"
	}
	if in.Channel == domain.ChannelAuto && in.To.Address != "" {
		details["channel"] = "auto requires a contact (contact_id or external_id), not a raw address"
	}
	if len(details) > 0 {
		return domain.ErrValidation.WithDetails(details)
	}
	return nil
}

// prepare resolves recipient, channel and content. It may return several
// messages (push to a contact with several devices).
func (s *Service) prepare(ctx context.Context, caller Caller, in SendInput) ([]prepared, error) {
	var contact *sqlcgen.Contact
	if in.To.ContactID != nil || in.To.ExternalID != "" {
		var err error
		if in.To.ContactID != nil {
			contact, err = s.contacts.Get(ctx, caller.Project.ID, *in.To.ContactID)
		} else {
			var c sqlcgen.Contact
			c, err = s.db.Queries.GetContactByExternalID(ctx, sqlcgen.GetContactByExternalIDParams{ProjectID: caller.Project.ID, ExternalID: in.To.ExternalID})
			if errors.Is(err, pgx.ErrNoRows) {
				err = domain.ErrInvalidRecipient.WithMessage("contact not found: " + in.To.ExternalID)
			}
			contact = &c
		}
		if err != nil {
			if errors.Is(err, domain.ErrNotFound) {
				return nil, domain.ErrInvalidRecipient.WithMessage("contact not found")
			}
			return nil, err
		}
	}

	channel := in.Channel
	if channel == domain.ChannelAuto {
		order := make([]domain.Channel, 0, len(caller.Project.AutoChannelOrder))
		for _, c := range caller.Project.AutoChannelOrder {
			order = append(order, domain.Channel(c))
		}
		if len(order) == 0 {
			order = domain.DefaultAutoOrder
		}
		channel = ""
		for _, c := range order {
			if addrs, _ := s.addressesFor(ctx, contact, c); len(addrs) > 0 {
				channel = c
				break
			}
		}
		if channel == "" {
			return nil, domain.ErrInvalidRecipient.WithMessage("contact has no reachable channel")
		}
	}

	var addresses []string
	if in.To.Address != "" {
		addr, err := normalizeAddress(channel, in.To.Address)
		if err != nil {
			return nil, err
		}
		addresses = []string{addr}
	} else {
		var err error
		if addresses, err = s.addressesFor(ctx, contact, channel); err != nil {
			return nil, err
		}
		if len(addresses) == 0 {
			return nil, domain.ErrInvalidRecipient.WithMessage("contact has no " + string(channel) + " address")
		}
	}

	// Content.
	locale := in.Locale
	if locale == "" && contact != nil {
		locale = domain.Locale(contact.Locale)
	}
	fallback := domain.Locale(caller.Project.DefaultLocale)
	data := mergeData(contact, in.Data)
	subject, body := in.Subject, in.Body
	var tplVersion *int32
	if in.Template != "" {
		tpl, err := s.templates.Resolve(ctx, caller.Project.ID, strings.ToLower(in.Template), channel, locale, fallback)
		if err != nil {
			if errors.Is(err, domain.ErrNotFound) {
				return nil, domain.ErrValidation.WithDetails(map[string]any{"template": err.Error()})
			}
			return nil, err
		}
		res, err := render.Render(channel, tpl.Subject, tpl.Body, data)
		if err != nil {
			if errors.Is(err, render.ErrMissingVars) {
				return nil, domain.ErrValidation.WithDetails(map[string]any{"data": err.Error()})
			}
			return nil, err
		}
		subject, body, tplVersion = res.Subject, res.Body, &tpl.Version
	} else if strings.Contains(body, "{{") || strings.Contains(subject, "{{") {
		// Ad hoc bodies may also use placeholders.
		res, err := render.Render(channel, subject, body, data)
		if err != nil {
			return nil, domain.ErrValidation.WithDetails(map[string]any{"body": err.Error()})
		}
		subject, body = res.Subject, res.Body
	}
	if channel == domain.ChannelPush && in.Title != "" {
		subject = in.Title
	}
	if channel == domain.ChannelEmail && strings.TrimSpace(subject) == "" {
		return nil, domain.ErrValidation.WithDetails(map[string]any{"subject": "required for email"})
	}

	meta := in.Metadata
	if meta == nil {
		meta = map[string]any{}
	}
	metaJSON, _ := json.Marshal(meta)
	out := make([]prepared, 0, len(addresses))
	for _, addr := range addresses {
		p := prepared{
			ID: ids.New(), Channel: channel, To: addr, TemplateKey: strings.ToLower(in.Template), TemplateVersion: tplVersion,
			Subject: subject, Body: body, Priority: in.Priority, ScheduledAt: in.ScheduledAt, Metadata: metaJSON,
		}
		if contact != nil {
			id := contact.ID
			p.ContactID = &id
		}
		out = append(out, p)
	}
	return out, nil
}

func (s *Service) addressesFor(ctx context.Context, c *sqlcgen.Contact, ch domain.Channel) ([]string, error) {
	if c == nil {
		return nil, nil
	}
	switch ch {
	case domain.ChannelSMS:
		if c.Phone != "" {
			return []string{c.Phone}, nil
		}
	case domain.ChannelEmail:
		if c.Email != "" {
			return []string{c.Email}, nil
		}
	case domain.ChannelTelegram:
		if c.TelegramChatID != "" {
			return []string{c.TelegramChatID}, nil
		}
	case domain.ChannelWhatsApp:
		if c.Whatsapp != "" {
			return []string{c.Whatsapp}, nil
		}
		if c.Phone != "" {
			return []string{c.Phone}, nil
		}
	case domain.ChannelSlack:
		if c.SlackID != "" {
			return []string{c.SlackID}, nil
		}
	case domain.ChannelPush:
		return s.contacts.ActiveTokens(ctx, c.ID)
	}
	return nil, nil
}

func normalizeAddress(ch domain.Channel, raw string) (string, error) {
	raw = strings.TrimSpace(raw)
	switch ch {
	case domain.ChannelSMS, domain.ChannelWhatsApp:
		p, err := phone.Normalize(raw, "")
		if err != nil {
			return "", domain.ErrInvalidRecipient.WithMessage("invalid phone number")
		}
		return p, nil
	case domain.ChannelSlack:
		if raw == "" || strings.ContainsAny(raw, " \n") {
			return "", domain.ErrInvalidRecipient.WithMessage("invalid slack channel or user id")
		}
		return raw, nil
	case domain.ChannelEmail:
		raw = strings.ToLower(raw)
		if _, err := mail.ParseAddress(raw); err != nil || strings.ContainsAny(raw, " <>") {
			return "", domain.ErrInvalidRecipient.WithMessage("invalid email address")
		}
		return raw, nil
	case domain.ChannelTelegram:
		if raw == "" || strings.ContainsAny(raw, " \n") {
			return "", domain.ErrInvalidRecipient.WithMessage("invalid telegram chat id")
		}
		return raw, nil
	case domain.ChannelPush:
		if len(raw) < 20 || strings.ContainsAny(raw, " \n") {
			return "", domain.ErrInvalidRecipient.WithMessage("invalid fcm token")
		}
		return raw, nil
	}
	return "", domain.ErrValidation.WithDetails(map[string]any{"channel": "unknown"})
}

// overlay returns shared data with per-recipient data on top.
func overlay(shared, own map[string]any) map[string]any {
	if len(shared) == 0 {
		return mergeData(nil, own)
	}
	out := make(map[string]any, len(shared)+len(own))
	for k, v := range shared {
		out[k] = v
	}
	for k, v := range own {
		out[k] = v
	}
	return out
}

// mergeData exposes contact attributes to templates under the data map,
// with explicit data winning. Contact fields are available as
// {{.contact.phone}} etc.
func mergeData(c *sqlcgen.Contact, data map[string]any) map[string]any {
	out := make(map[string]any, len(data)+1)
	if c != nil {
		var attrs map[string]any
		_ = json.Unmarshal(c.Attributes, &attrs)
		for k, v := range attrs {
			out[k] = v
		}
		out["contact"] = map[string]any{"id": c.ID.String(), "external_id": c.ExternalID, "name": c.Name, "phone": c.Phone, "email": c.Email, "locale": c.Locale}
		if _, ok := out["name"]; !ok && c.Name != "" {
			out["name"] = c.Name
		}
	}
	for k, v := range data {
		out[k] = v
	}
	return out
}

func (s *Service) checkQuota(ctx context.Context, caller Caller, n int) error {
	if caller.IsTest {
		return nil
	}
	p := caller.Project
	now := s.now().UTC()
	if p.DailyQuota > 0 {
		if err := s.quota(ctx, p.ID, "d:"+now.Format("2006-01-02"), p.DailyQuota, n, now.Add(-24*time.Hour).Truncate(24*time.Hour), 48*time.Hour, "daily"); err != nil {
			return err
		}
	}
	if p.MonthlyQuota > 0 {
		monthStart := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, time.UTC)
		if err := s.quota(ctx, p.ID, "m:"+now.Format("2006-01"), p.MonthlyQuota, n, monthStart, 32*24*time.Hour, "monthly"); err != nil {
			return err
		}
	}
	return nil
}

// quota increments a Redis counter seeded from the database on first use
// and rolls the increment back when the limit would be exceeded.
func (s *Service) quota(ctx context.Context, projectID uuid.UUID, period string, limit int64, n int, since time.Time, ttl time.Duration, name string) error {
	key := "quota:" + projectID.String() + ":" + period
	count, created, err := s.redis.Incr(ctx, key, int64(n), ttl)
	if err != nil {
		return err
	}
	if created {
		existing, err := s.db.Queries.CountMessagesSince(ctx, sqlcgen.CountMessagesSinceParams{ProjectID: projectID, Since: since})
		if err == nil && existing > 0 {
			c, _, _ := s.redis.Incr(ctx, key, existing, ttl)
			count = c
		}
	}
	if count > limit {
		_, _, _ = s.redis.Incr(ctx, key, -int64(n), ttl)
		return domain.ErrQuotaExceeded.WithMessage(fmt.Sprintf("%s quota of %d messages reached", name, limit)).
			WithDetails(map[string]any{"quota": name, "limit": limit})
	}
	return nil
}

func (s *Service) persistAndEnqueue(ctx context.Context, caller Caller, batchID *uuid.UUID, items []prepared, idemKeyStr string) ([]sqlcgen.Message, error) {
	out := make([]sqlcgen.Message, 0, len(items))
	err := s.db.WithTx(ctx, func(q *sqlcgen.Queries) error {
		for i, p := range items {
			var idem *string
			if idemKeyStr != "" && i == 0 {
				idem = &idemKeyStr
				if _, err := q.ReserveIdempotencyKey(ctx, sqlcgen.ReserveIdempotencyKeyParams{ProjectID: caller.Project.ID, IdempotencyKey: idemKeyStr, MessageID: p.ID}); err != nil {
					if errors.Is(err, pgx.ErrNoRows) {
						return domain.ErrDuplicateRequest
					}
					return err
				}
			}
			m, err := q.CreateMessage(ctx, sqlcgen.CreateMessageParams{
				ID: p.ID, ProjectID: caller.Project.ID, BatchID: batchID, Channel: sqlcgen.Channel(p.Channel), ToAddress: p.To, ContactID: p.ContactID,
				TemplateKey: p.TemplateKey, TemplateVersion: p.TemplateVersion, RenderedSubject: p.Subject, RenderedBody: p.Body,
				Status: sqlcgen.MessageStatusQueued, Priority: sqlcgen.MessagePriority(p.Priority), ScheduledAt: p.ScheduledAt,
				Metadata: p.Metadata, IdempotencyKey: idem, IsTest: caller.IsTest,
			})
			if err != nil {
				return err
			}
			payload, _ := json.Marshal(map[string]any{"scheduled_at": p.ScheduledAt, "priority": p.Priority})
			if _, err := q.CreateMessageEvent(ctx, sqlcgen.CreateMessageEventParams{MessageID: m.ID, Type: sqlcgen.EventTypeQueued, Payload: payload}); err != nil {
				return err
			}
			out = append(out, m)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	// Enqueue after commit so the worker always finds the rows.
	for i := range out {
		m := &out[i]
		if err := s.queue.EnqueueSend(ctx, m.ID, domain.Channel(m.Channel), domain.Priority(m.Priority), m.ScheduledAt); err != nil {
			_ = s.db.Queries.MarkMessageFailed(ctx, sqlcgen.MarkMessageFailedParams{ID: m.ID, ErrorCode: "enqueue_failed", ErrorMessage: err.Error()})
			return nil, err
		}
		if s.Events != nil {
			id := m.ID
			s.Events.Publish(ctx, events.Event{Type: "message.queued", ProjectID: m.ProjectID, MessageID: &id, Channel: string(m.Channel), Status: "queued", To: m.ToAddress})
		}
	}
	return out, nil
}

func (s *Service) replay(ctx context.Context, projectID uuid.UUID, key string) (*Result, error) {
	id, err := s.db.Queries.GetIdempotencyKey(ctx, sqlcgen.GetIdempotencyKeyParams{ProjectID: projectID, IdempotencyKey: key})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil //nolint:nilnil // no replay
		}
		return nil, err
	}
	return s.replayByID(ctx, projectID, id.String())
}

func (s *Service) replayByID(ctx context.Context, projectID uuid.UUID, idStr string) (*Result, error) {
	id, err := uuid.Parse(idStr)
	if err != nil {
		return nil, domain.ErrDuplicateRequest
	}
	m, err := s.db.Queries.GetMessage(ctx, sqlcgen.GetMessageParams{ID: id, ProjectID: projectID})
	if err != nil {
		return nil, domain.ErrDuplicateRequest
	}
	return &Result{Message: m, Duplicate: true}, nil
}

// Detail is a message with its timeline.
type Detail struct {
	Message sqlcgen.Message
	Events  []sqlcgen.MessageEvent
}

// Get returns a message and its events.
func (s *Service) Get(ctx context.Context, projectID, id uuid.UUID) (*Detail, error) {
	m, err := s.db.Queries.GetMessage(ctx, sqlcgen.GetMessageParams{ID: id, ProjectID: projectID})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrNotFound.WithMessage("message not found")
		}
		return nil, err
	}
	events, err := s.db.Queries.ListMessageEvents(ctx, id)
	if err != nil {
		return nil, err
	}
	return &Detail{Message: m, Events: events}, nil
}

// ListFilter is the query of GET /messages.
type ListFilter struct {
	Status  *domain.MessageStatus
	Channel *domain.Channel
	From    *time.Time
	To      *time.Time
	Cursor  *uuid.UUID
	Limit   int32
}

// List pages messages newest first with an id cursor.
func (s *Service) List(ctx context.Context, projectID uuid.UUID, f ListFilter) ([]sqlcgen.Message, *uuid.UUID, error) {
	if f.Limit <= 0 || f.Limit > 200 {
		f.Limit = 50
	}
	params := sqlcgen.ListMessagesParams{ProjectID: projectID, FromTs: f.From, ToTs: f.To, CursorID: f.Cursor, RowLimit: f.Limit + 1}
	if f.Status != nil {
		params.Status = sqlcgen.NullMessageStatus{MessageStatus: sqlcgen.MessageStatus(*f.Status), Valid: true}
	}
	if f.Channel != nil {
		params.Channel = sqlcgen.NullChannel{Channel: sqlcgen.Channel(*f.Channel), Valid: true}
	}
	rows, err := s.db.Queries.ListMessages(ctx, params)
	if err != nil {
		return nil, nil, err
	}
	var next *uuid.UUID
	if len(rows) > int(f.Limit) {
		rows = rows[:f.Limit]
		id := rows[len(rows)-1].ID
		next = &id
	}
	return rows, next, nil
}

// Cancel stops a queued (or scheduled) message.
func (s *Service) Cancel(ctx context.Context, projectID, id uuid.UUID) (*sqlcgen.Message, error) {
	m, err := s.db.Queries.GetMessage(ctx, sqlcgen.GetMessageParams{ID: id, ProjectID: projectID})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrNotFound.WithMessage("message not found")
		}
		return nil, err
	}
	if m.Status != sqlcgen.MessageStatusQueued {
		return nil, domain.ErrConflict.WithMessage("only queued messages can be cancelled (status is " + string(m.Status) + ")")
	}
	if c, ok := s.queue.(interface {
		CancelSend(context.Context, uuid.UUID, domain.Channel) (bool, error)
	}); ok {
		if _, err := c.CancelSend(ctx, id, domain.Channel(m.Channel)); err != nil {
			return nil, err
		}
	}
	n, err := s.db.Queries.CancelMessage(ctx, sqlcgen.CancelMessageParams{ID: id, ProjectID: projectID})
	if err != nil {
		return nil, err
	}
	if n == 0 {
		return nil, domain.ErrConflict.WithMessage("message is already being processed")
	}
	_, _ = s.db.Queries.CreateMessageEvent(ctx, sqlcgen.CreateMessageEventParams{MessageID: id, Type: sqlcgen.EventTypeCancelled, Payload: json.RawMessage(`{}`)})
	m.Status = sqlcgen.MessageStatusCancelled
	return &m, nil
}

// Resend re-queues a failed or cancelled message (admin action).
func (s *Service) Resend(ctx context.Context, projectID, id uuid.UUID) (*sqlcgen.Message, error) {
	m, err := s.db.Queries.ResetMessageForResend(ctx, sqlcgen.ResetMessageForResendParams{ID: id, ProjectID: projectID})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrConflict.WithMessage("only failed or cancelled messages can be resent")
		}
		return nil, err
	}
	_, _ = s.db.Queries.CreateMessageEvent(ctx, sqlcgen.CreateMessageEventParams{MessageID: id, Type: sqlcgen.EventTypeQueued, Payload: json.RawMessage(`{"resend":true}`)})
	if err := s.queue.EnqueueSend(ctx, m.ID, domain.Channel(m.Channel), domain.Priority(m.Priority), nil); err != nil {
		return nil, err
	}
	return &m, nil
}

// GetBatch returns a batch with refreshed counters.
func (s *Service) GetBatch(ctx context.Context, projectID, id uuid.UUID) (*sqlcgen.Batch, error) {
	if _, err := s.db.Queries.GetBatch(ctx, sqlcgen.GetBatchParams{ID: id, ProjectID: projectID}); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrNotFound.WithMessage("batch not found")
		}
		return nil, err
	}
	b, err := s.db.Queries.RefreshBatchCounters(ctx, id)
	if err != nil {
		return nil, err
	}
	return &b, nil
}

func idemKey(projectID uuid.UUID, key string) string { return projectID.String() + ":" + key }

func describe(r Recipient) string {
	switch {
	case r.Address != "":
		return r.Address
	case r.ContactID != nil:
		return "contact:" + r.ContactID.String()
	}
	return "external:" + r.ExternalID
}

func reason(err error) string {
	var de *domain.Error
	if errors.As(err, &de) {
		if len(de.Details) > 0 {
			b, _ := json.Marshal(de.Details)
			return de.Message + " " + string(b)
		}
		return de.Message
	}
	return "internal error"
}
