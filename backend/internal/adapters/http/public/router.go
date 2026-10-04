// Package public serves /api/v1, the API client applications call with an
// API key, plus provider callbacks.
package public

import (
	"time"

	"github.com/gofiber/fiber/v2"

	httpx "github.com/Esca6585dev/habarchy/backend/internal/adapters/http"
	"github.com/Esca6585dev/habarchy/backend/internal/adapters/http/admin"
	"github.com/Esca6585dev/habarchy/backend/internal/adapters/http/middleware"
	"github.com/Esca6585dev/habarchy/backend/internal/app/contacts"
	"github.com/Esca6585dev/habarchy/backend/internal/app/delivery"
	"github.com/Esca6585dev/habarchy/backend/internal/app/messages"
	"github.com/Esca6585dev/habarchy/backend/internal/app/otp"
	"github.com/Esca6585dev/habarchy/backend/internal/app/projects"
	"github.com/Esca6585dev/habarchy/backend/internal/app/providers"
	"github.com/Esca6585dev/habarchy/backend/internal/app/stats"
	"github.com/Esca6585dev/habarchy/backend/internal/app/templates"
	"github.com/Esca6585dev/habarchy/backend/internal/domain"
)

// Handlers groups the public API dependencies. Messages, OTP, Contacts,
// Providers and Delivery may be nil in tests that only exercise templates.
type Handlers struct {
	Projects           *projects.Service
	Templates          *templates.Service
	Messages           *messages.Service
	OTP                *otp.Service
	Contacts           *contacts.Service
	Providers          *providers.Service
	Delivery           *delivery.Service
	Stats              *stats.Service
	SignatureTolerance time.Duration
	// APIRatePerSec limits requests per API key (token bucket); 0 = off.
	APIRatePerSec float64
	Limiter       middleware.Limiter
}

// Register mounts the public API under /api/v1 and callbacks under
// /callbacks.
func (h *Handlers) Register(app fiber.Router) {
	r := app.Group("/api/v1", middleware.RequireAPIKey(h.Projects, middleware.APIKeyConfig{SignatureTolerance: h.SignatureTolerance}))
	if h.Limiter != nil && h.APIRatePerSec > 0 {
		r.Use(middleware.RateLimitByKey(h.Limiter, h.APIRatePerSec))
	}

	r.Get("/me", func(c *fiber.Ctx) error {
		k := middleware.Caller(c)
		return httpx.OK(c, fiber.Map{
			"project": fiber.Map{"id": k.Project.ID, "name": k.Project.Name, "slug": k.Project.Slug, "default_locale": k.Project.DefaultLocale},
			"key":     fiber.Map{"name": k.Key.Name, "prefix": k.Key.Prefix, "scopes": k.Key.Scopes, "test": k.IsTest()},
		})
	})

	t := r.Group("/templates", middleware.RequireScope(domain.ScopeTemplates))
	t.Get("/", h.listTemplates)
	t.Post("/", h.createTemplate)
	t.Post("/preview", h.previewTemplate)
	t.Get("/:template_id", h.getTemplate)
	t.Put("/:template_id", h.updateTemplate)
	t.Delete("/:template_id", h.deleteTemplate)

	if h.Messages != nil {
		send := middleware.RequireScope(domain.ScopeMessagesSend)
		read := middleware.RequireScope(domain.ScopeMessagesRead)
		r.Post("/messages", send, h.sendMessage)
		r.Post("/messages/batch", send, h.sendBatch)
		r.Get("/messages", read, h.listMessages)
		r.Get("/messages/:message_id", read, h.getMessage)
		r.Post("/messages/:message_id/cancel", send, h.cancelMessage)
		r.Get("/batches/:batch_id", read, h.getBatch)
	}
	if h.OTP != nil {
		o := r.Group("/otp", middleware.RequireScope(domain.ScopeOTP))
		o.Post("/send", h.otpSend)
		o.Post("/verify", h.otpVerify)
	}
	if h.Contacts != nil {
		ct := r.Group("/contacts", middleware.RequireScope(domain.ScopeContacts))
		ct.Get("/", h.listContacts)
		ct.Post("/", h.createContact)
		ct.Get("/:contact_id", h.getContact)
		ct.Put("/:contact_id", h.updateContact)
		ct.Delete("/:contact_id", h.deleteContact)

		d := r.Group("/devices", middleware.RequireScope(domain.ScopeDevices))
		d.Get("/", h.listDevices)
		d.Post("/", h.registerDevice)
		d.Delete("/:token", h.removeDevice)
	}
	if h.Stats != nil {
		r.Get("/usage", middleware.RequireScope(domain.ScopeUsage), h.usage)
	}
	if h.Delivery != nil && h.Providers != nil {
		app.Post("/callbacks/sms/:provider_id", h.smsCallback)
		app.Get("/callbacks/sms/:provider_id", h.smsCallback)
	}
}

func (h *Handlers) listTemplates(c *fiber.Ctx) error {
	list, err := h.Templates.List(c.UserContext(), middleware.Caller(c).Project.ID, c.Query("key"))
	if err != nil {
		return err
	}
	out := make([]admin.TemplateResponse, 0, len(list))
	for i := range list {
		out = append(out, admin.ToTemplate(&list[i]))
	}
	return httpx.OK(c, out)
}

func (h *Handlers) createTemplate(c *fiber.Ctx) error {
	var req admin.TemplateRequest
	if err := httpx.Bind(c, &req); err != nil {
		return err
	}
	tpl, err := h.Templates.Create(c.UserContext(), middleware.Caller(c).Project.ID, templates.Input{
		Key: req.Key, Channel: req.Channel, Locale: req.Locale, Subject: req.Subject, Body: req.Body,
	})
	if err != nil {
		return err
	}
	return httpx.Created(c, admin.ToTemplate(tpl))
}

func (h *Handlers) getTemplate(c *fiber.Ctx) error {
	id, err := httpx.ParamUUID(c, "template_id")
	if err != nil {
		return err
	}
	tpl, err := h.Templates.Get(c.UserContext(), middleware.Caller(c).Project.ID, id)
	if err != nil {
		return err
	}
	return httpx.OK(c, admin.ToTemplate(tpl))
}

func (h *Handlers) updateTemplate(c *fiber.Ctx) error {
	id, err := httpx.ParamUUID(c, "template_id")
	if err != nil {
		return err
	}
	req := admin.TemplateRequest{IsUpdate: true}
	if err := httpx.Bind(c, &req); err != nil {
		return err
	}
	tpl, err := h.Templates.Update(c.UserContext(), middleware.Caller(c).Project.ID, id, templates.Input{Subject: req.Subject, Body: req.Body, IsActive: req.IsActive})
	if err != nil {
		return err
	}
	return httpx.OK(c, admin.ToTemplate(tpl))
}

func (h *Handlers) deleteTemplate(c *fiber.Ctx) error {
	id, err := httpx.ParamUUID(c, "template_id")
	if err != nil {
		return err
	}
	if err := h.Templates.Delete(c.UserContext(), middleware.Caller(c).Project.ID, id); err != nil {
		return err
	}
	return httpx.NoContent(c)
}

func (h *Handlers) previewTemplate(c *fiber.Ctx) error {
	var req admin.PreviewRequest
	if err := httpx.Bind(c, &req); err != nil {
		return err
	}
	if req.Channel == "" {
		req.Channel = domain.ChannelSMS
	}
	p, err := templates.PreviewRaw(req.Channel, req.Subject, req.Body, req.Data)
	if err != nil {
		return err
	}
	return httpx.OK(c, p)
}

// usage serves GET /usage?from=&to=&group_by=channel|day for API clients.
func (h *Handlers) usage(c *fiber.Ctx) error {
	from, to := time.Now().UTC().AddDate(0, 0, -30), time.Now().UTC()
	details := map[string]any{}
	if s := c.Query("from"); s != "" {
		t, err := time.Parse("2006-01-02", s)
		if err != nil {
			details["from"] = "YYYY-MM-DD"
		}
		from = t
	}
	if s := c.Query("to"); s != "" {
		t, err := time.Parse("2006-01-02", s)
		if err != nil {
			details["to"] = "YYYY-MM-DD"
		}
		to = t
	}
	if len(details) > 0 {
		return domain.ErrValidation.WithDetails(details)
	}
	rows, err := h.Stats.Usage(c.UserContext(), caller(c).Project.ID, from, to, c.Query("group_by"))
	if err != nil {
		return err
	}
	return httpx.JSON(c, fiber.StatusOK, rows, fiber.Map{"from": from.Format("2006-01-02"), "to": to.Format("2006-01-02"), "group_by": c.Query("group_by")})
}
