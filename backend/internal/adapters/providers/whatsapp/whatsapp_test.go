package whatsapp_test

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Esca6585dev/habarchy/backend/internal/adapters/providers/whatsapp"
	"github.com/Esca6585dev/habarchy/backend/internal/ports"
)

func TestSendTextAndErrors(t *testing.T) {
	var got map[string]any
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Authorization") != "Bearer tok" || r.URL.Path != "/v20.0/123/messages" {
			t.Errorf("request %s %s", r.URL.Path, r.Header.Get("Authorization"))
		}
		_ = json.NewDecoder(r.Body).Decode(&got)
		if got["to"] == "99365000999" {
			w.WriteHeader(400)
			_, _ = w.Write([]byte(`{"error":{"message":"not a whatsapp user","code":131026}}`))
			return
		}
		if got["to"] == "99365000429" {
			w.WriteHeader(429)
			_, _ = w.Write([]byte(`{"error":{"message":"rate","code":130429}}`))
			return
		}
		_, _ = w.Write([]byte(`{"messages":[{"id":"wamid.abc"}]}`))
	}))
	defer srv.Close()
	p, err := whatsapp.New(whatsapp.Config{AccessToken: "tok", PhoneNumberID: "123", Endpoint: srv.URL}, srv.Client())
	if err != nil {
		t.Fatal(err)
	}
	res, err := p.Send(context.Background(), ports.ChatMessage{To: "+99365123456", Text: "Salam", Subject: "Habar"})
	if err != nil || res.ProviderMessageID != "wamid.abc" {
		t.Fatalf("%+v %v", res, err)
	}
	if got["type"] != "text" || got["text"].(map[string]any)["body"] != "*Habar*\nSalam" {
		t.Fatalf("payload %v", got)
	}
	_, err = p.Send(context.Background(), ports.ChatMessage{To: "+99365000999", Text: "x"})
	var pe *ports.ProviderError
	if !errors.As(err, &pe) || pe.Retryable || pe.Code != "invalid_recipient" {
		t.Fatalf("got %v", err)
	}
	_, err = p.Send(context.Background(), ports.ChatMessage{To: "+99365000429", Text: "x"})
	if !errors.As(err, &pe) || !pe.Retryable || pe.Code != "rate_limited" {
		t.Fatalf("got %v", err)
	}
	// Template message.
	_, _ = p.Send(context.Background(), ports.ChatMessage{To: "+99365123456", Text: "ignored", Extra: map[string]any{"wa_template": map[string]any{"name": "otp", "language": "tk"}}})
	if got["type"] != "template" || got["template"].(map[string]any)["name"] != "otp" {
		t.Fatalf("template payload %v", got)
	}
	if _, err := whatsapp.New(whatsapp.Config{}, nil); err == nil {
		t.Fatal("empty config accepted")
	}
}
