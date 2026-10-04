// Package sandbox is the provider used for test API keys: it accepts
// every message without contacting anything and reports it as delivered
// immediately. It lets client apps and the Flutter demo exercise the whole
// pipeline (queue, status, webhooks) with zero cost.
package sandbox

import (
	"context"
	"strings"

	"github.com/Esca6585dev/habarchy/backend/internal/ports"
	"github.com/Esca6585dev/habarchy/backend/pkg/ids"
)

// Provider implements every channel interface.
type Provider struct{}

// FailMarker in the recipient or text forces a permanent failure, so
// clients can test their error handling: e.g. send to "+99365000FAIL".
const FailMarker = "FAIL"

// RetryMarker forces a retryable failure.
const RetryMarker = "RETRY"

func (Provider) result(addr, text string) (*ports.SendResult, error) {
	switch {
	case strings.Contains(addr, FailMarker) || strings.Contains(text, "[sandbox:fail]"):
		return nil, &ports.ProviderError{Code: "sandbox_fail", Message: "forced failure", Retryable: false}
	case strings.Contains(addr, RetryMarker) || strings.Contains(text, "[sandbox:retry]"):
		return nil, &ports.ProviderError{Code: "sandbox_retry", Message: "forced retryable failure", Retryable: true}
	}
	return &ports.SendResult{ProviderMessageID: "sandbox-" + ids.New().String(), Raw: map[string]any{"sandbox": true}}, nil
}

// Send (SMS).
func (p Provider) Send(_ context.Context, m ports.SMSMessage) (*ports.SendResult, error) {
	return p.result(m.To, m.Text)
}

// SendEmail satisfies ports.EmailProvider through the Email wrapper.
type Email struct{ Provider }

// Send (email).
func (e Email) Send(_ context.Context, m ports.EmailMessage) (*ports.SendResult, error) {
	return e.result(m.To, m.Text+m.HTML)
}

// Push satisfies ports.PushProvider.
type Push struct{ Provider }

// Send (push).
func (p Push) Send(_ context.Context, m ports.PushMessage) (*ports.PushResult, error) {
	r, err := p.result(strings.Join(m.Tokens, ","), m.Body)
	if err != nil {
		return nil, err
	}
	return &ports.PushResult{SendResult: *r}, nil
}

// Telegram satisfies ports.TelegramProvider.
type Telegram struct{ Provider }

// Send (telegram).
func (t Telegram) Send(_ context.Context, m ports.TelegramMessage) (*ports.SendResult, error) {
	return t.result(m.ChatID, m.Text)
}
