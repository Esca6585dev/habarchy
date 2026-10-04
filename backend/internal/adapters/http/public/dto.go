package public

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"

	"github.com/Esca6585dev/habarchy/backend/internal/adapters/postgres/sqlcgen"
	"github.com/Esca6585dev/habarchy/backend/internal/app/messages"
	"github.com/Esca6585dev/habarchy/backend/internal/domain"
)

// RecipientRef accepts "to": "+993..." or "to": {"contact_id": "..."} or
// {"external_id": "..."}.
type RecipientRef struct {
	Address    string
	ContactID  *uuid.UUID
	ExternalID string
}

// UnmarshalJSON implements the dual shape.
func (r *RecipientRef) UnmarshalJSON(b []byte) error {
	var s string
	if err := json.Unmarshal(b, &s); err == nil {
		r.Address = s
		return nil
	}
	var o struct {
		ContactID  *uuid.UUID `json:"contact_id"`
		ExternalID string     `json:"external_id"`
		Address    string     `json:"address"`
	}
	if err := json.Unmarshal(b, &o); err != nil {
		return err
	}
	r.ContactID, r.ExternalID, r.Address = o.ContactID, o.ExternalID, o.Address
	return nil
}

func (r RecipientRef) toDomain() messages.Recipient {
	return messages.Recipient{Address: r.Address, ContactID: r.ContactID, ExternalID: r.ExternalID}
}

type sendRequest struct {
	Channel        domain.Channel  `json:"channel"`
	To             RecipientRef    `json:"to"`
	Template       string          `json:"template" validate:"max=64"`
	Data           map[string]any  `json:"data"`
	Subject        string          `json:"subject" validate:"max=998"`
	Title          string          `json:"title" validate:"max=200"`
	Body           string          `json:"body" validate:"max=65536"`
	Locale         domain.Locale   `json:"locale"`
	ScheduledAt    *time.Time      `json:"scheduled_at"`
	Priority       domain.Priority `json:"priority"`
	IdempotencyKey string          `json:"idempotency_key" validate:"max=128"`
	Metadata       map[string]any  `json:"metadata"`
}

type batchRequest struct {
	Channel        domain.Channel  `json:"channel"`
	Template       string          `json:"template" validate:"max=64"`
	Subject        string          `json:"subject" validate:"max=998"`
	Title          string          `json:"title" validate:"max=200"`
	Body           string          `json:"body" validate:"max=65536"`
	Locale         domain.Locale   `json:"locale"`
	ScheduledAt    *time.Time      `json:"scheduled_at"`
	Priority       domain.Priority `json:"priority"`
	IdempotencyKey string          `json:"idempotency_key" validate:"max=128"`
	Metadata       map[string]any  `json:"metadata"`
	Recipients     []struct {
		To   RecipientRef   `json:"to"`
		Data map[string]any `json:"data"`
	} `json:"recipients" validate:"required,min=1"`
}

// MessageResponse is the public shape of a message.
type MessageResponse struct {
	ID                uuid.UUID       `json:"id"`
	Status            string          `json:"status"`
	Channel           string          `json:"channel"`
	To                string          `json:"to"`
	ContactID         *uuid.UUID      `json:"contact_id,omitempty"`
	BatchID           *uuid.UUID      `json:"batch_id,omitempty"`
	Template          string          `json:"template,omitempty"`
	TemplateVersion   *int32          `json:"template_version,omitempty"`
	Subject           string          `json:"subject,omitempty"`
	Body              string          `json:"body"`
	Priority          string          `json:"priority"`
	ProviderID        *uuid.UUID      `json:"provider_id,omitempty"`
	ProviderMessageID string          `json:"provider_message_id,omitempty"`
	ErrorCode         string          `json:"error_code,omitempty"`
	ErrorMessage      string          `json:"error_message,omitempty"`
	Attempts          int32           `json:"attempts"`
	ScheduledAt       *time.Time      `json:"scheduled_at,omitempty"`
	SentAt            *time.Time      `json:"sent_at,omitempty"`
	DeliveredAt       *time.Time      `json:"delivered_at,omitempty"`
	CostMicros        int64           `json:"cost_micros"`
	Currency          string          `json:"currency"`
	Metadata          json.RawMessage `json:"metadata"`
	IsTest            bool            `json:"is_test"`
	CreatedAt         time.Time       `json:"created_at"`
	UpdatedAt         time.Time       `json:"updated_at"`
}

// ToMessage converts a row.
func ToMessage(m *sqlcgen.Message) MessageResponse {
	return MessageResponse{
		ID: m.ID, Status: string(m.Status), Channel: string(m.Channel), To: m.ToAddress, ContactID: m.ContactID, BatchID: m.BatchID,
		Template: m.TemplateKey, TemplateVersion: m.TemplateVersion, Subject: m.RenderedSubject, Body: m.RenderedBody, Priority: string(m.Priority),
		ProviderID: m.ProviderID, ProviderMessageID: m.ProviderMessageID, ErrorCode: m.ErrorCode, ErrorMessage: m.ErrorMessage, Attempts: m.Attempts,
		ScheduledAt: m.ScheduledAt, SentAt: m.SentAt, DeliveredAt: m.DeliveredAt, CostMicros: m.CostMicros, Currency: m.Currency,
		Metadata: m.Metadata, IsTest: m.IsTest, CreatedAt: m.CreatedAt, UpdatedAt: m.UpdatedAt,
	}
}

// EventResponse is one timeline entry.
type EventResponse struct {
	ID         uuid.UUID       `json:"id"`
	Type       string          `json:"type"`
	ProviderID *uuid.UUID      `json:"provider_id,omitempty"`
	Payload    json.RawMessage `json:"payload"`
	CreatedAt  time.Time       `json:"created_at"`
}

// ToEvents converts rows.
func ToEvents(rows []sqlcgen.MessageEvent) []EventResponse {
	out := make([]EventResponse, 0, len(rows))
	for i := range rows {
		r := &rows[i]
		out = append(out, EventResponse{ID: r.ID, Type: string(r.Type), ProviderID: r.ProviderID, Payload: r.Payload, CreatedAt: r.CreatedAt})
	}
	return out
}

// BatchResponse is the public shape of a batch.
type BatchResponse struct {
	ID          uuid.UUID  `json:"id"`
	Status      string     `json:"status"`
	Channel     string     `json:"channel"`
	Template    string     `json:"template,omitempty"`
	Total       int32      `json:"total"`
	Queued      int32      `json:"queued"`
	Sent        int32      `json:"sent"`
	Delivered   int32      `json:"delivered"`
	Failed      int32      `json:"failed"`
	CreatedAt   time.Time  `json:"created_at"`
	CompletedAt *time.Time `json:"completed_at,omitempty"`
}

// ToBatch converts a row.
func ToBatch(b *sqlcgen.Batch) BatchResponse {
	return BatchResponse{ID: b.ID, Status: b.Status, Channel: string(b.Channel), Template: b.TemplateKey, Total: b.Total, Queued: b.Queued,
		Sent: b.Sent, Delivered: b.Delivered, Failed: b.Failed, CreatedAt: b.CreatedAt, CompletedAt: b.CompletedAt}
}

type otpSendRequest struct {
	Channel  domain.Channel `json:"channel"`
	To       string         `json:"to" validate:"required"`
	Length   int            `json:"length" validate:"omitempty,min=4,max=10"`
	TTL      int            `json:"ttl" validate:"omitempty,min=30,max=3600"` // seconds
	Template string         `json:"template" validate:"max=64"`
	Locale   domain.Locale  `json:"locale"`
	Data     map[string]any `json:"data"`
}

type otpVerifyRequest struct {
	To   string `json:"to" validate:"required"`
	Code string `json:"code" validate:"required,min=4,max=10"`
}

type contactRequest struct {
	ExternalID     string         `json:"external_id" validate:"max=128"`
	Phone          string         `json:"phone" validate:"max=32"`
	Email          string         `json:"email" validate:"max=254"`
	TelegramChatID string         `json:"telegram_chat_id" validate:"max=64"`
	Locale         domain.Locale  `json:"locale"`
	Tags           []string       `json:"tags" validate:"max=50,dive,max=64"`
	Attributes     map[string]any `json:"attributes"`
}

// ContactResponse is the public shape of a contact.
type ContactResponse struct {
	ID             uuid.UUID       `json:"id"`
	ExternalID     string          `json:"external_id"`
	Phone          string          `json:"phone"`
	Email          string          `json:"email"`
	TelegramChatID string          `json:"telegram_chat_id"`
	Locale         string          `json:"locale"`
	Tags           []string        `json:"tags"`
	Attributes     json.RawMessage `json:"attributes"`
	CreatedAt      time.Time       `json:"created_at"`
	UpdatedAt      time.Time       `json:"updated_at"`
}

// ToContact converts a row.
func ToContact(c *sqlcgen.Contact) ContactResponse {
	return ContactResponse{ID: c.ID, ExternalID: c.ExternalID, Phone: c.Phone, Email: c.Email, TelegramChatID: c.TelegramChatID,
		Locale: c.Locale, Tags: c.Tags, Attributes: c.Attributes, CreatedAt: c.CreatedAt, UpdatedAt: c.UpdatedAt}
}

type deviceRequest struct {
	Token      string          `json:"token" validate:"required,min=20,max=4096"`
	Platform   domain.Platform `json:"platform" validate:"required"`
	AppVersion string          `json:"app_version" validate:"max=64"`
	ContactID  *uuid.UUID      `json:"contact_id"`
	ExternalID string          `json:"external_id" validate:"max=128"`
}

// DeviceResponse is the public shape of a device.
type DeviceResponse struct {
	ID         uuid.UUID  `json:"id"`
	ContactID  *uuid.UUID `json:"contact_id,omitempty"`
	Platform   string     `json:"platform"`
	TokenHint  string     `json:"token_hint"`
	AppVersion string     `json:"app_version"`
	IsActive   bool       `json:"is_active"`
	LastSeenAt time.Time  `json:"last_seen_at"`
	CreatedAt  time.Time  `json:"created_at"`
}

// ToDevice converts a row (the full token is never echoed back).
func ToDevice(d *sqlcgen.Device) DeviceResponse {
	hint := d.FcmToken
	if len(hint) > 8 {
		hint = hint[:4] + "…" + hint[len(hint)-4:]
	}
	return DeviceResponse{ID: d.ID, ContactID: d.ContactID, Platform: string(d.Platform), TokenHint: hint, AppVersion: d.AppVersion,
		IsActive: d.IsActive, LastSeenAt: d.LastSeenAt, CreatedAt: d.CreatedAt}
}
