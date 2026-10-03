package http

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/rs/zerolog"

	"github.com/Esca6585dev/habarchy/backend/internal/config"
	"github.com/Esca6585dev/habarchy/backend/internal/domain"
)

type fakePinger struct{ err error }

func (f fakePinger) Ping(context.Context) error { return f.err }

func testConfig(t *testing.T) *config.Config {
	t.Helper()
	t.Setenv("HABARCHY_DATABASE_URL", "postgres://x")
	t.Setenv("HABARCHY_JWT_SECRET", strings.Repeat("s", 32))
	t.Setenv("HABARCHY_MASTER_KEY", strings.Repeat("a", 64))
	cfg, err := config.Load()
	if err != nil {
		t.Fatal(err)
	}
	return cfg
}

func TestHealthAndReady(t *testing.T) {
	cfg := testConfig(t)
	app := NewServer(cfg, zerolog.Nop(), Deps{DB: fakePinger{}, Redis: fakePinger{errors.New("down")}})

	resp, err := app.Test(httptest.NewRequest("GET", "/healthz", nil))
	if err != nil || resp.StatusCode != 200 {
		t.Fatalf("healthz: %v %d", err, resp.StatusCode)
	}
	if resp.Header.Get("X-Request-Id") == "" {
		t.Fatal("request id header missing")
	}

	resp, err = app.Test(httptest.NewRequest("GET", "/readyz", nil))
	if err != nil || resp.StatusCode != 503 {
		t.Fatalf("readyz: %v %d", err, resp.StatusCode)
	}
	body, _ := io.ReadAll(resp.Body)
	var env struct {
		Data struct {
			Status string            `json:"status"`
			Checks map[string]string `json:"checks"`
		} `json:"data"`
	}
	if err := json.Unmarshal(body, &env); err != nil {
		t.Fatal(err)
	}
	if env.Data.Status != "degraded" || env.Data.Checks["postgres"] != "ok" || env.Data.Checks["redis"] != "down" {
		t.Fatalf("unexpected readyz body: %s", body)
	}

	// Client-supplied request ids are preserved.
	req := httptest.NewRequest("GET", "/healthz", nil)
	req.Header.Set("X-Request-Id", "abc-123")
	resp, _ = app.Test(req)
	if resp.Header.Get("X-Request-Id") != "abc-123" {
		t.Fatal("client request id not echoed")
	}
}

func TestErrorEnvelope(t *testing.T) {
	cfg := testConfig(t)
	app := NewServer(cfg, zerolog.Nop(), Deps{})
	app.Get("/boom", func(*fiberCtx) error {
		return domain.ErrQuotaExceeded.WithDetails(map[string]any{"limit": 100})
	})
	app.Get("/panic", func(*fiberCtx) error { panic("x") })

	resp, _ := app.Test(httptest.NewRequest("GET", "/boom", nil))
	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != 429 || !strings.Contains(string(body), `"code":"quota_exceeded"`) || !strings.Contains(string(body), `"limit":100`) {
		t.Fatalf("boom: %d %s", resp.StatusCode, body)
	}

	resp, _ = app.Test(httptest.NewRequest("GET", "/nowhere", nil))
	body, _ = io.ReadAll(resp.Body)
	if resp.StatusCode != 404 || !strings.Contains(string(body), `"code":"not_found"`) {
		t.Fatalf("404: %d %s", resp.StatusCode, body)
	}

	resp, _ = app.Test(httptest.NewRequest("GET", "/panic", nil))
	body, _ = io.ReadAll(resp.Body)
	if resp.StatusCode != 500 || !strings.Contains(string(body), `"code":"internal_error"`) || strings.Contains(string(body), "x\"") {
		t.Fatalf("panic: %d %s", resp.StatusCode, body)
	}
}

func TestLoggerStatusForErrors(t *testing.T) {
	cfg := testConfig(t)
	var buf strings.Builder
	log := zerolog.New(&buf)
	app := NewServer(cfg, log, Deps{})
	resp, _ := app.Test(httptest.NewRequest("GET", "/missing", nil))
	if resp.StatusCode != 404 {
		t.Fatalf("status %d", resp.StatusCode)
	}
	if !strings.Contains(buf.String(), `"status":404`) {
		t.Fatalf("request log should carry the real status: %s", buf.String())
	}
}
