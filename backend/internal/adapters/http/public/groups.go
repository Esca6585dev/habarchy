package public

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"

	httpx "github.com/Esca6585dev/habarchy/backend/internal/adapters/http"
	"github.com/Esca6585dev/habarchy/backend/internal/adapters/postgres/sqlcgen"
	"github.com/Esca6585dev/habarchy/backend/internal/app/groups"
)

// GroupResponse is the public shape of a contact group.
type GroupResponse struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	MemberCount int64     `json:"member_count"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// ToGroup converts a row (member count from the list query when known).
func ToGroup(g *sqlcgen.ContactGroup, members int64) GroupResponse {
	return GroupResponse{ID: g.ID, Name: g.Name, Description: g.Description, MemberCount: members, CreatedAt: g.CreatedAt, UpdatedAt: g.UpdatedAt}
}

// ToGroupRow converts a list row.
func ToGroupRow(r *sqlcgen.ListGroupsRow) GroupResponse {
	return GroupResponse{ID: r.ID, Name: r.Name, Description: r.Description, MemberCount: r.MemberCount, CreatedAt: r.CreatedAt, UpdatedAt: r.UpdatedAt}
}

// GroupRequest creates / renames a group.
type GroupRequest struct {
	Name        string `json:"name" validate:"required,max=100"`
	Description string `json:"description" validate:"max=500"`
}

// MembersRequest adds members: existing contacts by id or external_id, or
// new contacts created inline (name + phone/email/…).
type MembersRequest struct {
	ContactIDs  []uuid.UUID      `json:"contact_ids"`
	ExternalIDs []string         `json:"external_ids" validate:"max=1000,dive,max=128"`
	Contacts    []contactRequest `json:"contacts" validate:"max=1000,dive"`
}

// ToMembersInput converts the request.
func (r MembersRequest) ToMembersInput() groups.MembersInput {
	in := groups.MembersInput{ContactIDs: r.ContactIDs, ExternalIDs: r.ExternalIDs}
	for _, c := range r.Contacts {
		in.Contacts = append(in.Contacts, contactInput(c))
	}
	return in
}

func (h *Handlers) listGroups(c *fiber.Ctx) error {
	rows, err := h.Groups.List(c.UserContext(), caller(c).Project.ID)
	if err != nil {
		return err
	}
	out := make([]GroupResponse, 0, len(rows))
	for i := range rows {
		out = append(out, ToGroupRow(&rows[i]))
	}
	return httpx.OK(c, out)
}

func (h *Handlers) createGroup(c *fiber.Ctx) error {
	var req GroupRequest
	if err := httpx.Bind(c, &req); err != nil {
		return err
	}
	g, err := h.Groups.Create(c.UserContext(), caller(c).Project.ID, groups.Input{Name: req.Name, Description: req.Description})
	if err != nil {
		return err
	}
	return httpx.Created(c, ToGroup(g, 0))
}

func (h *Handlers) getGroup(c *fiber.Ctx) error {
	id, err := httpx.ParamUUID(c, "group_id")
	if err != nil {
		return err
	}
	g, err := h.Groups.Get(c.UserContext(), caller(c).Project.ID, id)
	if err != nil {
		return err
	}
	_, total, _ := h.Groups.Members(c.UserContext(), caller(c).Project.ID, id, 1, 0)
	return httpx.OK(c, ToGroup(g, total))
}

func (h *Handlers) updateGroup(c *fiber.Ctx) error {
	id, err := httpx.ParamUUID(c, "group_id")
	if err != nil {
		return err
	}
	var req GroupRequest
	if err := httpx.Bind(c, &req); err != nil {
		return err
	}
	g, err := h.Groups.Update(c.UserContext(), caller(c).Project.ID, id, groups.Input{Name: req.Name, Description: req.Description})
	if err != nil {
		return err
	}
	return httpx.OK(c, ToGroup(g, 0))
}

func (h *Handlers) deleteGroup(c *fiber.Ctx) error {
	id, err := httpx.ParamUUID(c, "group_id")
	if err != nil {
		return err
	}
	if err := h.Groups.Delete(c.UserContext(), caller(c).Project.ID, id); err != nil {
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
	rows, total, err := h.Groups.Members(c.UserContext(), caller(c).Project.ID, id, page.Limit, page.Offset)
	if err != nil {
		return err
	}
	out := make([]ContactResponse, 0, len(rows))
	for i := range rows {
		out = append(out, ToContact(&rows[i]))
	}
	return httpx.JSON(c, fiber.StatusOK, out, fiber.Map{"total": total, "limit": page.Limit, "offset": page.Offset})
}

func (h *Handlers) addGroupMembers(c *fiber.Ctx) error {
	id, err := httpx.ParamUUID(c, "group_id")
	if err != nil {
		return err
	}
	var req MembersRequest
	if err := httpx.Bind(c, &req); err != nil {
		return err
	}
	res, err := h.Groups.AddMembers(c.UserContext(), caller(c).Project.ID, id, req.ToMembersInput())
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
	if err := h.Groups.RemoveMember(c.UserContext(), caller(c).Project.ID, id, cid); err != nil {
		return err
	}
	return httpx.NoContent(c)
}
