// Package whatsapp sends messages through the Meta WhatsApp Business Cloud
// API (graph.facebook.com). Free-form text is allowed inside the 24 h
// customer-service window; outside it a pre-approved template is required
// (metadata.wa_template).
package whatsapp

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

// Config is the decrypted credentials JSON of a whatsapp_cloud provider.
type Config struct {
	AccessToken   string `json:"access_token"`    // permanent system-user token
	PhoneNumberID string `json:"phone_number_id"` // the sender's phone number id
	APIVersion    string `json:"api_version"`     // default v20.0
	Endpoint      string `json:"endpoint,omitempty"`
	PreviewURL    bool   `json:"preview_url"`
}

// Provider implements ports.ChatProvider.
type Provider struct {
	cfg    Config
	client *http.Client
}

// New validates cfg.
func New(cfg Config, client *http.Client) (*Provider, error) {
	if strings.TrimSpace(cfg.AccessToken) == "" || strings.TrimSpace(cfg.PhoneNumberID) == "" {
		return nil, errors.New("whatsapp: access_token and phone_number_id are required")
	}
	if cfg.APIVersion == "" {
		cfg.APIVersion = "v20.0"
	}
	if cfg.Endpoint == "" {
		cfg.Endpoint = "https://graph.facebook.com"
	}
	if client == nil {
		client = &http.Client{Timeout: 20 * time.Second}
	}
	return &Provider{cfg: cfg, client: client}, nil
}

// Send posts a text (or template) message.
func (p *Provider) Send(ctx context.Context, msg ports.ChatMessage) (*ports.SendResult, error) {
	to := strings.TrimPrefix(strings.TrimSpace(msg.To), "+")
	if to == "" {
		return nil, &ports.ProviderError{Code: "invalid_recipient", Message: "phone number missing"}
	}
	payload := map[string]any{"messaging_product": "whatsapp", "recipient_type": "individual", "to": to}
	if tpl, ok := msg.Extra["wa_template"].(map[string]any); ok && tpl["name"] != nil {
		lang := map[string]any{"code": "en"}
		if l, ok := tpl["language"].(string); ok && l != "" {
			lang["code"] = l
		}
		t := map[string]any{"name": tpl["name"], "language": lang}
		if comps, ok := tpl["components"]; ok {
			t["components"] = comps
		}
		payload["type"], payload["template"] = "template", t
	} else {
		text := msg.Text
		if msg.Subject != "" {
			text = "*" + msg.Subject + "*\n" + text
		}
		payload["type"], payload["text"] = "text", map[string]any{"body": text, "preview_url": p.cfg.PreviewURL}
	}
	body, _ := json.Marshal(payload)
	url := fmt.Sprintf("%s/%s/%s/messages", p.cfg.Endpoint, p.cfg.APIVersion, p.cfg.PhoneNumberID)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(body))
	if err != nil {
		return nil, &ports.ProviderError{Code: "bad_request", Message: err.Error()}
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+p.cfg.AccessToken)
	res, err := p.client.Do(req)
	if err != nil {
		return nil, &ports.ProviderError{Code: "network_error", Message: err.Error(), Retryable: true}
	}
	defer func() { _ = res.Body.Close() }()
	raw, _ := io.ReadAll(io.LimitReader(res.Body, 1<<20))

	var out struct {
		Messages []struct {
			ID string `json:"id"`
		} `json:"messages"`
		Error struct {
			Message string `json:"message"`
			Code    int    `json:"code"`
			Subcode int    `json:"error_subcode"`
		} `json:"error"`
	}
	_ = json.Unmarshal(raw, &out)
	rawMap := map[string]any{"status": res.StatusCode}
	if res.StatusCode >= 200 && res.StatusCode < 300 && len(out.Messages) > 0 {
		return &ports.SendResult{ProviderMessageID: out.Messages[0].ID, Raw: rawMap}, nil
	}
	rawMap["error"], rawMap["code"] = out.Error.Message, out.Error.Code
	code, retry := classify(res.StatusCode, out.Error.Code)
	msgText := out.Error.Message
	if msgText == "" {
		msgText = strings.TrimSpace(string(raw))
		if len(msgText) > 200 {
			msgText = msgText[:200]
		}
	}
	return nil, &ports.ProviderError{Code: code, Message: msgText, Retryable: retry, Raw: rawMap}
}

// classify maps Graph API error codes to Habarchy codes.
// https://developers.facebook.com/docs/whatsapp/cloud-api/support/error-codes
func classify(status, code int) (string, bool) {
	switch code {
	case 131026, 131030, 131047, 131051, 131052, 131053:
		return "invalid_recipient", false // not a WhatsApp user / not in allowed list / unsupported type / media
	case 131031, 131042, 133010, 130497:
		return "account_restricted", false
	case 190, 10, 200, 100:
		if code == 100 {
			return "bad_request", false
		}
		return "auth", false
	case 131056, 131057:
		return "pair_rate_limit", true
	case 4, 80007, 130429, 131048:
		return "rate_limited", true
	case 131000, 131016, 131021, 131045:
		return "provider_error", true
	}
	switch {
	case status == 401 || status == 403:
		return "auth", false
	case status == 400:
		return "bad_request", false
	case status == 429 || status >= 500:
		return "provider_error", true
	}
	return "provider_error", true
}
