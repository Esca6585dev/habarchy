package seed_test

import (
	"bytes"
	"context"
	"strings"
	"testing"

	"github.com/Esca6585dev/habarchy/backend/internal/adapters/postgres/pgtest"
	"github.com/Esca6585dev/habarchy/backend/internal/app/auth"
	"github.com/Esca6585dev/habarchy/backend/internal/app/projects"
	"github.com/Esca6585dev/habarchy/backend/internal/app/seed"
	"github.com/Esca6585dev/habarchy/backend/internal/app/templates"
	"github.com/Esca6585dev/habarchy/backend/internal/config"
	"github.com/Esca6585dev/habarchy/backend/pkg/crypto"
	"github.com/Esca6585dev/habarchy/backend/pkg/password"
)

func TestSeedFlowIsIdempotent(t *testing.T) {
	db := pgtest.Open(t)
	t.Setenv("HABARCHY_DATABASE_URL", "postgres://unused")
	t.Setenv("HABARCHY_JWT_SECRET", strings.Repeat("j", 32))
	t.Setenv("HABARCHY_MASTER_KEY", strings.Repeat("ab", 32))
	cfg, err := config.Load()
	if err != nil {
		t.Fatal(err)
	}
	cipher, _ := crypto.NewCipherFromString(cfg.Security.MasterKey)
	fast := password.Params{Memory: 8 * 1024, Iterations: 1, Parallelism: 1, SaltLength: 16, KeyLength: 32}
	authSvc := auth.New(db, cfg.Auth, cipher, auth.WithPasswordHasher(func(p string) (string, error) { return password.HashWithParams(p, fast) }))
	s := &seed.Seeder{DB: db, Auth: authSvc, Projects: projects.New(db, cipher), Templates: templates.New(db)}
	ctx := context.Background()

	opt := seed.Options{AdminEmail: "seed@example.com", AdminPassword: "seed-password-1", ProjectSlug: "seed-demo", LiveKey: true}
	r1, err := s.Run(ctx, opt)
	if err != nil {
		t.Fatal(err)
	}
	if !r1.AdminCreated || !r1.ProjectCreated || r1.TemplatesCreated != r1.TemplatesTotal || r1.TemplatesTotal != 18 {
		t.Fatalf("first run: %+v", r1)
	}
	if !strings.HasPrefix(r1.TestKey, "hb_test_") || !strings.HasPrefix(r1.LiveKey, "hb_live_") {
		t.Fatalf("keys: %q %q", r1.TestKey, r1.LiveKey)
	}

	// Second run reuses everything except keys.
	r2, err := s.Run(ctx, seed.Options{AdminEmail: "seed@example.com", ProjectSlug: "seed-demo"})
	if err != nil {
		t.Fatal(err)
	}
	if r2.AdminCreated || r2.ProjectCreated || r2.TemplatesCreated != 0 || r2.Project.ID != r1.Project.ID || r2.TestKey == r1.TestKey || r2.LiveKey != "" {
		t.Fatalf("second run: %+v", r2)
	}

	// The demo key authenticates and the templates render.
	caller, err := s.Projects.AuthenticateAPIKey(ctx, r2.TestKey, "127.0.0.1", projects.SignatureInput{})
	if err != nil || caller.Project.ID != r1.Project.ID {
		t.Fatalf("authenticate: %v", err)
	}
	var buf bytes.Buffer
	r2.Print(&buf)
	if !strings.Contains(buf.String(), r2.TestKey) || strings.Contains(buf.String(), "hb_live_") {
		t.Fatalf("print: %s", buf.String())
	}
}

func TestSampleTemplatesCoverAllLocales(t *testing.T) {
	seen := map[string]bool{}
	for _, in := range seed.SampleTemplates() {
		k := in.Key + "/" + string(in.Channel) + "/" + string(in.Locale)
		if seen[k] {
			t.Fatalf("duplicate %s", k)
		}
		seen[k] = true
		if in.Body == "" || (in.Channel == "email" && in.Subject == "") {
			t.Fatalf("incomplete %s", k)
		}
	}
	if len(seen) != 18 {
		t.Fatalf("got %d templates", len(seen))
	}
}
