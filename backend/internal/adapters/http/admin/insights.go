package admin

import (
	"bufio"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/valyala/fasthttp"

	httpx "github.com/Esca6585dev/habarchy/backend/internal/adapters/http"
	"github.com/Esca6585dev/habarchy/backend/internal/adapters/http/middleware"
	"github.com/Esca6585dev/habarchy/backend/internal/adapters/postgres/sqlcgen"
	"github.com/Esca6585dev/habarchy/backend/internal/app/events"
	"github.com/Esca6585dev/habarchy/backend/internal/app/stats"
	"github.com/Esca6585dev/habarchy/backend/internal/domain"
)

// parseWindow reads ?from=&to= (RFC3339 or YYYY-MM-DD) or ?days=N.
func parseWindow(c *fiber.Ctx, defaultDays int) (stats.Window, error) {
	if c.Query("from") == "" && c.Query("to") == "" {
		days := c.QueryInt("days", defaultDays)
		if days < 1 || days > 366 {
			days = defaultDays
		}
		return stats.LastDays(days), nil
	}
	parse := func(s string, end bool) (time.Time, error) {
		if t, err := time.Parse(time.RFC3339, s); err == nil {
			return t, nil
		}
		t, err := time.Parse("2006-01-02", s)
		if err != nil {
			return t, err
		}
		if end {
			t = t.Add(24 * time.Hour)
		}
		return t, nil
	}
	w := stats.LastDays(defaultDays)
	var err error
	if s := c.Query("from"); s != "" {
		if w.From, err = parse(s, false); err != nil {
			return w, domain.ErrValidation.WithDetails(map[string]any{"from": "RFC3339 or YYYY-MM-DD"})
		}
	}
	if s := c.Query("to"); s != "" {
		if w.To, err = parse(s, true); err != nil {
			return w, domain.ErrValidation.WithDetails(map[string]any{"to": "RFC3339 or YYYY-MM-DD"})
		}
	}
	if !w.To.After(w.From) {
		return w, domain.ErrValidation.WithDetails(map[string]any{"to": "must be after from"})
	}
	return w, nil
}

func (h *Handlers) dashboard(c *fiber.Ctx) error {
	w, err := parseWindow(c, 14)
	if err != nil {
		return err
	}
	d, err := h.Stats.Dashboard(c.UserContext(), membership(c).Project.ID, w)
	if err != nil {
		return err
	}
	return httpx.OK(c, d)
}

func (h *Handlers) usage(c *fiber.Ctx) error {
	w, err := parseWindow(c, 30)
	if err != nil {
		return err
	}
	groupBy := c.Query("group_by")
	rows, err := h.Stats.Usage(c.UserContext(), membership(c).Project.ID, w.From, w.To.Add(-time.Second), groupBy)
	if err != nil {
		return err
	}
	if c.Query("format") == "csv" {
		c.Set(fiber.HeaderContentType, "text/csv; charset=utf-8")
		c.Set(fiber.HeaderContentDisposition, `attachment; filename="habarchy-usage.csv"`)
		wr := csv.NewWriter(c)
		_ = wr.Write([]string{"day", "channel", "queued", "sent", "delivered", "failed", "cost", "currency"})
		for _, r := range rows {
			_ = wr.Write([]string{r.Day, r.Channel, strconv.FormatInt(r.Queued, 10), strconv.FormatInt(r.Sent, 10),
				strconv.FormatInt(r.Delivered, 10), strconv.FormatInt(r.Failed, 10), fmt.Sprintf("%.4f", float64(r.CostMicros)/1e6), r.Currency})
		}
		wr.Flush()
		return nil
	}
	return httpx.JSON(c, fiber.StatusOK, rows, fiber.Map{"from": w.From, "to": w.To, "group_by": groupBy})
}

func (h *Handlers) health(c *fiber.Ctx) error {
	hs, err := h.Stats.Health(c.UserContext(), membership(c).Project.ID)
	if err != nil {
		return err
	}
	return httpx.OK(c, hs)
}

func (h *Handlers) auditLogs(c *fiber.Ctx) error {
	pid := membership(c).Project.ID
	page := httpx.ParsePage(c, 50, 200)
	rows, err := h.DB.Queries.ListAuditLogs(c.UserContext(), sqlcgen.ListAuditLogsParams{ProjectID: &pid, RowLimit: page.Limit, RowOffset: page.Offset})
	if err != nil {
		return err
	}
	total, _ := h.DB.Queries.CountAuditLogs(c.UserContext(), &pid)
	return httpx.JSON(c, fiber.StatusOK, rows, fiber.Map{"total": total, "limit": page.Limit, "offset": page.Offset})
}

func (h *Handlers) listContacts(c *fiber.Ctx) error {
	page := httpx.ParsePage(c, 50, 200)
	params := sqlcgen.AdminListContactsParams{ProjectID: membership(c).Project.ID, RowLimit: page.Limit, RowOffset: page.Offset}
	if s := strings.TrimSpace(c.Query("search")); s != "" {
		params.Search = &s
	}
	rows, err := h.DB.Queries.AdminListContacts(c.UserContext(), params)
	if err != nil {
		return err
	}
	total, _ := h.DB.Queries.CountContacts(c.UserContext(), membership(c).Project.ID)
	return httpx.JSON(c, fiber.StatusOK, rows, fiber.Map{"total": total, "limit": page.Limit, "offset": page.Offset})
}

func (h *Handlers) listDevices(c *fiber.Ctx) error {
	page := httpx.ParsePage(c, 50, 200)
	rows, err := h.DB.Queries.ListDevices(c.UserContext(), sqlcgen.ListDevicesParams{ProjectID: membership(c).Project.ID, RowLimit: page.Limit, RowOffset: page.Offset})
	if err != nil {
		return err
	}
	type row struct {
		sqlcgen.Device
		FcmToken string `json:"fcm_token"` // masked override
	}
	out := make([]row, 0, len(rows))
	for _, d := range rows {
		tok := d.FcmToken
		if len(tok) > 12 {
			tok = tok[:6] + "…" + tok[len(tok)-4:]
		}
		out = append(out, row{Device: d, FcmToken: tok})
	}
	return httpx.JSON(c, fiber.StatusOK, out, fiber.Map{"limit": page.Limit, "offset": page.Offset})
}

// ---- users (global admin helpers for inviting members) ----

type createUserRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8,max=256"`
	FullName string `json:"full_name" validate:"max=100"`
}

func (h *Handlers) listUsers(c *fiber.Ctx) error {
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

func (h *Handlers) createUser(c *fiber.Ctx) error {
	var req createUserRequest
	if err := httpx.Bind(c, &req); err != nil {
		return err
	}
	u, err := h.Auth.CreateUser(c.UserContext(), req.Email, req.Password, req.FullName)
	if err != nil {
		return err
	}
	return httpx.Created(c, toUser(u))
}

func (h *Handlers) overview(c *fiber.Ctx) error {
	g, err := h.DB.Queries.GlobalCounts(c.UserContext())
	if err != nil {
		return err
	}
	return httpx.OK(c, fiber.Map{
		"projects": g.Projects, "users": g.Users, "messages_24h": g.Messages24h,
		"pending": g.Pending, "webhooks_pending": g.WebhooksPending,
	})
}

// ---- SSE ----

// stream sends live message events for the projects the user can see
// (or one project with ?project_id=). Auth: Bearer header, ?access_token=
// or the habarchy_access cookie, because EventSource cannot set headers.
func (h *Handlers) stream(c *fiber.Ctx) error {
	userID := middleware.UserID(c)
	var ids []uuid.UUID
	if s := c.Query("project_id"); s != "" {
		pid, err := uuid.Parse(s)
		if err != nil {
			return domain.ErrValidation.WithDetails(map[string]any{"project_id": "must be a uuid"})
		}
		if _, err := h.Projects.Authorize(c.UserContext(), userID, pid, domain.RoleViewer); err != nil {
			return err
		}
		ids = []uuid.UUID{pid}
	} else {
		list, err := h.Projects.ListForUser(c.UserContext(), userID)
		if err != nil {
			return err
		}
		for _, m := range list {
			ids = append(ids, m.Project.ID)
		}
	}
	if h.Redis == nil {
		return domain.ErrInternal.WithMessage("live stream unavailable")
	}

	c.Set(fiber.HeaderContentType, "text/event-stream")
	c.Set(fiber.HeaderCacheControl, "no-cache")
	c.Set(fiber.HeaderConnection, "keep-alive")
	c.Set("X-Accel-Buffering", "no")

	ctx, cancel := contextWithCancel(c)
	ch := events.Subscribe(ctx, h.Redis, ids)
	c.Context().SetBodyStreamWriter(fasthttp.StreamWriter(func(w *bufio.Writer) {
		defer cancel()
		_, _ = fmt.Fprintf(w, "event: ready\ndata: {\"projects\":%d}\n\n", len(ids))
		if err := w.Flush(); err != nil {
			return
		}
		ticker := time.NewTicker(20 * time.Second)
		defer ticker.Stop()
		for {
			select {
			case ev, ok := <-ch:
				if !ok {
					return
				}
				b, _ := json.Marshal(ev)
				_, _ = fmt.Fprintf(w, "event: %s\ndata: %s\n\n", ev.Type, b)
				if err := w.Flush(); err != nil {
					return
				}
			case <-ticker.C:
				_, _ = fmt.Fprint(w, ": keep-alive\n\n")
				if err := w.Flush(); err != nil {
					return
				}
			case <-ctx.Done():
				return
			}
		}
	}))
	return nil
}
