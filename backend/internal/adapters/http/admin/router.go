package admin

import (
	"net/http"

	"github.com/gofiber/adaptor/v2"
	"github.com/gofiber/fiber/v2"
	goredis "github.com/redis/go-redis/v9"

	"github.com/Esca6585dev/habarchy/backend/internal/adapters/http/middleware"
	"github.com/Esca6585dev/habarchy/backend/internal/adapters/postgres"
	"github.com/Esca6585dev/habarchy/backend/internal/app/auth"
	"github.com/Esca6585dev/habarchy/backend/internal/app/contactimport"
	"github.com/Esca6585dev/habarchy/backend/internal/app/contacts"
	"github.com/Esca6585dev/habarchy/backend/internal/app/groups"
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
	Contacts  *contacts.Service
	Groups    *groups.Service
	Imports   *contactimport.Service
	Messages  *messages.Service
	Webhooks  *webhooks.Service
	Stats     *stats.Service
	DB        *postgres.DB
	Redis     *goredis.Client // live events (SSE)
	// QueueUI is the asynqmon handler, mounted at /admin/queues behind
	// admin auth when set.
	QueueUI http.Handler
	// PublicURL is the API base URL phones should use (gateway pairing).
	PublicURL string
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

	if h.Imports != nil && h.Imports.Google != nil {
		// OAuth redirect target: authenticated by the signed state, not a JWT.
		r.Get("/integrations/google/callback", h.googleCallback)
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
		m.Post("/send", h.requireRole(domain.RoleDeveloper), h.sendMessages)
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
		if h.Contacts != nil {
			ct := p.Group("/contacts", h.requireRole(domain.RoleDeveloper))
			if h.Imports != nil {
				ct.Post("/import", h.importContacts)
				ct.Post("/import/carddav", h.importCardDAV)
				ct.Get("/import/google/url", h.googleImportURL)
			}
			ct.Post("/", h.createContact)
			ct.Get("/:contact_id", h.requireRole(domain.RoleViewer), h.getContact)
			ct.Put("/:contact_id", h.updateContact)
			ct.Delete("/:contact_id", h.deleteContact)
		}
		if h.Groups != nil {
			g := p.Group("/groups", h.requireRole(domain.RoleViewer))
			g.Get("/", h.listGroups)
			g.Post("/", h.requireRole(domain.RoleDeveloper), h.createGroup)
			g.Get("/:group_id", h.getGroup)
			g.Put("/:group_id", h.requireRole(domain.RoleDeveloper), h.updateGroup)
			g.Delete("/:group_id", h.requireRole(domain.RoleDeveloper), h.deleteGroup)
			g.Get("/:group_id/members", h.listGroupMembers)
			g.Post("/:group_id/members", h.requireRole(domain.RoleDeveloper), h.addGroupMembers)
			g.Delete("/:group_id/members/:contact_id", h.requireRole(domain.RoleDeveloper), h.removeGroupMember)
		}
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
		pv.Get("/:provider_id/pairing", h.gatewayPairing)
		p.Get("/inbound", h.requireRole(domain.RoleViewer), h.listInbound)
	}
}
