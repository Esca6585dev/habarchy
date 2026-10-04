package admin

import (
	"strconv"

	"github.com/gofiber/fiber/v2"

	httpx "github.com/Esca6585dev/habarchy/backend/internal/adapters/http"
	"github.com/Esca6585dev/habarchy/backend/internal/app/templates"
	"github.com/Esca6585dev/habarchy/backend/internal/domain"
)

func (h *Handlers) listTemplates(c *fiber.Ctx) error {
	list, err := h.Templates.List(c.UserContext(), membership(c).Project.ID, c.Query("key"))
	if err != nil {
		return err
	}
	out := make([]TemplateResponse, 0, len(list))
	for i := range list {
		out = append(out, ToTemplate(&list[i]))
	}
	return httpx.OK(c, out)
}

func (h *Handlers) createTemplate(c *fiber.Ctx) error {
	var req TemplateRequest
	if err := httpx.Bind(c, &req); err != nil {
		return err
	}
	tpl, err := h.Templates.Create(c.UserContext(), membership(c).Project.ID, templates.Input{
		Key: req.Key, Channel: req.Channel, Locale: req.Locale, Subject: req.Subject, Body: req.Body,
	})
	if err != nil {
		return err
	}
	return httpx.Created(c, ToTemplate(tpl))
}

func (h *Handlers) getTemplate(c *fiber.Ctx) error {
	id, err := httpx.ParamUUID(c, "template_id")
	if err != nil {
		return err
	}
	tpl, err := h.Templates.Get(c.UserContext(), membership(c).Project.ID, id)
	if err != nil {
		return err
	}
	return httpx.OK(c, ToTemplate(tpl))
}

func (h *Handlers) updateTemplate(c *fiber.Ctx) error {
	id, err := httpx.ParamUUID(c, "template_id")
	if err != nil {
		return err
	}
	req := TemplateRequest{IsUpdate: true}
	if err := httpx.Bind(c, &req); err != nil {
		return err
	}
	tpl, err := h.Templates.Update(c.UserContext(), membership(c).Project.ID, id, templates.Input{
		Subject: req.Subject, Body: req.Body, IsActive: req.IsActive,
	})
	if err != nil {
		return err
	}
	return httpx.OK(c, ToTemplate(tpl))
}

func (h *Handlers) deleteTemplate(c *fiber.Ctx) error {
	id, err := httpx.ParamUUID(c, "template_id")
	if err != nil {
		return err
	}
	if err := h.Templates.Delete(c.UserContext(), membership(c).Project.ID, id); err != nil {
		return err
	}
	return httpx.NoContent(c)
}

func (h *Handlers) templateVersions(c *fiber.Ctx) error {
	id, err := httpx.ParamUUID(c, "template_id")
	if err != nil {
		return err
	}
	rows, err := h.Templates.Versions(c.UserContext(), membership(c).Project.ID, id)
	if err != nil {
		return err
	}
	return httpx.OK(c, ToTemplateVersions(rows))
}

func (h *Handlers) restoreTemplateVersion(c *fiber.Ctx) error {
	id, err := httpx.ParamUUID(c, "template_id")
	if err != nil {
		return err
	}
	v, err := strconv.Atoi(c.Params("version"))
	if err != nil || v < 1 {
		return domain.ErrValidation.WithDetails(map[string]any{"version": "must be a positive integer"})
	}
	tpl, err := h.Templates.Restore(c.UserContext(), membership(c).Project.ID, id, int32(v)) //nolint:gosec // validated above
	if err != nil {
		return err
	}
	return httpx.OK(c, ToTemplate(tpl))
}

func (h *Handlers) previewTemplate(c *fiber.Ctx) error {
	var req PreviewRequest
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

func (h *Handlers) previewStoredTemplate(c *fiber.Ctx) error {
	id, err := httpx.ParamUUID(c, "template_id")
	if err != nil {
		return err
	}
	var req PreviewRequest
	if len(c.Body()) > 0 {
		if err := httpx.Bind(c, &req); err != nil {
			return err
		}
	}
	p, err := h.Templates.PreviewStored(c.UserContext(), membership(c).Project.ID, id, req.Data)
	if err != nil {
		return err
	}
	return httpx.OK(c, p)
}
