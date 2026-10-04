// Package seed installs demo data: an admin user (optional), a "demo"
// project, sample templates in tk/ru/en and fresh API keys. It is
// idempotent — running it twice reuses the project and templates and only
// issues new keys.
package seed

import (
	"context"
	"errors"
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"

	"github.com/Esca6585dev/habarchy/backend/internal/adapters/postgres"
	"github.com/Esca6585dev/habarchy/backend/internal/adapters/postgres/sqlcgen"
	"github.com/Esca6585dev/habarchy/backend/internal/app/auth"
	"github.com/Esca6585dev/habarchy/backend/internal/app/projects"
	"github.com/Esca6585dev/habarchy/backend/internal/app/templates"
	"github.com/Esca6585dev/habarchy/backend/internal/domain"
	"github.com/Esca6585dev/habarchy/backend/pkg/ids"
)

// Options control what the seeder creates.
type Options struct {
	// AdminEmail / AdminPassword create the admin when it does not exist yet.
	// With an empty email the most recently created user owns the project.
	AdminEmail    string
	AdminPassword string
	AdminName     string
	ProjectSlug   string // default "demo"
	ProjectName   string // default "Demo"
	// LiveKey also issues an hb_live_ key (needs real providers to deliver).
	LiveKey bool
}

// Result is what the caller prints. Plaintext keys are shown exactly once.
type Result struct {
	Admin            sqlcgen.User
	AdminCreated     bool
	Project          sqlcgen.Project
	ProjectCreated   bool
	TemplatesCreated int
	TemplatesTotal   int
	TestKey          string
	LiveKey          string
}

// Seeder wires the application services.
type Seeder struct {
	DB        *postgres.DB
	Auth      *auth.Service
	Projects  *projects.Service
	Templates *templates.Service
}

// Run performs the seeding.
func (s *Seeder) Run(ctx context.Context, opt Options) (*Result, error) {
	if opt.ProjectSlug == "" {
		opt.ProjectSlug = "demo"
	}
	if opt.ProjectName == "" {
		opt.ProjectName = "Demo"
	}
	res := &Result{}

	admin, created, err := s.admin(ctx, opt)
	if err != nil {
		return nil, err
	}
	res.Admin, res.AdminCreated = *admin, created

	project, created, err := s.project(ctx, admin.ID, opt)
	if err != nil {
		return nil, err
	}
	res.Project, res.ProjectCreated = *project, created

	existing, err := s.Templates.List(ctx, project.ID, "")
	if err != nil {
		return nil, err
	}
	have := map[string]bool{}
	for i := range existing {
		t := &existing[i]
		have[t.Key+"/"+string(t.Channel)+"/"+t.Locale] = true
	}
	for _, t := range SampleTemplates() {
		res.TemplatesTotal++
		if have[t.Key+"/"+string(t.Channel)+"/"+string(t.Locale)] {
			continue
		}
		if _, err := s.Templates.Create(ctx, project.ID, t); err != nil {
			return nil, fmt.Errorf("template %s/%s/%s: %w", t.Key, t.Channel, t.Locale, err)
		}
		res.TemplatesCreated++
	}

	for _, g := range []string{"Işdeşler", "Müşderiler"} {
		if _, err := s.DB.Queries.CreateGroup(ctx, sqlcgen.CreateGroupParams{ID: ids.New(), ProjectID: project.ID, Name: g}); err != nil && !postgres.IsUniqueViolation(err) {
			return nil, err
		}
	}

	stamp := time.Now().UTC().Format("2006-01-02 15:04")
	testKey, err := s.Projects.CreateAPIKey(ctx, admin.ID, project.ID, projects.APIKeyInput{Name: "seed test " + stamp, Live: false, Scopes: domain.AllScopes})
	if err != nil {
		return nil, err
	}
	res.TestKey = testKey.Plaintext
	if opt.LiveKey {
		liveKey, err := s.Projects.CreateAPIKey(ctx, admin.ID, project.ID, projects.APIKeyInput{Name: "seed live " + stamp, Live: true, Scopes: domain.AllScopes})
		if err != nil {
			return nil, err
		}
		res.LiveKey = liveKey.Plaintext
	}
	return res, nil
}

func (s *Seeder) admin(ctx context.Context, opt Options) (*sqlcgen.User, bool, error) {
	email := strings.ToLower(strings.TrimSpace(opt.AdminEmail))
	if email == "" {
		users, err := s.DB.Queries.ListUsers(ctx, sqlcgen.ListUsersParams{RowOffset: 0, RowLimit: 1})
		if err != nil {
			return nil, false, err
		}
		if len(users) == 0 {
			return nil, false, errors.New("no users yet: set HABARCHY_ADMIN_EMAIL and HABARCHY_ADMIN_PASSWORD")
		}
		return &users[0], false, nil
	}
	user, err := s.DB.Queries.GetUserByEmail(ctx, email)
	switch {
	case err == nil:
		return &user, false, nil
	case errors.Is(err, pgx.ErrNoRows):
		if opt.AdminPassword == "" {
			return nil, false, fmt.Errorf("user %s does not exist: set HABARCHY_ADMIN_PASSWORD to create it", email)
		}
		created, err := s.Auth.CreateUser(ctx, email, opt.AdminPassword, opt.AdminName)
		if err != nil {
			return nil, false, err
		}
		return created, true, nil
	default:
		return nil, false, err
	}
}

func (s *Seeder) project(ctx context.Context, ownerID uuid.UUID, opt Options) (*sqlcgen.Project, bool, error) {
	memberships, err := s.Projects.ListForUser(ctx, ownerID)
	if err != nil {
		return nil, false, err
	}
	for i := range memberships {
		if memberships[i].Project.Slug == opt.ProjectSlug {
			return &memberships[i].Project, false, nil
		}
	}
	p, err := s.Projects.Create(ctx, ownerID, projects.CreateInput{
		Name: opt.ProjectName, Slug: opt.ProjectSlug, DailyQuota: 1000, MonthlyQuota: 20000, DefaultLocale: domain.LocaleTK,
	})
	if err != nil {
		return nil, false, err
	}
	return p, true, nil
}

// SampleTemplates are the otp / welcome / password_reset templates in
// tk, ru and en for SMS and email (push uses the same keys with a title).
func SampleTemplates() []templates.Input {
	type tr struct{ tk, ru, en string }
	sms := map[string]tr{
		"otp": {
			tk: "{{.code}} — Habarçy tassyklama kody. {{.minutes}} minut dowam edýär. Hiç kime aýtmaň.",
			ru: "{{.code}} — код подтверждения Habarçy. Действует {{.minutes}} мин. Никому не сообщайте.",
			en: "{{.code}} is your Habarchy verification code. Valid for {{.minutes}} minutes. Do not share it.",
		},
		"welcome": {
			tk: "Hoş geldiňiz, {{.name}}! Hasabyňyz döredildi. Soraglar üçin: {{.support}}",
			ru: "Добро пожаловать, {{.name}}! Ваш аккаунт создан. Поддержка: {{.support}}",
			en: "Welcome, {{.name}}! Your account is ready. Support: {{.support}}",
		},
		"password_reset": {
			tk: "Paroly täzelemek üçin kod: {{.code}}. Siz soramadyk bolsaňyz, bu habary äsgermäň.",
			ru: "Код для сброса пароля: {{.code}}. Если это были не вы, проигнорируйте сообщение.",
			en: "Your password reset code is {{.code}}. If you did not request it, ignore this message.",
		},
	}
	emailSubject := map[string]tr{
		"otp":            {tk: "Tassyklama kody: {{.code}}", ru: "Код подтверждения: {{.code}}", en: "Your verification code: {{.code}}"},
		"welcome":        {tk: "Hoş geldiňiz, {{.name}}!", ru: "Добро пожаловать, {{.name}}!", en: "Welcome, {{.name}}!"},
		"password_reset": {tk: "Paroly täzelemek", ru: "Сброс пароля", en: "Reset your password"},
	}
	emailBody := map[string]tr{
		"otp": {
			tk: "<p>Salam!</p><p>Tassyklama kodyňyz: <b style=\"font-size:20px\">{{.code}}</b></p><p>Kod {{.minutes}} minut dowam edýär.</p>",
			ru: "<p>Здравствуйте!</p><p>Ваш код подтверждения: <b style=\"font-size:20px\">{{.code}}</b></p><p>Код действует {{.minutes}} мин.</p>",
			en: "<p>Hello!</p><p>Your verification code: <b style=\"font-size:20px\">{{.code}}</b></p><p>It is valid for {{.minutes}} minutes.</p>",
		},
		"welcome": {
			tk: "<p>Hoş geldiňiz, <b>{{.name}}</b>!</p><p>Hasabyňyz üstünlikli döredildi.</p><p>Soraglaryňyz bolsa: {{.support}}</p>",
			ru: "<p>Добро пожаловать, <b>{{.name}}</b>!</p><p>Ваш аккаунт успешно создан.</p><p>Поддержка: {{.support}}</p>",
			en: "<p>Welcome, <b>{{.name}}</b>!</p><p>Your account has been created.</p><p>Questions? {{.support}}</p>",
		},
		"password_reset": {
			tk: "<p>Paroly täzelemek üçin şu salga giriň:</p><p><a href=\"{{.link}}\">{{.link}}</a></p><p>Siz soramadyk bolsaňyz, bu haty äsgermäň.</p>",
			ru: "<p>Чтобы сбросить пароль, перейдите по ссылке:</p><p><a href=\"{{.link}}\">{{.link}}</a></p><p>Если это были не вы, проигнорируйте письмо.</p>",
			en: "<p>To reset your password open this link:</p><p><a href=\"{{.link}}\">{{.link}}</a></p><p>If you did not request it, ignore this email.</p>",
		},
	}
	locales := []struct {
		l    domain.Locale
		pick func(tr) string
	}{
		{domain.LocaleTK, func(t tr) string { return t.tk }},
		{domain.LocaleRU, func(t tr) string { return t.ru }},
		{domain.LocaleEN, func(t tr) string { return t.en }},
	}
	keys := []string{"otp", "welcome", "password_reset"}
	var out []templates.Input
	for _, k := range keys {
		for _, lc := range locales {
			out = append(out, templates.Input{Key: k, Channel: domain.ChannelSMS, Locale: lc.l, Body: lc.pick(sms[k])})
			out = append(out, templates.Input{Key: k, Channel: domain.ChannelEmail, Locale: lc.l, Subject: lc.pick(emailSubject[k]), Body: lc.pick(emailBody[k])})
		}
	}
	return out
}

// Print writes a human-readable summary (the only place keys are shown).
func (r *Result) Print(w io.Writer) {
	verb := func(created bool) string {
		if created {
			return "created"
		}
		return "exists"
	}
	var b strings.Builder
	fmt.Fprintf(&b, "admin user   %s (%s)\n", r.Admin.Email, verb(r.AdminCreated))
	fmt.Fprintf(&b, "project      %s  slug=%s  id=%s (%s)\n", r.Project.Name, r.Project.Slug, r.Project.ID, verb(r.ProjectCreated))
	fmt.Fprintf(&b, "templates    %d/%d created (otp, welcome, password_reset × sms,email × tk,ru,en)\n", r.TemplatesCreated, r.TemplatesTotal)
	fmt.Fprintf(&b, "test api key %s   (sandbox: messages succeed instantly, nothing is sent)\n", r.TestKey)
	if r.LiveKey != "" {
		fmt.Fprintf(&b, "live api key %s   (add providers in the admin panel first)\n", r.LiveKey)
	}
	b.WriteString("\n")
	fmt.Fprintf(&b, "try it:\n  curl -X POST $HABARCHY/api/v1/messages -H 'X-Api-Key: %s' -H 'Content-Type: application/json' \\\n", r.TestKey)
	b.WriteString(`    -d '{"channel":"sms","to":"+99365123456","template":"otp","data":{"code":"4821","minutes":5}}'` + "\n")
	_, _ = io.WriteString(w, b.String())
}
