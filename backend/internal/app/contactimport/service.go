package contactimport

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"github.com/Esca6585dev/habarchy/backend/internal/adapters/postgres"
	"github.com/Esca6585dev/habarchy/backend/internal/adapters/postgres/sqlcgen"
	"github.com/Esca6585dev/habarchy/backend/internal/app/audit"
	"github.com/Esca6585dev/habarchy/backend/internal/app/contacts"
	"github.com/Esca6585dev/habarchy/backend/internal/app/groups"
	"github.com/Esca6585dev/habarchy/backend/internal/domain"
)

// Service merges parsed contacts into a project.
type Service struct {
	db       *postgres.DB
	contacts *contacts.Service
	groups   *groups.Service
	// Google / CardDAV connectors (nil when not configured).
	Google  *GoogleConnector
	CardDAV *CardDAVClient
}

// New creates the service.
func New(db *postgres.DB, c *contacts.Service, g *groups.Service) *Service {
	return &Service{db: db, contacts: c, groups: g, CardDAV: NewCardDAVClient(nil)}
}

// Options control an import.
type Options struct {
	GroupID *uuid.UUID // add every imported contact to this group
	DryRun  bool       // parse and report only
	Source  string     // file name / "google" / "carddav" for the audit log
	Tag     string     // extra tag for every contact (e.g. "import-2026-10")
}

// Result summarises an import.
type Result struct {
	Total      int              `json:"total"`
	Created    int              `json:"created"`
	Updated    int              `json:"updated"`
	Unchanged  int              `json:"unchanged"`
	Skipped    int              `json:"skipped"`
	Errors     []RowError       `json:"errors"`
	Preview    []contacts.Input `json:"preview"`
	DryRun     bool             `json:"dry_run"`
	GroupAdded int64            `json:"group_added"`
}

// Import merges inputs into the project (and optionally a group).
func (s *Service) Import(ctx context.Context, projectID uuid.UUID, inputs []contacts.Input, opt Options) (*Result, error) {
	res := &Result{Total: len(inputs), Errors: []RowError{}, Preview: []contacts.Input{}, DryRun: opt.DryRun}
	if len(inputs) > MaxRows {
		return nil, domain.ErrValidation.WithMessage(fmt.Sprintf("too many contacts (max %d per import)", MaxRows))
	}
	if opt.GroupID != nil {
		if _, err := s.groups.Get(ctx, projectID, *opt.GroupID); err != nil {
			return nil, err
		}
	}
	for i := range inputs {
		if len(res.Preview) < 20 {
			res.Preview = append(res.Preview, inputs[i])
		}
	}
	if opt.DryRun {
		return res, nil
	}
	ids := make([]uuid.UUID, 0, len(inputs))
	for i, in := range inputs {
		if opt.Tag != "" {
			in.Tags = append(in.Tags, opt.Tag)
		}
		c, created, err := s.contacts.Merge(ctx, projectID, in)
		if err != nil {
			res.Skipped++
			res.Errors = append(res.Errors, RowError{Row: i + 1, Reason: reason(err), Line: describe(in)})
			continue
		}
		switch {
		case created:
			res.Created++
		case c.UpdatedAt.After(c.CreatedAt):
			res.Updated++
		default:
			res.Unchanged++
		}
		ids = append(ids, c.ID)
	}
	if opt.GroupID != nil && len(ids) > 0 {
		n, err := s.db.Queries.AddGroupMembers(ctx, sqlcgen.AddGroupMembersParams{GroupID: *opt.GroupID, ProjectID: projectID, ContactIds: ids})
		if err != nil {
			return nil, err
		}
		res.GroupAdded = n
	}
	audit.Record(ctx, s.db.Queries, &projectID, "contacts.import", "project", projectID.String(), map[string]any{
		"source": opt.Source, "total": res.Total, "created": res.Created, "updated": res.Updated, "skipped": res.Skipped,
	})
	return res, nil
}

func reason(err error) string {
	var de *domain.Error
	if ok := asDomain(err, &de); ok {
		if len(de.Details) > 0 {
			return fmt.Sprintf("%s: %v", de.Message, de.Details)
		}
		return de.Message
	}
	return err.Error()
}

func describe(in contacts.Input) string {
	for _, v := range []string{in.Name, in.Phone, in.Email, in.ExternalID} {
		if v != "" {
			return v
		}
	}
	return ""
}
