package fcm

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync/atomic"
	"testing"

	"github.com/Esca6585dev/habarchy/backend/internal/ports"
)

const fakeSA = `{"type":"service_account","project_id":"habarchy-demo","private_key":"x","client_email":"a@b"}`

func TestSendMixedTokens(t *testing.T) {
	var calls atomic.Int32
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		calls.Add(1)
		if !strings.HasSuffix(r.URL.Path, "/v1/projects/habarchy-demo/messages:send") {
			t.Errorf("path %s", r.URL.Path)
		}
		body, _ := io.ReadAll(r.Body)
		var req v1Request
		_ = json.Unmarshal(body, &req)
		switch req.Message.Token {
		case "dead":
			w.WriteHeader(404)
			_, _ = w.Write([]byte(`{"error":{"status":"NOT_FOUND","message":"Requested entity was not found.","details":[{"@type":"type.googleapis.com/google.firebase.fcm.v1.FcmError","errorCode":"UNREGISTERED"}]}}`))
		case "busy":
			w.WriteHeader(503)
			_, _ = w.Write([]byte(`{"error":{"status":"UNAVAILABLE","message":"try later"}}`))
		default:
			if req.Message.Notification == nil || req.Message.Notification.Title != "Salam" || req.Message.Data["k"] != "v" {
				t.Errorf("payload %s", body)
			}
			_, _ = w.Write([]byte(`{"name":"projects/habarchy-demo/messages/0:1"}`))
		}
	}))
	defer srv.Close()

	p, err := New(context.Background(), Config{ServiceAccount: json.RawMessage(fakeSA), Endpoint: srv.URL, Concurrency: 2}, srv.Client())
	if err != nil {
		t.Fatal(err)
	}
	res, err := p.Send(context.Background(), ports.PushMessage{Tokens: []string{"ok1", "dead", "busy", "ok2"}, Title: "Salam", Body: "b", Data: map[string]string{"k": "v"}})
	if err != nil {
		t.Fatalf("partial success must not error: %v", err)
	}
	if res.ProviderMessageID != "projects/habarchy-demo/messages/0:1" || len(res.InvalidTokens) != 1 || res.InvalidTokens[0] != "dead" || calls.Load() != 4 {
		t.Fatalf("%+v calls=%d", res, calls.Load())
	}

	// All tokens dead -> permanent failure.
	res, err = p.Send(context.Background(), ports.PushMessage{Tokens: []string{"dead"}, Body: "x"})
	var pe *ports.ProviderError
	if !errors.As(err, &pe) || pe.Retryable || pe.Code != "unregistered" || len(res.InvalidTokens) != 1 {
		t.Fatalf("expected permanent unregistered, got %v %+v", err, res)
	}
	// Only transient failures -> retryable.
	_, err = p.Send(context.Background(), ports.PushMessage{Tokens: []string{"busy"}, Body: "x"})
	if !errors.As(err, &pe) || !pe.Retryable {
		t.Fatalf("expected retryable, got %v", err)
	}
}

func TestConfig(t *testing.T) {
	if _, err := New(context.Background(), Config{}, http.DefaultClient); err == nil {
		t.Fatal("empty config accepted")
	}
	if _, err := New(context.Background(), Config{ServiceAccount: json.RawMessage(`{"type":"user"}`)}, http.DefaultClient); err == nil {
		t.Fatal("non service account accepted")
	}
}
