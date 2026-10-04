// Package ports defines the interfaces between the application core and
// the outside world. Adapters implement them; use cases depend on them.
package ports

import (
	"context"
	"time"

	"github.com/google/uuid"

	"github.com/Esca6585dev/habarchy/backend/internal/domain"
)

// SendResult is what a provider returns after a successful hand-over.
type SendResult struct {
	ProviderMessageID string
	CostMicros        int64
	Currency          string
	// Raw is the provider response kept for the message timeline. Never
	// include credentials in it.
	Raw map[string]any
}

// ProviderError is returned by providers so the worker can decide between
// retrying, falling back to the next provider or failing permanently.
type ProviderError struct {
	Code      string // short machine code, e.g. "invalid_number", "timeout"
	Message   string
	Retryable bool
	Raw       map[string]any
}

func (e *ProviderError) Error() string { return e.Code + ": " + e.Message }

// SMSMessage is the channel-specific payload for SMS providers.
type SMSMessage struct {
	To     string // E.164
	Text   string
	Sender string
	// Ref is the Habarchy message id, for providers that need to correlate
	// asynchronous results (the Android gateway outbox).
	Ref string
}

// EmailMessage is the payload for email providers.
type EmailMessage struct {
	To          string
	Subject     string
	HTML        string
	Text        string
	ReplyTo     string
	Attachments []EmailAttachment
}

// EmailAttachment is an inline (base64) or remote (URL) attachment.
type EmailAttachment struct {
	Filename    string
	ContentType string
	Content     []byte
	URL         string
}

// PushMessage is the payload for push providers.
type PushMessage struct {
	Tokens  []string
	Topic   string
	Title   string
	Body    string
	Data    map[string]string
	Android map[string]any
	APNs    map[string]any
}

// PushResult reports per-token outcomes so stale tokens can be disabled.
type PushResult struct {
	SendResult
	InvalidTokens []string
}

// TelegramMessage is the payload for Telegram providers.
type TelegramMessage struct {
	ChatID    string
	Text      string
	ParseMode string // Markdown | HTML | ""
	Buttons   [][]TelegramButton
}

// TelegramButton is an inline keyboard button.
type TelegramButton struct {
	Text string
	URL  string
	Data string
}

// ChatMessage is the payload for chat-style providers (WhatsApp, Slack).
type ChatMessage struct {
	To      string // phone (WhatsApp) or channel / user id (Slack)
	Text    string
	Subject string         // optional heading (Slack bold line)
	Extra   map[string]any // provider specific, from message metadata
}

// ChatProvider sends chat messages.
type ChatProvider interface {
	Send(ctx context.Context, msg ChatMessage) (*SendResult, error)
}

// SMSProvider sends SMS.
type SMSProvider interface {
	Send(ctx context.Context, msg SMSMessage) (*SendResult, error)
}

// EmailProvider sends email.
type EmailProvider interface {
	Send(ctx context.Context, msg EmailMessage) (*SendResult, error)
}

// PushProvider sends push notifications.
type PushProvider interface {
	Send(ctx context.Context, msg PushMessage) (*PushResult, error)
}

// TelegramProvider sends Telegram messages.
type TelegramProvider interface {
	Send(ctx context.Context, msg TelegramMessage) (*SendResult, error)
}

// Queue enqueues background work. Implemented by the asynq adapter.
type Queue interface {
	EnqueueSend(ctx context.Context, messageID uuid.UUID, channel domain.Channel, priority domain.Priority, at *time.Time) error
	EnqueueWebhook(ctx context.Context, deliveryID uuid.UUID) error
}

// Cipher encrypts secrets at rest.
type Cipher interface {
	Encrypt(plaintext, aad []byte) ([]byte, error)
	Decrypt(blob, aad []byte) ([]byte, error)
}
