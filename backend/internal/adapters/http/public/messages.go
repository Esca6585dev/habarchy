package public

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"

	httpx "github.com/Esca6585dev/habarchy/backend/internal/adapters/http"
	"github.com/Esca6585dev/habarchy/backend/internal/adapters/http/middleware"
	"github.com/Esca6585dev/habarchy/backend/internal/app/messages"
	"github.com/Esca6585dev/habarchy/backend/internal/domain"
)

func caller(c *fiber.Ctx) messages.Caller {
	k := middleware.Caller(c)
	return messages.Caller{Project: k.Project, IsTest: k.IsTest()}
}

func (h *Handlers) sendMessage(c *fiber.Ctx) error {
	var req sendRequest
	if err := httpx.Bind(c, &req); err != nil {
		return err
	}
	res, err := h.Messages.Send(c.UserContext(), caller(c), messages.SendInput{
		Channel: req.Channel, To: req.To.toDomain(), Template: req.Template, Data: req.Data, Subject: req.Subject, Title: req.Title,
		Body: req.Body, Locale: req.Locale, ScheduledAt: req.ScheduledAt, Priority: req.Priority, IdempotencyKey: req.IdempotencyKey, Metadata: req.Metadata,
	})
	if err != nil {
		return err
	}
	if res.Duplicate {
		return httpx.JSON(c, fiber.StatusOK, ToMessage(&res.Message), fiber.Map{"duplicate": true, "code": domain.CodeDuplicateRequest})
	}
	return httpx.JSON(c, fiber.StatusAccepted, fiber.Map{"id": res.Message.ID, "status": res.Message.Status, "channel": res.Message.Channel, "to": res.Message.ToAddress, "scheduled_at": res.Message.ScheduledAt}, nil)
}

func (h *Handlers) sendBatch(c *fiber.Ctx) error {
	var req batchRequest
	if err := httpx.Bind(c, &req); err != nil {
		return err
	}
	in := messages.BatchInput{
		Channel: req.Channel, Template: req.Template, Subject: req.Subject, Title: req.Title, Body: req.Body, Locale: req.Locale,
		ScheduledAt: req.ScheduledAt, Priority: req.Priority, IdempotencyKey: req.IdempotencyKey, Metadata: req.Metadata,
	}
	for _, r := range req.Recipients {
		in.Recipients = append(in.Recipients, messages.BatchRecipient{Recipient: r.To.toDomain(), Data: r.Data})
	}
	res, err := h.Messages.SendBatch(c.UserContext(), caller(c), in)
	if err != nil {
		return err
	}
	status := fiber.StatusAccepted
	meta := fiber.Map{"accepted": res.Accepted, "rejected": res.Rejected}
	if res.Duplicate {
		status, meta["duplicate"] = fiber.StatusOK, true
	}
	return httpx.JSON(c, status, ToBatch(&res.Batch), meta)
}

func (h *Handlers) getBatch(c *fiber.Ctx) error {
	id, err := httpx.ParamUUID(c, "batch_id")
	if err != nil {
		return err
	}
	b, err := h.Messages.GetBatch(c.UserContext(), caller(c).Project.ID, id)
	if err != nil {
		return err
	}
	return httpx.OK(c, ToBatch(b))
}

func (h *Handlers) getMessage(c *fiber.Ctx) error {
	id, err := httpx.ParamUUID(c, "message_id")
	if err != nil {
		return err
	}
	d, err := h.Messages.Get(c.UserContext(), caller(c).Project.ID, id)
	if err != nil {
		return err
	}
	return httpx.OK(c, fiber.Map{"message": ToMessage(&d.Message), "events": ToEvents(d.Events)})
}

func (h *Handlers) listMessages(c *fiber.Ctx) error {
	f, err := parseListFilter(c)
	if err != nil {
		return err
	}
	rows, next, err := h.Messages.List(c.UserContext(), caller(c).Project.ID, f)
	if err != nil {
		return err
	}
	out := make([]MessageResponse, 0, len(rows))
	for i := range rows {
		out = append(out, ToMessage(&rows[i]))
	}
	return httpx.JSON(c, fiber.StatusOK, out, fiber.Map{"next_cursor": next, "limit": f.Limit})
}

// ParseListFilter reads ?status=&channel=&from=&to=&cursor=&limit=.
func parseListFilter(c *fiber.Ctx) (messages.ListFilter, error) {
	f := messages.ListFilter{Limit: int32(c.QueryInt("limit", 50))} //nolint:gosec // bounded in service
	details := map[string]any{}
	if s := c.Query("status"); s != "" {
		st := domain.MessageStatus(s)
		if !st.Valid() {
			details["status"] = "unknown status"
		}
		f.Status = &st
	}
	if s := c.Query("channel"); s != "" {
		ch := domain.Channel(s)
		if !ch.Valid() {
			details["channel"] = "unknown channel"
		}
		f.Channel = &ch
	}
	for name, dst := range map[string]**time.Time{"from": &f.From, "to": &f.To} {
		if s := c.Query(name); s != "" {
			t, err := time.Parse(time.RFC3339, s)
			if err != nil {
				details[name] = "RFC3339 timestamp"
			}
			*dst = &t
		}
	}
	if s := c.Query("cursor"); s != "" {
		id, err := uuid.Parse(s)
		if err != nil {
			details["cursor"] = "must be a message id"
		}
		f.Cursor = &id
	}
	if len(details) > 0 {
		return f, domain.ErrValidation.WithDetails(details)
	}
	return f, nil
}

func (h *Handlers) cancelMessage(c *fiber.Ctx) error {
	id, err := httpx.ParamUUID(c, "message_id")
	if err != nil {
		return err
	}
	m, err := h.Messages.Cancel(c.UserContext(), caller(c).Project.ID, id)
	if err != nil {
		return err
	}
	return httpx.OK(c, ToMessage(m))
}

func (h *Handlers) otpSend(c *fiber.Ctx) error {
	var req otpSendRequest
	if err := httpx.Bind(c, &req); err != nil {
		return err
	}
	res, err := h.OTP.Send(c.UserContext(), caller(c), otpInput(req, c.IP()))
	if err != nil {
		return err
	}
	return httpx.JSON(c, fiber.StatusAccepted, res, nil)
}

func (h *Handlers) otpVerify(c *fiber.Ctx) error {
	var req otpVerifyRequest
	if err := httpx.Bind(c, &req); err != nil {
		return err
	}
	res, err := h.OTP.Verify(c.UserContext(), caller(c).Project.ID, req.To, req.Code)
	if err != nil {
		return err
	}
	return httpx.OK(c, res)
}
