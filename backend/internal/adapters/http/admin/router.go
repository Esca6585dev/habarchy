package admin

import (
	"github.com/gofiber/fiber/v2"

	"github.com/Esca6585dev/habarchy/backend/internal/adapters/http/middleware"
	"github.com/Esca6585dev/habarchy/backend/internal/app/auth"
	"github.com/Esca6585dev/habarchy/backend/internal/app/projects"
	"github.com/Esca6585dev/habarchy/backend/internal/app/templates"
	"github.com/Esca6585dev/habarchy/backend/internal/domain"
)

// Handlers groups the admin API dependencies.
type Handlers struct {
	Auth      *auth.Service
	Projects  *projects.Service
	Templates *templates.Service
}

// Register mounts the admin API under /api/admin.
func (h *Handlers) Register(app fiber.Router) {
	r := app.Group("/api/admin")

	r.Post("/auth/login", h.login)
	r.Post("/auth/refresh", h.refresh)
	r.Post("/auth/logout", h.logout)

	// Everything below requires a valid access token.
	r.Use(middleware.RequireJWT(h.Auth))

	r.Get("/me", h.me)
	r.Put("/me/password", h.changePassword)
	r.Post("/me/totp/setup", h.totpSetup)
	r.Post("/me/totp/confirm", h.totpConfirm)
	r.Post("/me/totp/disable", h.totpDisable)
	r.Post("/auth/logout-all", h.logoutAll)

	r.Get("/projects", h.listProjects)
	r.Post("/projects", h.createProject)
	p := r.Group("/projects/:project_id")
	p.Get("/", h.getProject)
	p.Patch("/", h.updateProject)
	p.Delete("/", h.deleteProject)

	p.Get("/members", h.listMembers)
	p.Put("/members", h.setMember)
	p.Delete("/members/:user_id", h.removeMember)

	p.Get("/api-keys", h.listAPIKeys)
	p.Post("/api-keys", h.createAPIKey)
	p.Delete("/api-keys/:key_id", h.revokeAPIKey)

	t := p.Group("/templates", h.requireRole(domain.RoleDeveloper))
	t.Get("/", h.listTemplates)
	t.Post("/", h.createTemplate)
	t.Post("/preview", h.previewTemplate)
	t.Get("/:template_id", h.getTemplate)
	t.Put("/:template_id", h.updateTemplate)
	t.Delete("/:template_id", h.deleteTemplate)
	t.Get("/:template_id/versions", h.templateVersions)
	t.Post("/:template_id/versions/:version/restore", h.restoreTemplateVersion)
	t.Post("/:template_id/preview", h.previewStoredTemplate)
}
