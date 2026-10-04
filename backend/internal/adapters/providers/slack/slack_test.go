package slack_test

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Esca6585dev/habarchy/backend/internal/adapters/providers/slack"
	"github.com/Esca6585dev/habarchy/backend/internal/ports"
)

func TestBotTokenAndWebhook(t *testing.T) {
	var got map[string]any
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_ = json.NewDecoder(r.Body).Decode(&got)
		switch r.URL.Path {
		case "/chat.postMessage":
			if r.Header.Get("Authorization") != "Bearer xoxb-1" {
				t.Errorf("auth %q", r.Header.Get("Authorization"))
			}
			if got["channel"] == "C404" {
				_, _ = w.Write([]byte(`{"ok":false,"error":"channel_not_found"}`))
				return
			}
			_, _ = w.Write([]byte(`{"ok":true,"channel":"C1","ts":"1700.1"}`))
		case "/hook":
			_, _ = w.Write([]byte("ok"))
		default:
			w.WriteHeader(404)
		}
	}))
	defer srv.Close()

	bot, err := slack.New(slack.Config{BotToken: "xoxb-1", Endpoint: srv.URL, DefaultChannel: "C1"}, srv.Client())
	if err != nil {
		t.Fatal(err)
	}
	res, err := bot.Send(context.Background(), ports.ChatMessage{Text: "hi", Subject: "Alert"})
	if err != nil || res.ProviderMessageID != "C1:1700.1" || got["channel"] != "C1" || got["text"] != "*Alert*\nhi" {
		t.Fatalf("%+v %v %v", res, err, got)
	}
	_, err = bot.Send(context.Background(), ports.ChatMessage{To: "C404", Text: "x"})
	var pe *ports.ProviderError
	if !errors.As(err, &pe) || pe.Retryable || pe.Code != "invalid_recipient" {
		t.Fatalf("got %v", err)
	}

	hook, err := slack.New(slack.Config{WebhookURL: "https://example.com/hook"}, srv.Client())
	if err != nil {
		t.Fatal(err)
	}
	_ = hook
	// Webhook URL must be https; point a second provider at the test server by rewriting transport.
	hook2, _ := slack.New(slack.Config{WebhookURL: "https://hooks.test/hook"}, &http.Client{Transport: rewrite{srv.URL, srv.Client().Transport}})
	res, err = hook2.Send(context.Background(), ports.ChatMessage{Text: "hello"})
	if err != nil || res.ProviderMessageID == "" || got["text"] != "hello" {
		t.Fatalf("%+v %v", res, err)
	}
	if _, err := slack.New(slack.Config{}, nil); err == nil {
		t.Fatal("empty config accepted")
	}
	if _, err := slack.New(slack.Config{BotToken: "nope"}, nil); err == nil {
		t.Fatal("bad token accepted")
	}
}

type rewrite struct {
	base string
	rt   http.RoundTripper
}

func (r rewrite) RoundTrip(req *http.Request) (*http.Response, error) {
	u := r.base + req.URL.Path
	nreq := req.Clone(req.Context())
	var err error
	if nreq.URL, err = req.URL.Parse(u); err != nil {
		return nil, err
	}
	return r.rt.RoundTrip(nreq)
}
