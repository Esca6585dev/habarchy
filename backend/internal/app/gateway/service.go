// Package gateway holds the use cases behind the Android SMS gateway: the
// phone authenticates with its pairing key, leases SMS from its outbox,
// reports results and heartbeats, and forwards SMS it received.
package gateway

import (
	"context"
	"encoding/json"
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"

	"github.com/Esca6585dev/habarchy/backend/internal/adapters/postgres"
	"github.com/Esca6585dev/habarchy/backend/internal/adapters/postgres/sqlcgen"
	"github.com/Esca6585dev/habarchy/backend/internal/adapters/providers/sms/androidgw"
	"github.com/Esca6585dev/habarchy/backend/internal/app/delivery"
	"github.com/Esca6585dev/habarchy/backend/internal/app/webhooks"
	"github.com/Esca6585dev/habarchy/backend/internal/domain"
	"github.com/Esca6585dev/habarchy/backend/internal/ports"
	"github.com/Esca6585dev/habarchy/backend/pkg/crypto"
	"github.com/Esca6585dev/habarchy/backend/pkg/ids"
)

// OnlineWindow is how recent a heartbeat must be for the phone to count
// as online on the health page.
const OnlineWindow = 90 * time.Second

// Service implements androidgw.Store and the phone-facing use cases.
type Service struct {
	db       *postgres.DB
	delivery *delivery.Service // nil in the worker is fine: receipts arrive via the API
	webhooks *webhooks.Service
	// LeasePoll is the interval for long-polling the outbox.
	LeasePoll time.Duration
}

// New creates the service.
func New(db *postgres.DB, d *delivery.Service, w *webhooks.Service) *Service {
	return &Service{db: db, delivery: d, webhooks: w, LeasePoll: time.Second}
}

// Caller is an authenticated phone.
type Caller struct {
	Provider sqlcgen.Provider
	Device   sqlcgen.GatewayDevice
}

// Authenticate resolves an X-Gateway-Key header.
func (s *Service) Authenticate(ctx context.Context, key string) (*Caller, error) {
	key = strings.TrimSpace(key)
	if len(key) < androidgw.MinKeyLen {
		return nil, domain.ErrUnauthorized.WithMessage("missing or malformed gateway key")
	}
	dev, err := s.db.Queries.GetGatewayDeviceByKeyHash(ctx, crypto.SHA256(key))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrUnauthorized.WithMessage("invalid gateway key")
		}
		return nil, err
	}
	prov, err := s.db.Queries.GetProviderByID(ctx, dev.ProviderID)
	if err != nil {
		return nil, err
	}
	if !prov.IsActive {
		return nil, domain.ErrUnauthorized.WithMessage("provider is disabled")
	}
	project, err := s.db.Queries.GetProject(ctx, prov.ProjectID)
	if err != nil {
		return nil, err
	}
	if project.Status != sqlcgen.ProjectStatusActive {
		return nil, domain.ErrUnauthorized.WithMessage("project is " + string(project.Status))
	}
	return &Caller{Provider: prov, Device: dev}, nil
}

// Lease hands up to limit pending SMS to the phone, long-polling for up
// to wait when the outbox is empty.
func (s *Service) Lease(ctx context.Context, providerID uuid.UUID, limit int, wait time.Duration) ([]sqlcgen.GatewayOutbox, error) {
	if limit <= 0 || limit > 20 {
		limit = 5
	}
	deadline := time.Now().Add(wait)
	for {
		rows, err := s.db.Queries.LeaseGatewayOutbox(ctx, sqlcgen.LeaseGatewayOutboxParams{ProviderID: providerID, RowLimit: int32(limit)}) //nolint:gosec // bounded above
		if err != nil {
			return nil, err
		}
		if len(rows) > 0 || time.Now().After(deadline) {
			return rows, nil
		}
		select {
		case <-ctx.Done():
			return []sqlcgen.GatewayOutbox{}, nil
		case <-time.After(s.LeasePoll):
		}
	}
}

// Report is a phone's result for one outbox row.
type Report struct {
	Status       string // sent | failed | delivered
	ErrorCode    string
	ErrorMessage string
	Parts        int
}

// ReportResult applies the phone's report. "sent"/"failed" wake the
// waiting worker; a later "delivered" (or a "failed" after "sent") is a
// delivery receipt for the message itself.
func (s *Service) ReportResult(ctx context.Context, caller *Caller, outboxID uuid.UUID, r Report) (*sqlcgen.GatewayOutbox, bool, error) {
	r.Status = strings.ToLower(strings.TrimSpace(r.Status))
	if r.Status != "sent" && r.Status != "failed" && r.Status != "delivered" {
		return nil, false, domain.ErrValidation.WithDetails(map[string]any{"status": "sent, failed or delivered"})
	}
	row, err := s.db.Queries.GetGatewayOutboxForProvider(ctx, sqlcgen.GetGatewayOutboxForProviderParams{ID: outboxID, ProviderID: caller.Provider.ID})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, false, domain.ErrNotFound
		}
		return nil, false, err
	}
	prev := row.Status
	applied := false
	switch prev {
	case "pending", "leased":
		// First report: the worker is (or was) waiting on this row.
		applied = true
	case "sent":
		// Receipt after the worker already marked the message sent.
		applied = r.Status != "sent"
	default:
		// expired / failed / delivered: record only.
	}
	if applied {
		if r.Status == "failed" && r.ErrorCode == "" {
			r.ErrorCode = "generic_failure"
		}
		row, err = s.db.Queries.SetGatewayOutboxResult(ctx, sqlcgen.SetGatewayOutboxResultParams{
			ID: outboxID, Status: r.Status, ErrorCode: r.ErrorCode, ErrorMessage: truncate(r.ErrorMessage, 500), Parts: int32(min(r.Parts, 255)), //nolint:gosec // bounded
		})
		if err != nil {
			return nil, false, err
		}
	}
	if prev == "sent" && applied && s.delivery != nil {
		raw := map[string]any{"gateway": "android", "outbox_id": outboxID.String()}
		delivered := r.Status == "delivered"
		status := r.Status
		if !delivered {
			status = r.ErrorCode
		}
		if err := s.delivery.HandleReceipt(ctx, caller.Provider.ID, outboxID.String(), delivered, status, raw); err != nil {
			return nil, false, err
		}
	}
	return &row, applied, nil
}

// Heartbeat stores the phone's state for the health page.
func (s *Service) Heartbeat(ctx context.Context, caller *Caller, info map[string]any) error {
	if info == nil {
		info = map[string]any{}
	}
	delete(info, "gateway_key")
	b, err := json.Marshal(info)
	if err != nil {
		return domain.ErrValidation.WithMessage("device info must be a JSON object")
	}
	return s.db.Queries.TouchGatewayDevice(ctx, sqlcgen.TouchGatewayDeviceParams{ProviderID: caller.Provider.ID, DeviceInfo: b})
}

// PendingCount is how many SMS wait for the phone.
func (s *Service) PendingCount(ctx context.Context, providerID uuid.UUID) (int64, error) {
	return s.db.Queries.CountGatewayOutboxPending(ctx, providerID)
}

// Inbound records an SMS the phone received and fires the sms.inbound webhook.
func (s *Service) Inbound(ctx context.Context, caller *Caller, from, text string, receivedAt time.Time) (*sqlcgen.GatewayInbound, error) {
	from, text = strings.TrimSpace(from), strings.TrimSpace(text)
	if from == "" || text == "" {
		return nil, domain.ErrValidation.WithDetails(map[string]any{"from": "required", "text": "required"})
	}
	if receivedAt.IsZero() {
		receivedAt = time.Now()
	}
	row, err := s.db.Queries.InsertGatewayInbound(ctx, sqlcgen.InsertGatewayInboundParams{
		ID: ids.New(), ProviderID: caller.Provider.ID, ProjectID: caller.Provider.ProjectID, FromAddress: from, Text: text, ReceivedAt: receivedAt,
	})
	if err != nil {
		return nil, err
	}
	if s.webhooks != nil {
		_ = s.webhooks.Emit(ctx, caller.Provider.ProjectID, domain.WebhookSMSInbound, nil, nil, map[string]any{
			"id": row.ID, "provider_id": row.ProviderID, "provider": caller.Provider.Name, "from": row.FromAddress, "text": row.Text, "received_at": row.ReceivedAt,
		})
	}
	return &row, nil
}

// Online reports whether a heartbeat arrived recently.
func Online(dev *sqlcgen.GatewayDevice, now time.Time) bool {
	return dev != nil && dev.LastSeenAt != nil && now.Sub(*dev.LastSeenAt) <= OnlineWindow
}

// ---- androidgw.Store ----

// Enqueue inserts an outbox row.
func (s *Service) Enqueue(ctx context.Context, providerID, messageID uuid.UUID, to, text string, simSlot int, expiresAt time.Time) (uuid.UUID, error) {
	row, err := s.db.Queries.EnqueueGatewayOutbox(ctx, sqlcgen.EnqueueGatewayOutboxParams{
		ID: ids.New(), ProviderID: providerID, MessageID: messageID, ToAddress: to, Text: text, SimSlot: int32(simSlot), ExpiresAt: expiresAt, //nolint:gosec // -1..1
	})
	if err != nil {
		return uuid.Nil, err
	}
	return row.ID, nil
}

// Outcome reads the row state.
func (s *Service) Outcome(ctx context.Context, outboxID uuid.UUID) (androidgw.Outcome, error) {
	row, err := s.db.Queries.GetGatewayOutbox(ctx, outboxID)
	if err != nil {
		return androidgw.Outcome{}, err
	}
	return androidgw.Outcome{Status: row.Status, ErrorCode: row.ErrorCode, ErrorMessage: row.ErrorMessage, Parts: int(row.Parts)}, nil
}

// Expire gives up on a row the phone did not handle in time.
func (s *Service) Expire(ctx context.Context, outboxID uuid.UUID) error {
	_, err := s.db.Queries.ExpireGatewayOutbox(ctx, outboxID)
	return err
}

// Factory builds the SMS provider for an android_sms row (set on
// providers.Service.GatewayFactory by the worker).
func (s *Service) Factory(providerID uuid.UUID, creds json.RawMessage) (ports.SMSProvider, error) {
	cfg, err := androidgw.ParseConfig(creds)
	if err != nil {
		return nil, err
	}
	return androidgw.New(providerID, cfg, s), nil
}

func truncate(v string, n int) string {
	if len(v) <= n {
		return v
	}
	return v[:n]
}
