package admin

import (
	"encoding/json"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"

	httpx "github.com/Esca6585dev/habarchy/backend/internal/adapters/http"
	"github.com/Esca6585dev/habarchy/backend/internal/adapters/postgres/sqlcgen"
	"github.com/Esca6585dev/habarchy/backend/internal/domain"
)

// MessageRow is the admin shape of a message (same as public plus flags).
type MessageRow struct {
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
	IdempotencyKey    *string         `json:"idempotency_key,omitempty"`
	IsTest            bool            `json:"is_test"`
	CreatedAt         time.Time       `json:"created_at"`
	UpdatedAt         time.Time       `json:"updated_at"`
}

func toMessageRow(m *sqlcgen.Message) MessageRow {
	return MessageRow{
		ID: m.ID, Status: string(m.Status), Channel: string(m.Channel), To: m.ToAddress, ContactID: m.ContactID, BatchID: m.BatchID,
		Template: m.TemplateKey, TemplateVersion: m.TemplateVersion, Subject: m.RenderedSubject, Body: m.RenderedBody, Priority: string(m.Priority),
		ProviderID: m.ProviderID, ProviderMessageID: m.ProviderMessageID, ErrorCode: m.ErrorCode, ErrorMessage: m.ErrorMessage, Attempts: m.Attempts,
		ScheduledAt: m.ScheduledAt, SentAt: m.SentAt, DeliveredAt: m.DeliveredAt, CostMicros: m.CostMicros, Currency: m.Currency,
		Metadata: m.Metadata, IdempotencyKey: m.IdempotencyKey, IsTest: m.IsTest, CreatedAt: m.CreatedAt, UpdatedAt: m.UpdatedAt,
	}
}

// listMessages is the admin message log: filters + search + cursor.
func (h *Handlers) listMessages(c *fiber.Ctx) error {
	pid := membership(c).Project.ID
	limit := int32(c.QueryInt("limit", 50)) //nolint:gosec // bounded below
	if limit <= 0 || limit > 200 {
		limit = 50
	}
	params := sqlcgen.AdminListMessagesParams{ProjectID: pid, RowLimit: limit + 1}
	details := map[string]any{}
	if s := c.Query("status"); s != "" {
		if !domain.MessageStatus(s).Valid() {
			details["status"] = "unknown status"
		}
		params.Status = sqlcgen.NullMessageStatus{MessageStatus: sqlcgen.MessageStatus(s), Valid: true}
	}
	if s := c.Query("channel"); s != "" {
		if !domain.Channel(s).Valid() {
			details["channel"] = "unknown channel"
		}
		params.Channel = sqlcgen.NullChannel{Channel: sqlcgen.Channel(s), Valid: true}
	}
	for name, dst := range map[string]**time.Time{"from": &params.FromTs, "to": &params.ToTs} {
		if s := c.Query(name); s != "" {
			t, err := time.Parse(time.RFC3339, s)
			if err != nil {
				details[name] = "RFC3339 timestamp"
			}
			*dst = &t
		}
	}
	for name, dst := range map[string]**uuid.UUID{"cursor": &params.CursorID, "batch_id": &params.BatchID, "contact_id": &params.ContactID} {
		if s := c.Query(name); s != "" {
			id, err := uuid.Parse(s)
			if err != nil {
				details[name] = "must be a uuid"
			}
			*dst = &id
		}
	}
	if s := strings.TrimSpace(c.Query("search")); s != "" {
		params.Search = &s
	}
	includeTest := c.QueryBool("include_test", true)
	params.IncludeTest = &includeTest
	if len(details) > 0 {
		return domain.ErrValidation.WithDetails(details)
	}
	rows, err := h.DB.Queries.AdminListMessages(c.UserContext(), params)
	if err != nil {
		return err
	}
	var next *uuid.UUID
	if len(rows) > int(limit) {
		rows = rows[:limit]
		id := rows[len(rows)-1].ID
		next = &id
	}
	out := make([]MessageRow, 0, len(rows))
	for i := range rows {
		out = append(out, toMessageRow(&rows[i]))
	}
	return httpx.JSON(c, fiber.StatusOK, out, fiber.Map{"next_cursor": next, "limit": limit})
}

// getMessage returns the message, its timeline (with raw provider
// responses) and webhook deliveries about it.
func (h *Handlers) getMessage(c *fiber.Ctx) error {
	id, err := httpx.ParamUUID(c, "message_id")
	if err != nil {
		return err
	}
	d, err := h.Messages.Get(c.UserContext(), membership(c).Project.ID, id)
	if err != nil {
		return err
	}
	hooks, err := h.DB.Queries.ListWebhookDeliveriesFiltered(c.UserContext(), sqlcgen.ListWebhookDeliveriesFilteredParams{
		ProjectID: membership(c).Project.ID, MessageID: &id, RowLimit: 50,
	})
	if err != nil {
		return err
	}
	return httpx.OK(c, fiber.Map{"message": toMessageRow(&d.Message), "events": d.Events, "webhooks": toWebhookRows(hooks)})
}

// resendMessage re-queues a failed or cancelled message (developer role).
func (h *Handlers) resendMessage(c *fiber.Ctx) error {
	id, err := httpx.ParamUUID(c, "message_id")
	if err != nil {
		return err
	}
	m, err := h.Messages.Resend(c.UserContext(), membership(c).Project.ID, id)
	if err != nil {
		return err
	}
	return httpx.JSON(c, fiber.StatusAccepted, toMessageRow(m), nil)
}

func (h *Handlers) cancelMessage(c *fiber.Ctx) error {
	id, err := httpx.ParamUUID(c, "message_id")
	if err != nil {
		return err
	}
	m, err := h.Messages.Cancel(c.UserContext(), membership(c).Project.ID, id)
	if err != nil {
		return err
	}
	return httpx.OK(c, toMessageRow(m))
}

func (h *Handlers) listBatches(c *fiber.Ctx) error {
	page := httpx.ParsePage(c, 50, 200)
	rows, err := h.DB.Queries.ListBatches(c.UserContext(), sqlcgen.ListBatchesParams{ProjectID: membership(c).Project.ID, RowLimit: page.Limit, RowOffset: page.Offset})
	if err != nil {
		return err
	}
	return httpx.OK(c, rows)
}

func (h *Handlers) getBatch(c *fiber.Ctx) error {
	id, err := httpx.ParamUUID(c, "batch_id")
	if err != nil {
		return err
	}
	b, err := h.Messages.GetBatch(c.UserContext(), membership(c).Project.ID, id)
	if err != nil {
		return err
	}
	return httpx.OK(c, b)
}

// WebhookRow is the admin shape of a webhook delivery.
type WebhookRow struct {
	ID           uuid.UUID       `json:"id"`
	Event        string          `json:"event"`
	URL          string          `json:"url"`
	MessageID    *uuid.UUID      `json:"message_id,omitempty"`
	BatchID      *uuid.UUID      `json:"batch_id,omitempty"`
	Payload      json.RawMessage `json:"payload"`
	ResponseCode *int32          `json:"response_code"`
	ResponseBody string          `json:"response_body"`
	Attempts     int32           `json:"attempts"`
	NextRetryAt  *time.Time      `json:"next_retry_at"`
	DeliveredAt  *time.Time      `json:"delivered_at"`
	CreatedAt    time.Time       `json:"created_at"`
	Status       string          `json:"status"` // delivered | pending | failed
}

func toWebhookRows(rows []sqlcgen.WebhookDelivery) []WebhookRow {
	out := make([]WebhookRow, 0, len(rows))
	for i := range rows {
		out = append(out, toWebhookRow(&rows[i]))
	}
	return out
}

func toWebhookRow(w *sqlcgen.WebhookDelivery) WebhookRow {
	status := "pending"
	switch {
	case w.DeliveredAt != nil:
		status = "delivered"
	case w.NextRetryAt == nil && w.Attempts > 0:
		status = "failed"
	}
	return WebhookRow{ID: w.ID, Event: w.Event, URL: w.Url, MessageID: w.MessageID, BatchID: w.BatchID, Payload: w.Payload,
		ResponseCode: w.ResponseCode, ResponseBody: w.ResponseBody, Attempts: w.Attempts, NextRetryAt: w.NextRetryAt,
		DeliveredAt: w.DeliveredAt, CreatedAt: w.CreatedAt, Status: status}
}

func (h *Handlers) listWebhooks(c *fiber.Ctx) error {
	page := httpx.ParsePage(c, 50, 200)
	params := sqlcgen.ListWebhookDeliveriesFilteredParams{ProjectID: membership(c).Project.ID, RowLimit: page.Limit, RowOffset: page.Offset}
	if e := c.Query("event"); e != "" {
		params.Event = &e
	}
	if c.Query("failed") == "true" {
		t := true
		params.OnlyFailed = &t
	}
	rows, err := h.DB.Queries.ListWebhookDeliveriesFiltered(c.UserContext(), params)
	if err != nil {
		return err
	}
	return httpx.JSON(c, fiber.StatusOK, toWebhookRows(rows), fiber.Map{"limit": page.Limit, "offset": page.Offset})
}

func (h *Handlers) getWebhook(c *fiber.Ctx) error {
	id, err := httpx.ParamUUID(c, "delivery_id")
	if err != nil {
		return err
	}
	w, err := h.DB.Queries.GetWebhookDelivery(c.UserContext(), sqlcgen.GetWebhookDeliveryParams{ID: id, ProjectID: membership(c).Project.ID})
	if err != nil {
		return domain.ErrNotFound.WithMessage("webhook delivery not found")
	}
	return httpx.OK(c, toWebhookRow(&w))
}

func (h *Handlers) resendWebhook(c *fiber.Ctx) error {
	id, err := httpx.ParamUUID(c, "delivery_id")
	if err != nil {
		return err
	}
	w, err := h.Webhooks.Resend(c.UserContext(), membership(c).Project.ID, id)
	if err != nil {
		return err
	}
	return httpx.JSON(c, fiber.StatusAccepted, toWebhookRow(w), nil)
}
