package public

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"

	httpx "github.com/Esca6585dev/habarchy/backend/internal/adapters/http"
	"github.com/Esca6585dev/habarchy/backend/internal/adapters/providers/sms/httpgeneric"
	"github.com/Esca6585dev/habarchy/backend/internal/domain"
)

// smsCallback receives HTTP delivery reports for an http_sms provider:
// POST|GET /callbacks/sms/{provider_id}. The provider's DLR config maps the
// request fields to a message id and a status. No API key is required
// (gateways call it), so an unknown provider id answers 404 and nothing
// else is revealed.
func (h *Handlers) smsCallback(c *fiber.Ctx) error {
	pid, err := uuid.Parse(c.Params("provider_id"))
	if err != nil {
		return domain.ErrNotFound
	}
	prov, err := h.Providers.GetByID(c.UserContext(), pid)
	if err != nil || prov.Type != "http_sms" {
		return domain.ErrNotFound
	}
	creds, err := h.Providers.Decrypt(prov)
	if err != nil {
		return err
	}
	var cfg httpgeneric.Config
	if err := json.Unmarshal(creds, &cfg); err != nil {
		return err
	}
	if cfg.DLR.MessageIDParam == "" || cfg.DLR.StatusParam == "" {
		return domain.ErrValidation.WithMessage("provider has no dlr mapping configured")
	}

	values := map[string]string{}
	c.Context().QueryArgs().VisitAll(func(k, v []byte) { values[string(k)] = string(v) })
	ct := string(c.Request().Header.ContentType())
	switch {
	case strings.HasPrefix(ct, fiber.MIMEApplicationJSON):
		var body map[string]any
		if json.Unmarshal(c.Body(), &body) == nil {
			flatten("", body, values)
		}
	case strings.HasPrefix(ct, fiber.MIMEApplicationForm), strings.HasPrefix(ct, fiber.MIMEMultipartForm):
		c.Context().PostArgs().VisitAll(func(k, v []byte) { values[string(k)] = string(v) })
	}

	msgID, delivered := cfg.DLR.ParseDLR(values)
	if msgID == "" {
		return domain.ErrValidation.WithMessage("delivery report has no message id")
	}
	if delivered == nil {
		// Intermediate state (ENROUTE, ACCEPTD...): acknowledge, ignore.
		return httpx.OK(c, fiber.Map{"received": true, "applied": false})
	}
	raw := map[string]any{}
	for k, v := range values {
		raw[k] = v
	}
	if err := h.Delivery.HandleReceipt(c.UserContext(), pid, msgID, *delivered, values[cfg.DLR.StatusParam], raw); err != nil {
		return err
	}
	return httpx.OK(c, fiber.Map{"received": true, "applied": true})
}

// flatten turns nested JSON into dot-keys so "result.id" style params work.
func flatten(prefix string, in map[string]any, out map[string]string) {
	for k, v := range in {
		key := k
		if prefix != "" {
			key = prefix + "." + k
		}
		switch x := v.(type) {
		case map[string]any:
			flatten(key, x, out)
		case nil:
		default:
			out[key] = fmt.Sprint(x)
		}
	}
}
