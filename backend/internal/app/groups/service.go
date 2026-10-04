// Package groups manages contact groups (colleagues, classmates,
// customers…) used to broadcast a message to many contacts at once.
package groups

import (
	"context"
	"errors"
	"strings"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"

	"github.com/Esca6585dev/habarchy/backend/internal/adapters/postgres"
	"github.com/Esca6585dev/habarchy/backend/internal/adapters/postgres/sqlcgen"
	"github.com/Esca6585dev/habarchy/backend/internal/app/audit"
	"github.com/Esca6585dev/habarchy/backend/internal/app/contacts"
	"github.com/Esca6585dev/habarchy/backend/internal/domain"
	"github.com/Esca6585dev/habarchy/backend/pkg/ids"
)

// Service holds the group use cases.
type Service struct {
	db       *postgres.DB
	contacts *contacts.Service
}

// New creates the service.
func New(db *postgres.DB, c *contacts.Service) *Service { return &Service{db: db, contacts: c} }

// Input is the payload for Create and Update.
type Input struct {
	Name        string
	Description string
}

func (in *Input) normalize() error {
	in.Name = strings.TrimSpace(in.Name)
	in.Description = strings.TrimSpace(in.Description)
	if in.Name == "" || len(in.Name) > 100 {
		return domain.ErrValidation.WithDetails(map[string]any{"name": "required, max 100 characters"})
	}
	return nil
}

// Create adds a group.
func (s *Service) Create(ctx context.Context, projectID uuid.UUID, in Input) (*sqlcgen.ContactGroup, error) {
	if err := in.normalize(); err != nil {
		return nil, err
	}
	g, err := s.db.Queries.CreateGroup(ctx, sqlcgen.CreateGroupParams{ID: ids.New(), ProjectID: projectID, Name: in.Name, Description: in.Description})
	if err != nil {
		if postgres.IsUniqueViolation(err) {
			return nil, domain.ErrConflict.WithMessage("a group with this name already exists")
		}
		return nil, err
	}
	audit.Record(ctx, s.db.Queries, &projectID, "group.create", "group", g.ID.String(), map[string]any{"name": g.Name})
	return &g, nil
}

// Get returns one group.
func (s *Service) Get(ctx context.Context, projectID, id uuid.UUID) (*sqlcgen.ContactGroup, error) {
	g, err := s.db.Queries.GetGroup(ctx, sqlcgen.GetGroupParams{ID: id, ProjectID: projectID})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrNotFound.WithMessage("group not found")
		}
		return nil, err
	}
	return &g, nil
}

// List returns the project's groups with member counts.
func (s *Service) List(ctx context.Context, projectID uuid.UUID) ([]sqlcgen.ListGroupsRow, error) {
	return s.db.Queries.ListGroups(ctx, projectID)
}

// Update renames a group.
func (s *Service) Update(ctx context.Context, projectID, id uuid.UUID, in Input) (*sqlcgen.ContactGroup, error) {
	if err := in.normalize(); err != nil {
		return nil, err
	}
	g, err := s.db.Queries.UpdateGroup(ctx, sqlcgen.UpdateGroupParams{ID: id, ProjectID: projectID, Name: in.Name, Description: in.Description})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrNotFound.WithMessage("group not found")
		}
		if postgres.IsUniqueViolation(err) {
			return nil, domain.ErrConflict.WithMessage("a group with this name already exists")
		}
		return nil, err
	}
	return &g, nil
}

// Delete removes a group (contacts stay).
func (s *Service) Delete(ctx context.Context, projectID, id uuid.UUID) error {
	n, err := s.db.Queries.DeleteGroup(ctx, sqlcgen.DeleteGroupParams{ID: id, ProjectID: projectID})
	if err != nil {
		return err
	}
	if n == 0 {
		return domain.ErrNotFound.WithMessage("group not found")
	}
	audit.Record(ctx, s.db.Queries, &projectID, "group.delete", "group", id.String(), nil)
	return nil
}

// MembersInput adds members by contact id, by external id, or by creating
// new contacts inline (e.g. a phone number typed in the app).
type MembersInput struct {
	ContactIDs  []uuid.UUID
	ExternalIDs []string
	Contacts    []contacts.Input
}

// AddResult reports what happened.
type AddResult struct {
	Added           int64    `json:"added"`
	CreatedContacts int      `json:"created_contacts"`
	NotFound        []string `json:"not_found"`
}

// AddMembers resolves and adds contacts to the group.
func (s *Service) AddMembers(ctx context.Context, projectID, groupID uuid.UUID, in MembersInput) (*AddResult, error) {
	if _, err := s.Get(ctx, projectID, groupID); err != nil {
		return nil, err
	}
	res := &AddResult{NotFound: []string{}}
	idSet := map[uuid.UUID]bool{}
	for _, id := range in.ContactIDs {
		idSet[id] = true
	}
	for _, ext := range in.ExternalIDs {
		c, err := s.db.Queries.GetContactByExternalID(ctx, sqlcgen.GetContactByExternalIDParams{ProjectID: projectID, ExternalID: strings.TrimSpace(ext)})
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				res.NotFound = append(res.NotFound, ext)
				continue
			}
			return nil, err
		}
		idSet[c.ID] = true
	}
	for _, cin := range in.Contacts {
		c, err := s.findOrCreate(ctx, projectID, cin)
		if err != nil {
			return nil, err
		}
		if c.created {
			res.CreatedContacts++
		}
		idSet[c.id] = true
	}
	if len(idSet) == 0 {
		if len(res.NotFound) > 0 {
			return res, nil
		}
		return nil, domain.ErrValidation.WithDetails(map[string]any{"members": "contact_ids, external_ids or contacts required"})
	}
	all := make([]uuid.UUID, 0, len(idSet))
	for id := range idSet {
		all = append(all, id)
	}
	n, err := s.db.Queries.AddGroupMembers(ctx, sqlcgen.AddGroupMembersParams{GroupID: groupID, ProjectID: projectID, ContactIds: all})
	if err != nil {
		return nil, err
	}
	res.Added = n
	return res, nil
}

type found struct {
	id      uuid.UUID
	created bool
}

// findOrCreate matches an inline contact by external_id, phone or email
// before creating it, so typing the same number twice does not duplicate.
func (s *Service) findOrCreate(ctx context.Context, projectID uuid.UUID, in contacts.Input) (found, error) {
	if ext := strings.TrimSpace(in.ExternalID); ext != "" {
		if c, err := s.db.Queries.GetContactByExternalID(ctx, sqlcgen.GetContactByExternalIDParams{ProjectID: projectID, ExternalID: ext}); err == nil {
			return found{id: c.ID}, nil
		}
	}
	if p := strings.TrimSpace(in.Phone); p != "" {
		if c, err := s.db.Queries.GetContactByPhone(ctx, sqlcgen.GetContactByPhoneParams{ProjectID: projectID, Phone: normalizePhone(p)}); err == nil {
			return found{id: c.ID}, nil
		}
	}
	if e := strings.TrimSpace(in.Email); e != "" {
		if c, err := s.db.Queries.GetContactByEmail(ctx, sqlcgen.GetContactByEmailParams{ProjectID: projectID, Email: e}); err == nil {
			return found{id: c.ID}, nil
		}
	}
	c, err := s.contacts.Create(ctx, projectID, in)
	if err != nil {
		return found{}, err
	}
	return found{id: c.ID, created: true}, nil
}

// RemoveMember detaches one contact.
func (s *Service) RemoveMember(ctx context.Context, projectID, groupID, contactID uuid.UUID) error {
	if _, err := s.Get(ctx, projectID, groupID); err != nil {
		return err
	}
	n, err := s.db.Queries.RemoveGroupMember(ctx, sqlcgen.RemoveGroupMemberParams{GroupID: groupID, ContactID: contactID})
	if err != nil {
		return err
	}
	if n == 0 {
		return domain.ErrNotFound.WithMessage("contact is not a member")
	}
	return nil
}

// Members lists a page of the group's contacts and the total.
func (s *Service) Members(ctx context.Context, projectID, groupID uuid.UUID, limit, offset int32) ([]sqlcgen.Contact, int64, error) {
	if _, err := s.Get(ctx, projectID, groupID); err != nil {
		return nil, 0, err
	}
	rows, err := s.db.Queries.ListGroupMembers(ctx, sqlcgen.ListGroupMembersParams{GroupID: groupID, RowLimit: limit, RowOffset: offset})
	if err != nil {
		return nil, 0, err
	}
	total, err := s.db.Queries.CountGroupMembers(ctx, groupID)
	if err != nil {
		return nil, 0, err
	}
	return rows, total, nil
}

// GroupsOf lists the groups a contact belongs to.
func (s *Service) GroupsOf(ctx context.Context, contactID uuid.UUID) ([]sqlcgen.ContactGroup, error) {
	return s.db.Queries.ListGroupsForContact(ctx, contactID)
}
