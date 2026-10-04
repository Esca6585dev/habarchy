// Package worker wires asynq task handlers to the application services
// and manages long-lived SMPP sessions.
package worker

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"sync"

	"github.com/google/uuid"
	"github.com/hibiken/asynq"
	"github.com/rs/zerolog"

	"github.com/Esca6585dev/habarchy/backend/internal/adapters/providers/sms/smpp"
	"github.com/Esca6585dev/habarchy/backend/internal/app/delivery"
	"github.com/Esca6585dev/habarchy/backend/internal/app/webhooks"
	"github.com/Esca6585dev/habarchy/backend/internal/ports"
	"github.com/Esca6585dev/habarchy/backend/internal/queue"
)

// NewMux registers every task type.
func NewMux(d *delivery.Service, w *webhooks.Service) *asynq.ServeMux {
	mux := asynq.NewServeMux()
	for _, t := range []string{queue.TaskSendSMS, queue.TaskSendEmail, queue.TaskSendPush, queue.TaskSendTelegram, queue.TaskSendWhatsApp, queue.TaskSendSlack} {
		mux.HandleFunc(t, d.HandleTask)
	}
	mux.HandleFunc(queue.TaskWebhook, func(ctx context.Context, t *asynq.Task) error {
		var p queue.WebhookPayload
		if err := json.Unmarshal(t.Payload(), &p); err != nil {
			return fmt.Errorf("bad payload (%s): %w", err.Error(), asynq.SkipRetry)
		}
		return w.Deliver(ctx, p.DeliveryID)
	})
	return mux
}

// IsFailure tells asynq which errors count as failures in its metrics.
// Retryable outcomes are expected behaviour, not failures.
func IsFailure(err error) bool {
	return !errors.Is(err, delivery.ErrRetry) && !errors.Is(err, webhooks.ErrRetry)
}

// SMPPPool keeps one bound session per SMPP provider.
type SMPPPool struct {
	log      zerolog.Logger
	delivery *delivery.Service
	mu       sync.Mutex
	sessions map[uuid.UUID]*smpp.Session
}

// NewSMPPPool creates the pool.
func NewSMPPPool(log zerolog.Logger, d *delivery.Service) *SMPPPool {
	return &SMPPPool{log: log, delivery: d, sessions: map[uuid.UUID]*smpp.Session{}}
}

// Factory matches providers.Service.SMPPFactory. The providers service
// caches the result per provider updated_at, so a credentials change
// yields a new bind; the old session is closed here.
func (p *SMPPPool) Factory(ctx context.Context, id uuid.UUID, creds json.RawMessage) (ports.SMSProvider, error) {
	cfg, err := smpp.ParseConfig(creds)
	if err != nil {
		return nil, err
	}
	p.mu.Lock()
	if old, ok := p.sessions[id]; ok {
		_ = old.Close()
		delete(p.sessions, id)
	}
	p.mu.Unlock()

	log := p.log.With().Str("provider_id", id.String()).Str("smsc", cfg.Host).Logger()
	sess, err := smpp.Dial(ctx, cfg, log, func(ctx context.Context, r smpp.Receipt) {
		if !r.Final {
			return
		}
		if err := p.delivery.HandleReceipt(ctx, id, r.MessageID, r.Delivered, r.Status, map[string]any{"dlr": r.Text}); err != nil {
			log.Error().Err(err).Msg("apply smpp receipt")
		}
	})
	if err != nil {
		return nil, err
	}
	p.mu.Lock()
	p.sessions[id] = sess
	p.mu.Unlock()
	return sess, nil
}

// Close unbinds every session.
func (p *SMPPPool) Close() {
	p.mu.Lock()
	defer p.mu.Unlock()
	for id, s := range p.sessions {
		_ = s.Close()
		delete(p.sessions, id)
	}
}
