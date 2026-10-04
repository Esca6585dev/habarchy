// Package telegram sends messages through the Telegram Bot API.
package telegram

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/Esca6585dev/habarchy/backend/internal/ports"
)

// Config is the decrypted credentials JSON of a telegram_bot provider.
type Config struct {
	BotToken       string `json:"bot_token"`
	DefaultChatID  string `json:"default_chat_id"`
	ParseMode      string `json:"parse_mode"` // Markdown | MarkdownV2 | HTML | ""
	DisablePreview bool   `json:"disable_web_page_preview"`
	Endpoint       string `json:"endpoint,omitempty"` // tests
}

// Provider implements ports.TelegramProvider.
type Provider struct {
	cfg    Config
	client *http.Client
}

// New validates cfg.
func New(cfg Config, client *http.Client) (*Provider, error) {
	if cfg.BotToken == "" || !strings.Contains(cfg.BotToken, ":") {
		return nil, errors.New("telegram: bot_token is required (format 123456:ABC...)")
	}
	switch cfg.ParseMode {
	case "", "Markdown", "MarkdownV2", "HTML":
	default:
		return nil, errors.New("telegram: parse_mode must be Markdown, MarkdownV2 or HTML")
	}
	if cfg.Endpoint == "" {
		cfg.Endpoint = "https://api.telegram.org"
	}
	if client == nil {
		client = &http.Client{Timeout: 15 * time.Second}
	}
	return &Provider{cfg: cfg, client: client}, nil
}

type sendMessageReq struct {
	ChatID                string          `json:"chat_id"`
	Text                  string          `json:"text"`
	ParseMode             string          `json:"parse_mode,omitempty"`
	DisableWebPagePreview bool            `json:"disable_web_page_preview,omitempty"`
	ReplyMarkup           *inlineKeyboard `json:"reply_markup,omitempty"`
}

type inlineKeyboard struct {
	InlineKeyboard [][]inlineButton `json:"inline_keyboard"`
}

type inlineButton struct {
	Text         string `json:"text"`
	URL          string `json:"url,omitempty"`
	CallbackData string `json:"callback_data,omitempty"`
}

// Send implements ports.TelegramProvider.
func (p *Provider) Send(ctx context.Context, msg ports.TelegramMessage) (*ports.SendResult, error) {
	chatID := msg.ChatID
	if chatID == "" {
		chatID = p.cfg.DefaultChatID
	}
	if chatID == "" {
		return nil, &ports.ProviderError{Code: "invalid_recipient", Message: "chat_id missing"}
	}
	parseMode := msg.ParseMode
	if parseMode == "" {
		parseMode = p.cfg.ParseMode
	}
	req := sendMessageReq{ChatID: chatID, Text: msg.Text, ParseMode: parseMode, DisableWebPagePreview: p.cfg.DisablePreview}
	if len(msg.Buttons) > 0 {
		kb := &inlineKeyboard{}
		for _, row := range msg.Buttons {
			var r []inlineButton
			for _, b := range row {
				r = append(r, inlineButton{Text: b.Text, URL: b.URL, CallbackData: b.Data})
			}
			kb.InlineKeyboard = append(kb.InlineKeyboard, r)
		}
		req.ReplyMarkup = kb
	}
	body, _ := json.Marshal(req)
	url := fmt.Sprintf("%s/bot%s/sendMessage", p.cfg.Endpoint, p.cfg.BotToken)
	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(body))
	if err != nil {
		return nil, &ports.ProviderError{Code: "bad_request", Message: err.Error()}
	}
	httpReq.Header.Set("Content-Type", "application/json")
	res, err := p.client.Do(httpReq)
	if err != nil {
		return nil, &ports.ProviderError{Code: "network_error", Message: err.Error(), Retryable: true}
	}
	defer func() { _ = res.Body.Close() }()
	raw, _ := io.ReadAll(io.LimitReader(res.Body, 1<<20))

	var tg struct {
		OK          bool   `json:"ok"`
		Description string `json:"description"`
		ErrorCode   int    `json:"error_code"`
		Result      struct {
			MessageID int64 `json:"message_id"`
		} `json:"result"`
		Parameters struct {
			RetryAfter int `json:"retry_after"`
		} `json:"parameters"`
	}
	_ = json.Unmarshal(raw, &tg)
	// The bot token is part of the URL; never include the URL in Raw.
	info := map[string]any{"status": res.StatusCode, "description": tg.Description}
	if tg.OK {
		return &ports.SendResult{ProviderMessageID: strconv.FormatInt(tg.Result.MessageID, 10), Raw: info}, nil
	}
	switch {
	case tg.ErrorCode == 429 || res.StatusCode == 429:
		info["retry_after"] = tg.Parameters.RetryAfter
		return nil, &ports.ProviderError{Code: "rate_limited", Message: tg.Description, Retryable: true, Raw: info}
	case tg.ErrorCode == 401:
		return nil, &ports.ProviderError{Code: "telegram_auth", Message: tg.Description, Raw: info}
	case tg.ErrorCode == 400 || tg.ErrorCode == 403:
		// chat not found, bot blocked by the user, bad markup...
		return nil, &ports.ProviderError{Code: "invalid_recipient", Message: tg.Description, Raw: info}
	case res.StatusCode >= 500:
		return nil, &ports.ProviderError{Code: "provider_unavailable", Message: tg.Description, Retryable: true, Raw: info}
	}
	return nil, &ports.ProviderError{Code: "telegram_error", Message: tg.Description, Raw: info}
}
