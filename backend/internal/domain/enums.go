// Package domain holds the core entities, enumerations and errors of
// Habarchy. It has no dependencies on frameworks, databases or transports.
package domain

import "fmt"

// Channel is a delivery channel supported by Habarchy.
type Channel string

const (
	ChannelSMS      Channel = "sms"
	ChannelEmail    Channel = "email"
	ChannelPush     Channel = "push"
	ChannelTelegram Channel = "telegram"
	// ChannelAuto is only valid on the public API: Habarchy picks the best
	// channel for the contact. It is never stored on a message.
	ChannelAuto Channel = "auto"
)

// AllChannels lists the channels a message can actually be delivered on.
var AllChannels = []Channel{ChannelSMS, ChannelEmail, ChannelPush, ChannelTelegram}

// DefaultAutoOrder is the order "auto" tries channels when a project has not
// configured its own.
var DefaultAutoOrder = []Channel{ChannelPush, ChannelSMS, ChannelEmail}

// Valid reports whether c is a deliverable channel (auto is excluded).
func (c Channel) Valid() bool {
	switch c {
	case ChannelSMS, ChannelEmail, ChannelPush, ChannelTelegram:
		return true
	}
	return false
}

// ParseChannel converts a string into a Channel, accepting "auto".
func ParseChannel(s string) (Channel, error) {
	c := Channel(s)
	if c == ChannelAuto || c.Valid() {
		return c, nil
	}
	return "", fmt.Errorf("%w: unknown channel %q", ErrValidation, s)
}

func (c Channel) String() string { return string(c) }

// MessageStatus is the lifecycle state of a message.
type MessageStatus string

const (
	StatusQueued     MessageStatus = "queued"
	StatusProcessing MessageStatus = "processing"
	StatusSent       MessageStatus = "sent"
	StatusDelivered  MessageStatus = "delivered"
	StatusFailed     MessageStatus = "failed"
	StatusCancelled  MessageStatus = "cancelled"
)

// Valid reports whether s is a known status.
func (s MessageStatus) Valid() bool {
	switch s {
	case StatusQueued, StatusProcessing, StatusSent, StatusDelivered, StatusFailed, StatusCancelled:
		return true
	}
	return false
}

// Terminal reports whether no further transitions are expected.
func (s MessageStatus) Terminal() bool {
	return s == StatusDelivered || s == StatusFailed || s == StatusCancelled
}

// CanTransitionTo enforces the allowed status graph.
func (s MessageStatus) CanTransitionTo(next MessageStatus) bool {
	switch s {
	case StatusQueued:
		return next == StatusProcessing || next == StatusCancelled || next == StatusFailed
	case StatusProcessing:
		// processing -> queued means a retry was scheduled.
		return next == StatusSent || next == StatusFailed || next == StatusQueued
	case StatusSent:
		return next == StatusDelivered || next == StatusFailed
	}
	return false
}

func (s MessageStatus) String() string { return string(s) }

// Priority controls queue ordering.
type Priority string

const (
	PriorityHigh   Priority = "high"
	PriorityNormal Priority = "normal"
	PriorityLow    Priority = "low"
)

// Valid reports whether p is a known priority.
func (p Priority) Valid() bool {
	return p == PriorityHigh || p == PriorityNormal || p == PriorityLow
}

func (p Priority) String() string { return string(p) }

// ProviderType identifies a concrete provider implementation.
type ProviderType string

const (
	ProviderHTTPSMS     ProviderType = "http_sms"
	ProviderSMPP        ProviderType = "smpp"
	ProviderSMTP        ProviderType = "smtp"
	ProviderFCM         ProviderType = "fcm"
	ProviderTelegramBot ProviderType = "telegram_bot"
)

// Channel returns the channel a provider type serves.
func (t ProviderType) Channel() Channel {
	switch t {
	case ProviderHTTPSMS, ProviderSMPP:
		return ChannelSMS
	case ProviderSMTP:
		return ChannelEmail
	case ProviderFCM:
		return ChannelPush
	case ProviderTelegramBot:
		return ChannelTelegram
	}
	return ""
}

// Valid reports whether t is a known provider type.
func (t ProviderType) Valid() bool { return t.Channel() != "" }

func (t ProviderType) String() string { return string(t) }

// Role is a project membership role, ordered from most to least privileged.
type Role string

const (
	RoleOwner     Role = "owner"
	RoleAdmin     Role = "admin"
	RoleDeveloper Role = "developer"
	RoleViewer    Role = "viewer"
)

var roleRank = map[Role]int{RoleOwner: 4, RoleAdmin: 3, RoleDeveloper: 2, RoleViewer: 1}

// Valid reports whether r is a known role.
func (r Role) Valid() bool { _, ok := roleRank[r]; return ok }

// AtLeast reports whether r grants at least the privileges of minimum.
func (r Role) AtLeast(minimum Role) bool { return roleRank[r] >= roleRank[minimum] }

func (r Role) String() string { return string(r) }

// Platform is a device platform for push tokens.
type Platform string

const (
	PlatformAndroid Platform = "android"
	PlatformIOS     Platform = "ios"
	PlatformWeb     Platform = "web"
)

// Valid reports whether p is a known platform.
func (p Platform) Valid() bool {
	return p == PlatformAndroid || p == PlatformIOS || p == PlatformWeb
}

func (p Platform) String() string { return string(p) }

// Locale is a supported UI / template locale.
type Locale string

const (
	LocaleTK Locale = "tk"
	LocaleRU Locale = "ru"
	LocaleEN Locale = "en"
)

// DefaultLocale is used when neither contact nor request specify one.
const DefaultLocale = LocaleTK

// Valid reports whether l is a supported locale.
func (l Locale) Valid() bool { return l == LocaleTK || l == LocaleRU || l == LocaleEN }

func (l Locale) String() string { return string(l) }

// ProjectStatus is the lifecycle state of a tenant.
type ProjectStatus string

const (
	ProjectActive    ProjectStatus = "active"
	ProjectSuspended ProjectStatus = "suspended"
	ProjectArchived  ProjectStatus = "archived"
)

// Valid reports whether s is a known project status.
func (s ProjectStatus) Valid() bool {
	return s == ProjectActive || s == ProjectSuspended || s == ProjectArchived
}

// EventType is a message timeline event.
type EventType string

const (
	EventQueued           EventType = "queued"
	EventAttempt          EventType = "attempt"
	EventProviderResponse EventType = "provider_response"
	EventSent             EventType = "sent"
	EventDelivered        EventType = "delivered"
	EventFailed           EventType = "failed"
	EventCancelled        EventType = "cancelled"
	EventWebhookSent      EventType = "webhook_sent"
)

// WebhookEvent is an outbound webhook event name.
type WebhookEvent string

const (
	WebhookMessageSent      WebhookEvent = "message.sent"
	WebhookMessageDelivered WebhookEvent = "message.delivered"
	WebhookMessageFailed    WebhookEvent = "message.failed"
	WebhookBatchCompleted   WebhookEvent = "batch.completed"
)

// APIKeyScope restricts what an API key may do.
type APIKeyScope string

const (
	ScopeMessagesSend APIKeyScope = "messages:send"
	ScopeMessagesRead APIKeyScope = "messages:read"
	ScopeOTP          APIKeyScope = "otp"
	ScopeTemplates    APIKeyScope = "templates"
	ScopeContacts     APIKeyScope = "contacts"
	ScopeDevices      APIKeyScope = "devices"
	ScopeUsage        APIKeyScope = "usage"
)

// AllScopes is the full scope set granted to a key created without an
// explicit list.
var AllScopes = []APIKeyScope{
	ScopeMessagesSend, ScopeMessagesRead, ScopeOTP, ScopeTemplates, ScopeContacts, ScopeDevices, ScopeUsage,
}

// Valid reports whether s is a known scope.
func (s APIKeyScope) Valid() bool {
	for _, k := range AllScopes {
		if k == s {
			return true
		}
	}
	return false
}

// API key prefixes. The environment is encoded in the key so test traffic
// can never hit live providers by mistake.
const (
	APIKeyPrefixLive = "hb_live_"
	APIKeyPrefixTest = "hb_test_"
)
