// Package contacts manages a project's address book and push devices.
package contacts

import (
	"context"
	"encoding/json"
	"errors"
	"net/mail"
	"strings"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"

	"github.com/Esca6585dev/habarchy/backend/internal/adapters/postgres"
	"github.com/Esca6585dev/habarchy/backend/internal/adapters/postgres/sqlcgen"
	"github.com/Esca6585dev/habarchy/backend/internal/domain"
	"github.com/Esca6585dev/habarchy/backend/pkg/phone"
)

// Service holds the contact use cases.
type Service struct {
	db *postgres.DB
}

// New creates the service.
func New(db *postgres.DB) *Service { return &Service{db: db} }

// Input is the payload for Create and Update.
type Input struct {
	ExternalID     string
	Name           string
	Phone          string
	Email          string
	WhatsApp       string // E.164; empty = same as Phone
	TelegramChatID string
	SlackID        string
	Locale         domain.Locale
	Tags           []string
	Attributes     map[string]any
}

func (in *Input) normalize() error {
	details := map[string]any{}
	in.ExternalID = strings.TrimSpace(in.ExternalID)
	if in.Phone = strings.TrimSpace(in.Phone); in.Phone != "" {
		p, err := phone.Normalize(in.Phone, "")
		if err != nil {
			details["phone"] = "invalid phone number"
		}
		in.Phone = p
	}
	if in.Email = strings.ToLower(strings.TrimSpace(in.Email)); in.Email != "" {
		if _, err := mail.ParseAddress(in.Email); err != nil || strings.ContainsAny(in.Email, " <>") {
			details["email"] = "invalid email address"
		}
	}
	in.TelegramChatID = strings.TrimSpace(in.TelegramChatID)
	in.Name = strings.TrimSpace(in.Name)
	in.SlackID = strings.TrimSpace(in.SlackID)
	if in.WhatsApp = strings.TrimSpace(in.WhatsApp); in.WhatsApp != "" {
		p, err := phone.Normalize(in.WhatsApp, "")
		if err != nil {
			details["whatsapp"] = "invalid phone number"
		}
		in.WhatsApp = p
	}
	if in.Locale == "" {
		in.Locale = domain.DefaultLocale
	}
	if !in.Locale.Valid() {
		details["locale"] = "tk, ru or en"
	}
	if in.Tags == nil {
		in.Tags = []string{}
	}
	if in.Attributes == nil {
		in.Attributes = map[string]any{}
	}
	if in.Phone == "" && in.Email == "" && in.TelegramChatID == "" && in.ExternalID == "" && in.WhatsApp == "" && in.SlackID == "" {
		details["contact"] = "at least one of external_id, phone, email, whatsapp, telegram_chat_id, slack_id is required"
	}
	if len(details) > 0 {
		return domain.ErrValidation.WithDetails(details)
	}
	return nil
}

// Create adds a contact.
func (s *Service) Create(ctx context.Context, projectID uuid.UUID, in Input) (*sqlcgen.Contact, error) {
	if err := in.normalize(); err != nil {
		return nil, err
	}
	attrs, _ := json.Marshal(in.Attributes)
	c, err := s.db.Queries.CreateContact(ctx, sqlcgen.CreateContactParams{
		ProjectID: projectID, ExternalID: in.ExternalID, Name: in.Name, Phone: in.Phone, Email: in.Email, Whatsapp: in.WhatsApp,
		TelegramChatID: in.TelegramChatID, SlackID: in.SlackID, Locale: string(in.Locale), Tags: in.Tags, Attributes: attrs,
	})
	if err != nil {
		if postgres.IsUniqueViolation(err) {
			return nil, domain.ErrConflict.WithMessage("a contact with this external_id already exists")
		}
		return nil, err
	}
	return &c, nil
}

// Upsert creates or updates by external_id (used by clients syncing users).
func (s *Service) Upsert(ctx context.Context, projectID uuid.UUID, in Input) (*sqlcgen.Contact, bool, error) {
	if in.ExternalID != "" {
		existing, err := s.db.Queries.GetContactByExternalID(ctx, sqlcgen.GetContactByExternalIDParams{ProjectID: projectID, ExternalID: strings.TrimSpace(in.ExternalID)})
		if err == nil {
			c, err := s.Update(ctx, projectID, existing.ID, in)
			return c, false, err
		}
		if !errors.Is(err, pgx.ErrNoRows) {
			return nil, false, err
		}
	}
	c, err := s.Create(ctx, projectID, in)
	return c, true, err
}

// Get returns one contact.
func (s *Service) Get(ctx context.Context, projectID, id uuid.UUID) (*sqlcgen.Contact, error) {
	c, err := s.db.Queries.GetContact(ctx, sqlcgen.GetContactParams{ID: id, ProjectID: projectID})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrNotFound.WithMessage("contact not found")
		}
		return nil, err
	}
	return &c, nil
}

// List filters by tag and free-text search.
func (s *Service) List(ctx context.Context, projectID uuid.UUID, tag, search string, limit, offset int32) ([]sqlcgen.Contact, error) {
	var tagPtr, searchPtr *string
	if tag != "" {
		tagPtr = &tag
	}
	if search != "" {
		searchPtr = &search
	}
	return s.db.Queries.ListContacts(ctx, sqlcgen.ListContactsParams{ProjectID: projectID, Tag: tagPtr, Search: searchPtr, RowLimit: limit, RowOffset: offset})
}

// Update replaces a contact's address fields. Tags and attributes that are
// omitted (nil) are kept, so clients syncing a phone number do not wipe
// the attributes another system wrote.
func (s *Service) Update(ctx context.Context, projectID, id uuid.UUID, in Input) (*sqlcgen.Contact, error) {
	keepTags, keepAttrs := in.Tags == nil, in.Attributes == nil
	if err := in.normalize(); err != nil {
		return nil, err
	}
	attrs, _ := json.Marshal(in.Attributes)
	if keepTags || keepAttrs {
		cur, err := s.Get(ctx, projectID, id)
		if err != nil {
			return nil, err
		}
		if keepTags {
			in.Tags = cur.Tags
		}
		if keepAttrs {
			attrs = cur.Attributes
		}
	}
	c, err := s.db.Queries.UpdateContact(ctx, sqlcgen.UpdateContactParams{
		ID: id, ProjectID: projectID, ExternalID: in.ExternalID, Name: in.Name, Phone: in.Phone, Email: in.Email, Whatsapp: in.WhatsApp,
		TelegramChatID: in.TelegramChatID, SlackID: in.SlackID, Locale: string(in.Locale), Tags: in.Tags, Attributes: attrs,
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrNotFound.WithMessage("contact not found")
		}
		if postgres.IsUniqueViolation(err) {
			return nil, domain.ErrConflict.WithMessage("a contact with this external_id already exists")
		}
		return nil, err
	}
	return &c, nil
}

// Delete soft-deletes a contact: the row is kept (deleted_at set) and hidden
// from every normal query, so the project keeps it as a record and can
// restore it. Devices keep existing. Use Purge for permanent removal.
func (s *Service) Delete(ctx context.Context, projectID, id uuid.UUID) error {
	n, err := s.db.Queries.DeleteContact(ctx, sqlcgen.DeleteContactParams{ID: id, ProjectID: projectID})
	if err != nil {
		return err
	}
	if n == 0 {
		return domain.ErrNotFound.WithMessage("contact not found")
	}
	return nil
}

// Restore brings a soft-deleted contact back.
func (s *Service) Restore(ctx context.Context, projectID, id uuid.UUID) (*sqlcgen.Contact, error) {
	c, err := s.db.Queries.RestoreContact(ctx, sqlcgen.RestoreContactParams{ID: id, ProjectID: projectID})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrNotFound.WithMessage("contact not in trash")
		}
		if postgres.IsUniqueViolation(err) {
			return nil, domain.ErrConflict.WithMessage("another active contact now uses this external_id")
		}
		return nil, err
	}
	return &c, nil
}

// ListDeleted returns the project's trashed contacts and the total.
func (s *Service) ListDeleted(ctx context.Context, projectID uuid.UUID, search string, limit, offset int32) ([]sqlcgen.Contact, int64, error) {
	params := sqlcgen.ListDeletedContactsParams{ProjectID: projectID, RowLimit: limit, RowOffset: offset}
	if search = strings.TrimSpace(search); search != "" {
		params.Search = &search
	}
	rows, err := s.db.Queries.ListDeletedContacts(ctx, params)
	if err != nil {
		return nil, 0, err
	}
	total, err := s.db.Queries.CountDeletedContacts(ctx, projectID)
	return rows, total, err
}

// Purge permanently deletes one trashed contact.
func (s *Service) Purge(ctx context.Context, projectID, id uuid.UUID) error {
	n, err := s.db.Queries.PurgeContact(ctx, sqlcgen.PurgeContactParams{ID: id, ProjectID: projectID})
	if err != nil {
		return err
	}
	if n == 0 {
		return domain.ErrNotFound.WithMessage("contact not in trash")
	}
	return nil
}

// PurgeAll permanently deletes every trashed contact of a project.
func (s *Service) PurgeAll(ctx context.Context, projectID uuid.UUID) (int64, error) {
	return s.db.Queries.PurgeDeletedContacts(ctx, projectID)
}

// DeviceInput registers or refreshes an FCM token.
type DeviceInput struct {
	Token      string
	Platform   domain.Platform
	AppVersion string
	ContactID  *uuid.UUID
	ExternalID string // alternative to ContactID
}

// RegisterDevice upserts a device by token. When ExternalID is given and no
// such contact exists, one is created so push can target it later.
func (s *Service) RegisterDevice(ctx context.Context, projectID uuid.UUID, in DeviceInput) (*sqlcgen.Device, error) {
	in.Token = strings.TrimSpace(in.Token)
	if len(in.Token) < 20 {
		return nil, domain.ErrValidation.WithDetails(map[string]any{"token": "fcm token required"})
	}
	if !in.Platform.Valid() {
		return nil, domain.ErrValidation.WithDetails(map[string]any{"platform": "android, ios or web"})
	}
	contactID := in.ContactID
	if contactID == nil && in.ExternalID != "" {
		c, _, err := s.Upsert(ctx, projectID, Input{ExternalID: in.ExternalID})
		if err != nil {
			return nil, err
		}
		contactID = &c.ID
	}
	if contactID != nil {
		if _, err := s.Get(ctx, projectID, *contactID); err != nil {
			return nil, err
		}
	}
	d, err := s.db.Queries.UpsertDevice(ctx, sqlcgen.UpsertDeviceParams{
		ProjectID: projectID, ContactID: contactID, Platform: sqlcgen.DevicePlatform(in.Platform), FcmToken: in.Token, AppVersion: in.AppVersion,
	})
	if err != nil {
		return nil, err
	}
	return &d, nil
}

// RemoveDevice deletes a token (app uninstall / logout).
func (s *Service) RemoveDevice(ctx context.Context, projectID uuid.UUID, token string) error {
	n, err := s.db.Queries.DeleteDeviceByToken(ctx, sqlcgen.DeleteDeviceByTokenParams{FcmToken: token, ProjectID: projectID})
	if err != nil {
		return err
	}
	if n == 0 {
		return domain.ErrNotFound.WithMessage("device not found")
	}
	return nil
}

// DisableDevice marks a token inactive after FCM reports UNREGISTERED.
func (s *Service) DisableDevice(ctx context.Context, token string) error {
	_, err := s.db.Queries.DisableDeviceByToken(ctx, token)
	return err
}

// ListDevices lists a project's devices.
func (s *Service) ListDevices(ctx context.Context, projectID uuid.UUID, limit, offset int32) ([]sqlcgen.Device, error) {
	return s.db.Queries.ListDevices(ctx, sqlcgen.ListDevicesParams{ProjectID: projectID, RowLimit: limit, RowOffset: offset})
}

// ActiveTokens returns the active FCM tokens of a contact.
func (s *Service) ActiveTokens(ctx context.Context, contactID uuid.UUID) ([]string, error) {
	devs, err := s.db.Queries.ListActiveDevicesForContact(ctx, &contactID)
	if err != nil {
		return nil, err
	}
	out := make([]string, 0, len(devs))
	for _, d := range devs {
		out = append(out, d.FcmToken)
	}
	return out, nil
}

// Merge finds an existing contact by external_id, phone or e-mail and fills
// its empty fields from in; otherwise it creates the contact. Used by
// imports and group quick-add so the same person is never duplicated.
func (s *Service) Merge(ctx context.Context, projectID uuid.UUID, in Input) (*sqlcgen.Contact, bool, error) {
	if err := in.normalize(); err != nil {
		return nil, false, err
	}
	var existing *sqlcgen.Contact
	if in.ExternalID != "" {
		if c, err := s.db.Queries.GetContactByExternalID(ctx, sqlcgen.GetContactByExternalIDParams{ProjectID: projectID, ExternalID: in.ExternalID}); err == nil {
			existing = &c
		}
	}
	if existing == nil && in.Phone != "" {
		if c, err := s.db.Queries.GetContactByPhone(ctx, sqlcgen.GetContactByPhoneParams{ProjectID: projectID, Phone: in.Phone}); err == nil {
			existing = &c
		}
	}
	if existing == nil && in.Email != "" {
		if c, err := s.db.Queries.GetContactByEmail(ctx, sqlcgen.GetContactByEmailParams{ProjectID: projectID, Email: in.Email}); err == nil {
			existing = &c
		}
	}
	if existing == nil {
		c, err := s.Create(ctx, projectID, in)
		return c, true, err
	}
	merged := Input{
		ExternalID: pick(existing.ExternalID, in.ExternalID), Name: pick(existing.Name, in.Name), Phone: pick(existing.Phone, in.Phone),
		Email: pick(existing.Email, in.Email), WhatsApp: pick(existing.Whatsapp, in.WhatsApp), TelegramChatID: pick(existing.TelegramChatID, in.TelegramChatID),
		SlackID: pick(existing.SlackID, in.SlackID), Locale: domain.Locale(existing.Locale), Tags: unionTags(existing.Tags, in.Tags),
	}
	unchanged := merged.ExternalID == existing.ExternalID && merged.Name == existing.Name && merged.Phone == existing.Phone && merged.Email == existing.Email &&
		merged.WhatsApp == existing.Whatsapp && merged.TelegramChatID == existing.TelegramChatID && merged.SlackID == existing.SlackID && len(merged.Tags) == len(existing.Tags)
	if unchanged {
		return existing, false, nil // nothing new
	}
	merged.Attributes = nil // keep stored attributes
	c, err := s.Update(ctx, projectID, existing.ID, merged)
	return c, false, err
}

func pick(current, incoming string) string {
	if current != "" {
		return current
	}
	return incoming
}

func unionTags(a, b []string) []string {
	seen := map[string]bool{}
	out := make([]string, 0, len(a)+len(b))
	for _, t := range append(append([]string{}, a...), b...) {
		if t != "" && !seen[t] {
			seen[t] = true
			out = append(out, t)
		}
	}
	return out
}
