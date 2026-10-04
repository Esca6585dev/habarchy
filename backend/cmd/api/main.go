// Command api runs the Habarchy HTTP server (public API + admin API).
package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	httpadapter "github.com/Esca6585dev/habarchy/backend/internal/adapters/http"
	"github.com/Esca6585dev/habarchy/backend/internal/adapters/http/admin"
	"github.com/Esca6585dev/habarchy/backend/internal/adapters/http/public"
	"github.com/Esca6585dev/habarchy/backend/internal/adapters/postgres"
	"github.com/Esca6585dev/habarchy/backend/internal/adapters/redis"
	"github.com/Esca6585dev/habarchy/backend/internal/app/auth"
	"github.com/Esca6585dev/habarchy/backend/internal/app/contacts"
	"github.com/Esca6585dev/habarchy/backend/internal/app/delivery"
	"github.com/Esca6585dev/habarchy/backend/internal/app/messages"
	"github.com/Esca6585dev/habarchy/backend/internal/app/otp"
	"github.com/Esca6585dev/habarchy/backend/internal/app/projects"
	"github.com/Esca6585dev/habarchy/backend/internal/app/providers"
	"github.com/Esca6585dev/habarchy/backend/internal/app/templates"
	"github.com/Esca6585dev/habarchy/backend/internal/app/webhooks"
	"github.com/Esca6585dev/habarchy/backend/internal/config"
	"github.com/Esca6585dev/habarchy/backend/internal/queue"
	"github.com/Esca6585dev/habarchy/backend/pkg/crypto"
	"github.com/Esca6585dev/habarchy/backend/pkg/logger"
)

func main() {
	genKey := flag.Bool("genkey", false, "print a new HABARCHY_MASTER_KEY and exit")
	migrateOnly := flag.Bool("migrate", false, "apply database migrations and exit")
	showEnv := flag.Bool("env", false, "print documented environment variables and exit")
	createAdmin := flag.Bool("create-admin", false, "create an admin user from HABARCHY_ADMIN_EMAIL / HABARCHY_ADMIN_PASSWORD and exit")
	flag.Parse()

	switch {
	case *genKey:
		key, err := crypto.GenerateKey()
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		fmt.Println(key)
		return
	case *showEnv:
		_ = config.Usage()
		return
	}

	if err := run(*migrateOnly, *createAdmin); err != nil {
		fmt.Fprintln(os.Stderr, "api:", err)
		os.Exit(1)
	}
}

func run(migrateOnly, createAdmin bool) error {
	cfg, err := config.Load()
	if err != nil {
		return err
	}
	log := logger.New("api", cfg.LogLevel, cfg.LogPretty)

	cipher, err := crypto.NewCipherFromString(cfg.Security.MasterKey)
	if err != nil {
		return err
	}

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	db, err := postgres.Connect(ctx, cfg.Database)
	if err != nil {
		return err
	}
	defer db.Close()

	if migrateOnly || cfg.Database.AutoMigrate {
		if err := db.Migrate(ctx, log); err != nil {
			return err
		}
		v, _ := db.MigrationVersion(ctx)
		log.Info().Int64("version", v).Msg("database schema up to date")
		if migrateOnly {
			return nil
		}
	}

	authSvc := auth.New(db, cfg.Auth, cipher)
	projectSvc := projects.New(db, cipher)
	templateSvc := templates.New(db)

	if createAdmin {
		email, pass := os.Getenv("HABARCHY_ADMIN_EMAIL"), os.Getenv("HABARCHY_ADMIN_PASSWORD")
		if email == "" || pass == "" {
			return errors.New("set HABARCHY_ADMIN_EMAIL and HABARCHY_ADMIN_PASSWORD")
		}
		user, err := authSvc.CreateUser(ctx, email, pass, os.Getenv("HABARCHY_ADMIN_NAME"))
		if err != nil {
			return err
		}
		log.Info().Str("email", user.Email).Str("id", user.ID.String()).Msg("admin user created")
		return nil
	}

	rdb, err := redis.Connect(ctx, cfg.Redis)
	if err != nil {
		return err
	}
	defer func() { _ = rdb.Close() }()

	q := queue.NewClient(queue.RedisOpt(rdb.Options()), cfg.Queue.MaxRetry, cfg.Queue.WebhookRetry)
	defer func() { _ = q.Close() }()

	contactSvc := contacts.New(db)
	providerSvc := providers.New(db, cipher)
	webhookSvc := webhooks.New(db, cipher, q, cfg.Queue.WebhookRetry, nil)
	webhookSvc.AllowPrivate = cfg.Security.WebhookAllowPrivate
	messageSvc := messages.New(db, rdb, q, contactSvc, templateSvc, messages.Limits{
		IdempotencyTTL: cfg.Security.IdempotencyTTL, BatchMaxRecipients: cfg.Limits.BatchMaxRecipients,
	})
	otpSvc := otp.New(rdb, messageSvc, otp.Limits{
		Length: cfg.Limits.OTPLength, TTL: cfg.Limits.OTPTTL, MaxAttempts: cfg.Limits.OTPMaxAttempts,
		PerAddressHour: cfg.Limits.OTPPerAddressRate, PerIPHour: cfg.Limits.OTPPerIPRate,
	})
	// The API only applies receipts (callbacks); it never sends.
	deliverySvc := delivery.New(db, rdb, providerSvc, webhookSvc, contactSvc, log, cfg.Queue.MaxRetry)

	app := httpadapter.NewServer(cfg, log, httpadapter.Deps{DB: db, Redis: rdb})
	(&admin.Handlers{Auth: authSvc, Projects: projectSvc, Templates: templateSvc, Providers: providerSvc}).Register(app)
	(&public.Handlers{
		Projects: projectSvc, Templates: templateSvc, Messages: messageSvc, OTP: otpSvc, Contacts: contactSvc,
		Providers: providerSvc, Delivery: deliverySvc, SignatureTolerance: cfg.Security.SignatureTolerance,
		APIRatePerSec: cfg.Limits.APIRatePerSec, Limiter: rdb,
	}).Register(app)

	errCh := make(chan error, 1)
	go func() {
		log.Info().Str("addr", cfg.HTTP.Addr).Str("env", cfg.Env).Msg("http server listening")
		errCh <- app.Listen(cfg.HTTP.Addr)
	}()

	select {
	case err := <-errCh:
		return err
	case <-ctx.Done():
	}

	log.Info().Msg("shutting down")
	if err := app.ShutdownWithTimeout(cfg.HTTP.ShutdownTimeout); err != nil && !errors.Is(err, context.DeadlineExceeded) {
		return err
	}
	time.Sleep(50 * time.Millisecond)
	return nil
}
