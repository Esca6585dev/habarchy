package admin

import (
	"encoding/json"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"

	httpx "github.com/Esca6585dev/habarchy/backend/internal/adapters/http"
	"github.com/Esca6585dev/habarchy/backend/internal/adapters/postgres/sqlcgen"
	"github.com/Esca6585dev/habarchy/backend/internal/app/contacts"
	"github.com/Esca6585dev/habarchy/backend/internal/app/groups"
	"github.com/Esca6585dev/habarchy/backend/internal/app/messages"
	"github.com/Esca6585dev/habarchy/backend/internal/domain"
)

// ---- contacts (admin) ----

// ContactResponse mirrors the public contact shape (snake_case).
type ContactResponse struct {
	ID             uuid.UUID       `json:"id"`
	ExternalID     string          `json:"external_id"`
	Name           string          `json:"name"`
	Phone          string          `json:"phone"`
	Email          string          `json:"email"`
	WhatsApp       string          `json:"whatsapp"`
	TelegramChatID string          `json:"telegram_chat_id"`
	SlackID        string          `json:"slack_id"`
	Locale         string          `json:"locale"`
	Tags           []string        `json:"tags"`
	Attributes     json.RawMessage `json:"attributes"`
	CreatedAt      time.Time       `json:"created_at"`
	UpdatedAt      time.Time       `json:"updated_at"`
}

func toContact(c *sqlcgen.Contact) ContactResponse {
	return ContactResponse{ID: c.ID, ExternalID: c.ExternalID, Name: c.Name, Phone: c.Phone, Email: c.Email, WhatsApp: c.Whatsapp,
		TelegramChatID: c.TelegramChatID, SlackID: c.SlackID, Locale: c.Locale, Tags: c.Tags, Attributes: c.Attributes, CreatedAt: c.CreatedAt, UpdatedAt: c.UpdatedAt}
}

type contactRequest struct {
	ExternalID     string         `json:"external_id" validate:"max=128"`
	Name           string         `json:"name" validate:"max=200"`
	Phone          string         `json:"phone" validate:"max=32"`
	Email          string         `json:"email" validate:"max=254"`
	WhatsApp       string         `json:"whatsapp" validate:"max=32"`
	TelegramChatID string         `json:"telegram_chat_id" validate:"max=64"`
	SlackID        string         `json:"slack_id" validate:"max=64"`
	Locale         domain.Locale  `json:"locale"`
	Tags           []string       `json:"tags" validate:"max=50,dive,max=64"`
	Attributes     map[string]any `json:"attributes"`
}

func (r contactRequest) input() contacts.Input {
	return contacts.Input{ExternalID: r.ExternalID, Name: r.Name, Phone: r.Phone, Email: r.Email, WhatsApp: r.WhatsApp,
		TelegramChatID: r.TelegramChatID, SlackID: r.SlackID, Locale: r.Locale, Tags: r.Tags, Attributes: r.Attributes}
}

func (h *Handlers) createContact(c *fiber.Ctx) error {
	var req contactRequest
	if err := httpx.Bind(c, &req); err != nil {
		return err
	}
	ct, created, err := h.Contacts.Upsert(c.UserContext(), membership(c).Project.ID, req.input())
	if err != nil {
		return err
	}
	status := fiber.StatusOK
	if created {
		status = fiber.StatusCreated
	}
	return httpx.JSON(c, status, toContact(ct), nil)
}

func (h *Handlers) getContact(c *fiber.Ctx) error {
	id, err := httpx.ParamUUID(c, "contact_id")
	if err != nil {
		return err
	}
	ct, err := h.Contacts.Get(c.UserContext(), membership(c).Project.ID, id)
	if err != nil {
		return err
	}
	out := fiber.Map{"contact": toContact(ct)}
	if h.Groups != nil {
		gs, _ := h.Groups.GroupsOf(c.UserContext(), id)
		list := make([]GroupResponse, 0, len(gs))
		for i := range gs {
			list = append(list, toGroup(&gs[i], 0))
		}
		out["groups"] = list
	}
	return httpx.OK(c, out)
}

func (h *Handlers) updateContact(c *fiber.Ctx) error {
	id, err := httpx.ParamUUID(c, "contact_id")
	if err != nil {
		return err
	}
	var req contactRequest
	if err := httpx.Bind(c, &req); err != nil {
		return err
	}
	ct, err := h.Contacts.Update(c.UserContext(), membership(c).Project.ID, id, req.input())
	if err != nil {
		return err
	}
	return httpx.OK(c, toContact(ct))
}

func (h *Handlers) deleteContact(c *fiber.Ctx) error {
	id, err := httpx.ParamUUID(c, "contact_id")
	if err != nil {
		return err
	}
	if err := h.Contacts.Delete(c.UserContext(), membership(c).Project.ID, id); err != nil {
		return err
	}
	return httpx.NoContent(c)
}

// ---- groups (admin) ----

// GroupResponse is a contact group.
type GroupResponse struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	MemberCount int64     `json:"member_count"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func toGroup(g *sqlcgen.ContactGroup, members int64) GroupResponse {
	return GroupResponse{ID: g.ID, Name: g.Name, Description: g.Description, MemberCount: members, CreatedAt: g.CreatedAt, UpdatedAt: g.UpdatedAt}
}

type groupRequest struct {
	Name        string `json:"name" validate:"required,max=100"`
	Description string `json:"description" validate:"max=500"`
}

type membersRequest struct {
	ContactIDs  []uuid.UUID      `json:"contact_ids"`
	ExternalIDs []string         `json:"external_ids" validate:"max=1000,dive,max=128"`
	Contacts    []contactRequest `json:"contacts" validate:"max=1000,dive"`
}

func (h *Handlers) listGroups(c *fiber.Ctx) error {
	rows, err := h.Groups.List(c.UserContext(), membership(c).Project.ID)
	if err != nil {
		return err
	}
	out := make([]GroupResponse, 0, len(rows))
	for i := range rows {
		r := &rows[i]
		out = append(out, GroupResponse{ID: r.ID, Name: r.Name, Description: r.Description, MemberCount: r.MemberCount, CreatedAt: r.CreatedAt, UpdatedAt: r.UpdatedAt})
	}
	return httpx.OK(c, out)
}

func (h *Handlers) createGroup(c *fiber.Ctx) error {
	var req groupRequest
	if err := httpx.Bind(c, &req); err != nil {
		return err
	}
	g, err := h.Groups.Create(c.UserContext(), membership(c).Project.ID, groups.Input{Name: req.Name, Description: req.Description})
	if err != nil {
		return err
	}
	return httpx.Created(c, toGroup(g, 0))
}

func (h *Handlers) getGroup(c *fiber.Ctx) error {
	id, err := httpx.ParamUUID(c, "group_id")
	if err != nil {
		return err
	}
	g, err := h.Groups.Get(c.UserContext(), membership(c).Project.ID, id)
	if err != nil {
		return err
	}
	_, total, _ := h.Groups.Members(c.UserContext(), membership(c).Project.ID, id, 1, 0)
	return httpx.OK(c, toGroup(g, total))
}

func (h *Handlers) updateGroup(c *fiber.Ctx) error {
	id, err := httpx.ParamUUID(c, "group_id")
	if err != nil {
		return err
	}
	var req groupRequest
	if err := httpx.Bind(c, &req); err != nil {
		return err
	}
	g, err := h.Groups.Update(c.UserContext(), membership(c).Project.ID, id, groups.Input{Name: req.Name, Description: req.Description})
	if err != nil {
		return err
	}
	return httpx.OK(c, toGroup(g, 0))
}

func (h *Handlers) deleteGroup(c *fiber.Ctx) error {
	id, err := httpx.ParamUUID(c, "group_id")
	if err != nil {
		return err
	}
	if err := h.Groups.Delete(c.UserContext(), membership(c).Project.ID, id); err != nil {
		return err
	}
	return httpx.NoContent(c)
}

func (h *Handlers) listGroupMembers(c *fiber.Ctx) error {
	id, err := httpx.ParamUUID(c, "group_id")
	if err != nil {
		return err
	}
	page := httpx.ParsePage(c, 50, 500)
	rows, total, err := h.Groups.Members(c.UserContext(), membership(c).Project.ID, id, page.Limit, page.Offset)
	if err != nil {
		return err
	}
	out := make([]ContactResponse, 0, len(rows))
	for i := range rows {
		out = append(out, toContact(&rows[i]))
	}
	return httpx.JSON(c, fiber.StatusOK, out, fiber.Map{"total": total, "limit": page.Limit, "offset": page.Offset})
}

func (h *Handlers) addGroupMembers(c *fiber.Ctx) error {
	id, err := httpx.ParamUUID(c, "group_id")
	if err != nil {
		return err
	}
	var req membersRequest
	if err := httpx.Bind(c, &req); err != nil {
		return err
	}
	in := groups.MembersInput{ContactIDs: req.ContactIDs, ExternalIDs: req.ExternalIDs}
	for _, ct := range req.Contacts {
		in.Contacts = append(in.Contacts, ct.input())
	}
	res, err := h.Groups.AddMembers(c.UserContext(), membership(c).Project.ID, id, in)
	if err != nil {
		return err
	}
	return httpx.OK(c, res)
}

func (h *Handlers) removeGroupMember(c *fiber.Ctx) error {
	id, err := httpx.ParamUUID(c, "group_id")
	if err != nil {
		return err
	}
	cid, err := httpx.ParamUUID(c, "contact_id")
	if err != nil {
		return err
	}
	if err := h.Groups.RemoveMember(c.UserContext(), membership(c).Project.ID, id, cid); err != nil {
		return err
	}
	return httpx.NoContent(c)
}

// ---- compose & send from the admin panel / app ----

type sendRequest struct {
	Channel     domain.Channel  `json:"channel"`
	Template    string          `json:"template" validate:"max=64"`
	Data        map[string]any  `json:"data"`
	Subject     string          `json:"subject" validate:"max=998"`
	Title       string          `json:"title" validate:"max=200"`
	Body        string          `json:"body" validate:"max=65536"`
	Locale      domain.Locale   `json:"locale"`
	ScheduledAt *time.Time      `json:"scheduled_at"`
	Priority    domain.Priority `json:"priority"`
	Metadata    map[string]any  `json:"metadata"`
	GroupIDs    []uuid.UUID     `json:"group_ids"`
	ContactIDs  []uuid.UUID     `json:"contact_ids"`
	To          []string        `json:"to" validate:"max=1000,dive,max=254"` // raw addresses
	// IsTest routes through the sandbox (nothing is really sent).
	IsTest bool `json:"is_test"`
}

// sendMessages composes a batch to groups, contacts and/or raw addresses.
func (h *Handlers) sendMessages(c *fiber.Ctx) error {
	var req sendRequest
	if err := httpx.Bind(c, &req); err != nil {
		return err
	}
	if len(req.GroupIDs) == 0 && len(req.ContactIDs) == 0 && len(req.To) == 0 {
		return domain.ErrValidation.WithDetails(map[string]any{"recipients": "group_ids, contact_ids or to required"})
	}
	in := messages.BatchInput{
		Channel: req.Channel, Template: req.Template, Subject: req.Subject, Title: req.Title, Body: req.Body, Locale: req.Locale,
		ScheduledAt: req.ScheduledAt, Priority: req.Priority, Metadata: req.Metadata, GroupIDs: req.GroupIDs, Data: req.Data,
	}
	for _, id := range req.ContactIDs {
		cid := id
		in.Recipients = append(in.Recipients, messages.BatchRecipient{Recipient: messages.Recipient{ContactID: &cid}})
	}
	for _, addr := range req.To {
		in.Recipients = append(in.Recipients, messages.BatchRecipient{Recipient: messages.Recipient{Address: addr}})
	}
	res, err := h.Messages.SendBatch(c.UserContext(), messages.Caller{Project: membership(c).Project, IsTest: req.IsTest}, in)
	if err != nil {
		return err
	}
	status := fiber.StatusAccepted
	meta := fiber.Map{"accepted": res.Accepted, "rejected": res.Rejected}
	if res.Duplicate {
		status, meta["duplicate"] = fiber.StatusOK, true
	}
	return httpx.JSON(c, status, fiber.Map{
		"id": res.Batch.ID, "status": res.Batch.Status, "channel": res.Batch.Channel, "template": res.Batch.TemplateKey,
		"total": res.Batch.Total, "queued": res.Batch.Queued, "sent": res.Batch.Sent, "delivered": res.Batch.Delivered, "failed": res.Batch.Failed,
		"created_at": res.Batch.CreatedAt, "completed_at": res.Batch.CompletedAt,
	}, meta)
}
