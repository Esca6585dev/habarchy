package domain

import (
	"time"

	"github.com/google/uuid"
)

// User is an admin-panel user. Passwords are argon2id hashes; TOTP secrets
// are stored encrypted and never returned by the API.
type User struct {
	ID            uuid.UUID
	Email         string
	PasswordHash  string
	FullName      string
	IsActive      bool
	TOTPEnabled   bool
	TOTPSecretEnc []byte
	LastLoginAt   *time.Time
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

// Project is a tenant. Every API key, provider, template, contact and
// message belongs to exactly one project.
type Project struct {
	ID               uuid.UUID
	Name             string
	Slug             string
	Status           ProjectStatus
	DailyQuota       int64 // 0 = unlimited
	MonthlyQuota     int64 // 0 = unlimited
	WebhookURL       string
	WebhookSecretEnc []byte
	DefaultLocale    Locale
	AllowedIPs       []string  // CIDRs; empty = any
	AutoChannelOrder []Channel // order used by channel=auto
	CreatedAt        time.Time
	UpdatedAt        time.Time
}

// ProjectMember links a user to a project with a role.
type ProjectMember struct {
	ProjectID uuid.UUID
	UserID    uuid.UUID
	Role      Role
	CreatedAt time.Time
}

// APIKey authenticates a client application. Only the SHA-256 hash is
// stored; the plaintext is shown once at creation.
type APIKey struct {
	ID          uuid.UUID
	ProjectID   uuid.UUID
	Name        string
	Prefix      string // hb_live_ or hb_test_
	Hint        string // last 4 chars, for display
	KeyHash     []byte
	Scopes      []APIKeyScope
	IPAllowlist []string
	LastUsedAt  *time.Time
	ExpiresAt   *time.Time
	RevokedAt   *time.Time
	CreatedAt   time.Time
}

// IsLive reports whether the key may reach real providers.
func (k APIKey) IsLive() bool { return k.Prefix == APIKeyPrefixLive }

// Active reports whether the key can still be used right now.
func (k APIKey) Active(now time.Time) bool {
	if k.RevokedAt != nil {
		return false
	}
	return k.ExpiresAt == nil || now.Before(*k.ExpiresAt)
}

// HasScope reports whether the key grants s.
func (k APIKey) HasScope(s APIKeyScope) bool {
	for _, have := range k.Scopes {
		if have == s {
			return true
		}
	}
	return false
}

// Provider is a configured delivery backend for one channel of a project.
// Credentials are AES-256-GCM encrypted JSON; the shape depends on Type.
type Provider struct {
	ID              uuid.UUID
	ProjectID       uuid.UUID
	Name            string
	Channel         Channel
	Type            ProviderType
	Priority        int // lower = tried first
	IsActive        bool
	CredentialsEnc  []byte
	RateLimitPerSec int // 0 = unlimited
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

// Template is the active version of a message template for one
// (project, key, channel, locale). Older versions live in TemplateVersion.
type Template struct {
	ID           uuid.UUID
	ProjectID    uuid.UUID
	Key          string
	Channel      Channel
	Locale       Locale
	Subject      string
	Body         string
	RequiredVars []string
	Version      int
	IsActive     bool
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

// TemplateVersion is an immutable historical snapshot of a template.
type TemplateVersion struct {
	ID           uuid.UUID
	TemplateID   uuid.UUID
	Version      int
	Subject      string
	Body         string
	RequiredVars []string
	CreatedBy    *uuid.UUID
	CreatedAt    time.Time
}

// Contact is an address-book entry of a project.
type Contact struct {
	ID             uuid.UUID
	ProjectID      uuid.UUID
	ExternalID     string
	Phone          string // E.164
	Email          string
	TelegramChatID string
	Locale         Locale
	Tags           []string
	Attributes     map[string]any
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

// Device is a push-capable device of a contact.
type Device struct {
	ID         uuid.UUID
	ProjectID  uuid.UUID
	ContactID  *uuid.UUID
	Platform   Platform
	FCMToken   string
	AppVersion string
	LastSeenAt time.Time
	IsActive   bool
	CreatedAt  time.Time
}

// Message is the central entity: one delivery attempt chain to one address.
type Message struct {
	ID                uuid.UUID // UUID v7, time-ordered
	ProjectID         uuid.UUID
	BatchID           *uuid.UUID
	Channel           Channel
	ToAddress         string // normalized phone / lowercase email / token / chat id
	ContactID         *uuid.UUID
	TemplateKey       string
	TemplateVersion   *int
	RenderedSubject   string
	RenderedBody      string
	Status            MessageStatus
	Priority          Priority
	ProviderID        *uuid.UUID
	ProviderMessageID string
	ErrorCode         string
	ErrorMessage      string
	Attempts          int
	ScheduledAt       *time.Time
	SentAt            *time.Time
	DeliveredAt       *time.Time
	CostMicros        int64
	Currency          string
	Metadata          map[string]any
	IdempotencyKey    string
	IsTest            bool
	CreatedAt         time.Time
	UpdatedAt         time.Time
}

// MessageEvent is one entry of a message's timeline.
type MessageEvent struct {
	ID         uuid.UUID
	MessageID  uuid.UUID
	Type       EventType
	ProviderID *uuid.UUID
	Payload    map[string]any
	CreatedAt  time.Time
}

// Batch groups messages created by one POST /messages/batch call.
type Batch struct {
	ID             uuid.UUID
	ProjectID      uuid.UUID
	TemplateKey    string
	Channel        Channel
	Total          int
	Queued         int
	Sent           int
	Delivered      int
	Failed         int
	Status         string // processing | completed
	IdempotencyKey string
	CreatedAt      time.Time
	CompletedAt    *time.Time
}

// WebhookDelivery is one outbound webhook call and its retry state.
type WebhookDelivery struct {
	ID           uuid.UUID
	ProjectID    uuid.UUID
	MessageID    *uuid.UUID
	BatchID      *uuid.UUID
	Event        WebhookEvent
	URL          string
	Payload      []byte
	Signature    string
	ResponseCode *int
	ResponseBody string
	Attempts     int
	NextRetryAt  *time.Time
	DeliveredAt  *time.Time
	CreatedAt    time.Time
}

// AuditLog records who did what in the admin panel.
type AuditLog struct {
	ID         uuid.UUID
	ProjectID  *uuid.UUID
	UserID     *uuid.UUID
	Action     string
	EntityType string
	EntityID   string
	Changes    map[string]any
	IP         string
	UserAgent  string
	CreatedAt  time.Time
}

// UsageDaily is the aggregated counter row for (project, channel, day).
type UsageDaily struct {
	ProjectID  uuid.UUID
	Day        time.Time
	Channel    Channel
	Queued     int64
	Sent       int64
	Delivered  int64
	Failed     int64
	CostMicros int64
	Currency   string
}
