package admin

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"

	httpx "github.com/Esca6585dev/habarchy/backend/internal/adapters/http"
	"github.com/Esca6585dev/habarchy/backend/internal/adapters/postgres/sqlcgen"
	"github.com/Esca6585dev/habarchy/backend/internal/domain"
)

// Platform endpoints give the owner (platform_admin) a cross-tenant view:
// every tenant, every user and the whole message log — who sent what, to
// whom. They are mounted behind RequirePlatformAdmin.

// platformTenants lists every tenant (project).
func (h *Handlers) platformTenants(c *fiber.Ctx) error {
	page := httpx.ParsePage(c, 50, 200)
	rows, err := h.DB.Queries.ListProjects(c.UserContext(), sqlcgen.ListProjectsParams{RowLimit: page.Limit, RowOffset: page.Offset})
	if err != nil {
		return err
	}
	out := make([]fiber.Map, 0, len(rows))
	for i := range rows {
		p := &rows[i]
		out = append(out, fiber.Map{
			"id": p.ID, "name": p.Name, "slug": p.Slug, "status": p.Status,
			"daily_quota": p.DailyQuota, "monthly_quota": p.MonthlyQuota,
			"default_locale": p.DefaultLocale, "created_at": p.CreatedAt,
		})
	}
	total, _ := h.DB.Queries.CountProjects(c.UserContext())
	return httpx.JSON(c, fiber.StatusOK, out, fiber.Map{"total": total})
}

// platformUsers lists every user — admins and phone-login clients.
func (h *Handlers) platformUsers(c *fiber.Ctx) error {
	page := httpx.ParsePage(c, 50, 200)
	rows, err := h.DB.Queries.ListUsers(c.UserContext(), sqlcgen.ListUsersParams{RowLimit: page.Limit, RowOffset: page.Offset})
	if err != nil {
		return err
	}
	out := make([]UserResponse, 0, len(rows))
	for i := range rows {
		out = append(out, toUser(&rows[i]))
	}
	total, _ := h.DB.Queries.CountUsers(c.UserContext())
	return httpx.JSON(c, fiber.StatusOK, out, fiber.Map{"total": total})
}

// platformMessages is the cross-tenant message log: who sent what, to whom,
// through which tenant. Keyset pagination on message id, newest first.
func (h *Handlers) platformMessages(c *fiber.Ctx) error {
	limit := int32(c.QueryInt("limit", 50)) //nolint:gosec // bounded below
	if limit <= 0 || limit > 200 {
		limit = 50
	}
	params := sqlcgen.ListAllMessagesParams{RowLimit: limit + 1}
	details := map[string]any{}
	if s := c.Query("status"); s != "" {
		if !domain.MessageStatus(s).Valid() {
			details["status"] = "unknown status"
		}
		params.Status = sqlcgen.NullMessageStatus{MessageStatus: sqlcgen.MessageStatus(s), Valid: true}
	}
	if s := c.Query("channel"); s != "" {
		if !domain.Channel(s).Valid() {
			details["channel"] = "unknown channel"
		}
		params.Channel = sqlcgen.NullChannel{Channel: sqlcgen.Channel(s), Valid: true}
	}
	if s := c.Query("cursor"); s != "" {
		id, err := uuid.Parse(s)
		if err != nil {
			details["cursor"] = "must be a uuid"
		}
		params.CursorID = &id
	}
	if len(details) > 0 {
		return domain.ErrValidation.WithDetails(details)
	}
	rows, err := h.DB.Queries.ListAllMessages(c.UserContext(), params)
	if err != nil {
		return err
	}
	var next *uuid.UUID
	if len(rows) > int(limit) {
		rows = rows[:limit]
		id := rows[len(rows)-1].Message.ID
		next = &id
	}
	out := make([]fiber.Map, 0, len(rows))
	for i := range rows {
		m := toMessageRow(&rows[i].Message)
		out = append(out, fiber.Map{
			"message": m,
			"tenant":  fiber.Map{"id": rows[i].Message.ProjectID, "name": rows[i].ProjectName, "slug": rows[i].ProjectSlug},
		})
	}
	return httpx.JSON(c, fiber.StatusOK, out, fiber.Map{"next_cursor": next, "limit": limit})
}
