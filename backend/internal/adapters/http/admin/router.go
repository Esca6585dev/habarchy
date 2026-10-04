package admin

import (
	"net/http"

	"github.com/gofiber/adaptor/v2"
	"github.com/gofiber/fiber/v2"
	goredis "github.com/redis/go-redis/v9"

	"github.com/Esca6585dev/habarchy/backend/internal/adapters/http/middleware"
	"github.com/Esca6585dev/habarchy/backend/internal/adapters/postgres"
	"github.com/Esca6585dev/habarchy/backend/internal/app/auth"
	"github.com/Esca6585dev/habarchy/backend/internal/app/messages"
	"github.com/Esca6585dev/habarchy/backend/internal/app/projects"
	"github.com/Esca6585dev/habarchy/backend/internal/app/providers"
	"github.com/Esca6585dev/habarchy/backend/internal/app/stats"
	"github.com/Esca6585dev/habarchy/backend/internal/app/templates"
	"github.com/Esca6585dev/habarchy/backend/internal/app/webhooks"
	"github.com/Esca6585dev/habarchy/backend/internal/domain"
)

// Handlers groups the admin API dependencies. Optional fields may be nil
// in tests; the routes that need them are then not mounted.
type Handlers struct {
	Auth      *auth.Service
	Projects  *projects.Service
	Templates *templates.Service
	Providers *providers.Service
	Messages  *messages.Service
	Webhooks  *webhooks.Service
	Stats     *stats.Service
	DB        *postgres.DB
	Redis     *goredis.Client // live events (SSE)
	// QueueUI is the asynqmon handler, mounted at /admin/queues behind
	// admin auth when set.
	QueueUI http.Handler
}

// Register mounts the admin API under /api/admin.
func (h *Handlers) Register(app fiber.Router) {
	r := app.Group("/api/admin")

	r.Post("/auth/login", h.login)
	r.Post("/auth/refresh", h.refresh)
	r.Post("/auth/logout", h.logout)

	// SSE and the queue UI accept the token from a query param / cookie.
	if h.Redis != nil && h.Projects != nil {
		r.Get("/stream", middleware.RequireJWTFlexible(h.Auth), h.stream)
	}
	if h.QueueUI != nil {
		app.Use("/admin/queues", middleware.RequireJWTFlexible(h.Auth), adaptor.HTTPHandler(h.QueueUI))
	}

	// Everything below requires a valid access token.
	r.Use(middleware.RequireJWT(h.Auth))
	if h.DB != nil {
		r.Get("/users", h.listUsers)
		r.Post("/users", h.createUser)
		r.Get("/overview", h.overview)
	}

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

	if h.Messages != nil && h.DB != nil {
		m := p.Group("/messages", h.requireRole(domain.RoleViewer))
		m.Get("/", h.listMessages)
		m.Get("/:message_id", h.getMessage)
		m.Post("/:message_id/resend", h.requireRole(domain.RoleDeveloper), h.resendMessage)
		m.Post("/:message_id/cancel", h.requireRole(domain.RoleDeveloper), h.cancelMessage)
		b := p.Group("/batches", h.requireRole(domain.RoleViewer))
		b.Get("/", h.listBatches)
		b.Get("/:batch_id", h.getBatch)
	}
	if h.Webhooks != nil && h.DB != nil {
		w := p.Group("/webhooks", h.requireRole(domain.RoleViewer))
		w.Get("/", h.listWebhooks)
		w.Get("/:delivery_id", h.getWebhook)
		w.Post("/:delivery_id/resend", h.requireRole(domain.RoleDeveloper), h.resendWebhook)
	}
	if h.Stats != nil && h.DB != nil {
		v := p.Group("/", h.requireRole(domain.RoleViewer))
		v.Get("/dashboard", h.dashboard)
		v.Get("/usage", h.usage)
		v.Get("/health", h.health)
		v.Get("/contacts", h.listContacts)
		v.Get("/devices", h.listDevices)
		p.Get("/audit-logs", h.requireRole(domain.RoleAdmin), h.auditLogs)
	}

	if h.Providers != nil {
		pv := p.Group("/providers", h.requireRole(domain.RoleAdmin))
		pv.Get("/", h.listProviders)
		pv.Post("/", h.createProvider)
		pv.Get("/:provider_id", h.getProvider)
		pv.Put("/:provider_id", h.updateProvider)
		pv.Delete("/:provider_id", h.deleteProvider)
		pv.Post("/:provider_id/test", h.testSendProvider)
	}
}
