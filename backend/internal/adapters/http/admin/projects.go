package admin

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"

	httpx "github.com/Esca6585dev/habarchy/backend/internal/adapters/http"
	"github.com/Esca6585dev/habarchy/backend/internal/adapters/http/middleware"
	"github.com/Esca6585dev/habarchy/backend/internal/app/projects"
	"github.com/Esca6585dev/habarchy/backend/internal/domain"
)

const localMembership = "membership"

// requireRole authorizes the :project_id route group and caches the
// membership in locals.
func (h *Handlers) requireRole(minRole domain.Role) fiber.Handler {
	return func(c *fiber.Ctx) error {
		pid, err := httpx.ParamUUID(c, "project_id")
		if err != nil {
			return err
		}
		m, err := h.Projects.Authorize(c.UserContext(), middleware.UserID(c), pid, minRole)
		if err != nil {
			return err
		}
		c.Locals(localMembership, m)
		return c.Next()
	}
}

func membership(c *fiber.Ctx) *projects.Membership {
	m, _ := c.Locals(localMembership).(*projects.Membership)
	return m
}

func projectID(c *fiber.Ctx) (uuid.UUID, error) { return httpx.ParamUUID(c, "project_id") }

func (h *Handlers) listProjects(c *fiber.Ctx) error {
	list, err := h.Projects.ListForUser(c.UserContext(), middleware.UserID(c))
	if err != nil {
		return err
	}
	out := make([]ProjectResponse, 0, len(list))
	for i := range list {
		out = append(out, toProject(&list[i].Project, list[i].Role))
	}
	return httpx.OK(c, out)
}

func (h *Handlers) createProject(c *fiber.Ctx) error {
	var req createProjectRequest
	if err := httpx.Bind(c, &req); err != nil {
		return err
	}
	p, err := h.Projects.Create(c.UserContext(), middleware.UserID(c), projects.CreateInput{
		Name: req.Name, Slug: req.Slug, DailyQuota: req.DailyQuota, MonthlyQuota: req.MonthlyQuota, DefaultLocale: req.DefaultLocale,
	})
	if err != nil {
		return err
	}
	return httpx.Created(c, toProject(p, domain.RoleOwner))
}

func (h *Handlers) getProject(c *fiber.Ctx) error {
	pid, err := projectID(c)
	if err != nil {
		return err
	}
	m, err := h.Projects.Authorize(c.UserContext(), middleware.UserID(c), pid, domain.RoleViewer)
	if err != nil {
		return err
	}
	return httpx.OK(c, toProject(&m.Project, m.Role))
}

func (h *Handlers) updateProject(c *fiber.Ctx) error {
	pid, err := projectID(c)
	if err != nil {
		return err
	}
	var req updateProjectRequest
	if err := httpx.Bind(c, &req); err != nil {
		return err
	}
	p, err := h.Projects.Update(c.UserContext(), middleware.UserID(c), pid, projects.UpdateInput{
		Name: req.Name, Status: req.Status, DailyQuota: req.DailyQuota, MonthlyQuota: req.MonthlyQuota,
		WebhookURL: req.WebhookURL, WebhookSecret: req.WebhookSecret, DefaultLocale: req.DefaultLocale,
		AllowedIPs: req.AllowedIPs, AutoChannelOrder: req.AutoChannelOrder,
	})
	if err != nil {
		return err
	}
	return httpx.OK(c, toProject(p, ""))
}

func (h *Handlers) deleteProject(c *fiber.Ctx) error {
	pid, err := projectID(c)
	if err != nil {
		return err
	}
	if err := h.Projects.Delete(c.UserContext(), middleware.UserID(c), pid); err != nil {
		return err
	}
	return httpx.NoContent(c)
}

func (h *Handlers) listMembers(c *fiber.Ctx) error {
	pid, err := projectID(c)
	if err != nil {
		return err
	}
	list, err := h.Projects.ListMembers(c.UserContext(), middleware.UserID(c), pid)
	if err != nil {
		return err
	}
	return httpx.OK(c, list)
}

func (h *Handlers) setMember(c *fiber.Ctx) error {
	pid, err := projectID(c)
	if err != nil {
		return err
	}
	var req setMemberRequest
	if err := httpx.Bind(c, &req); err != nil {
		return err
	}
	m, err := h.Projects.SetMember(c.UserContext(), middleware.UserID(c), pid, req.Email, req.Role)
	if err != nil {
		return err
	}
	return httpx.OK(c, m)
}

func (h *Handlers) removeMember(c *fiber.Ctx) error {
	pid, err := projectID(c)
	if err != nil {
		return err
	}
	uid, err := httpx.ParamUUID(c, "user_id")
	if err != nil {
		return err
	}
	if err := h.Projects.RemoveMember(c.UserContext(), middleware.UserID(c), pid, uid); err != nil {
		return err
	}
	return httpx.NoContent(c)
}

func (h *Handlers) listAPIKeys(c *fiber.Ctx) error {
	pid, err := projectID(c)
	if err != nil {
		return err
	}
	keys, err := h.Projects.ListAPIKeys(c.UserContext(), middleware.UserID(c), pid)
	if err != nil {
		return err
	}
	out := make([]APIKeyResponse, 0, len(keys))
	for i := range keys {
		out = append(out, toAPIKey(&keys[i], ""))
	}
	return httpx.OK(c, out)
}

func (h *Handlers) createAPIKey(c *fiber.Ctx) error {
	pid, err := projectID(c)
	if err != nil {
		return err
	}
	var req createAPIKeyRequest
	if err := httpx.Bind(c, &req); err != nil {
		return err
	}
	created, err := h.Projects.CreateAPIKey(c.UserContext(), middleware.UserID(c), pid, projects.APIKeyInput{
		Name: req.Name, Live: req.Live, Scopes: req.Scopes, IPAllowlist: req.IPAllowlist, ExpiresAt: req.ExpiresAt, RequireSignature: req.RequireSignature,
	})
	if err != nil {
		return err
	}
	// The plaintext key appears in this response only.
	return httpx.JSON(c, fiber.StatusCreated, toAPIKey(&created.Key, created.Plaintext), fiber.Map{"notice": "store this key now; it cannot be shown again"})
}

func (h *Handlers) revokeAPIKey(c *fiber.Ctx) error {
	pid, err := projectID(c)
	if err != nil {
		return err
	}
	kid, err := httpx.ParamUUID(c, "key_id")
	if err != nil {
		return err
	}
	if err := h.Projects.RevokeAPIKey(c.UserContext(), middleware.UserID(c), pid, kid); err != nil {
		return err
	}
	return httpx.NoContent(c)
}
