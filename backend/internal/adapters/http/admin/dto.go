// Package admin serves /api/admin for the web and mobile admin apps.
package admin

import (
	"time"

	"github.com/google/uuid"

	"github.com/Esca6585dev/habarchy/backend/internal/adapters/postgres/sqlcgen"
	"github.com/Esca6585dev/habarchy/backend/internal/domain"
)

// ---- auth ----

type loginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
	TOTPCode string `json:"totp_code"`
}

type refreshRequest struct {
	RefreshToken string `json:"refresh_token" validate:"required"`
}

type changePasswordRequest struct {
	CurrentPassword string `json:"current_password" validate:"required"`
	NewPassword     string `json:"new_password" validate:"required,min=8,max=256"`
}

type totpCodeRequest struct {
	Code string `json:"code" validate:"required,len=6,numeric"`
}

type passwordRequest struct {
	Password string `json:"password" validate:"required"`
}

// UserResponse is the public shape of a user.
type UserResponse struct {
	ID          uuid.UUID  `json:"id"`
	Email       string     `json:"email"`
	FullName    string     `json:"full_name"`
	IsActive    bool       `json:"is_active"`
	TOTPEnabled bool       `json:"totp_enabled"`
	LastLoginAt *time.Time `json:"last_login_at"`
	CreatedAt   time.Time  `json:"created_at"`
}

func toUser(u *sqlcgen.User) UserResponse {
	return UserResponse{ID: u.ID, Email: u.Email, FullName: u.FullName, IsActive: u.IsActive, TOTPEnabled: u.TotpEnabled, LastLoginAt: u.LastLoginAt, CreatedAt: u.CreatedAt}
}

// ---- projects ----

type createProjectRequest struct {
	Name          string        `json:"name" validate:"required,min=2,max=100"`
	Slug          string        `json:"slug" validate:"max=60"`
	DailyQuota    int64         `json:"daily_quota" validate:"min=0"`
	MonthlyQuota  int64         `json:"monthly_quota" validate:"min=0"`
	DefaultLocale domain.Locale `json:"default_locale"`
}

type updateProjectRequest struct {
	Name             *string               `json:"name" validate:"omitempty,min=2,max=100"`
	Status           *domain.ProjectStatus `json:"status"`
	DailyQuota       *int64                `json:"daily_quota" validate:"omitempty,min=0"`
	MonthlyQuota     *int64                `json:"monthly_quota" validate:"omitempty,min=0"`
	WebhookURL       *string               `json:"webhook_url" validate:"omitempty,max=2048"`
	WebhookSecret    *string               `json:"webhook_secret" validate:"omitempty,max=256"`
	DefaultLocale    *domain.Locale        `json:"default_locale"`
	AllowedIPs       *[]string             `json:"allowed_ips"`
	AutoChannelOrder *[]domain.Channel     `json:"auto_channel_order"`
}

// ProjectResponse is the public shape of a project.
type ProjectResponse struct {
	ID               uuid.UUID         `json:"id"`
	Name             string            `json:"name"`
	Slug             string            `json:"slug"`
	Status           string            `json:"status"`
	DailyQuota       int64             `json:"daily_quota"`
	MonthlyQuota     int64             `json:"monthly_quota"`
	WebhookURL       string            `json:"webhook_url"`
	HasWebhookSecret bool              `json:"has_webhook_secret"`
	DefaultLocale    string            `json:"default_locale"`
	AllowedIPs       []string          `json:"allowed_ips"`
	AutoChannelOrder []sqlcgen.Channel `json:"auto_channel_order"`
	Role             domain.Role       `json:"role,omitempty"`
	CreatedAt        time.Time         `json:"created_at"`
	UpdatedAt        time.Time         `json:"updated_at"`
}

func toProject(p *sqlcgen.Project, role domain.Role) ProjectResponse {
	return ProjectResponse{
		ID: p.ID, Name: p.Name, Slug: p.Slug, Status: string(p.Status), DailyQuota: p.DailyQuota, MonthlyQuota: p.MonthlyQuota,
		WebhookURL: p.WebhookUrl, HasWebhookSecret: len(p.WebhookSecretEnc) > 0, DefaultLocale: p.DefaultLocale,
		AllowedIPs: p.AllowedIps, AutoChannelOrder: p.AutoChannelOrder, Role: role, CreatedAt: p.CreatedAt, UpdatedAt: p.UpdatedAt,
	}
}

type setMemberRequest struct {
	Email string      `json:"email" validate:"required,email"`
	Role  domain.Role `json:"role" validate:"required"`
}

// ---- api keys ----

type createAPIKeyRequest struct {
	Name             string               `json:"name" validate:"required,min=1,max=100"`
	Live             bool                 `json:"live"`
	Scopes           []domain.APIKeyScope `json:"scopes"`
	IPAllowlist      []string             `json:"ip_allowlist"`
	ExpiresAt        *time.Time           `json:"expires_at"`
	RequireSignature bool                 `json:"require_signature"`
}

// APIKeyResponse never includes the hash. Key is set only on creation.
type APIKeyResponse struct {
	ID               uuid.UUID  `json:"id"`
	Name             string     `json:"name"`
	Prefix           string     `json:"prefix"`
	Hint             string     `json:"hint"`
	Key              string     `json:"key,omitempty"`
	Scopes           []string   `json:"scopes"`
	IPAllowlist      []string   `json:"ip_allowlist"`
	RequireSignature bool       `json:"require_signature"`
	LastUsedAt       *time.Time `json:"last_used_at"`
	ExpiresAt        *time.Time `json:"expires_at"`
	RevokedAt        *time.Time `json:"revoked_at"`
	CreatedAt        time.Time  `json:"created_at"`
}

func toAPIKey(k *sqlcgen.ApiKey, plaintext string) APIKeyResponse {
	return APIKeyResponse{
		ID: k.ID, Name: k.Name, Prefix: k.Prefix, Hint: k.Hint, Key: plaintext, Scopes: k.Scopes, IPAllowlist: k.IpAllowlist,
		RequireSignature: k.RequireSignature, LastUsedAt: k.LastUsedAt, ExpiresAt: k.ExpiresAt, RevokedAt: k.RevokedAt, CreatedAt: k.CreatedAt,
	}
}

// ---- templates ----

// TemplateRequest is shared by admin and public APIs.
type TemplateRequest struct {
	Key      string         `json:"key" validate:"required_without=IsUpdate,max=64"`
	Channel  domain.Channel `json:"channel"`
	Locale   domain.Locale  `json:"locale"`
	Subject  string         `json:"subject" validate:"max=998"`
	Body     string         `json:"body" validate:"max=65536"`
	IsActive *bool          `json:"is_active"`
	IsUpdate bool           `json:"-"`
}

// TemplateResponse is the public shape of a template.
type TemplateResponse struct {
	ID           uuid.UUID `json:"id"`
	Key          string    `json:"key"`
	Channel      string    `json:"channel"`
	Locale       string    `json:"locale"`
	Subject      string    `json:"subject"`
	Body         string    `json:"body"`
	RequiredVars []string  `json:"required_vars"`
	Version      int32     `json:"version"`
	IsActive     bool      `json:"is_active"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// ToTemplate converts a row.
func ToTemplate(t *sqlcgen.Template) TemplateResponse {
	return TemplateResponse{
		ID: t.ID, Key: t.Key, Channel: string(t.Channel), Locale: t.Locale, Subject: t.Subject, Body: t.Body,
		RequiredVars: t.RequiredVars, Version: t.Version, IsActive: t.IsActive, CreatedAt: t.CreatedAt, UpdatedAt: t.UpdatedAt,
	}
}

// TemplateVersionResponse is one history entry.
type TemplateVersionResponse struct {
	Version      int32      `json:"version"`
	Subject      string     `json:"subject"`
	Body         string     `json:"body"`
	RequiredVars []string   `json:"required_vars"`
	CreatedBy    *uuid.UUID `json:"created_by"`
	CreatedAt    time.Time  `json:"created_at"`
}

// ToTemplateVersions converts rows.
func ToTemplateVersions(rows []sqlcgen.TemplateVersion) []TemplateVersionResponse {
	out := make([]TemplateVersionResponse, 0, len(rows))
	for i := range rows {
		r := &rows[i]
		out = append(out, TemplateVersionResponse{Version: r.Version, Subject: r.Subject, Body: r.Body, RequiredVars: r.RequiredVars, CreatedBy: r.CreatedBy, CreatedAt: r.CreatedAt})
	}
	return out
}

// PreviewRequest renders either a stored template (by id in the path) or
// ad hoc content.
type PreviewRequest struct {
	Channel domain.Channel `json:"channel"`
	Subject string         `json:"subject" validate:"max=998"`
	Body    string         `json:"body" validate:"max=65536"`
	Data    map[string]any `json:"data"`
}
