// Package smtp sends email through any SMTP server.
package smtp

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/wneessen/go-mail"

	"github.com/Esca6585dev/habarchy/backend/internal/ports"
)

// Config is the decrypted credentials JSON of an smtp provider.
type Config struct {
	Host       string `json:"host"`
	Port       int    `json:"port"`
	TLSMode    string `json:"tls_mode"` // none | starttls (default) | tls
	Username   string `json:"username"`
	Password   string `json:"password"`
	FromName   string `json:"from_name"`
	FromEmail  string `json:"from_email"`
	ReplyTo    string `json:"reply_to"`
	TimeoutSec int    `json:"timeout_sec"`
	// AuthType overrides auto-detection: plain | login | cram-md5 | none.
	AuthType string `json:"auth_type"`
}

// Provider implements ports.EmailProvider.
type Provider struct {
	cfg  Config
	http *http.Client
}

// New validates cfg.
func New(cfg Config) (*Provider, error) {
	if cfg.Host == "" || cfg.FromEmail == "" {
		return nil, errors.New("smtp: host and from_email are required")
	}
	if cfg.Port == 0 {
		cfg.Port = 587
	}
	cfg.TLSMode = strings.ToLower(cfg.TLSMode)
	switch cfg.TLSMode {
	case "", "starttls":
		cfg.TLSMode = "starttls"
	case "none", "tls":
	default:
		return nil, errors.New("smtp: tls_mode must be none, starttls or tls")
	}
	if cfg.TimeoutSec <= 0 {
		cfg.TimeoutSec = 30
	}
	return &Provider{cfg: cfg, http: &http.Client{Timeout: 20 * time.Second}}, nil
}

func (p *Provider) client() (*mail.Client, error) {
	opts := []mail.Option{
		mail.WithPort(p.cfg.Port),
		mail.WithTimeout(time.Duration(p.cfg.TimeoutSec) * time.Second),
	}
	switch p.cfg.TLSMode {
	case "none":
		opts = append(opts, mail.WithTLSPortPolicy(mail.NoTLS))
	case "tls":
		opts = append(opts, mail.WithSSLPort(false))
	default:
		opts = append(opts, mail.WithTLSPortPolicy(mail.TLSMandatory))
	}
	if p.cfg.Username != "" {
		opts = append(opts, mail.WithUsername(p.cfg.Username), mail.WithPassword(p.cfg.Password))
		switch strings.ToLower(p.cfg.AuthType) {
		case "login":
			opts = append(opts, mail.WithSMTPAuth(mail.SMTPAuthLogin))
		case "cram-md5":
			opts = append(opts, mail.WithSMTPAuth(mail.SMTPAuthCramMD5))
		case "plain":
			opts = append(opts, mail.WithSMTPAuth(mail.SMTPAuthPlain))
		default:
			opts = append(opts, mail.WithSMTPAuth(mail.SMTPAuthAutoDiscover))
		}
	}
	return mail.NewClient(p.cfg.Host, opts...)
}

// Send implements ports.EmailProvider.
func (p *Provider) Send(ctx context.Context, msg ports.EmailMessage) (*ports.SendResult, error) {
	m := mail.NewMsg()
	if err := m.FromFormat(p.cfg.FromName, p.cfg.FromEmail); err != nil {
		return nil, perm("invalid_from", err.Error())
	}
	if err := m.To(msg.To); err != nil {
		return nil, perm("invalid_recipient", err.Error())
	}
	if msg.ReplyTo != "" {
		_ = m.ReplyTo(msg.ReplyTo)
	} else if p.cfg.ReplyTo != "" {
		_ = m.ReplyTo(p.cfg.ReplyTo)
	}
	m.Subject(msg.Subject)
	switch {
	case msg.HTML != "" && msg.Text != "":
		m.SetBodyString(mail.TypeTextPlain, msg.Text)
		m.AddAlternativeString(mail.TypeTextHTML, msg.HTML)
	case msg.HTML != "":
		m.SetBodyString(mail.TypeTextHTML, msg.HTML)
	default:
		m.SetBodyString(mail.TypeTextPlain, msg.Text)
	}
	for _, a := range msg.Attachments {
		content := a.Content
		if len(content) == 0 && a.URL != "" {
			b, err := p.fetch(ctx, a.URL)
			if err != nil {
				return nil, temp("attachment_fetch", err.Error())
			}
			content = b
		}
		name := a.Filename
		if name == "" {
			name = "attachment"
		}
		if err := m.AttachReader(name, strings.NewReader(string(content))); err != nil {
			return nil, perm("attachment_error", err.Error())
		}
	}
	id := m.GetMessageID()

	c, err := p.client()
	if err != nil {
		return nil, perm("smtp_config", err.Error())
	}
	if err := c.DialAndSendWithContext(ctx, m); err != nil {
		return nil, classify(err)
	}
	return &ports.SendResult{ProviderMessageID: id, Raw: map[string]any{"message_id": id}}, nil
}

func (p *Provider) fetch(ctx context.Context, url string) ([]byte, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	res, err := p.http.Do(req)
	if err != nil {
		return nil, err
	}
	defer func() { _ = res.Body.Close() }()
	if res.StatusCode != 200 {
		return nil, fmt.Errorf("attachment %s: http %d", url, res.StatusCode)
	}
	return io.ReadAll(io.LimitReader(res.Body, 10<<20))
}

// classify maps SMTP failures: permanent 5xx reply codes about the
// recipient are not retried, everything else is.
func classify(err error) *ports.ProviderError {
	s := err.Error()
	var se *mail.SendError
	if errors.As(err, &se) {
		if se.IsTemp() {
			return temp("smtp_temporary", s)
		}
		if strings.Contains(s, "550") || strings.Contains(s, "553") || strings.Contains(s, "recipient") {
			return perm("invalid_recipient", s)
		}
	}
	if strings.Contains(s, "535") || strings.Contains(s, "authentication") {
		return perm("smtp_auth", s)
	}
	return temp("smtp_error", s)
}

func perm(code, msg string) *ports.ProviderError {
	return &ports.ProviderError{Code: code, Message: msg, Retryable: false}
}

func temp(code, msg string) *ports.ProviderError {
	return &ports.ProviderError{Code: code, Message: msg, Retryable: true}
}
