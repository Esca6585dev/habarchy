// Package public serves /api/v1, the API client applications call with an
// API key.
package public

import (
	"time"

	"github.com/gofiber/fiber/v2"

	httpx "github.com/Esca6585dev/habarchy/backend/internal/adapters/http"
	"github.com/Esca6585dev/habarchy/backend/internal/adapters/http/admin"
	"github.com/Esca6585dev/habarchy/backend/internal/adapters/http/middleware"
	"github.com/Esca6585dev/habarchy/backend/internal/app/projects"
	"github.com/Esca6585dev/habarchy/backend/internal/app/templates"
	"github.com/Esca6585dev/habarchy/backend/internal/domain"
)

// Handlers groups the public API dependencies.
type Handlers struct {
	Projects           *projects.Service
	Templates          *templates.Service
	SignatureTolerance time.Duration
}

// Register mounts the public API under /api/v1.
func (h *Handlers) Register(app fiber.Router) {
	r := app.Group("/api/v1", middleware.RequireAPIKey(h.Projects, middleware.APIKeyConfig{SignatureTolerance: h.SignatureTolerance}))

	r.Get("/me", func(c *fiber.Ctx) error {
		caller := middleware.Caller(c)
		return httpx.OK(c, fiber.Map{
			"project": fiber.Map{"id": caller.Project.ID, "name": caller.Project.Name, "slug": caller.Project.Slug},
			"key":     fiber.Map{"name": caller.Key.Name, "prefix": caller.Key.Prefix, "scopes": caller.Key.Scopes, "test": caller.IsTest()},
		})
	})

	t := r.Group("/templates", middleware.RequireScope(domain.ScopeTemplates))
	t.Get("/", h.listTemplates)
	t.Post("/", h.createTemplate)
	t.Post("/preview", h.previewTemplate)
	t.Get("/:template_id", h.getTemplate)
	t.Put("/:template_id", h.updateTemplate)
	t.Delete("/:template_id", h.deleteTemplate)
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
