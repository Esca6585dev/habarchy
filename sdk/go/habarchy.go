// Package habarchy is a minimal client for the Habarchy public API
// (https://github.com/Esca6585dev/habarchy): SendMessage, SendBatch,
// GetMessage, SendOTP, VerifyOTP, RegisterDevice, with optional HMAC
// request signing.
package habarchy

import (
	"bytes"
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

// Client talks to one Habarchy project (one API key).
type Client struct {
	BaseURL string
	APIKey  string
	// Sign adds X-Timestamp / X-Signature to every request (required when
	// the key was created with require_signature).
	Sign bool
	HTTP *http.Client
	Now  func() time.Time
}

// Option configures the client.
type Option func(*Client)

// WithSigning enables HMAC request signatures.
func WithSigning() Option { return func(c *Client) { c.Sign = true } }

// WithHTTPClient replaces the HTTP client (timeouts, proxies, tests).
func WithHTTPClient(h *http.Client) Option { return func(c *Client) { c.HTTP = h } }

// New creates a client. baseURL is e.g. "https://habarchy.example.tm".
func New(baseURL, apiKey string, opts ...Option) *Client {
	c := &Client{BaseURL: strings.TrimRight(baseURL, "/"), APIKey: apiKey, HTTP: &http.Client{Timeout: 20 * time.Second}, Now: time.Now}
	for _, o := range opts {
		o(c)
	}
	return c
}

// Error is a non-2xx response from the API.
type Error struct {
	Status  int            `json:"-"`
	Code    string         `json:"code"`
	Message string         `json:"message"`
	Details map[string]any `json:"details,omitempty"`
}

func (e *Error) Error() string {
	return fmt.Sprintf("habarchy: %s (%d): %s", e.Code, e.Status, e.Message)
}

// Recipient is "to": either an address or a contact reference.
type Recipient struct {
	Address    string `json:"address,omitempty"`
	ContactID  string `json:"contact_id,omitempty"`
	ExternalID string `json:"external_id,omitempty"`
}

// To addresses a phone number (E.164), e-mail, Telegram chat id or FCM token.
func To(address string) Recipient { return Recipient{Address: address} }

// ToContact addresses a stored contact by id.
func ToContact(id string) Recipient { return Recipient{ContactID: id} }

// ToExternal addresses a contact by your own user id (external_id).
func ToExternal(id string) Recipient { return Recipient{ExternalID: id} }

// MarshalJSON encodes a bare address as a string.
func (r Recipient) MarshalJSON() ([]byte, error) {
	if r.ContactID == "" && r.ExternalID == "" {
		return json.Marshal(r.Address)
	}
	type raw Recipient
	return json.Marshal(raw(r))
}

// SendMessageRequest is POST /api/v1/messages.
type SendMessageRequest struct {
	Channel        string         `json:"channel"` // sms | email | push | telegram | auto
	To             Recipient      `json:"to"`
	Template       string         `json:"template,omitempty"`
	Data           map[string]any `json:"data,omitempty"`
	Subject        string         `json:"subject,omitempty"`
	Title          string         `json:"title,omitempty"`
	Body           string         `json:"body,omitempty"`
	Locale         string         `json:"locale,omitempty"`
	ScheduledAt    *time.Time     `json:"scheduled_at,omitempty"`
	Priority       string         `json:"priority,omitempty"` // high | normal | low
	IdempotencyKey string         `json:"idempotency_key,omitempty"`
	Metadata       map[string]any `json:"metadata,omitempty"`
}

// Accepted is the 202 response of SendMessage.
type Accepted struct {
	ID          string     `json:"id"`
	Status      string     `json:"status"`
	Channel     string     `json:"channel"`
	To          string     `json:"to"`
	ScheduledAt *time.Time `json:"scheduled_at,omitempty"`
	// Duplicate is true when the idempotency key was seen before (200).
	Duplicate bool `json:"-"`
}

// Message is the full message resource.
type Message struct {
	ID                string          `json:"id"`
	Status            string          `json:"status"`
	Channel           string          `json:"channel"`
	To                string          `json:"to"`
	ContactID         string          `json:"contact_id,omitempty"`
	BatchID           string          `json:"batch_id,omitempty"`
	Template          string          `json:"template,omitempty"`
	Subject           string          `json:"subject,omitempty"`
	Body              string          `json:"body"`
	Priority          string          `json:"priority"`
	ProviderID        string          `json:"provider_id,omitempty"`
	ProviderMessageID string          `json:"provider_message_id,omitempty"`
	ErrorCode         string          `json:"error_code,omitempty"`
	ErrorMessage      string          `json:"error_message,omitempty"`
	Attempts          int             `json:"attempts"`
	ScheduledAt       *time.Time      `json:"scheduled_at,omitempty"`
	SentAt            *time.Time      `json:"sent_at,omitempty"`
	DeliveredAt       *time.Time      `json:"delivered_at,omitempty"`
	CostMicros        int64           `json:"cost_micros"`
	Currency          string          `json:"currency"`
	Metadata          json.RawMessage `json:"metadata"`
	IsTest            bool            `json:"is_test"`
	CreatedAt         time.Time       `json:"created_at"`
	UpdatedAt         time.Time       `json:"updated_at"`
}

// Event is one timeline entry.
type Event struct {
	ID         string          `json:"id"`
	Type       string          `json:"type"`
	ProviderID string          `json:"provider_id,omitempty"`
	Payload    json.RawMessage `json:"payload"`
	CreatedAt  time.Time       `json:"created_at"`
}

// MessageDetail is GET /api/v1/messages/{id}.
type MessageDetail struct {
	Message Message `json:"message"`
	Events  []Event `json:"events"`
}

// BatchRecipient is one entry of SendBatch.
type BatchRecipient struct {
	To   Recipient      `json:"to"`
	Data map[string]any `json:"data,omitempty"`
}

// SendBatchRequest is POST /api/v1/messages/batch.
type SendBatchRequest struct {
	Channel        string           `json:"channel"`
	Template       string           `json:"template,omitempty"`
	Subject        string           `json:"subject,omitempty"`
	Title          string           `json:"title,omitempty"`
	Body           string           `json:"body,omitempty"`
	Locale         string           `json:"locale,omitempty"`
	ScheduledAt    *time.Time       `json:"scheduled_at,omitempty"`
	Priority       string           `json:"priority,omitempty"`
	IdempotencyKey string           `json:"idempotency_key,omitempty"`
	Metadata       map[string]any   `json:"metadata,omitempty"`
	Recipients     []BatchRecipient `json:"recipients"`
}

// Batch is the batch resource.
type Batch struct {
	ID          string     `json:"id"`
	Status      string     `json:"status"`
	Channel     string     `json:"channel"`
	Template    string     `json:"template,omitempty"`
	Total       int        `json:"total"`
	Queued      int        `json:"queued"`
	Sent        int        `json:"sent"`
	Delivered   int        `json:"delivered"`
	Failed      int        `json:"failed"`
	CreatedAt   time.Time  `json:"created_at"`
	CompletedAt *time.Time `json:"completed_at,omitempty"`
	// Accepted / Rejected come from the response meta of SendBatch.
	Accepted int `json:"-"`
	Rejected int `json:"-"`
}

// SendOTPRequest is POST /api/v1/otp/send.
type SendOTPRequest struct {
	Channel  string         `json:"channel,omitempty"` // sms (default) | email | telegram
	To       string         `json:"to"`
	Length   int            `json:"length,omitempty"`
	TTL      int            `json:"ttl,omitempty"` // seconds
	Template string         `json:"template,omitempty"`
	Locale   string         `json:"locale,omitempty"`
	Data     map[string]any `json:"data,omitempty"`
}

// OTPSent is the response of SendOTP.
type OTPSent struct {
	MessageID string `json:"message_id"`
	To        string `json:"to"`
	Channel   string `json:"channel"`
	ExpiresIn int    `json:"expires_in"`
	Length    int    `json:"length"`
}

// OTPVerified is the response of VerifyOTP.
type OTPVerified struct {
	Verified          bool `json:"verified"`
	AttemptsRemaining int  `json:"attempts_remaining"`
}

// RegisterDeviceRequest is POST /api/v1/devices.
type RegisterDeviceRequest struct {
	Token      string `json:"token"`
	Platform   string `json:"platform"` // android | ios | web
	AppVersion string `json:"app_version,omitempty"`
	ContactID  string `json:"contact_id,omitempty"`
	ExternalID string `json:"external_id,omitempty"`
}

// Device is the device resource (the token is never echoed back).
type Device struct {
	ID         string    `json:"id"`
	ContactID  string    `json:"contact_id,omitempty"`
	Platform   string    `json:"platform"`
	TokenHint  string    `json:"token_hint"`
	AppVersion string    `json:"app_version"`
	IsActive   bool      `json:"is_active"`
	LastSeenAt time.Time `json:"last_seen_at"`
	CreatedAt  time.Time `json:"created_at"`
}

// SendMessage queues one message (202). With an idempotency key a repeated
// call returns the original message with Duplicate=true.
func (c *Client) SendMessage(ctx context.Context, req SendMessageRequest) (*Accepted, error) {
	var out Accepted
	meta, err := c.do(ctx, http.MethodPost, "/api/v1/messages", req, &out)
	if err != nil {
		return nil, err
	}
	out.Duplicate, _ = meta["duplicate"].(bool)
	return &out, nil
}

// SendBatch queues up to batch_max_recipients messages at once.
func (c *Client) SendBatch(ctx context.Context, req SendBatchRequest) (*Batch, error) {
	var out Batch
	meta, err := c.do(ctx, http.MethodPost, "/api/v1/messages/batch", req, &out)
	if err != nil {
		return nil, err
	}
	out.Accepted = metaInt(meta, "accepted")
	out.Rejected = metaInt(meta, "rejected")
	return &out, nil
}

// GetMessage returns the message and its timeline.
func (c *Client) GetMessage(ctx context.Context, id string) (*MessageDetail, error) {
	var out MessageDetail
	if _, err := c.do(ctx, http.MethodGet, "/api/v1/messages/"+url.PathEscape(id), nil, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// CancelMessage cancels a queued / scheduled message.
func (c *Client) CancelMessage(ctx context.Context, id string) (*Message, error) {
	var out Message
	if _, err := c.do(ctx, http.MethodPost, "/api/v1/messages/"+url.PathEscape(id)+"/cancel", nil, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// GetBatch returns batch counters.
func (c *Client) GetBatch(ctx context.Context, id string) (*Batch, error) {
	var out Batch
	if _, err := c.do(ctx, http.MethodGet, "/api/v1/batches/"+url.PathEscape(id), nil, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// SendOTP generates and sends a one-time code.
func (c *Client) SendOTP(ctx context.Context, req SendOTPRequest) (*OTPSent, error) {
	var out OTPSent
	if _, err := c.do(ctx, http.MethodPost, "/api/v1/otp/send", req, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// VerifyOTP checks a code the user typed.
func (c *Client) VerifyOTP(ctx context.Context, to, code string) (*OTPVerified, error) {
	var out OTPVerified
	if _, err := c.do(ctx, http.MethodPost, "/api/v1/otp/verify", map[string]string{"to": to, "code": code}, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// RegisterDevice stores an FCM token for push messages.
func (c *Client) RegisterDevice(ctx context.Context, req RegisterDeviceRequest) (*Device, error) {
	var out Device
	if _, err := c.do(ctx, http.MethodPost, "/api/v1/devices", req, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// Signature computes the X-Signature value: hex(HMAC-SHA256(apiKey,
// timestamp "\n" METHOD "\n" path "\n" hex(sha256(body)))).
func Signature(apiKey, timestamp, method, path string, body []byte) string {
	sum := sha256.Sum256(body)
	msg := timestamp + "\n" + strings.ToUpper(method) + "\n" + path + "\n" + hex.EncodeToString(sum[:])
	mac := hmac.New(sha256.New, []byte(apiKey))
	mac.Write([]byte(msg))
	return hex.EncodeToString(mac.Sum(nil))
}

func (c *Client) do(ctx context.Context, method, path string, in, out any) (map[string]any, error) {
	var body []byte
	if in != nil {
		var err error
		if body, err = json.Marshal(in); err != nil {
			return nil, err
		}
	}
	req, err := http.NewRequestWithContext(ctx, method, c.BaseURL+path, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	req.Header.Set("X-Api-Key", c.APIKey)
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", "habarchy-go/1.0")
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	if c.Sign {
		ts := strconv.FormatInt(c.Now().Unix(), 10)
		req.Header.Set("X-Timestamp", ts)
		req.Header.Set("X-Signature", Signature(c.APIKey, ts, method, path, body))
	}
	res, err := c.HTTP.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	raw, err := io.ReadAll(io.LimitReader(res.Body, 4<<20))
	if err != nil {
		return nil, err
	}
	var env struct {
		Data  json.RawMessage `json:"data"`
		Meta  map[string]any  `json:"meta"`
		Error *Error          `json:"error"`
	}
	if len(raw) > 0 {
		if err := json.Unmarshal(raw, &env); err != nil {
			if res.StatusCode >= 300 {
				return nil, &Error{Status: res.StatusCode, Code: "http_error", Message: strings.TrimSpace(string(raw))}
			}
			return nil, fmt.Errorf("habarchy: decode response: %w", err)
		}
	}
	if res.StatusCode >= 300 {
		if env.Error == nil {
			env.Error = &Error{Code: "http_error", Message: res.Status}
		}
		env.Error.Status = res.StatusCode
		return nil, env.Error
	}
	if out != nil && len(env.Data) > 0 {
		if err := json.Unmarshal(env.Data, out); err != nil {
			return nil, fmt.Errorf("habarchy: decode data: %w", err)
		}
	}
	return env.Meta, nil
}

func metaInt(m map[string]any, k string) int {
	if f, ok := m[k].(float64); ok {
		return int(f)
	}
	return 0
}

// IsCode reports whether err is an API error with the given code.
func IsCode(err error, code string) bool {
	var e *Error
	return errors.As(err, &e) && e.Code == code
}
