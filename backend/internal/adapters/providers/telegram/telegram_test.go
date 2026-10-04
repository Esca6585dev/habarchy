package telegram

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

func TestSend(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !strings.HasPrefix(r.URL.Path, "/bot123:ABC/sendMessage") {
			t.Errorf("path %s", r.URL.Path)
		}
		body, _ := io.ReadAll(r.Body)
		var req sendMessageReq
		_ = json.Unmarshal(body, &req)
		switch req.ChatID {
		case "blocked":
			w.WriteHeader(403)
			_, _ = w.Write([]byte(`{"ok":false,"error_code":403,"description":"Forbidden: bot was blocked by the user"}`))
		case "flood":
			w.WriteHeader(429)
			_, _ = w.Write([]byte(`{"ok":false,"error_code":429,"description":"Too Many Requests","parameters":{"retry_after":5}}`))
		default:
			if req.ParseMode != "HTML" || req.ReplyMarkup == nil || req.ReplyMarkup.InlineKeyboard[0][0].URL != "https://tds.gov.tm" {
				t.Errorf("payload %s", body)
			}
			_, _ = w.Write([]byte(`{"ok":true,"result":{"message_id":42}}`))
		}
	}))
	defer srv.Close()
	p, err := New(Config{BotToken: "123:ABC", ParseMode: "HTML", Endpoint: srv.URL}, srv.Client())
	if err != nil {
		t.Fatal(err)
	}
	res, err := p.Send(context.Background(), ports.TelegramMessage{ChatID: "777", Text: "<b>Salam</b>", Buttons: [][]ports.TelegramButton{{{Text: "Portal", URL: "https://tds.gov.tm"}}}})
	if err != nil || res.ProviderMessageID != "42" {
		t.Fatalf("%+v %v", res, err)
	}
	var pe *ports.ProviderError
	_, err = p.Send(context.Background(), ports.TelegramMessage{ChatID: "blocked", Text: "x"})
	if !errors.As(err, &pe) || pe.Retryable || pe.Code != "invalid_recipient" {
		t.Fatalf("blocked: %v", err)
	}
	_, err = p.Send(context.Background(), ports.TelegramMessage{ChatID: "flood", Text: "x"})
	if !errors.As(err, &pe) || !pe.Retryable || pe.Raw["retry_after"] != 5 {
		t.Fatalf("flood: %v", err)
	}
	if _, err := New(Config{BotToken: "nocolon"}, nil); err == nil {
		t.Fatal("bad token accepted")
	}
}
