package gateway_test

import (
	"context"
	"encoding/json"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/rs/zerolog"

	httpx "github.com/Esca6585dev/habarchy/backend/internal/adapters/http"
	gatewayhttp "github.com/Esca6585dev/habarchy/backend/internal/adapters/http/gateway"
	"github.com/Esca6585dev/habarchy/backend/internal/adapters/postgres/pgtest"
	"github.com/Esca6585dev/habarchy/backend/internal/adapters/providers/sms/androidgw"
	"github.com/Esca6585dev/habarchy/backend/internal/app/auth"
	"github.com/Esca6585dev/habarchy/backend/internal/app/gateway"
	"github.com/Esca6585dev/habarchy/backend/internal/app/projects"
	"github.com/Esca6585dev/habarchy/backend/internal/app/providers"
	"github.com/Esca6585dev/habarchy/backend/internal/app/stats"
	"github.com/Esca6585dev/habarchy/backend/internal/config"
	"github.com/Esca6585dev/habarchy/backend/internal/domain"
	"github.com/Esca6585dev/habarchy/backend/internal/ports"
	"github.com/Esca6585dev/habarchy/backend/pkg/crypto"
	"github.com/Esca6585dev/habarchy/backend/pkg/password"
)

type env struct {
	t         *testing.T
	app       *fiber.App
	gateway   *gateway.Service
	providers *providers.Service
	stats     *stats.Service
	projectID uuid.UUID
	prov      uuid.UUID
	key       string
}

func newEnv(t *testing.T) *env {
	t.Helper()
	db := pgtest.Open(t)
	t.Setenv("HABARCHY_DATABASE_URL", "postgres://unused")
	t.Setenv("HABARCHY_JWT_SECRET", strings.Repeat("j", 32))
	t.Setenv("HABARCHY_MASTER_KEY", strings.Repeat("ab", 32))
	cfg, err := config.Load()
	if err != nil {
		t.Fatal(err)
	}
	cipher, _ := crypto.NewCipherFromString(cfg.Security.MasterKey)
	ctx := context.Background()
	fast := password.Params{Memory: 8 * 1024, Iterations: 1, Parallelism: 1, SaltLength: 16, KeyLength: 32}
	authSvc := auth.New(db, cfg.Auth, cipher, auth.WithPasswordHasher(func(p string) (string, error) { return password.HashWithParams(p, fast) }))
	user, err := authSvc.CreateUser(ctx, "gw-"+uuid.NewString()[:8]+"@example.com", "password-123", "GW")
	if err != nil {
		t.Fatal(err)
	}
	projectSvc := projects.New(db, cipher)
	project, err := projectSvc.Create(ctx, user.ID, projects.CreateInput{Name: "GW " + uuid.NewString()[:6]})
	if err != nil {
		t.Fatal(err)
	}
	providerSvc := providers.New(db, cipher)
	// No gateway_key given: the service generates one.
	prov, err := providerSvc.Create(ctx, project.ID, providers.Input{Name: "Phone", Type: domain.ProviderAndroidSMS, Credentials: json.RawMessage(`{"sim_slot":1,"timeout_sec":10}`)})
	if err != nil {
		t.Fatal(err)
	}
	settings, err := providerSvc.Redacted(prov)
	if err != nil {
		t.Fatal(err)
	}
	key, _ := settings["gateway_key"].(string)
	if !strings.HasPrefix(key, "gw_") || len(key) < 20 {
		t.Fatalf("generated key %q", key)
	}
	gw := gateway.New(db, nil, nil)
	gw.LeasePoll = 20 * time.Millisecond
	providerSvc.GatewayFactory = gw.Factory

	app := httpx.NewServer(cfg, zerolog.Nop(), httpx.Deps{})
	(&gatewayhttp.Handlers{Gateway: gw, MaxWait: 2 * time.Second}).Register(app)
	return &env{t: t, app: app, gateway: gw, providers: providerSvc, stats: stats.New(db, nil), projectID: project.ID, prov: prov.ID, key: key}
}

func (e *env) call(method, path, body string, key string) (int, map[string]any) {
	e.t.Helper()
	req := httptest.NewRequest(method, path, strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	if key != "" {
		req.Header.Set("X-Gateway-Key", key)
	}
	res, err := e.app.Test(req, 10_000)
	if err != nil {
		e.t.Fatal(err)
	}
	var out map[string]any
	_ = json.NewDecoder(res.Body).Decode(&out)
	return res.StatusCode, out
}

func TestGatewayFlowPhoneSendsAndReportsDelivered(t *testing.T) {
	e := newEnv(t)
	ctx := context.Background()

	if st, _ := e.call("GET", "/api/gateway/v1/me", "", "wrong-key-wrong-key"); st != 401 {
		t.Fatalf("wrong key: %d", st)
	}
	st, me := e.call("GET", "/api/gateway/v1/me", "", e.key)
	if st != 200 || me["data"].(map[string]any)["provider"] != "Phone" {
		t.Fatalf("me: %d %v", st, me)
	}

	// Empty outbox: long poll returns quickly with wait=0.
	st, body := e.call("GET", "/api/gateway/v1/outbox?wait=0", "", e.key)
	if st != 200 || len(body["data"].([]any)) != 0 {
		t.Fatalf("empty outbox: %d %v", st, body)
	}

	// Worker side: the provider adapter enqueues and waits.
	prov, _ := e.providers.GetByID(ctx, e.prov)
	adapter, err := e.providers.Build(ctx, prov)
	if err != nil {
		t.Fatal(err)
	}
	sms := adapter.(*androidgw.Provider)
	sms.Poll = 20 * time.Millisecond
	msgID := uuid.New()
	type sendOut struct {
		res *ports.SendResult
		err error
	}
	done := make(chan sendOut, 1)
	go func() {
		r, err := sms.Send(ctx, ports.SMSMessage{To: "+99365123456", Text: "Salam", Ref: msgID.String()})
		done <- sendOut{r, err}
	}()

	// Phone side: lease, send, report.
	st, body = e.call("GET", "/api/gateway/v1/outbox?wait=2&limit=5", "", e.key)
	items := body["data"].([]any)
	if st != 200 || len(items) != 1 {
		t.Fatalf("lease: %d %v", st, body)
	}
	item := items[0].(map[string]any)
	if item["to"] != "+99365123456" || item["text"] != "Salam" || item["sim_slot"].(float64) != 1 || item["message_id"] != msgID.String() {
		t.Fatalf("item %v", item)
	}
	// Second lease must not hand the same row out again.
	if _, body2 := e.call("GET", "/api/gateway/v1/outbox?wait=0", "", e.key); len(body2["data"].([]any)) != 0 {
		t.Fatalf("row leased twice: %v", body2)
	}
	outboxID := item["id"].(string)
	st, body = e.call("POST", "/api/gateway/v1/outbox/"+outboxID+"/result", `{"status":"sent","parts":1}`, e.key)
	if st != 200 || body["data"].(map[string]any)["applied"] != true {
		t.Fatalf("report sent: %d %v", st, body)
	}
	out := <-done
	if out.err != nil || out.res.ProviderMessageID != outboxID {
		t.Fatalf("worker result: %+v %v", out.res, out.err)
	}

	// Delivery report later: recorded on the row (no delivery service in this test).
	st, body = e.call("POST", "/api/gateway/v1/outbox/"+outboxID+"/result", `{"status":"delivered"}`, e.key)
	if st != 200 || body["data"].(map[string]any)["status"] != "delivered" {
		t.Fatalf("report delivered: %d %v", st, body)
	}
	// Duplicate / wrong id.
	if st, _ := e.call("POST", "/api/gateway/v1/outbox/"+uuid.NewString()+"/result", `{"status":"sent"}`, e.key); st != 404 {
		t.Fatalf("unknown outbox id: %d", st)
	}
	if st, _ := e.call("POST", "/api/gateway/v1/outbox/"+outboxID+"/result", `{"status":"bogus"}`, e.key); st != 400 {
		t.Fatalf("bad status: %d", st)
	}

	// Heartbeat shows up on the health page as online.
	st, _ = e.call("POST", "/api/gateway/v1/heartbeat", `{"battery":77,"network":"LTE","app_version":"0.1.0","gateway_key":"must-not-be-stored"}`, e.key)
	if st != 200 {
		t.Fatalf("heartbeat: %d", st)
	}
	h, err := e.stats.Health(ctx, e.projectID)
	if err != nil {
		t.Fatal(err)
	}
	var found bool
	for _, p := range h.Providers {
		if p.ID == e.prov {
			found = true
			if p.Gateway == nil || !p.Gateway.Online || strings.Contains(string(p.Gateway.Info), "must-not-be-stored") || !strings.Contains(string(p.Gateway.Info), "LTE") {
				t.Fatalf("gateway health %+v", p.Gateway)
			}
			if p.Status != "idle" { // no messages counted against this provider in the test
				t.Fatalf("status %s", p.Status)
			}
		}
	}
	if !found {
		t.Fatal("provider missing from health")
	}

	// Inbound SMS is stored.
	st, body = e.call("POST", "/api/gateway/v1/inbound", `{"from":"+99365000111","text":"Salam, bu test"}`, e.key)
	if st != 201 || body["data"].(map[string]any)["id"] == nil {
		t.Fatalf("inbound: %d %v", st, body)
	}
}

func TestGatewayPhoneFailureAndTimeout(t *testing.T) {
	e := newEnv(t)
	ctx := context.Background()
	prov, _ := e.providers.GetByID(ctx, e.prov)
	adapter, _ := e.providers.Build(ctx, prov)
	sms := adapter.(*androidgw.Provider)
	sms.Poll = 20 * time.Millisecond

	// Phone reports a permanent failure.
	done := make(chan error, 1)
	go func() {
		_, err := sms.Send(ctx, ports.SMSMessage{To: "+99365123456", Text: "x", Ref: uuid.NewString()})
		done <- err
	}()
	_, body := e.call("GET", "/api/gateway/v1/outbox?wait=2", "", e.key)
	id := body["data"].([]any)[0].(map[string]any)["id"].(string)
	e.call("POST", "/api/gateway/v1/outbox/"+id+"/result", `{"status":"failed","error_code":"null_pdu","error_message":"bad"}`, e.key)
	err := <-done
	pe, ok := err.(*ports.ProviderError) //nolint:errorlint // direct return
	if !ok || pe.Retryable || pe.Code != "gateway_null_pdu" {
		t.Fatalf("got %v", err)
	}

	// Phone offline: the worker gives up after the context deadline and the row expires.
	cctx, cancel := context.WithTimeout(ctx, 150*time.Millisecond)
	defer cancel()
	_, err = sms.Send(cctx, ports.SMSMessage{To: "+99365123456", Text: "y", Ref: uuid.NewString()})
	pe, ok = err.(*ports.ProviderError) //nolint:errorlint // direct return
	if !ok || !pe.Retryable {
		t.Fatalf("offline: %v", err)
	}
	if _, body := e.call("GET", "/api/gateway/v1/outbox?wait=0", "", e.key); len(body["data"].([]any)) != 0 {
		t.Fatalf("expired row leased: %v", body)
	}

	// Rotating credentials re-registers the key; the old one stops working.
	if _, err := e.providers.Update(ctx, e.projectID, e.prov, providers.Input{Credentials: json.RawMessage(`{"gateway_key":"new-key-0123456789abcdef"}`)}); err != nil {
		t.Fatal(err)
	}
	if st, _ := e.call("GET", "/api/gateway/v1/me", "", e.key); st != 401 {
		t.Fatalf("old key still works: %d", st)
	}
	if st, _ := e.call("GET", "/api/gateway/v1/me", "", "new-key-0123456789abcdef"); st != 200 {
		t.Fatalf("new key rejected: %d", st)
	}
	// Disabled provider cannot poll.
	off := false
	if _, err := e.providers.Update(ctx, e.projectID, e.prov, providers.Input{IsActive: &off}); err != nil {
		t.Fatal(err)
	}
	if st, _ := e.call("GET", "/api/gateway/v1/me", "", "new-key-0123456789abcdef"); st != 401 {
		t.Fatalf("disabled provider accepted: %d", st)
	}
}
