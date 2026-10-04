// Package e2e exercises the admin and public APIs end to end against a
// real Postgres (HABARCHY_TEST_DATABASE_URL).
package e2e

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/pquerna/otp/totp"
	"github.com/rs/zerolog"

	httpx "github.com/Esca6585dev/habarchy/backend/internal/adapters/http"
	"github.com/Esca6585dev/habarchy/backend/internal/adapters/http/admin"
	"github.com/Esca6585dev/habarchy/backend/internal/adapters/http/public"
	"github.com/Esca6585dev/habarchy/backend/internal/adapters/postgres"
	"github.com/Esca6585dev/habarchy/backend/internal/adapters/postgres/pgtest"
	"github.com/Esca6585dev/habarchy/backend/internal/app/auth"
	"github.com/Esca6585dev/habarchy/backend/internal/app/projects"
	"github.com/Esca6585dev/habarchy/backend/internal/app/templates"
	"github.com/Esca6585dev/habarchy/backend/internal/config"
	"github.com/Esca6585dev/habarchy/backend/pkg/crypto"
	"github.com/Esca6585dev/habarchy/backend/pkg/password"
)

type env struct {
	t         *testing.T
	app       *fiber.App
	auth      *auth.Service
	projects  *projects.Service
	templates *templates.Service
	db        *postgres.DB
	cfg       *config.Config
	now       time.Time
	clock     func() time.Time
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
	e := &env{t: t, now: time.Now(), db: db, cfg: cfg}
	e.clock = func() time.Time { return e.now }
	fast := password.Params{Memory: 8 * 1024, Iterations: 1, Parallelism: 1, SaltLength: 16, KeyLength: 32}
	e.auth = auth.New(db, cfg.Auth, cipher, auth.WithClock(e.clock), auth.WithPasswordHasher(func(p string) (string, error) {
		return password.HashWithParams(p, fast)
	}))
	projectSvc := projects.New(db, cipher)
	templateSvc := templates.New(db)
	e.projects, e.templates = projectSvc, templateSvc

	e.app = httpx.NewServer(cfg, zerolog.Nop(), httpx.Deps{})
	(&admin.Handlers{Auth: e.auth, Projects: projectSvc, Templates: templateSvc}).Register(e.app)
	(&public.Handlers{Projects: projectSvc, Templates: templateSvc, SignatureTolerance: cfg.Security.SignatureTolerance}).Register(e.app)
	return e
}

type resp struct {
	Status int
	Data   json.RawMessage
	Meta   json.RawMessage
	Error  *struct {
		Code    string         `json:"code"`
		Message string         `json:"message"`
		Details map[string]any `json:"details"`
	}
	Raw []byte
}

func (e *env) do(method, path string, body any, headers map[string]string) resp {
	e.t.Helper()
	var buf bytes.Buffer
	if body != nil {
		if b, ok := body.([]byte); ok {
			buf.Write(b)
		} else if err := json.NewEncoder(&buf).Encode(body); err != nil {
			e.t.Fatal(err)
		}
	}
	req := httptest.NewRequest(method, path, &buf)
	req.Header.Set("Content-Type", "application/json")
	for k, v := range headers {
		req.Header.Set(k, v)
	}
	res, err := e.app.Test(req, 10_000)
	if err != nil {
		e.t.Fatal(err)
	}
	raw, _ := io.ReadAll(res.Body)
	out := resp{Status: res.StatusCode, Raw: raw}
	if len(raw) > 0 {
		if err := json.Unmarshal(raw, &out); err != nil {
			e.t.Fatalf("%s %s: bad json %q", method, path, raw)
		}
	}
	return out
}

func (e *env) expect(r resp, status int) resp {
	e.t.Helper()
	if r.Status != status {
		e.t.Fatalf("expected %d, got %d: %s", status, r.Status, r.Raw)
	}
	return r
}

func decode[T any](t *testing.T, raw json.RawMessage) T {
	t.Helper()
	var v T
	if err := json.Unmarshal(raw, &v); err != nil {
		t.Fatalf("decode: %v (%s)", err, raw)
	}
	return v
}

func bearer(token string) map[string]string {
	return map[string]string{"Authorization": "Bearer " + token}
}

func TestAdminAndPublicFlow(t *testing.T) {
	e := newEnv(t)
	ctx := context.Background()

	if _, err := e.auth.CreateUser(ctx, "owner@habarchy.tm", "owner-pass-123", "Owner"); err != nil {
		t.Fatal(err)
	}
	if _, err := e.auth.CreateUser(ctx, "dev@habarchy.tm", "dev-pass-1234", "Dev"); err != nil {
		t.Fatal(err)
	}
	if _, err := e.auth.CreateUser(ctx, "owner@habarchy.tm", "x-pass-12345", ""); err == nil {
		t.Fatal("duplicate email must fail")
	}

	// ---- login ----
	e.expect(e.do("POST", "/api/admin/auth/login", map[string]any{"email": "owner@habarchy.tm", "password": "wrong-pass"}, nil), 401)
	r := e.expect(e.do("POST", "/api/admin/auth/login", map[string]any{"email": "not-an-email"}, nil), 400)
	if r.Error.Code != "validation_failed" || r.Error.Details["email"] == nil || r.Error.Details["password"] == nil {
		t.Fatalf("validation details: %s", r.Raw)
	}
	type loginOut struct {
		Tokens auth.Tokens        `json:"tokens"`
		User   admin.UserResponse `json:"user"`
	}
	login := decode[loginOut](t, e.expect(e.do("POST", "/api/admin/auth/login", map[string]any{"email": "Owner@habarchy.tm", "password": "owner-pass-123"}, nil), 200).Data)
	if login.Tokens.AccessToken == "" || login.Tokens.RefreshToken == "" || login.User.Email != "owner@habarchy.tm" {
		t.Fatalf("login: %+v", login)
	}
	owner := bearer(login.Tokens.AccessToken)
	devLogin := decode[loginOut](t, e.expect(e.do("POST", "/api/admin/auth/login", map[string]any{"email": "dev@habarchy.tm", "password": "dev-pass-1234"}, nil), 200).Data)
	dev := bearer(devLogin.Tokens.AccessToken)

	e.expect(e.do("GET", "/api/admin/me", nil, nil), 401)
	e.expect(e.do("GET", "/api/admin/me", nil, bearer("garbage")), 401)
	me := decode[admin.UserResponse](t, e.expect(e.do("GET", "/api/admin/me", nil, owner), 200).Data)
	if me.Email != "owner@habarchy.tm" || me.TOTPEnabled {
		t.Fatalf("me: %+v", me)
	}

	// ---- refresh rotation + reuse detection ----
	type tokOut struct {
		Tokens auth.Tokens `json:"tokens"`
	}
	rot := decode[tokOut](t, e.expect(e.do("POST", "/api/admin/auth/refresh", map[string]any{"refresh_token": login.Tokens.RefreshToken}, nil), 200).Data)
	if rot.Tokens.RefreshToken == login.Tokens.RefreshToken {
		t.Fatal("refresh token must rotate")
	}
	// Reusing the old token revokes the whole family, including the new one.
	e.expect(e.do("POST", "/api/admin/auth/refresh", map[string]any{"refresh_token": login.Tokens.RefreshToken}, nil), 401)
	e.expect(e.do("POST", "/api/admin/auth/refresh", map[string]any{"refresh_token": rot.Tokens.RefreshToken}, nil), 401)
	// Access tokens expire.
	e.now = e.now.Add(16 * time.Minute)
	e.expect(e.do("GET", "/api/admin/me", nil, owner), 401)
	e.now = e.now.Add(-16 * time.Minute)

	// ---- projects ----
	e.expect(e.do("POST", "/api/admin/projects", map[string]any{"name": "X"}, owner), 400)
	proj := decode[admin.ProjectResponse](t, e.expect(e.do("POST", "/api/admin/projects", map[string]any{"name": "TDS Portal", "daily_quota": 1000}, owner), 201).Data)
	if proj.Slug != "tds-portal" || proj.Role != "owner" || proj.DefaultLocale != "tk" {
		t.Fatalf("project: %+v", proj)
	}
	e.expect(e.do("POST", "/api/admin/projects", map[string]any{"name": "Other", "slug": "tds-portal"}, owner), 409)
	pbase := "/api/admin/projects/" + proj.ID.String()

	// Dev is not a member yet: project is invisible to them.
	e.expect(e.do("GET", pbase, nil, dev), 404)
	list := decode[[]admin.ProjectResponse](t, e.expect(e.do("GET", "/api/admin/projects", nil, dev), 200).Data)
	if len(list) != 0 {
		t.Fatalf("dev should see no projects, got %d", len(list))
	}

	// Add dev as developer; developers cannot manage members or keys.
	e.expect(e.do("PUT", pbase+"/members", map[string]any{"email": "dev@habarchy.tm", "role": "developer"}, owner), 200)
	e.expect(e.do("PUT", pbase+"/members", map[string]any{"email": "nobody@habarchy.tm", "role": "viewer"}, owner), 404)
	e.expect(e.do("PUT", pbase+"/members", map[string]any{"email": "dev@habarchy.tm", "role": "owner"}, dev), 403)
	members := decode[[]projects.Member](t, e.expect(e.do("GET", pbase+"/members", nil, dev), 200).Data)
	if len(members) != 2 {
		t.Fatalf("members: %+v", members)
	}
	e.expect(e.do("POST", pbase+"/api-keys", map[string]any{"name": "x"}, dev), 403)
	// The last owner cannot be removed.
	e.expect(e.do("DELETE", pbase+"/members/"+me.ID.String(), nil, owner), 409)

	// Update settings, including an encrypted webhook secret.
	upd := decode[admin.ProjectResponse](t, e.expect(e.do("PATCH", pbase, map[string]any{
		"webhook_url": "https://tds.gov.tm/hooks/habarchy", "webhook_secret": "whsec_123", "allowed_ips": []string{"10.0.0.0/8", "0.0.0.0/0"},
		"auto_channel_order": []string{"sms", "push"},
	}, owner), 200).Data)
	if !upd.HasWebhookSecret || upd.WebhookURL == "" || len(upd.AutoChannelOrder) != 2 {
		t.Fatalf("update: %+v", upd)
	}
	e.expect(e.do("PATCH", pbase, map[string]any{"allowed_ips": []string{"nope"}}, owner), 400)
	e.expect(e.do("PATCH", pbase, map[string]any{"status": "bogus"}, owner), 400)

	// ---- API keys ----
	created := decode[admin.APIKeyResponse](t, e.expect(e.do("POST", pbase+"/api-keys", map[string]any{
		"name": "laravel", "scopes": []string{"templates", "messages:send"},
	}, owner), 201).Data)
	if !strings.HasPrefix(created.Key, "hb_test_") || created.Hint != created.Key[len(created.Key)-4:] {
		t.Fatalf("api key: %+v", created)
	}
	keys := decode[[]admin.APIKeyResponse](t, e.expect(e.do("GET", pbase+"/api-keys", nil, dev), 200).Data)
	if len(keys) != 1 || keys[0].Key != "" {
		t.Fatalf("list must not expose keys: %+v", keys)
	}
	noScope := decode[admin.APIKeyResponse](t, e.expect(e.do("POST", pbase+"/api-keys", map[string]any{"name": "ro", "scopes": []string{"usage"}}, owner), 201).Data)
	signed := decode[admin.APIKeyResponse](t, e.expect(e.do("POST", pbase+"/api-keys", map[string]any{"name": "signed", "require_signature": true}, owner), 201).Data)
	e.expect(e.do("POST", pbase+"/api-keys", map[string]any{"name": "bad", "scopes": []string{"root"}}, owner), 400)

	// ---- public API ----
	apiKey := map[string]string{"X-Api-Key": created.Key}
	e.expect(e.do("GET", "/api/v1/me", nil, nil), 401)
	e.expect(e.do("GET", "/api/v1/me", nil, map[string]string{"X-Api-Key": "hb_test_nope"}), 401)
	e.expect(e.do("GET", "/api/v1/me", nil, apiKey), 200)
	r = e.expect(e.do("GET", "/api/v1/templates", nil, map[string]string{"X-Api-Key": noScope.Key}), 403)
	if r.Error.Code != "forbidden_scope" {
		t.Fatalf("scope error: %s", r.Raw)
	}

	// Signature required: missing -> 401, valid -> 200, wrong body -> 401.
	e.expect(e.do("GET", "/api/v1/me", nil, map[string]string{"X-Api-Key": signed.Key}), 401)
	ts := fmt.Sprint(time.Now().Unix())
	body := []byte(`{"key":"otp","channel":"sms","body":"Kod: {{.code}}"}`)
	sigHeaders := map[string]string{
		"X-Api-Key": signed.Key, "X-Timestamp": ts,
		"X-Signature": projects.Sign(signed.Key, ts, "POST", "/api/v1/templates", body),
	}
	tpl := decode[admin.TemplateResponse](t, e.expect(e.do("POST", "/api/v1/templates", body, sigHeaders), 201).Data)
	if tpl.Version != 1 || len(tpl.RequiredVars) != 1 || tpl.RequiredVars[0] != "code" || tpl.Locale != "tk" {
		t.Fatalf("template: %+v", tpl)
	}
	e.expect(e.do("POST", "/api/v1/templates", []byte(`{"key":"otp2","channel":"sms","body":"x"}`), sigHeaders), 401)

	// Duplicate (key, channel, locale) -> 409; other locale -> 201.
	e.expect(e.do("POST", "/api/v1/templates", map[string]any{"key": "otp", "channel": "sms", "body": "Код: {{.code}}"}, apiKey), 409)
	e.expect(e.do("POST", "/api/v1/templates", map[string]any{"key": "otp", "channel": "sms", "locale": "ru", "body": "Код: {{.code}}"}, apiKey), 201)
	e.expect(e.do("POST", "/api/v1/templates", map[string]any{"key": "otp", "channel": "email", "body": "no subject"}, apiKey), 400)
	e.expect(e.do("POST", "/api/v1/templates", map[string]any{"key": "otp", "channel": "push", "body": "{{.broken"}, apiKey), 400)

	// Preview with and without data.
	prev := decode[templates.Preview](t, e.expect(e.do("POST", "/api/v1/templates/preview", map[string]any{"body": "Salam {{.name}}, kod: {{.code}}", "data": map[string]any{"code": "4821"}}, apiKey), 200).Data)
	if prev.Body != "Salam {{name}}, kod: 4821" || len(prev.MissingVars) != 1 {
		t.Fatalf("preview: %+v", prev)
	}

	// ---- admin templates: update bumps version, history, restore ----
	tbase := pbase + "/templates/" + tpl.ID.String()
	v2 := decode[admin.TemplateResponse](t, e.expect(e.do("PUT", tbase, map[string]any{"body": "Siziň koduňyz: {{.code}}. {{.minutes}} min."}, dev), 200).Data)
	if v2.Version != 2 || len(v2.RequiredVars) != 2 {
		t.Fatalf("v2: %+v", v2)
	}
	versions := decode[[]admin.TemplateVersionResponse](t, e.expect(e.do("GET", tbase+"/versions", nil, dev), 200).Data)
	if len(versions) != 2 || versions[0].Version != 2 || versions[1].Body != "Kod: {{.code}}" {
		t.Fatalf("versions: %+v", versions)
	}
	v3 := decode[admin.TemplateResponse](t, e.expect(e.do("POST", tbase+"/versions/1/restore", nil, dev), 200).Data)
	if v3.Version != 3 || v3.Body != "Kod: {{.code}}" {
		t.Fatalf("restore: %+v", v3)
	}
	e.expect(e.do("POST", tbase+"/versions/99/restore", nil, dev), 404)
	sp := decode[templates.Preview](t, e.expect(e.do("POST", tbase+"/preview", map[string]any{"data": map[string]any{"code": "1111"}}, dev), 200).Data)
	if sp.Body != "Kod: 1111" {
		t.Fatalf("stored preview: %+v", sp)
	}
	all := decode[[]admin.TemplateResponse](t, e.expect(e.do("GET", pbase+"/templates?key=otp", nil, dev), 200).Data)
	if len(all) != 2 {
		t.Fatalf("list by key: %d", len(all))
	}

	// ---- revoke key -> public API rejects it ----
	e.expect(e.do("DELETE", pbase+"/api-keys/"+created.ID.String(), nil, owner), 204)
	e.expect(e.do("DELETE", pbase+"/api-keys/"+created.ID.String(), nil, owner), 404)
	e.expect(e.do("GET", "/api/v1/me", nil, apiKey), 401)

	// ---- suspended project blocks API keys ----
	e.expect(e.do("PATCH", pbase, map[string]any{"status": "suspended"}, owner), 200)
	e.expect(e.do("GET", "/api/v1/me", nil, map[string]string{"X-Api-Key": noScope.Key}), 401)

	// ---- password change revokes sessions; TOTP flow ----
	e.expect(e.do("PUT", "/api/admin/me/password", map[string]any{"current_password": "nope", "new_password": "new-pass-12345"}, owner), 401)
	e.expect(e.do("PUT", "/api/admin/me/password", map[string]any{"current_password": "owner-pass-123", "new_password": "new-pass-12345"}, owner), 204)
	login = decode[loginOut](t, e.expect(e.do("POST", "/api/admin/auth/login", map[string]any{"email": "owner@habarchy.tm", "password": "new-pass-12345"}, nil), 200).Data)
	owner = bearer(login.Tokens.AccessToken)

	setup := decode[auth.TOTPSetup](t, e.expect(e.do("POST", "/api/admin/me/totp/setup", nil, owner), 200).Data)
	if setup.Secret == "" || !strings.HasPrefix(setup.URL, "otpauth://totp/") {
		t.Fatalf("totp setup: %+v", setup)
	}
	e.expect(e.do("POST", "/api/admin/me/totp/confirm", map[string]any{"code": "000000"}, owner), 401)
	code, _ := totp.GenerateCode(setup.Secret, e.now)
	e.expect(e.do("POST", "/api/admin/me/totp/confirm", map[string]any{"code": code}, owner), 204)
	// Login now needs the code.
	r = e.expect(e.do("POST", "/api/admin/auth/login", map[string]any{"email": "owner@habarchy.tm", "password": "new-pass-12345"}, nil), 401)
	if r.Error.Details["totp_required"] != true {
		t.Fatalf("expected totp_required: %s", r.Raw)
	}
	code, _ = totp.GenerateCode(setup.Secret, e.now)
	e.expect(e.do("POST", "/api/admin/auth/login", map[string]any{"email": "owner@habarchy.tm", "password": "new-pass-12345", "totp_code": code}, nil), 200)
	e.expect(e.do("POST", "/api/admin/me/totp/disable", map[string]any{"password": "new-pass-12345"}, owner), 204)
	e.expect(e.do("POST", "/api/admin/auth/login", map[string]any{"email": "owner@habarchy.tm", "password": "new-pass-12345"}, nil), 200)

	// ---- delete project: owner only ----
	e.expect(e.do("DELETE", pbase, nil, dev), 403)
	e.expect(e.do("DELETE", pbase, nil, owner), 204)
	e.expect(e.do("GET", pbase, nil, owner), 404)
	_ = http.StatusOK
}
