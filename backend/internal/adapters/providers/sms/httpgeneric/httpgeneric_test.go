package httpgeneric

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/Esca6585dev/habarchy/backend/internal/ports"
)

func TestJSONPostProvider(t *testing.T) {
	var got map[string]any
	var auth string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		auth = r.Header.Get("Authorization")
		body, _ := io.ReadAll(r.Body)
		_ = json.Unmarshal(body, &got)
		if got["to"] == "+99365000000" {
			w.WriteHeader(500)
			_, _ = w.Write([]byte(`{"error":"gateway down"}`))
			return
		}
		if got["to"] == "+99365111111" {
			w.WriteHeader(400)
			_, _ = w.Write([]byte(`{"result":{"status":"INVALID_NUMBER"}}`))
			return
		}
		_, _ = w.Write([]byte(`{"result":{"status":"OK","id":123456}}`))
	}))
	defer srv.Close()

	p, err := New(Config{
		URL:          srv.URL + "/send",
		Headers:      map[string]string{"Authorization": "Bearer secret-token"},
		BodyTemplate: `{"to":"{{.To}}","text":{{.TextJSON}},"from":"{{.Sender}}"}`,
		Sender:       "HABARCHY",
		Success:      Matcher{JSONPath: "result.status", JSONEquals: "OK"},
		MessageID:    Extractor{JSONPath: "result.id"},
	}, srv.Client())
	if err != nil {
		t.Fatal(err)
	}
	res, err := p.Send(context.Background(), ports.SMSMessage{To: "+99365123456", Text: `Kod: "4821" täze`})
	if err != nil {
		t.Fatal(err)
	}
	if res.ProviderMessageID != "123456" || auth != "Bearer secret-token" || got["text"] != `Kod: "4821" täze` || got["from"] != "HABARCHY" {
		t.Fatalf("res=%+v got=%v auth=%q", res, got, auth)
	}

	_, err = p.Send(context.Background(), ports.SMSMessage{To: "+99365000000", Text: "x"})
	var pe *ports.ProviderError
	if !errors.As(err, &pe) || !pe.Retryable {
		t.Fatalf("500 should be retryable, got %v", err)
	}
	_, err = p.Send(context.Background(), ports.SMSMessage{To: "+99365111111", Text: "x"})
	if !errors.As(err, &pe) || pe.Retryable {
		t.Fatalf("400 should be permanent, got %v", err)
	}
}

func TestGetProviderWithRegex(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Query().Get("msg") != "Salam dünýä" || r.URL.Query().Get("phone") != "99365123456" {
			w.WriteHeader(400)
			return
		}
		_, _ = w.Write([]byte("OK: id=abc-789"))
	}))
	defer srv.Close()
	if _, err := New(Config{
		URL:     srv.URL + "/api?phone={{.To | trimPlus}}&msg={{.Text | urlquery}}",
		Method:  "GET",
		Success: Matcher{Regex: `^OK`},
	}, srv.Client()); err == nil {
		t.Fatal("unknown template func trimPlus must be rejected")
	}
	p, err := New(Config{
		URL:       srv.URL + "/api?phone={{slice .To 1}}&msg={{.Text | urlquery}}",
		Method:    "GET",
		Success:   Matcher{Regex: `^OK`},
		MessageID: Extractor{Regex: `id=([\w-]+)`},
	}, srv.Client())
	if err != nil {
		t.Fatal(err)
	}
	res, err := p.Send(context.Background(), ports.SMSMessage{To: "+99365123456", Text: "Salam dünýä"})
	if err != nil || res.ProviderMessageID != "abc-789" {
		t.Fatalf("%+v %v", res, err)
	}
}

func TestDLRParse(t *testing.T) {
	c := DLRConfig{MessageIDParam: "msgid", StatusParam: "status", DeliveredValues: []string{"DELIVRD"}, FailedValues: []string{"UNDELIV", "EXPIRED"}}
	id, d := c.ParseDLR(map[string]string{"msgid": "1", "status": "delivrd"})
	if id != "1" || d == nil || !*d {
		t.Fatal("expected delivered")
	}
	_, d = c.ParseDLR(map[string]string{"msgid": "1", "status": "EXPIRED"})
	if d == nil || *d {
		t.Fatal("expected failed")
	}
	_, d = c.ParseDLR(map[string]string{"msgid": "1", "status": "ENROUTE"})
	if d != nil {
		t.Fatal("expected unknown")
	}
}

func TestConfigValidation(t *testing.T) {
	if _, err := New(Config{}, nil); err == nil {
		t.Fatal("empty url accepted")
	}
	if _, err := New(Config{URL: "http://x", Method: "DELETE"}, nil); err == nil {
		t.Fatal("bad method accepted")
	}
	if _, err := New(Config{URL: "http://x", BodyTemplate: "{{.Oops"}, nil); err == nil || !strings.Contains(err.Error(), "body_template") {
		t.Fatalf("bad template accepted: %v", err)
	}
}
