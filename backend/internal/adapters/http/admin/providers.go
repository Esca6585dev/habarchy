package admin

import (
	"encoding/json"
	"net/url"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"

	httpx "github.com/Esca6585dev/habarchy/backend/internal/adapters/http"
	"github.com/Esca6585dev/habarchy/backend/internal/adapters/postgres/sqlcgen"
	"github.com/Esca6585dev/habarchy/backend/internal/app/providers"
	"github.com/Esca6585dev/habarchy/backend/internal/domain"
	"github.com/Esca6585dev/habarchy/backend/internal/ports"
)

type providerRequest struct {
	Name            string              `json:"name" validate:"required_without=IsUpdate,max=100"`
	Type            domain.ProviderType `json:"type"`
	Priority        *int                `json:"priority" validate:"omitempty,min=0,max=1000"`
	IsActive        *bool               `json:"is_active"`
	Credentials     json.RawMessage     `json:"credentials"`
	RateLimitPerSec *int                `json:"rate_limit_per_sec" validate:"omitempty,min=0,max=10000"`
	IsUpdate        bool                `json:"-"`
}

// ProviderResponse never includes raw credentials; Redacted shows masked
// non-secret settings so the form can be re-edited.
type ProviderResponse struct {
	ID              uuid.UUID      `json:"id"`
	Name            string         `json:"name"`
	Channel         string         `json:"channel"`
	Type            string         `json:"type"`
	Priority        int32          `json:"priority"`
	IsActive        bool           `json:"is_active"`
	RateLimitPerSec int32          `json:"rate_limit_per_sec"`
	Settings        map[string]any `json:"settings,omitempty"`
	CreatedAt       time.Time      `json:"created_at"`
	UpdatedAt       time.Time      `json:"updated_at"`
}

func (h *Handlers) toProvider(p *sqlcgen.Provider, withSettings bool) ProviderResponse {
	out := ProviderResponse{ID: p.ID, Name: p.Name, Channel: string(p.Channel), Type: string(p.Type), Priority: p.Priority,
		IsActive: p.IsActive, RateLimitPerSec: p.RateLimitPerSec, CreatedAt: p.CreatedAt, UpdatedAt: p.UpdatedAt}
	if withSettings {
		if s, err := h.Providers.Redacted(p); err == nil {
			out.Settings = s
		}
	}
	return out
}

func (h *Handlers) listProviders(c *fiber.Ctx) error {
	rows, err := h.Providers.List(c.UserContext(), membership(c).Project.ID)
	if err != nil {
		return err
	}
	out := make([]ProviderResponse, 0, len(rows))
	for i := range rows {
		out = append(out, h.toProvider(&rows[i], false))
	}
	return httpx.OK(c, out)
}

func (h *Handlers) createProvider(c *fiber.Ctx) error {
	var req providerRequest
	if err := httpx.Bind(c, &req); err != nil {
		return err
	}
	p, err := h.Providers.Create(c.UserContext(), membership(c).Project.ID, providers.Input{
		Name: req.Name, Type: req.Type, Priority: req.Priority, IsActive: req.IsActive, Credentials: req.Credentials, RateLimitPerSec: req.RateLimitPerSec,
	})
	if err != nil {
		return err
	}
	return httpx.Created(c, h.toProvider(p, true))
}

func (h *Handlers) getProvider(c *fiber.Ctx) error {
	id, err := httpx.ParamUUID(c, "provider_id")
	if err != nil {
		return err
	}
	p, err := h.Providers.Get(c.UserContext(), membership(c).Project.ID, id)
	if err != nil {
		return err
	}
	return httpx.OK(c, h.toProvider(p, true))
}

func (h *Handlers) updateProvider(c *fiber.Ctx) error {
	id, err := httpx.ParamUUID(c, "provider_id")
	if err != nil {
		return err
	}
	req := providerRequest{IsUpdate: true}
	if err := httpx.Bind(c, &req); err != nil {
		return err
	}
	p, err := h.Providers.Update(c.UserContext(), membership(c).Project.ID, id, providers.Input{
		Name: req.Name, Priority: req.Priority, IsActive: req.IsActive, Credentials: req.Credentials, RateLimitPerSec: req.RateLimitPerSec,
	})
	if err != nil {
		return err
	}
	return httpx.OK(c, h.toProvider(p, true))
}

func (h *Handlers) deleteProvider(c *fiber.Ctx) error {
	id, err := httpx.ParamUUID(c, "provider_id")
	if err != nil {
		return err
	}
	if err := h.Providers.Delete(c.UserContext(), membership(c).Project.ID, id); err != nil {
		return err
	}
	return httpx.NoContent(c)
}

type testSendRequest struct {
	To      string `json:"to" validate:"required"`
	Text    string `json:"text" validate:"max=2000"`
	Subject string `json:"subject" validate:"max=200"`
}

// testSendProvider calls the provider directly (bypassing the queue) so an
// admin can check credentials from the panel. Result and raw response are
// returned; nothing is stored as a message.
func (h *Handlers) testSendProvider(c *fiber.Ctx) error {
	id, err := httpx.ParamUUID(c, "provider_id")
	if err != nil {
		return err
	}
	var req testSendRequest
	if err := httpx.Bind(c, &req); err != nil {
		return err
	}
	p, err := h.Providers.Get(c.UserContext(), membership(c).Project.ID, id)
	if err != nil {
		return err
	}
	adapter, err := h.Providers.Build(c.UserContext(), p)
	if err != nil {
		return domain.ErrValidation.WithMessage("provider configuration error: " + err.Error())
	}
	text := req.Text
	if text == "" {
		text = "Habarchy test message / Habarçy synag habary"
	}
	subject := req.Subject
	if subject == "" {
		subject = "Habarchy test"
	}
	start := time.Now()
	var res *ports.SendResult
	switch domain.Channel(p.Channel) {
	case domain.ChannelSMS:
		res, err = adapter.(ports.SMSProvider).Send(c.UserContext(), ports.SMSMessage{To: req.To, Text: text})
	case domain.ChannelEmail:
		res, err = adapter.(ports.EmailProvider).Send(c.UserContext(), ports.EmailMessage{To: req.To, Subject: subject, Text: text})
	case domain.ChannelTelegram:
		res, err = adapter.(ports.TelegramProvider).Send(c.UserContext(), ports.TelegramMessage{ChatID: req.To, Text: text})
	case domain.ChannelWhatsApp, domain.ChannelSlack:
		res, err = adapter.(ports.ChatProvider).Send(c.UserContext(), ports.ChatMessage{To: req.To, Text: text, Subject: subject})
	case domain.ChannelPush:
		var pr *ports.PushResult
		pr, err = adapter.(ports.PushProvider).Send(c.UserContext(), ports.PushMessage{Tokens: []string{req.To}, Title: subject, Body: text})
		if pr != nil {
			res = &pr.SendResult
		}
	}
	out := fiber.Map{"ok": err == nil, "duration_ms": time.Since(start).Milliseconds()}
	if res != nil {
		out["provider_message_id"], out["raw"] = res.ProviderMessageID, res.Raw
	}
	if err != nil {
		if pe, ok := err.(*ports.ProviderError); ok { //nolint:errorlint // direct type from adapters
			out["error"] = fiber.Map{"code": pe.Code, "message": pe.Message, "retryable": pe.Retryable, "raw": pe.Raw}
		} else {
			out["error"] = fiber.Map{"code": "error", "message": err.Error()}
		}
	}
	return httpx.OK(c, out)
}

// gatewayPairing returns what the Habarchy Gateway app needs to pair with
// an android_sms provider: the API URL and the key (admin role only).
func (h *Handlers) gatewayPairing(c *fiber.Ctx) error {
	id, err := httpx.ParamUUID(c, "provider_id")
	if err != nil {
		return err
	}
	p, err := h.Providers.Get(c.UserContext(), membership(c).Project.ID, id)
	if err != nil {
		return err
	}
	if domain.ProviderType(p.Type) != domain.ProviderAndroidSMS {
		return domain.ErrValidation.WithMessage("pairing is only available for android_sms providers")
	}
	raw, err := h.Providers.Decrypt(p)
	if err != nil {
		return err
	}
	var creds struct {
		GatewayKey string `json:"gateway_key"`
		SimSlot    int    `json:"sim_slot"`
	}
	if err := json.Unmarshal(raw, &creds); err != nil {
		return err
	}
	var dev *sqlcgen.GatewayDevice
	if d, err := h.DB.Queries.GetGatewayDevice(c.UserContext(), p.ID); err == nil {
		dev = &d
	}
	q := url.Values{"url": {h.PublicURL}, "key": {creds.GatewayKey}, "name": {p.Name}}
	out := fiber.Map{
		"provider_id": p.ID, "name": p.Name, "api_url": h.PublicURL, "gateway_key": creds.GatewayKey, "sim_slot": creds.SimSlot,
		"qr": "habarchy://gateway?" + q.Encode(),
	}
	if dev != nil {
		out["last_seen_at"] = dev.LastSeenAt
		out["online"] = dev.LastSeenAt != nil && time.Since(*dev.LastSeenAt) <= 90*time.Second
		out["device_info"] = dev.DeviceInfo
	}
	return httpx.OK(c, out)
}

// listInbound lists SMS received by the project's gateway phones.
func (h *Handlers) listInbound(c *fiber.Ctx) error {
	pid := membership(c).Project.ID
	page := httpx.ParsePage(c, 50, 200)
	rows, err := h.DB.Queries.ListGatewayInbound(c.UserContext(), sqlcgen.ListGatewayInboundParams{ProjectID: pid, RowLimit: page.Limit, RowOffset: page.Offset})
	if err != nil {
		return err
	}
	total, _ := h.DB.Queries.CountGatewayInbound(c.UserContext(), pid)
	return httpx.JSON(c, fiber.StatusOK, rows, fiber.Map{"total": total, "limit": page.Limit, "offset": page.Offset})
}
