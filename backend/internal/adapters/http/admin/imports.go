package admin

import (
	"fmt"
	"html"
	"net/url"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"

	httpx "github.com/Esca6585dev/habarchy/backend/internal/adapters/http"
	"github.com/Esca6585dev/habarchy/backend/internal/adapters/http/middleware"
	"github.com/Esca6585dev/habarchy/backend/internal/app/contactimport"
	"github.com/Esca6585dev/habarchy/backend/internal/domain"
)

// MaxImportBytes caps uploaded contact files.
const MaxImportBytes = 10 << 20

// importOptions reads group_id / dry_run / tag from a multipart or JSON form.
func importOptions(c *fiber.Ctx, source string) (contactimport.Options, error) {
	opt := contactimport.Options{Source: source, DryRun: c.FormValue("dry_run") == "true" || c.Query("dry_run") == "true", Tag: strings.TrimSpace(c.FormValue("tag"))}
	if g := strings.TrimSpace(c.FormValue("group_id")); g != "" {
		id, err := uuid.Parse(g)
		if err != nil {
			return opt, domain.ErrValidation.WithDetails(map[string]any{"group_id": "invalid uuid"})
		}
		opt.GroupID = &id
	}
	return opt, nil
}

// ImportFile is shared by the admin and public endpoints.
func ImportFile(c *fiber.Ctx, svc *contactimport.Service, projectID uuid.UUID) error {
	fh, err := c.FormFile("file")
	if err != nil {
		return domain.ErrValidation.WithDetails(map[string]any{"file": "multipart field 'file' required (.xlsx, .csv, .vcf, .docx, .txt)"})
	}
	if fh.Size > MaxImportBytes {
		return domain.ErrValidation.WithDetails(map[string]any{"file": "max 10 MB"})
	}
	f, err := fh.Open()
	if err != nil {
		return err
	}
	defer func() { _ = f.Close() }()
	inputs, rowErrs, err := contactimport.Parse(fh.Filename, f)
	if err != nil {
		return domain.ErrValidation.WithDetails(map[string]any{"file": err.Error()})
	}
	opt, err := importOptions(c, fh.Filename)
	if err != nil {
		return err
	}
	res, err := svc.Import(c.UserContext(), projectID, inputs, opt)
	if err != nil {
		return err
	}
	res.Errors = append(res.Errors, rowErrs...)
	return httpx.OK(c, res)
}

func (h *Handlers) importContacts(c *fiber.Ctx) error {
	return ImportFile(c, h.Imports, membership(c).Project.ID)
}

type cardDAVRequest struct {
	ServerURL string `json:"server_url" validate:"max=500"`
	Username  string `json:"username" validate:"required,max=200"`
	Password  string `json:"password" validate:"required,max=200"`
	GroupID   string `json:"group_id"`
	DryRun    bool   `json:"dry_run"`
	Tag       string `json:"tag" validate:"max=64"`
}

// importCardDAV pulls contacts from iCloud (Apple ID + app-specific
// password) or any CardDAV server. Credentials are used once, not stored.
func (h *Handlers) importCardDAV(c *fiber.Ctx) error {
	var req cardDAVRequest
	if err := httpx.Bind(c, &req); err != nil {
		return err
	}
	inputs, err := h.Imports.CardDAV.Fetch(c.UserContext(), req.ServerURL, req.Username, req.Password)
	if err != nil {
		return domain.ErrValidation.WithMessage(err.Error())
	}
	opt := contactimport.Options{Source: "carddav:" + hostOf(req.ServerURL), DryRun: req.DryRun, Tag: req.Tag}
	if req.GroupID != "" {
		id, err := uuid.Parse(req.GroupID)
		if err != nil {
			return domain.ErrValidation.WithDetails(map[string]any{"group_id": "invalid uuid"})
		}
		opt.GroupID = &id
	}
	res, err := h.Imports.Import(c.UserContext(), membership(c).Project.ID, inputs, opt)
	if err != nil {
		return err
	}
	return httpx.OK(c, res)
}

func hostOf(raw string) string {
	if raw == "" {
		return "icloud"
	}
	if u, err := url.Parse(raw); err == nil && u.Host != "" {
		return u.Host
	}
	return raw
}

// googleImportURL returns the Google consent URL for this project.
func (h *Handlers) googleImportURL(c *fiber.Ctx) error {
	g := h.Imports.Google
	if !g.Configured() {
		return domain.ErrValidation.WithMessage("Google import is not configured (HABARCHY_GOOGLE_CLIENT_ID / HABARCHY_GOOGLE_CLIENT_SECRET)")
	}
	st := contactimport.State{ProjectID: membership(c).Project.ID, UserID: middleware.UserID(c), ReturnTo: strings.TrimSpace(c.Query("return_to"))}
	if gid := strings.TrimSpace(c.Query("group_id")); gid != "" {
		id, err := uuid.Parse(gid)
		if err != nil {
			return domain.ErrValidation.WithDetails(map[string]any{"group_id": "invalid uuid"})
		}
		st.GroupID = &id
	}
	return httpx.OK(c, fiber.Map{"url": g.AuthorizationURL(g.SignState(st))})
}

// googleCallback finishes the OAuth round-trip: verifies the signed state,
// re-checks the user's project role, imports, then redirects back (web) or
// renders a small page (mobile).
func (h *Handlers) googleCallback(c *fiber.Ctx) error {
	g := h.Imports.Google
	if !g.Configured() {
		return domain.ErrNotFound
	}
	st, err := g.ParseState(c.Query("state"))
	if err != nil {
		return domain.ErrValidation.WithMessage("invalid or expired state, start the import again")
	}
	if e := c.Query("error"); e != "" {
		return renderImportPage(c, st.ReturnTo, "Google: "+e, nil)
	}
	if _, err := h.Projects.Authorize(c.UserContext(), st.UserID, st.ProjectID, domain.RoleDeveloper); err != nil {
		return err
	}
	tok, err := g.Exchange(c.UserContext(), c.Query("code"))
	if err != nil {
		return renderImportPage(c, st.ReturnTo, err.Error(), nil)
	}
	inputs, err := g.FetchContacts(c.UserContext(), tok)
	if err != nil {
		return renderImportPage(c, st.ReturnTo, err.Error(), nil)
	}
	res, err := h.Imports.Import(c.UserContext(), st.ProjectID, inputs, contactimport.Options{Source: "google", GroupID: st.GroupID})
	if err != nil {
		return renderImportPage(c, st.ReturnTo, err.Error(), nil)
	}
	return renderImportPage(c, st.ReturnTo, "", res)
}

// renderImportPage redirects to return_to with the summary in the query, or
// shows a minimal HTML page when there is nowhere to return (mobile).
func renderImportPage(c *fiber.Ctx, returnTo, errMsg string, res *contactimport.Result) error {
	if returnTo != "" && (strings.HasPrefix(returnTo, "/") || strings.HasPrefix(returnTo, "http")) {
		q := url.Values{}
		if errMsg != "" {
			q.Set("import_error", errMsg)
		} else if res != nil {
			q.Set("imported", fmt.Sprint(res.Created+res.Updated))
			q.Set("created", fmt.Sprint(res.Created))
			q.Set("skipped", fmt.Sprint(res.Skipped))
		}
		sep := "?"
		if strings.Contains(returnTo, "?") {
			sep = "&"
		}
		return c.Redirect(returnTo+sep+q.Encode(), fiber.StatusFound)
	}
	title, body := "Google Contacts import", ""
	if errMsg != "" {
		body = "<p style='color:#b91c1c'>" + html.EscapeString(errMsg) + "</p>"
	} else if res != nil {
		body = fmt.Sprintf("<p>%d contacts: <b>%d</b> created, <b>%d</b> updated, %d unchanged, %d skipped.</p><p>You can close this page and return to the app.</p>", res.Total, res.Created, res.Updated, res.Unchanged, res.Skipped)
	}
	c.Set(fiber.HeaderContentType, fiber.MIMETextHTMLCharsetUTF8)
	return c.SendString("<!doctype html><html><head><meta charset='utf-8'><meta name='viewport' content='width=device-width'><title>" + title + "</title></head><body style='font-family:system-ui;padding:24px;max-width:480px;margin:auto'><h2>Habarçy · " + title + "</h2>" + body + "</body></html>")
}
