package e2e

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/hibiken/asynq"
	"github.com/rs/zerolog"

	httpx "github.com/Esca6585dev/habarchy/backend/internal/adapters/http"
	"github.com/Esca6585dev/habarchy/backend/internal/adapters/http/admin"
	"github.com/Esca6585dev/habarchy/backend/internal/adapters/http/public"
	rds "github.com/Esca6585dev/habarchy/backend/internal/adapters/redis"
	"github.com/Esca6585dev/habarchy/backend/internal/app/auth"
	"github.com/Esca6585dev/habarchy/backend/internal/app/contacts"
	"github.com/Esca6585dev/habarchy/backend/internal/app/delivery"
	"github.com/Esca6585dev/habarchy/backend/internal/app/messages"
	"github.com/Esca6585dev/habarchy/backend/internal/app/otp"
	"github.com/Esca6585dev/habarchy/backend/internal/app/providers"
	"github.com/Esca6585dev/habarchy/backend/internal/app/webhooks"
	"github.com/Esca6585dev/habarchy/backend/internal/config"
	"github.com/Esca6585dev/habarchy/backend/internal/queue"
	"github.com/Esca6585dev/habarchy/backend/internal/worker"
	"github.com/Esca6585dev/habarchy/backend/pkg/crypto"
)

// flowEnv runs the API and an in-process asynq worker against real
// Postgres and Redis (HABARCHY_TEST_DATABASE_URL, HABARCHY_TEST_REDIS_URL).
type flowEnv struct {
	*env
	rdb      *rds.Client
	q        *queue.Client
	srv      *asynq.Server
	hooks    *hookSink
	gateway  *httptest.Server
	gwCalls  *gatewayLog
	apiKey   string
	project  admin.ProjectResponse
	adminTok string
}

type gatewayLog struct {
	mu    sync.Mutex
	calls []map[string]any
	fail  bool // when true the gateway answers 503
}

type hookSink struct {
	mu     sync.Mutex
	events []webhooks.Event
	bodies [][]byte
	sigs   []string
	secret []byte
	srv    *httptest.Server
}

func (h *hookSink) wait(t *testing.T, typ string) []webhooks.Event {
	return h.waitN(t, typ, 1)
}

func (h *hookSink) waitN(t *testing.T, typ string, n int) []webhooks.Event {
	t.Helper()
	deadline := time.Now().Add(15 * time.Second)
	for time.Now().Before(deadline) {
		h.mu.Lock()
		var got []webhooks.Event
		for _, e := range h.events {
			if string(e.Type) == typ {
				got = append(got, e)
			}
		}
		h.mu.Unlock()
		if len(got) >= n {
			return got
		}
		time.Sleep(50 * time.Millisecond)
	}
	h.mu.Lock()
	defer h.mu.Unlock()
	t.Fatalf("timeout waiting for %d %s webhooks; got %d events: %+v", n, typ, len(h.events), h.events)
	return nil
}

func newFlowEnv(t *testing.T) *flowEnv {
	t.Helper()
	redisURL := os.Getenv("HABARCHY_TEST_REDIS_URL")
	if redisURL == "" {
		t.Skip("HABARCHY_TEST_REDIS_URL not set")
	}
	base := newEnv(t) // migrated Postgres + admin app with auth/projects/templates

	ctx := context.Background()
	rdb, err := rds.Connect(ctx, config.Redis{URL: redisURL})
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = rdb.Close() })
	// Clean slate: asynq keys and our own state from previous runs.
	_ = rdb.Raw().FlushDB(ctx).Err()

	cfg := base.cfg
	cipher, _ := crypto.NewCipherFromString(cfg.Security.MasterKey)
	q := queue.NewClient(queue.RedisOpt(rdb.Options()), 3, 8)
	t.Cleanup(func() { _ = q.Close() })

	sink := &hookSink{secret: []byte("whsec_test")}
	sink.srv = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		var ev webhooks.Event
		_ = json.Unmarshal(body, &ev)
		sink.mu.Lock()
		sink.events = append(sink.events, ev)
		sink.bodies = append(sink.bodies, body)
		sink.sigs = append(sink.sigs, r.Header.Get(webhooks.HeaderSignature))
		sink.mu.Unlock()
		w.WriteHeader(200)
	}))
	t.Cleanup(sink.srv.Close)

	gw := &gatewayLog{}
	gateway := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		var call map[string]any
		_ = json.Unmarshal(body, &call)
		gw.mu.Lock()
		gw.calls = append(gw.calls, call)
		fail := gw.fail
		gw.mu.Unlock()
		if fail {
			w.WriteHeader(503)
			_, _ = w.Write([]byte(`{"status":"busy"}`))
			return
		}
		if strings.HasSuffix(fmt.Sprint(call["to"]), "999") {
			w.WriteHeader(400)
			_, _ = w.Write([]byte(`{"status":"INVALID"}`))
			return
		}
		fmt.Fprintf(w, `{"status":"OK","id":"gw-%d"}`, len(gw.calls))
	}))
	t.Cleanup(gateway.Close)

	contactSvc := contacts.New(base.db)
	providerSvc := providers.New(base.db, cipher)
	webhookSvc := webhooks.New(base.db, cipher, q, 8, nil)
	webhookSvc.AllowPrivate = true
	messageSvc := messages.New(base.db, rdb, q, contactSvc, base.templates, messages.Limits{IdempotencyTTL: time.Hour, BatchMaxRecipients: 1000})
	otpSvc := otp.New(rdb, messageSvc, otp.Limits{Length: 6, TTL: 5 * time.Minute, MaxAttempts: 3, PerAddressHour: 3, PerIPHour: 100})
	deliverySvc := delivery.New(base.db, rdb, providerSvc, webhookSvc, contactSvc, zerolog.Nop(), 3)

	app := httpx.NewServer(cfg, zerolog.Nop(), httpx.Deps{})
	(&admin.Handlers{Auth: base.auth, Projects: base.projects, Templates: base.templates, Providers: providerSvc}).Register(app)
	(&public.Handlers{Projects: base.projects, Templates: base.templates, Messages: messageSvc, OTP: otpSvc, Contacts: contactSvc,
		Providers: providerSvc, Delivery: deliverySvc, SignatureTolerance: cfg.Security.SignatureTolerance}).Register(app)
	base.app = app

	srv := asynq.NewServer(queue.RedisOpt(rdb.Options()), asynq.Config{
		Concurrency: 4, Queues: queue.Weights(), LogLevel: asynq.ErrorLevel, IsFailure: worker.IsFailure,
		RetryDelayFunc: func(int, error, *asynq.Task) time.Duration { return 300 * time.Millisecond }, // fast retries in tests
	})
	if err := srv.Start(worker.NewMux(deliverySvc, webhookSvc)); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(srv.Shutdown)

	fe := &flowEnv{env: base, rdb: rdb, q: q, srv: srv, hooks: sink, gateway: gateway, gwCalls: gw}
	fe.bootstrap(t)
	return fe
}

// bootstrap creates an admin, a project with webhook + API key, a live
// http_sms provider pointing at the fake gateway and an otp template.
func (f *flowEnv) bootstrap(t *testing.T) {
	t.Helper()
	ctx := context.Background()
	if _, err := f.auth.CreateUser(ctx, "flow@habarchy.tm", "flow-pass-1234", "Flow"); err != nil {
		t.Fatal(err)
	}
	type loginOut struct {
		Tokens auth.Tokens `json:"tokens"`
	}
	login := decode[loginOut](t, f.expect(f.do("POST", "/api/admin/auth/login", map[string]any{"email": "flow@habarchy.tm", "password": "flow-pass-1234"}, nil), 200).Data)
	f.adminTok = login.Tokens.AccessToken
	tok := bearer(f.adminTok)
	f.project = decode[admin.ProjectResponse](t, f.expect(f.do("POST", "/api/admin/projects", map[string]any{"name": "Flow", "daily_quota": 6}, tok), 201).Data)
	pbase := "/api/admin/projects/" + f.project.ID.String()
	f.expect(f.do("PATCH", pbase, map[string]any{"webhook_url": f.hooks.srv.URL + "/hook", "webhook_secret": string(f.hooks.secret)}, tok), 200)
	key := decode[admin.APIKeyResponse](t, f.expect(f.do("POST", pbase+"/api-keys", map[string]any{"name": "live", "live": true}, tok), 201).Data)
	f.apiKey = key.Key
	prov := decode[admin.ProviderResponse](t, f.expect(f.do("POST", pbase+"/providers", map[string]any{
		"name": "fake-gateway", "type": "http_sms", "priority": 10, "rate_limit_per_sec": 100,
		"credentials": map[string]any{
			"url": f.gateway.URL + "/send", "body_template": `{"to":"{{.To}}","text":{{.TextJSON}}}`,
			"success": map[string]any{"json_path": "status", "json_equals": "OK"}, "message_id": map[string]any{"json_path": "id"},
			"dlr": map[string]any{"message_id_param": "msgid", "status_param": "status", "delivered_values": []string{"DELIVRD"}, "failed_values": []string{"UNDELIV"}},
		},
	}, tok), 201).Data)
	if prov.Settings["url"] == nil || prov.Settings["headers"] != nil {
		t.Fatalf("provider settings: %+v", prov.Settings)
	}
	f.expect(f.do("POST", pbase+"/templates", map[string]any{"key": "otp", "channel": "sms", "body": "Kod: {{.code}} ({{.minutes}} min)"}, tok), 201)
}

func (f *flowEnv) key() map[string]string { return map[string]string{"X-Api-Key": f.apiKey} }

func (f *flowEnv) waitStatus(t *testing.T, id string, want string) public.MessageResponse {
	t.Helper()
	deadline := time.Now().Add(15 * time.Second)
	var last public.MessageResponse
	for time.Now().Before(deadline) {
		r := f.expect(f.do("GET", "/api/v1/messages/"+id, nil, f.key()), 200)
		var out struct {
			Message public.MessageResponse `json:"message"`
		}
		_ = json.Unmarshal(r.Data, &out)
		last = out.Message
		if last.Status == want {
			return last
		}
		time.Sleep(50 * time.Millisecond)
	}
	t.Fatalf("message %s: wanted status %s, last %+v", id, want, last)
	return last
}

func TestSendFlow(t *testing.T) {
	f := newFlowEnv(t)
	key := f.key()
	pid := f.project.ID.String()

	// ---- validation ----
	f.expect(f.do("POST", "/api/v1/messages", map[string]any{"channel": "sms", "to": "+99365123456"}, key), 400)
	f.expect(f.do("POST", "/api/v1/messages", map[string]any{"channel": "sms", "to": "abc", "body": "x"}, key), 400)
	f.expect(f.do("POST", "/api/v1/messages", map[string]any{"channel": "sms", "to": "+99365123456", "template": "missing", "data": map[string]any{}}, key), 400)
	f.expect(f.do("POST", "/api/v1/messages", map[string]any{"channel": "sms", "to": "+99365123456", "template": "otp", "data": map[string]any{"code": "1"}}, key), 400) // minutes missing

	// ---- happy path: template send -> worker -> gateway -> sent -> webhook ----
	type accepted struct {
		ID     uuid.UUID `json:"id"`
		Status string    `json:"status"`
	}
	acc := decode[accepted](t, f.expect(f.do("POST", "/api/v1/messages", map[string]any{
		"channel": "sms", "to": "65 12 34 56", "template": "otp", "data": map[string]any{"code": "4821", "minutes": 5},
		"idempotency_key": "req-1", "metadata": map[string]any{"order": 7},
	}, key), 202).Data)
	if acc.Status != "queued" {
		t.Fatalf("accepted: %+v", acc)
	}
	sent := f.waitStatus(t, acc.ID.String(), "sent")
	if sent.To != "+99365123456" || sent.Body != "Kod: 4821 (5 min)" || sent.ProviderMessageID != "gw-1" || sent.Attempts != 1 {
		t.Fatalf("sent: %+v", sent)
	}
	ev := f.hooks.wait(t, "message.sent")
	data := ev[0].Data.(map[string]any)
	if data["message_id"] != acc.ID.String() || data["metadata"].(map[string]any)["order"] != float64(7) {
		t.Fatalf("webhook data: %+v", data)
	}
	f.hooks.mu.Lock()
	if !webhooks.Verify(f.hooks.secret, f.hooks.sigs[0], f.hooks.bodies[0], 5*time.Minute, time.Now()) {
		t.Fatal("webhook signature invalid")
	}
	f.hooks.mu.Unlock()

	// Idempotent replay returns the original message with 200.
	r := f.expect(f.do("POST", "/api/v1/messages", map[string]any{"channel": "sms", "to": "+99365123456", "body": "other", "idempotency_key": "req-1"}, key), 200)
	var meta map[string]any
	_ = json.Unmarshal(r.Meta, &meta)
	if meta["duplicate"] != true || !strings.Contains(string(r.Data), acc.ID.String()) {
		t.Fatalf("replay: %s %s", r.Meta, r.Data)
	}

	// ---- DLR callback -> delivered -> webhook ----
	provs := decode[[]admin.ProviderResponse](t, f.expect(f.do("GET", "/api/admin/projects/"+pid+"/providers", nil, bearer(f.adminTok)), 200).Data)
	cb := "/callbacks/sms/" + provs[0].ID.String()
	f.expect(f.do("POST", cb, map[string]any{"msgid": "gw-1", "status": "ENROUTE"}, nil), 200)
	f.expect(f.do("POST", cb, map[string]any{"msgid": "gw-1", "status": "DELIVRD"}, nil), 200)
	delivered := f.waitStatus(t, acc.ID.String(), "delivered")
	if delivered.DeliveredAt == nil {
		t.Fatal("delivered_at missing")
	}
	f.hooks.wait(t, "message.delivered")
	f.expect(f.do("POST", "/callbacks/sms/"+uuid.NewString(), map[string]any{"msgid": "x", "status": "DELIVRD"}, nil), 404)

	detail := f.expect(f.do("GET", "/api/v1/messages/"+acc.ID.String(), nil, key), 200)
	var d struct {
		Events []public.EventResponse `json:"events"`
	}
	_ = json.Unmarshal(detail.Data, &d)
	types := []string{}
	for _, e := range d.Events {
		types = append(types, e.Type)
	}
	if strings.Join(types, ",") != "queued,sent,delivered" {
		t.Fatalf("timeline: %v", types)
	}

	// ---- permanent provider error -> failed, no retry ----
	bad := decode[accepted](t, f.expect(f.do("POST", "/api/v1/messages", map[string]any{"channel": "sms", "to": "+99365000999", "body": "x"}, key), 202).Data)
	failed := f.waitStatus(t, bad.ID.String(), "failed")
	if failed.Attempts != 1 || failed.ErrorCode != "provider_rejected" {
		t.Fatalf("failed: %+v", failed)
	}
	f.hooks.wait(t, "message.failed")

	// ---- retryable error -> retries, then succeeds once gateway recovers ----
	f.gwCalls.mu.Lock()
	f.gwCalls.fail = true
	f.gwCalls.mu.Unlock()
	rt := decode[accepted](t, f.expect(f.do("POST", "/api/v1/messages", map[string]any{"channel": "sms", "to": "+99365111111", "body": "retry me"}, key), 202).Data)
	// Wait until the gateway has rejected the first attempt, then recover.
	deadline := time.Now().Add(10 * time.Second)
	for {
		f.gwCalls.mu.Lock()
		n := 0
		for _, c := range f.gwCalls.calls {
			if c["to"] == "+99365111111" {
				n++
			}
		}
		f.gwCalls.mu.Unlock()
		if n >= 1 || time.Now().After(deadline) {
			break
		}
		time.Sleep(20 * time.Millisecond)
	}
	f.gwCalls.mu.Lock()
	f.gwCalls.fail = false
	f.gwCalls.mu.Unlock()
	rsent := f.waitStatus(t, rt.ID.String(), "sent")
	if rsent.Attempts < 2 {
		t.Fatalf("expected retries, attempts=%d", rsent.Attempts)
	}

	// ---- scheduled message can be cancelled ----
	sch := decode[accepted](t, f.expect(f.do("POST", "/api/v1/messages", map[string]any{"channel": "sms", "to": "+99365222222", "body": "later", "scheduled_at": time.Now().Add(time.Hour).Format(time.RFC3339)}, key), 202).Data)
	canc := decode[public.MessageResponse](t, f.expect(f.do("POST", "/api/v1/messages/"+sch.ID.String()+"/cancel", nil, key), 200).Data)
	if canc.Status != "cancelled" {
		t.Fatalf("cancel: %+v", canc)
	}
	f.expect(f.do("POST", "/api/v1/messages/"+acc.ID.String()+"/cancel", nil, key), 409)

	// ---- list with filters and cursor ----
	list := f.expect(f.do("GET", "/api/v1/messages?limit=2", nil, key), 200)
	var page []public.MessageResponse
	_ = json.Unmarshal(list.Data, &page)
	_ = json.Unmarshal(list.Meta, &meta)
	if len(page) != 2 || meta["next_cursor"] == nil {
		t.Fatalf("page: %d %v", len(page), meta)
	}
	list = f.expect(f.do("GET", "/api/v1/messages?status=failed", nil, key), 200)
	_ = json.Unmarshal(list.Data, &page)
	if len(page) != 1 || page[0].ID != bad.ID {
		t.Fatalf("filter failed: %+v", page)
	}
	f.expect(f.do("GET", "/api/v1/messages?status=bogus", nil, key), 400)

	// ---- batch with contacts, one invalid, completion webhook ----
	c1 := decode[public.ContactResponse](t, f.expect(f.do("POST", "/api/v1/contacts", map[string]any{"external_id": "u1", "phone": "+99365333333", "attributes": map[string]any{"name": "Aman"}}, key), 201).Data)
	f.expect(f.do("POST", "/api/v1/contacts", map[string]any{"external_id": "u1", "phone": "+99365333334"}, key), 200) // upsert
	batch := f.expect(f.do("POST", "/api/v1/messages/batch", map[string]any{
		"channel": "sms", "body": "Salam {{.name}}!",
		"recipients": []map[string]any{
			{"to": map[string]any{"contact_id": c1.ID}, "data": map[string]any{}},
			{"to": "+99365444444", "data": map[string]any{"name": "Merdan"}},
			{"to": "bad-number", "data": map[string]any{"name": "X"}},
		},
	}, key), 202)
	var b public.BatchResponse
	_ = json.Unmarshal(batch.Data, &b)
	meta = map[string]any{}
	_ = json.Unmarshal(batch.Meta, &meta)
	if b.Total != 2 || meta["accepted"] != float64(2) || len(meta["rejected"].([]any)) != 1 {
		t.Fatalf("batch: %+v meta=%v", b, meta)
	}
	f.hooks.wait(t, "batch.completed")
	got := decode[public.BatchResponse](t, f.expect(f.do("GET", "/api/v1/batches/"+b.ID.String(), nil, key), 200).Data)
	if got.Status != "completed" || got.Sent != 2 || got.Queued != 0 {
		t.Fatalf("batch status: %+v", got)
	}
	// Contact attributes were available to the template.
	f.gwCalls.mu.Lock()
	texts := []string{}
	for _, c := range f.gwCalls.calls {
		texts = append(texts, fmt.Sprint(c["text"]))
	}
	f.gwCalls.mu.Unlock()
	if !contains(texts, "Salam Aman!") || !contains(texts, "Salam Merdan!") {
		t.Fatalf("gateway texts: %v", texts)
	}

	// ---- quota: daily_quota=6 live messages counted so far: 1+1+1+1(cancelled not counted)+2 = 6 ----
	r = f.expect(f.do("POST", "/api/v1/messages", map[string]any{"channel": "sms", "to": "+99365555555", "body": "over"}, key), 429)
	if r.Error.Code != "quota_exceeded" {
		t.Fatalf("quota: %s", r.Raw)
	}

	// ---- OTP: test key bypasses quota and providers (sandbox) ----
	tkey := decode[admin.APIKeyResponse](t, f.expect(f.do("POST", "/api/admin/projects/"+pid+"/api-keys", map[string]any{"name": "test"}, bearer(f.adminTok)), 201).Data)
	tk := map[string]string{"X-Api-Key": tkey.Key}
	otpRes := decode[otp.SendResult](t, f.expect(f.do("POST", "/api/v1/otp/send", map[string]any{"channel": "sms", "to": "+99365777777"}, tk), 202).Data)
	if otpRes.Length != 6 || otpRes.ExpiresIn != 300 {
		t.Fatalf("otp send: %+v", otpRes)
	}
	om := f.waitStatus(t, otpRes.MessageID.String(), "delivered") // sandbox delivers immediately
	if !om.IsTest || !strings.HasPrefix(om.Body, "Kod: ") {
		t.Fatalf("otp message: %+v", om)
	}
	code := strings.TrimSpace(strings.Split(strings.TrimPrefix(om.Body, "Kod: "), " ")[0])
	r = f.expect(f.do("POST", "/api/v1/otp/verify", map[string]any{"to": "+99365777777", "code": "000000"}, tk), 401)
	if r.Error.Details["attempts_remaining"] != float64(2) {
		t.Fatalf("otp wrong: %s", r.Raw)
	}
	vr := decode[otp.VerifyResult](t, f.expect(f.do("POST", "/api/v1/otp/verify", map[string]any{"to": "99365777777", "code": code}, tk), 200).Data)
	if !vr.Verified {
		t.Fatalf("otp verify: %+v", vr)
	}
	f.expect(f.do("POST", "/api/v1/otp/verify", map[string]any{"to": "+99365777777", "code": code}, tk), 401) // consumed
	// per-address rate limit: 3/hour
	f.expect(f.do("POST", "/api/v1/otp/send", map[string]any{"to": "+99365777777"}, tk), 202)
	f.expect(f.do("POST", "/api/v1/otp/send", map[string]any{"to": "+99365777777"}, tk), 202)
	f.expect(f.do("POST", "/api/v1/otp/send", map[string]any{"to": "+99365777777"}, tk), 429)

	// ---- devices ----
	dev := decode[public.DeviceResponse](t, f.expect(f.do("POST", "/api/v1/devices", map[string]any{"token": strings.Repeat("t", 40), "platform": "android", "external_id": "u1"}, key), 201).Data)
	if dev.ContactID == nil || *dev.ContactID != c1.ID || strings.Contains(dev.TokenHint, strings.Repeat("t", 20)) {
		t.Fatalf("device: %+v", dev)
	}
	// auto channel: project order default push,sms,email -> contact has a device -> push via... no push provider -> failed provider_unavailable
	auto := decode[accepted](t, f.expect(f.do("POST", "/api/v1/messages", map[string]any{"channel": "auto", "to": map[string]any{"external_id": "u1"}, "body": "hi", "title": "T"}, tk), 202).Data)
	am := f.waitStatus(t, auto.ID.String(), "delivered") // test key -> sandbox push
	if am.Channel != "push" || am.Subject != "T" {
		t.Fatalf("auto: %+v", am)
	}
	f.expect(f.do("DELETE", "/api/v1/devices/"+strings.Repeat("t", 40), nil, key), 204)

	// ---- admin test-send hits the provider directly ----
	ts := f.expect(f.do("POST", "/api/admin/projects/"+pid+"/providers/"+provs[0].ID.String()+"/test", map[string]any{"to": "+99365888888", "text": "ping"}, bearer(f.adminTok)), 200)
	if !strings.Contains(string(ts.Data), `"ok":true`) {
		t.Fatalf("test send: %s", ts.Data)
	}

	// ---- webhook log has deliveries with 200 responses ----
	n := 0
	for _, s := range []string{"message.sent", "message.delivered", "message.failed", "batch.completed"} {
		n += len(f.hooks.waitN(t, s, 1))
	}
	if n < 4 {
		t.Fatalf("webhooks: %d", n)
	}
	_ = fiber.StatusOK
}

func contains(list []string, s string) bool {
	for _, x := range list {
		if x == s {
			return true
		}
	}
	return false
}
