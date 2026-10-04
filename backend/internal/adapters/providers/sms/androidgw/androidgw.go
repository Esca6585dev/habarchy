// Package androidgw is the SMS provider backed by an Android phone running
// the Habarchy Gateway app. Send places the message in the phone's outbox
// and waits until the phone reports sent or failed (or the timeout
// passes), so retries and fallback to the next provider keep working.
package androidgw

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"

	"github.com/Esca6585dev/habarchy/backend/internal/ports"
)

// Config is the provider's credentials JSON.
type Config struct {
	// GatewayKey pairs the phone with this provider (X-Gateway-Key header).
	GatewayKey string `json:"gateway_key"`
	// SimSlot selects the SIM on dual-SIM phones: -1 = default, 0 or 1.
	SimSlot int `json:"sim_slot"`
	// TimeoutSec is how long the worker waits for the phone to send
	// (10–120, default 45). Afterwards the next provider is tried.
	TimeoutSec int `json:"timeout_sec"`
	// InboundEnabled lets the phone forward SMS it receives
	// (POST /api/gateway/v1/inbound). Default true; the phone also has its
	// own switch and keeps its receiver disabled unless both are on.
	InboundEnabled *bool `json:"inbound_enabled,omitempty"`
}

// Inbound reports the effective inbound setting.
func (c Config) Inbound() bool { return c.InboundEnabled == nil || *c.InboundEnabled }

// MinKeyLen is the minimum gateway key length accepted.
const MinKeyLen = 16

// ParseConfig validates credentials and applies defaults.
func ParseConfig(creds json.RawMessage) (Config, error) {
	var cfg Config
	if err := json.Unmarshal(creds, &cfg); err != nil {
		return cfg, err
	}
	cfg.GatewayKey = strings.TrimSpace(cfg.GatewayKey)
	if len(cfg.GatewayKey) < MinKeyLen {
		return cfg, fmt.Errorf("android_sms: gateway_key must have at least %d characters", MinKeyLen)
	}
	if cfg.SimSlot < -1 || cfg.SimSlot > 1 {
		return cfg, errors.New("android_sms: sim_slot must be -1, 0 or 1")
	}
	if cfg.TimeoutSec == 0 {
		cfg.TimeoutSec = 45
	}
	if cfg.TimeoutSec < 10 || cfg.TimeoutSec > 120 {
		return cfg, errors.New("android_sms: timeout_sec must be 10-120")
	}
	return cfg, nil
}

// Outcome is what the store reports back for an outbox row.
type Outcome struct {
	Status       string // pending | leased | sent | failed | delivered | expired
	ErrorCode    string
	ErrorMessage string
	Parts        int
}

// Store persists outbox rows; implemented by app/gateway on Postgres.
type Store interface {
	Enqueue(ctx context.Context, providerID, messageID uuid.UUID, to, text string, simSlot int, expiresAt time.Time) (uuid.UUID, error)
	Outcome(ctx context.Context, outboxID uuid.UUID) (Outcome, error)
	Expire(ctx context.Context, outboxID uuid.UUID) error
}

// Provider implements ports.SMSProvider.
type Provider struct {
	id    uuid.UUID
	cfg   Config
	store Store
	// Poll is how often the outcome is checked while waiting.
	Poll time.Duration
}

// New builds the provider.
func New(providerID uuid.UUID, cfg Config, store Store) *Provider {
	return &Provider{id: providerID, cfg: cfg, store: store, Poll: time.Second}
}

// Send queues the SMS for the phone and waits for its report.
func (p *Provider) Send(ctx context.Context, m ports.SMSMessage) (*ports.SendResult, error) {
	msgID, err := uuid.Parse(m.Ref)
	if err != nil {
		return nil, &ports.ProviderError{Code: "gateway_ref", Message: "android_sms needs the message id as Ref", Retryable: false}
	}
	timeout := time.Duration(p.cfg.TimeoutSec) * time.Second
	expires := time.Now().Add(timeout)
	id, err := p.store.Enqueue(ctx, p.id, msgID, m.To, m.Text, p.cfg.SimSlot, expires)
	if err != nil {
		return nil, &ports.ProviderError{Code: "gateway_enqueue", Message: err.Error(), Retryable: true}
	}

	ticker := time.NewTicker(p.Poll)
	defer ticker.Stop()
	deadline := time.NewTimer(timeout)
	defer deadline.Stop()
	graceUsed := false
	for {
		select {
		case <-ctx.Done():
			_ = p.store.Expire(context.WithoutCancel(ctx), id)
			return nil, &ports.ProviderError{Code: "timeout", Message: "cancelled while waiting for the phone", Retryable: true}
		case <-deadline.C:
			out, err := p.store.Outcome(ctx, id)
			if err == nil && out.Status == "leased" && !graceUsed {
				// The phone took the message right before the deadline:
				// give it a moment instead of double-sending via fallback.
				graceUsed = true
				deadline.Reset(15 * time.Second)
				continue
			}
			_ = p.store.Expire(context.WithoutCancel(ctx), id)
			return nil, &ports.ProviderError{Code: "gateway_offline", Message: "phone did not send within " + timeout.String(), Retryable: true,
				Raw: map[string]any{"outbox_id": id.String()}}
		case <-ticker.C:
			out, err := p.store.Outcome(ctx, id)
			if err != nil {
				continue
			}
			switch out.Status {
			case "sent", "delivered":
				return &ports.SendResult{ProviderMessageID: id.String(), Raw: map[string]any{"gateway": "android", "parts": out.Parts}}, nil
			case "failed":
				return nil, &ports.ProviderError{Code: code(out.ErrorCode), Message: out.ErrorMessage, Retryable: Retryable(out.ErrorCode),
					Raw: map[string]any{"outbox_id": id.String(), "phone_error": out.ErrorCode}}
			case "expired":
				return nil, &ports.ProviderError{Code: "gateway_offline", Message: "outbox entry expired", Retryable: true}
			}
		}
	}
}

func code(c string) string {
	if c == "" {
		return "gateway_failed"
	}
	return "gateway_" + c
}

// Retryable classifies Android SmsManager result codes reported by the
// phone. Unknown codes are retried (another provider may succeed).
func Retryable(phoneCode string) bool {
	switch strings.ToLower(phoneCode) {
	case "null_pdu", "invalid_number", "blocked", "short_code_not_allowed", "short_code_never_allowed":
		return false
	}
	return true
}
