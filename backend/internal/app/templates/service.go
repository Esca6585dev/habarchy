// Package templates manages message templates and their version history.
package templates

import (
	"context"
	"errors"
	"regexp"
	"strings"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"

	"github.com/Esca6585dev/habarchy/backend/internal/adapters/postgres"
	"github.com/Esca6585dev/habarchy/backend/internal/adapters/postgres/sqlcgen"
	"github.com/Esca6585dev/habarchy/backend/internal/app/audit"
	"github.com/Esca6585dev/habarchy/backend/internal/domain"
	"github.com/Esca6585dev/habarchy/backend/pkg/render"
)

// Service holds the template use cases. Authorization is the caller's job
// (admin: project role check; public API: scope check), so every method
// takes an already-trusted projectID.
type Service struct {
	db *postgres.DB
}

// New creates the service.
func New(db *postgres.DB) *Service { return &Service{db: db} }

var keyRe = regexp.MustCompile(`^[a-z0-9][a-z0-9_.-]{0,63}$`)

// Input is the payload for Create and Update.
type Input struct {
	Key      string
	Channel  domain.Channel
	Locale   domain.Locale
	Subject  string
	Body     string
	IsActive *bool
}

func (in *Input) validate(forCreate bool) ([]string, error) {
	details := map[string]any{}
	in.Key = strings.ToLower(strings.TrimSpace(in.Key))
	if forCreate {
		if !keyRe.MatchString(in.Key) {
			details["key"] = "lowercase letters, digits, _ . - (max 64)"
		}
		if !in.Channel.Valid() {
			details["channel"] = "sms, email, push or telegram"
		}
		if in.Locale == "" {
			in.Locale = domain.DefaultLocale
		}
		if !in.Locale.Valid() {
			details["locale"] = "tk, ru or en"
		}
	}
	if strings.TrimSpace(in.Body) == "" {
		details["body"] = "required"
	}
	if in.Channel == domain.ChannelSMS && in.Subject != "" {
		details["subject"] = "sms templates have no subject"
	}
	if in.Channel == domain.ChannelEmail && strings.TrimSpace(in.Subject) == "" {
		details["subject"] = "required for email"
	}
	if len(details) > 0 {
		return nil, domain.ErrValidation.WithDetails(details)
	}
	vars, err := render.Validate(in.Subject, in.Body)
	if err != nil {
		return nil, domain.ErrValidation.WithDetails(map[string]any{"body": strings.TrimPrefix(err.Error(), domain.ErrValidation.Error()+": ")})
	}
	return vars, nil
}

// Create adds a template and its first version.
func (s *Service) Create(ctx context.Context, projectID uuid.UUID, in Input) (*sqlcgen.Template, error) {
	vars, err := in.validate(true)
	if err != nil {
		return nil, err
	}
	var tpl sqlcgen.Template
	err = s.db.WithTx(ctx, func(q *sqlcgen.Queries) error {
		var err error
		tpl, err = q.CreateTemplate(ctx, sqlcgen.CreateTemplateParams{
			ProjectID: projectID, Key: in.Key, Channel: sqlcgen.Channel(in.Channel), Locale: string(in.Locale),
			Subject: in.Subject, Body: in.Body, RequiredVars: vars,
		})
		if err != nil {
			if postgres.IsUniqueViolation(err) {
				return domain.ErrConflict.WithMessage("template with this key, channel and locale already exists")
			}
			return err
		}
		if err := s.snapshot(ctx, q, tpl); err != nil {
			return err
		}
		audit.Record(ctx, q, &projectID, "template.create", "template", tpl.ID.String(), map[string]any{"key": tpl.Key, "channel": tpl.Channel, "locale": tpl.Locale})
		return nil
	})
	if err != nil {
		return nil, err
	}
	return &tpl, nil
}

// Get returns one template.
func (s *Service) Get(ctx context.Context, projectID, id uuid.UUID) (*sqlcgen.Template, error) {
	tpl, err := s.db.Queries.GetTemplate(ctx, sqlcgen.GetTemplateParams{ID: id, ProjectID: projectID})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrNotFound.WithMessage("template not found")
		}
		return nil, err
	}
	return &tpl, nil
}

// List returns all templates of a project, optionally filtered by key.
func (s *Service) List(ctx context.Context, projectID uuid.UUID, key string) ([]sqlcgen.Template, error) {
	if key != "" {
		return s.db.Queries.ListTemplatesByKey(ctx, sqlcgen.ListTemplatesByKeyParams{ProjectID: projectID, Key: strings.ToLower(key)})
	}
	return s.db.Queries.ListTemplates(ctx, projectID)
}

// Update saves new content, bumps the version and snapshots it.
func (s *Service) Update(ctx context.Context, projectID, id uuid.UUID, in Input) (*sqlcgen.Template, error) {
	current, err := s.Get(ctx, projectID, id)
	if err != nil {
		return nil, err
	}
	in.Channel = domain.Channel(current.Channel)
	if in.Body == "" {
		in.Body = current.Body
	}
	vars, err := in.validate(false)
	if err != nil {
		return nil, err
	}
	active := current.IsActive
	if in.IsActive != nil {
		active = *in.IsActive
	}
	var tpl sqlcgen.Template
	err = s.db.WithTx(ctx, func(q *sqlcgen.Queries) error {
		var err error
		tpl, err = q.UpdateTemplate(ctx, sqlcgen.UpdateTemplateParams{
			ID: id, ProjectID: projectID, Subject: in.Subject, Body: in.Body, RequiredVars: vars, IsActive: active,
		})
		if err != nil {
			return err
		}
		if err := s.snapshot(ctx, q, tpl); err != nil {
			return err
		}
		audit.Record(ctx, q, &projectID, "template.update", "template", tpl.ID.String(), map[string]any{"version": tpl.Version, "is_active": tpl.IsActive})
		return nil
	})
	if err != nil {
		return nil, err
	}
	return &tpl, nil
}

// Delete removes a template and its history.
func (s *Service) Delete(ctx context.Context, projectID, id uuid.UUID) error {
	return s.db.WithTx(ctx, func(q *sqlcgen.Queries) error {
		n, err := q.DeleteTemplate(ctx, sqlcgen.DeleteTemplateParams{ID: id, ProjectID: projectID})
		if err != nil {
			return err
		}
		if n == 0 {
			return domain.ErrNotFound.WithMessage("template not found")
		}
		audit.Record(ctx, q, &projectID, "template.delete", "template", id.String(), nil)
		return nil
	})
}

// Versions lists the history, newest first.
func (s *Service) Versions(ctx context.Context, projectID, id uuid.UUID) ([]sqlcgen.TemplateVersion, error) {
	if _, err := s.Get(ctx, projectID, id); err != nil {
		return nil, err
	}
	return s.db.Queries.ListTemplateVersions(ctx, id)
}

// Restore makes an old version the current content (as a new version).
func (s *Service) Restore(ctx context.Context, projectID, id uuid.UUID, version int32) (*sqlcgen.Template, error) {
	if _, err := s.Get(ctx, projectID, id); err != nil {
		return nil, err
	}
	v, err := s.db.Queries.GetTemplateVersion(ctx, sqlcgen.GetTemplateVersionParams{TemplateID: id, Version: version})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrNotFound.WithMessage("version not found")
		}
		return nil, err
	}
	return s.Update(ctx, projectID, id, Input{Subject: v.Subject, Body: v.Body})
}

// Preview renders a template (stored or ad hoc) with sample data.
type Preview struct {
	Subject      string   `json:"subject"`
	Body         string   `json:"body"`
	RequiredVars []string `json:"required_vars"`
	MissingVars  []string `json:"missing_vars"`
}

// PreviewStored renders a saved template with data.
func (s *Service) PreviewStored(ctx context.Context, projectID, id uuid.UUID, data map[string]any) (*Preview, error) {
	tpl, err := s.Get(ctx, projectID, id)
	if err != nil {
		return nil, err
	}
	return PreviewRaw(domain.Channel(tpl.Channel), tpl.Subject, tpl.Body, data)
}

// PreviewRaw renders unsaved content with data. Missing variables are
// reported instead of failing so the editor can show them live.
func PreviewRaw(channel domain.Channel, subject, body string, data map[string]any) (*Preview, error) {
	vars, err := render.Validate(subject, body)
	if err != nil {
		return nil, err
	}
	if data == nil {
		data = map[string]any{}
	}
	p := &Preview{RequiredVars: vars, MissingVars: []string{}}
	filled := make(map[string]any, len(vars))
	for _, v := range vars {
		if val, ok := data[v]; ok {
			filled[v] = val
		} else {
			p.MissingVars = append(p.MissingVars, v)
			filled[v] = "{{" + v + "}}"
		}
	}
	res, err := render.Render(channel, subject, body, filled)
	if err != nil {
		return nil, err
	}
	p.Subject, p.Body = res.Subject, res.Body
	return p, nil
}

// Resolve finds the active template for (key, channel, locale), falling
// back to the project default locale and then to any locale. Used by the
// send flow in step 3.
func (s *Service) Resolve(ctx context.Context, projectID uuid.UUID, key string, channel domain.Channel, locale, fallback domain.Locale) (*sqlcgen.Template, error) {
	try := []domain.Locale{locale, fallback, domain.LocaleTK, domain.LocaleRU, domain.LocaleEN}
	seen := map[domain.Locale]bool{}
	for _, l := range try {
		if l == "" || seen[l] {
			continue
		}
		seen[l] = true
		tpl, err := s.db.Queries.FindTemplate(ctx, sqlcgen.FindTemplateParams{ProjectID: projectID, Key: key, Channel: sqlcgen.Channel(channel), Locale: string(l)})
		if err == nil {
			return &tpl, nil
		}
		if !errors.Is(err, pgx.ErrNoRows) {
			return nil, err
		}
	}
	return nil, domain.ErrNotFound.WithMessage("no active template " + key + " for channel " + string(channel))
}

func (s *Service) snapshot(ctx context.Context, q *sqlcgen.Queries, tpl sqlcgen.Template) error {
	var createdBy *uuid.UUID
	if a, ok := audit.ActorFrom(ctx); ok {
		id := a.UserID
		createdBy = &id
	}
	_, err := q.CreateTemplateVersion(ctx, sqlcgen.CreateTemplateVersionParams{
		TemplateID: tpl.ID, Version: tpl.Version, Subject: tpl.Subject, Body: tpl.Body, RequiredVars: tpl.RequiredVars, CreatedBy: createdBy,
	})
	return err
}
