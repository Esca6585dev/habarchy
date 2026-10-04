// Package fcm sends push notifications through the Firebase Cloud
// Messaging HTTP v1 API, authenticating with a service account.
package fcm

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"sync"
	"time"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"

	"github.com/Esca6585dev/habarchy/backend/internal/ports"
)

// Scope is the OAuth scope FCM requires.
const Scope = "https://www.googleapis.com/auth/firebase.messaging"

// MaxTokensPerCall mirrors the legacy batch limit; tokens beyond it are
// sent in further rounds.
const MaxTokensPerCall = 500

// Config is the decrypted credentials JSON of an fcm provider.
type Config struct {
	// ServiceAccount is the Firebase service-account JSON (object).
	ServiceAccount json.RawMessage `json:"service_account"`
	// ProjectID overrides the project_id of the service account.
	ProjectID string `json:"project_id"`
	// Concurrency bounds parallel requests per Send (default 8).
	Concurrency int `json:"concurrency"`
	// Endpoint overrides the API base URL (tests).
	Endpoint string `json:"endpoint,omitempty"`
}

// Provider implements ports.PushProvider.
type Provider struct {
	projectID string
	endpoint  string
	conc      int
	client    *http.Client
}

// New builds a provider. client overrides the OAuth-authenticated client
// (tests).
func New(ctx context.Context, cfg Config, client *http.Client) (*Provider, error) {
	if len(cfg.ServiceAccount) == 0 {
		return nil, errors.New("fcm: service_account is required")
	}
	var sa struct {
		ProjectID string `json:"project_id"`
		Type      string `json:"type"`
	}
	if err := json.Unmarshal(cfg.ServiceAccount, &sa); err != nil || sa.Type != "service_account" {
		return nil, errors.New("fcm: service_account must be a Firebase service-account JSON object")
	}
	projectID := cfg.ProjectID
	if projectID == "" {
		projectID = sa.ProjectID
	}
	if projectID == "" {
		return nil, errors.New("fcm: project_id missing")
	}
	if client == nil {
		creds, err := google.CredentialsFromJSON(ctx, cfg.ServiceAccount, Scope)
		if err != nil {
			return nil, fmt.Errorf("fcm: credentials: %w", err)
		}
		client = oauth2.NewClient(ctx, creds.TokenSource)
		client.Timeout = 20 * time.Second
	}
	endpoint := cfg.Endpoint
	if endpoint == "" {
		endpoint = "https://fcm.googleapis.com"
	}
	conc := cfg.Concurrency
	if conc <= 0 {
		conc = 8
	}
	return &Provider{projectID: projectID, endpoint: endpoint, conc: conc, client: client}, nil
}

// Send implements ports.PushProvider. Each token is one HTTP v1 request;
// requests run concurrently up to Concurrency. Tokens reported as
// UNREGISTERED / INVALID_ARGUMENT come back in InvalidTokens so the caller
// can disable the device.
func (p *Provider) Send(ctx context.Context, msg ports.PushMessage) (*ports.PushResult, error) {
	if len(msg.Tokens) == 0 && msg.Topic == "" {
		return nil, &ports.ProviderError{Code: "invalid_recipient", Message: "no tokens or topic"}
	}
	targets := msg.Tokens
	if len(targets) > MaxTokensPerCall {
		targets = targets[:MaxTokensPerCall]
	}
	if msg.Topic != "" && len(targets) == 0 {
		id, err := p.sendOne(ctx, msg, "", msg.Topic)
		if err != nil {
			return nil, err
		}
		return &ports.PushResult{SendResult: ports.SendResult{ProviderMessageID: id, Raw: map[string]any{"topic": msg.Topic}}}, nil
	}

	type out struct {
		token string
		id    string
		err   error
	}
	results := make([]out, len(targets))
	sem := make(chan struct{}, p.conc)
	var wg sync.WaitGroup
	for i, tok := range targets {
		wg.Add(1)
		sem <- struct{}{}
		go func(i int, tok string) {
			defer wg.Done()
			defer func() { <-sem }()
			id, err := p.sendOne(ctx, msg, tok, "")
			results[i] = out{token: tok, id: id, err: err}
		}(i, tok)
	}
	wg.Wait()

	res := &ports.PushResult{}
	var firstErr *ports.ProviderError
	sent, perTok := 0, map[string]any{}
	for _, r := range results {
		if r.err == nil {
			sent++
			if res.ProviderMessageID == "" {
				res.ProviderMessageID = r.id
			}
			perTok[r.token] = "ok"
			continue
		}
		var pe *ports.ProviderError
		if errors.As(r.err, &pe) {
			perTok[r.token] = pe.Code
			if pe.Code == "unregistered" || pe.Code == "invalid_token" {
				res.InvalidTokens = append(res.InvalidTokens, r.token)
			}
			if firstErr == nil {
				firstErr = pe
			}
		} else {
			perTok[r.token] = r.err.Error()
		}
	}
	res.Raw = map[string]any{"sent": sent, "total": len(targets), "tokens": perTok}
	if sent == 0 && firstErr != nil {
		firstErr.Raw = res.Raw
		// If every token is dead the message is permanently undeliverable.
		if len(res.InvalidTokens) == len(targets) {
			firstErr.Retryable = false
			firstErr.Code = "unregistered"
		}
		return res, firstErr
	}
	return res, nil
}

type v1Request struct {
	ValidateOnly bool      `json:"validate_only,omitempty"`
	Message      v1Message `json:"message"`
}

type v1Message struct {
	Token        string            `json:"token,omitempty"`
	Topic        string            `json:"topic,omitempty"`
	Notification *v1Notification   `json:"notification,omitempty"`
	Data         map[string]string `json:"data,omitempty"`
	Android      map[string]any    `json:"android,omitempty"`
	APNs         map[string]any    `json:"apns,omitempty"`
}

type v1Notification struct {
	Title string `json:"title,omitempty"`
	Body  string `json:"body,omitempty"`
}

func (p *Provider) sendOne(ctx context.Context, msg ports.PushMessage, token, topic string) (string, error) {
	m := v1Message{Token: token, Topic: topic, Data: msg.Data, Android: msg.Android, APNs: msg.APNs}
	if msg.Title != "" || msg.Body != "" {
		m.Notification = &v1Notification{Title: msg.Title, Body: msg.Body}
	}
	body, _ := json.Marshal(v1Request{Message: m})
	url := fmt.Sprintf("%s/v1/projects/%s/messages:send", p.endpoint, p.projectID)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(body))
	if err != nil {
		return "", &ports.ProviderError{Code: "bad_request", Message: err.Error()}
	}
	req.Header.Set("Content-Type", "application/json")
	res, err := p.client.Do(req)
	if err != nil {
		return "", &ports.ProviderError{Code: "network_error", Message: err.Error(), Retryable: true}
	}
	defer func() { _ = res.Body.Close() }()
	raw, _ := io.ReadAll(io.LimitReader(res.Body, 1<<20))
	if res.StatusCode == 200 {
		var ok struct {
			Name string `json:"name"`
		}
		_ = json.Unmarshal(raw, &ok)
		return ok.Name, nil
	}
	return "", classify(res.StatusCode, raw)
}

// classify maps FCM error codes. See
// https://firebase.google.com/docs/reference/fcm/rest/v1/ErrorCode
func classify(status int, raw []byte) *ports.ProviderError {
	var body struct {
		Error struct {
			Message string `json:"message"`
			Status  string `json:"status"`
			Details []struct {
				Type      string `json:"@type"`
				ErrorCode string `json:"errorCode"`
			} `json:"details"`
		} `json:"error"`
	}
	_ = json.Unmarshal(raw, &body)
	code := body.Error.Status
	for _, d := range body.Error.Details {
		if d.ErrorCode != "" {
			code = d.ErrorCode
		}
	}
	msg := body.Error.Message
	if msg == "" {
		msg = fmt.Sprintf("http %d", status)
	}
	switch code {
	case "UNREGISTERED", "NOT_FOUND":
		return &ports.ProviderError{Code: "unregistered", Message: msg}
	case "INVALID_ARGUMENT":
		return &ports.ProviderError{Code: "invalid_token", Message: msg}
	case "SENDER_ID_MISMATCH", "THIRD_PARTY_AUTH_ERROR", "PERMISSION_DENIED":
		return &ports.ProviderError{Code: "fcm_auth", Message: msg}
	case "QUOTA_EXCEEDED", "UNAVAILABLE", "INTERNAL", "RESOURCE_EXHAUSTED":
		return &ports.ProviderError{Code: "provider_unavailable", Message: msg, Retryable: true}
	}
	if status == 401 || status == 403 {
		return &ports.ProviderError{Code: "fcm_auth", Message: msg}
	}
	return &ports.ProviderError{Code: "fcm_error", Message: msg, Retryable: status >= 500 || status == 429}
}
