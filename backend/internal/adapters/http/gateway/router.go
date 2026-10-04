// Package gateway serves /api/gateway/v1, the endpoints the Habarchy
// Gateway Android app calls with its pairing key (X-Gateway-Key).
package gateway

import (
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"

	httpx "github.com/Esca6585dev/habarchy/backend/internal/adapters/http"
	"github.com/Esca6585dev/habarchy/backend/internal/adapters/postgres/sqlcgen"
	"github.com/Esca6585dev/habarchy/backend/internal/app/gateway"
	"github.com/Esca6585dev/habarchy/backend/internal/domain"
)

// Handlers groups dependencies.
type Handlers struct {
	Gateway *gateway.Service
	// MaxWait caps the long-poll duration (default 25s).
	MaxWait time.Duration
}

const localCaller = "gateway.caller"

// Register mounts the routes.
func (h *Handlers) Register(app fiber.Router) {
	if h.MaxWait == 0 {
		h.MaxWait = 25 * time.Second
	}
	g := app.Group("/api/gateway/v1", h.requireKey)
	g.Get("/me", h.me)
	g.Get("/outbox", h.lease)
	g.Post("/outbox/:outbox_id/result", h.result)
	g.Post("/heartbeat", h.heartbeat)
	g.Post("/inbound", h.inbound)
}

func (h *Handlers) requireKey(c *fiber.Ctx) error {
	caller, err := h.Gateway.Authenticate(c.UserContext(), c.Get("X-Gateway-Key"))
	if err != nil {
		return err
	}
	c.Locals(localCaller, caller)
	return c.Next()
}

func caller(c *fiber.Ctx) *gateway.Caller {
	v, _ := c.Locals(localCaller).(*gateway.Caller)
	return v
}

// OutboxItem is what the phone receives.
type OutboxItem struct {
	ID        uuid.UUID `json:"id"`
	MessageID uuid.UUID `json:"message_id"`
	To        string    `json:"to"`
	Text      string    `json:"text"`
	SimSlot   int32     `json:"sim_slot"`
	ExpiresAt time.Time `json:"expires_at"`
}

func toItem(r *sqlcgen.GatewayOutbox) OutboxItem {
	return OutboxItem{ID: r.ID, MessageID: r.MessageID, To: r.ToAddress, Text: r.Text, SimSlot: r.SimSlot, ExpiresAt: r.ExpiresAt}
}

func (h *Handlers) me(c *fiber.Ctx) error {
	cl := caller(c)
	pending, _ := h.Gateway.PendingCount(c.UserContext(), cl.Provider.ID)
	cfg := h.Gateway.Config(cl)
	return httpx.OK(c, fiber.Map{
		"provider_id": cl.Provider.ID, "provider": cl.Provider.Name, "project_id": cl.Provider.ProjectID,
		"pending": pending, "last_seen_at": cl.Device.LastSeenAt, "server_time": time.Now().UTC(),
		"inbound_enabled": cfg.Inbound(), "sim_slot": cfg.SimSlot,
	})
}

func (h *Handlers) lease(c *fiber.Ctx) error {
	wait, _ := strconv.Atoi(c.Query("wait", "20"))
	limit, _ := strconv.Atoi(c.Query("limit", "5"))
	d := time.Duration(wait) * time.Second
	if d < 0 {
		d = 0
	}
	if d > h.MaxWait {
		d = h.MaxWait
	}
	rows, err := h.Gateway.Lease(c.UserContext(), caller(c).Provider.ID, limit, d)
	if err != nil {
		return err
	}
	out := make([]OutboxItem, 0, len(rows))
	for i := range rows {
		out = append(out, toItem(&rows[i]))
	}
	return httpx.OK(c, out)
}

type resultRequest struct {
	Status       string `json:"status" validate:"required,oneof=sent failed delivered"`
	ErrorCode    string `json:"error_code" validate:"max=64"`
	ErrorMessage string `json:"error_message" validate:"max=500"`
	Parts        int    `json:"parts" validate:"min=0,max=255"`
}

func (h *Handlers) result(c *fiber.Ctx) error {
	id, err := httpx.ParamUUID(c, "outbox_id")
	if err != nil {
		return err
	}
	var req resultRequest
	if err := httpx.Bind(c, &req); err != nil {
		return err
	}
	row, applied, err := h.Gateway.ReportResult(c.UserContext(), caller(c), id, gateway.Report{Status: req.Status, ErrorCode: req.ErrorCode, ErrorMessage: req.ErrorMessage, Parts: req.Parts})
	if err != nil {
		return err
	}
	return httpx.OK(c, fiber.Map{"id": row.ID, "status": row.Status, "applied": applied})
}

func (h *Handlers) heartbeat(c *fiber.Ctx) error {
	var info map[string]any
	if len(c.Body()) > 0 {
		if err := c.BodyParser(&info); err != nil {
			return domain.ErrValidation.WithMessage("body must be a JSON object")
		}
	}
	if err := h.Gateway.Heartbeat(c.UserContext(), caller(c), info); err != nil {
		return err
	}
	pending, _ := h.Gateway.PendingCount(c.UserContext(), caller(c).Provider.ID)
	return httpx.OK(c, fiber.Map{"ok": true, "pending": pending, "server_time": time.Now().UTC(), "inbound_enabled": h.Gateway.Config(caller(c)).Inbound()})
}

type inboundRequest struct {
	From       string     `json:"from" validate:"required,max=32"`
	Text       string     `json:"text" validate:"required,max=4096"`
	ReceivedAt *time.Time `json:"received_at"`
}

func (h *Handlers) inbound(c *fiber.Ctx) error {
	var req inboundRequest
	if err := httpx.Bind(c, &req); err != nil {
		return err
	}
	var at time.Time
	if req.ReceivedAt != nil {
		at = *req.ReceivedAt
	}
	row, err := h.Gateway.Inbound(c.UserContext(), caller(c), req.From, req.Text, at)
	if err != nil {
		return err
	}
	return httpx.Created(c, fiber.Map{"id": row.ID, "received_at": row.ReceivedAt})
}
