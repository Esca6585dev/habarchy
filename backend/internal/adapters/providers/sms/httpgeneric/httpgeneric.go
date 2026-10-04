// Package httpgeneric is a fully configurable HTTP SMS provider: URL,
// method, headers and body are templates, success and message-id
// extraction are declarative. Most local SMS gateways can be integrated
// from the admin panel without code.
package httpgeneric

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"text/template"
	"time"

	"github.com/Esca6585dev/habarchy/backend/internal/ports"
)

// Config is the decrypted credentials JSON of an http_sms provider.
type Config struct {
	URL          string            `json:"url"`
	Method       string            `json:"method"`        // GET | POST (default POST)
	Headers      map[string]string `json:"headers"`       // values are templates
	ContentType  string            `json:"content_type"`  // application/json (default) | application/x-www-form-urlencoded | text/plain
	BodyTemplate string            `json:"body_template"` // Go template with {{.To}} {{.Text}} {{.Sender}} {{.TextJSON}}
	Sender       string            `json:"sender"`
	TimeoutSec   int               `json:"timeout_sec"`
	Success      Matcher           `json:"success"`
	MessageID    Extractor         `json:"message_id"`
	// Retryable lists HTTP status codes treated as temporary (default
	// 408, 429, 5xx).
	RetryableStatus []int `json:"retryable_status"`
	// DLR configures inbound delivery reports on POST /callbacks/sms/{id}.
	DLR DLRConfig `json:"dlr"`
}

// Matcher decides whether a response means "accepted". All configured
// conditions must hold. With nothing configured, any 2xx is success.
type Matcher struct {
	StatusCodes []int  `json:"status_codes"`
	JSONPath    string `json:"json_path"`   // dot path, e.g. "result.status"
	JSONEquals  string `json:"json_equals"` // expected value as string
	Regex       string `json:"regex"`       // must match body
}

// Extractor pulls the provider message id out of the response.
type Extractor struct {
	JSONPath string `json:"json_path"`
	Regex    string `json:"regex"` // first capture group
	Header   string `json:"header"`
}

// DLRConfig maps a delivery-report callback to a message and status.
type DLRConfig struct {
	MessageIDParam  string   `json:"message_id_param"` // form/query/json field holding the provider message id
	StatusParam     string   `json:"status_param"`
	DeliveredValues []string `json:"delivered_values"` // e.g. ["DELIVRD","delivered","1"]
	FailedValues    []string `json:"failed_values"`
}

// Provider implements ports.SMSProvider over HTTP.
type Provider struct {
	cfg    Config
	client *http.Client
	body   *template.Template
	url    *template.Template
	hdrs   map[string]*template.Template
	succRe *regexp.Regexp
	idRe   *regexp.Regexp
}

// New validates cfg and builds the provider. client may be nil.
func New(cfg Config, client *http.Client) (*Provider, error) {
	if cfg.URL == "" {
		return nil, errors.New("http_sms: url is required")
	}
	cfg.Method = strings.ToUpper(cfg.Method)
	if cfg.Method == "" {
		cfg.Method = http.MethodPost
	}
	if cfg.Method != http.MethodGet && cfg.Method != http.MethodPost && cfg.Method != http.MethodPut {
		return nil, errors.New("http_sms: method must be GET, POST or PUT")
	}
	if cfg.ContentType == "" {
		cfg.ContentType = "application/json"
	}
	if cfg.TimeoutSec <= 0 {
		cfg.TimeoutSec = 15
	}
	if cfg.Method != http.MethodGet && cfg.BodyTemplate == "" {
		return nil, errors.New("http_sms: body_template is required for POST/PUT")
	}
	p := &Provider{cfg: cfg, client: client, hdrs: map[string]*template.Template{}}
	if p.client == nil {
		p.client = &http.Client{Timeout: time.Duration(cfg.TimeoutSec) * time.Second}
	}
	var err error
	if p.url, err = template.New("url").Funcs(funcs).Parse(cfg.URL); err != nil {
		return nil, fmt.Errorf("http_sms: url template: %w", err)
	}
	if cfg.BodyTemplate != "" {
		if p.body, err = template.New("body").Funcs(funcs).Parse(cfg.BodyTemplate); err != nil {
			return nil, fmt.Errorf("http_sms: body_template: %w", err)
		}
	}
	for k, v := range cfg.Headers {
		if p.hdrs[k], err = template.New("h").Funcs(funcs).Parse(v); err != nil {
			return nil, fmt.Errorf("http_sms: header %s: %w", k, err)
		}
	}
	if cfg.Success.Regex != "" {
		if p.succRe, err = regexp.Compile(cfg.Success.Regex); err != nil {
			return nil, fmt.Errorf("http_sms: success.regex: %w", err)
		}
	}
	if cfg.MessageID.Regex != "" {
		if p.idRe, err = regexp.Compile(cfg.MessageID.Regex); err != nil {
			return nil, fmt.Errorf("http_sms: message_id.regex: %w", err)
		}
	}
	return p, nil
}

var funcs = template.FuncMap{
	"urlquery": func(s string) string { return strings.ReplaceAll(templateURLQuery(s), "+", "%20") },
	"json":     func(s string) string { b, _ := json.Marshal(s); return string(b) },
}

// vars is the data passed to the templates.
type vars struct {
	To       string
	Text     string
	TextJSON string // JSON-escaped text, including quotes
	Sender   string
	Length   int
}

// Send implements ports.SMSProvider.
func (p *Provider) Send(ctx context.Context, msg ports.SMSMessage) (*ports.SendResult, error) {
	sender := msg.Sender
	if sender == "" {
		sender = p.cfg.Sender
	}
	textJSON, _ := json.Marshal(msg.Text)
	v := vars{To: msg.To, Text: msg.Text, TextJSON: string(textJSON), Sender: sender, Length: len([]rune(msg.Text))}

	url, err := render(p.url, v)
	if err != nil {
		return nil, perm("template_error", err.Error())
	}
	var body io.Reader
	if p.body != nil {
		b, err := render(p.body, v)
		if err != nil {
			return nil, perm("template_error", err.Error())
		}
		body = strings.NewReader(b)
	}
	req, err := http.NewRequestWithContext(ctx, p.cfg.Method, url, body)
	if err != nil {
		return nil, perm("bad_request", err.Error())
	}
	if body != nil {
		req.Header.Set("Content-Type", p.cfg.ContentType)
	}
	req.Header.Set("Accept", "application/json, text/plain, */*")
	for k, t := range p.hdrs {
		val, err := render(t, v)
		if err != nil {
			return nil, perm("template_error", err.Error())
		}
		req.Header.Set(k, val)
	}

	res, err := p.client.Do(req)
	if err != nil {
		return nil, temp("network_error", err.Error())
	}
	defer func() { _ = res.Body.Close() }()
	raw, _ := io.ReadAll(io.LimitReader(res.Body, 1<<20))
	rawInfo := map[string]any{"status": res.StatusCode, "body": truncate(string(raw), 2000)}

	if !p.matches(res, raw) {
		code := "provider_rejected"
		if p.retryable(res.StatusCode) {
			return nil, &ports.ProviderError{Code: "provider_unavailable", Message: fmt.Sprintf("http %d", res.StatusCode), Retryable: true, Raw: rawInfo}
		}
		return nil, &ports.ProviderError{Code: code, Message: fmt.Sprintf("http %d: %s", res.StatusCode, truncate(string(raw), 200)), Retryable: false, Raw: rawInfo}
	}
	return &ports.SendResult{ProviderMessageID: p.extractID(res, raw), Raw: rawInfo}, nil
}

func (p *Provider) retryable(status int) bool {
	if len(p.cfg.RetryableStatus) > 0 {
		for _, s := range p.cfg.RetryableStatus {
			if s == status {
				return true
			}
		}
		return false
	}
	return status == 408 || status == 429 || status >= 500
}

func (p *Provider) matches(res *http.Response, raw []byte) bool {
	m := p.cfg.Success
	if len(m.StatusCodes) > 0 {
		ok := false
		for _, s := range m.StatusCodes {
			if s == res.StatusCode {
				ok = true
			}
		}
		if !ok {
			return false
		}
	} else if res.StatusCode < 200 || res.StatusCode > 299 {
		return false
	}
	if m.JSONPath != "" {
		val, ok := JSONPath(raw, m.JSONPath)
		if !ok {
			return false
		}
		if m.JSONEquals != "" && fmt.Sprint(val) != m.JSONEquals {
			return false
		}
	}
	if p.succRe != nil && !p.succRe.Match(raw) {
		return false
	}
	return true
}

func (p *Provider) extractID(res *http.Response, raw []byte) string {
	e := p.cfg.MessageID
	if e.Header != "" {
		if h := res.Header.Get(e.Header); h != "" {
			return h
		}
	}
	if e.JSONPath != "" {
		if v, ok := JSONPath(raw, e.JSONPath); ok {
			switch x := v.(type) {
			case float64:
				return strconv.FormatFloat(x, 'f', -1, 64)
			default:
				return fmt.Sprint(x)
			}
		}
	}
	if p.idRe != nil {
		if m := p.idRe.FindSubmatch(raw); len(m) > 1 {
			return string(m[1])
		}
	}
	return ""
}

// JSONPath resolves a dot path ("a.b.0.c") in a JSON document.
func JSONPath(raw []byte, path string) (any, bool) {
	var doc any
	if err := json.Unmarshal(raw, &doc); err != nil {
		return nil, false
	}
	cur := doc
	for _, part := range strings.Split(path, ".") {
		switch node := cur.(type) {
		case map[string]any:
			v, ok := node[part]
			if !ok {
				return nil, false
			}
			cur = v
		case []any:
			i, err := strconv.Atoi(part)
			if err != nil || i < 0 || i >= len(node) {
				return nil, false
			}
			cur = node[i]
		default:
			return nil, false
		}
	}
	return cur, true
}

// ParseDLR interprets a delivery report callback. values holds the merged
// query/form/JSON fields. It returns the provider message id and whether
// the report means delivered (true), failed (false) or unknown (nil).
func (c DLRConfig) ParseDLR(values map[string]string) (messageID string, delivered *bool) {
	messageID = values[c.MessageIDParam]
	status := values[c.StatusParam]
	for _, v := range c.DeliveredValues {
		if strings.EqualFold(v, status) {
			t := true
			return messageID, &t
		}
	}
	for _, v := range c.FailedValues {
		if strings.EqualFold(v, status) {
			f := false
			return messageID, &f
		}
	}
	return messageID, nil
}

func render(t *template.Template, v vars) (string, error) {
	var buf bytes.Buffer
	if err := t.Execute(&buf, v); err != nil {
		return "", err
	}
	return buf.String(), nil
}

func perm(code, msg string) *ports.ProviderError {
	return &ports.ProviderError{Code: code, Message: msg, Retryable: false}
}

func temp(code, msg string) *ports.ProviderError {
	return &ports.ProviderError{Code: code, Message: msg, Retryable: true}
}

func truncate(s string, n int) string {
	if len(s) <= n {
		return s
	}
	return s[:n] + "…"
}
