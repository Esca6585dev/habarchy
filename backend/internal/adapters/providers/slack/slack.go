// Package slack posts messages to Slack, either with a bot token
// (chat.postMessage, recipient = channel or user id) or through an
// incoming webhook URL (recipient ignored, fixed channel).
package slack

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/Esca6585dev/habarchy/backend/internal/ports"
)

// Config is the decrypted credentials JSON of a slack provider.
type Config struct {
	BotToken       string `json:"bot_token"`       // xoxb-… (chat:write scope)
	WebhookURL     string `json:"webhook_url"`     // alternative: incoming webhook
	DefaultChannel string `json:"default_channel"` // used when the message has no recipient
	Endpoint       string `json:"endpoint,omitempty"`
}

// Provider implements ports.ChatProvider.
type Provider struct {
	cfg    Config
	client *http.Client
}

// New validates cfg.
func New(cfg Config, client *http.Client) (*Provider, error) {
	cfg.BotToken, cfg.WebhookURL = strings.TrimSpace(cfg.BotToken), strings.TrimSpace(cfg.WebhookURL)
	if cfg.BotToken == "" && cfg.WebhookURL == "" {
		return nil, errors.New("slack: bot_token or webhook_url is required")
	}
	if cfg.BotToken != "" && !strings.HasPrefix(cfg.BotToken, "xox") {
		return nil, errors.New("slack: bot_token should start with xoxb-")
	}
	if cfg.WebhookURL != "" && !strings.HasPrefix(cfg.WebhookURL, "https://") {
		return nil, errors.New("slack: webhook_url must be https")
	}
	if cfg.Endpoint == "" {
		cfg.Endpoint = "https://slack.com/api"
	}
	if client == nil {
		client = &http.Client{Timeout: 15 * time.Second}
	}
	return &Provider{cfg: cfg, client: client}, nil
}

// Send posts the message.
func (p *Provider) Send(ctx context.Context, msg ports.ChatMessage) (*ports.SendResult, error) {
	text := msg.Text
	if msg.Subject != "" {
		text = "*" + msg.Subject + "*\n" + text
	}
	if p.cfg.BotToken == "" {
		return p.webhook(ctx, text)
	}
	channel := strings.TrimSpace(msg.To)
	if channel == "" {
		channel = p.cfg.DefaultChannel
	}
	if channel == "" {
		return nil, &ports.ProviderError{Code: "invalid_recipient", Message: "slack channel missing"}
	}
	payload := map[string]any{"channel": channel, "text": text, "mrkdwn": true}
	if blocks, ok := msg.Extra["slack_blocks"]; ok {
		payload["blocks"] = blocks
	}
	body, _ := json.Marshal(payload)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, p.cfg.Endpoint+"/chat.postMessage", bytes.NewReader(body))
	if err != nil {
		return nil, &ports.ProviderError{Code: "bad_request", Message: err.Error()}
	}
	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	req.Header.Set("Authorization", "Bearer "+p.cfg.BotToken)
	res, err := p.client.Do(req)
	if err != nil {
		return nil, &ports.ProviderError{Code: "network_error", Message: err.Error(), Retryable: true}
	}
	defer func() { _ = res.Body.Close() }()
	raw, _ := io.ReadAll(io.LimitReader(res.Body, 1<<20))
	var out struct {
		OK    bool   `json:"ok"`
		Error string `json:"error"`
		TS    string `json:"ts"`
		Chan  string `json:"channel"`
	}
	_ = json.Unmarshal(raw, &out)
	rawMap := map[string]any{"status": res.StatusCode, "error": out.Error}
	if res.StatusCode == 429 {
		return nil, &ports.ProviderError{Code: "rate_limited", Message: "slack rate limit", Retryable: true, Raw: rawMap}
	}
	if out.OK {
		return &ports.SendResult{ProviderMessageID: out.Chan + ":" + out.TS, Raw: rawMap}, nil
	}
	code, retry := classify(out.Error)
	return nil, &ports.ProviderError{Code: code, Message: fmt.Sprintf("slack: %s", nz(out.Error, res.Status)), Retryable: retry, Raw: rawMap}
}

func (p *Provider) webhook(ctx context.Context, text string) (*ports.SendResult, error) {
	body, _ := json.Marshal(map[string]any{"text": text})
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, p.cfg.WebhookURL, bytes.NewReader(body))
	if err != nil {
		return nil, &ports.ProviderError{Code: "bad_request", Message: err.Error()}
	}
	req.Header.Set("Content-Type", "application/json")
	res, err := p.client.Do(req)
	if err != nil {
		return nil, &ports.ProviderError{Code: "network_error", Message: err.Error(), Retryable: true}
	}
	defer func() { _ = res.Body.Close() }()
	raw, _ := io.ReadAll(io.LimitReader(res.Body, 64<<10))
	rawMap := map[string]any{"status": res.StatusCode, "body": string(raw)}
	switch {
	case res.StatusCode >= 200 && res.StatusCode < 300:
		return &ports.SendResult{ProviderMessageID: fmt.Sprintf("webhook-%d", time.Now().UnixNano()), Raw: rawMap}, nil
	case res.StatusCode == 404 || res.StatusCode == 403 || res.StatusCode == 410:
		return nil, &ports.ProviderError{Code: "auth", Message: "slack webhook rejected: " + strings.TrimSpace(string(raw)), Retryable: false, Raw: rawMap}
	case res.StatusCode == 400:
		return nil, &ports.ProviderError{Code: "bad_request", Message: strings.TrimSpace(string(raw)), Retryable: false, Raw: rawMap}
	}
	return nil, &ports.ProviderError{Code: "provider_error", Message: res.Status, Retryable: true, Raw: rawMap}
}

func classify(e string) (string, bool) {
	switch e {
	case "channel_not_found", "user_not_found", "is_archived", "not_in_channel", "channel_is_archived":
		return "invalid_recipient", false
	case "invalid_auth", "not_authed", "account_inactive", "token_revoked", "missing_scope", "no_permission":
		return "auth", false
	case "msg_too_long", "no_text", "invalid_blocks", "invalid_arguments":
		return "bad_request", false
	case "ratelimited", "rate_limited", "service_unavailable", "internal_error", "fatal_error":
		return "provider_error", true
	}
	return "provider_error", true
}

func nz(a, b string) string {
	if a != "" {
		return a
	}
	return b
}
