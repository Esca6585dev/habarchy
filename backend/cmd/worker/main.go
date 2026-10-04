// Command worker consumes the asynq queues: it delivers messages through
// the configured providers (with fallback and retries) and posts webhooks.
package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/hibiken/asynq"

	"github.com/Esca6585dev/habarchy/backend/internal/adapters/postgres"
	"github.com/Esca6585dev/habarchy/backend/internal/adapters/redis"
	"github.com/Esca6585dev/habarchy/backend/internal/app/contacts"
	"github.com/Esca6585dev/habarchy/backend/internal/app/delivery"
	"github.com/Esca6585dev/habarchy/backend/internal/app/events"
	"github.com/Esca6585dev/habarchy/backend/internal/app/providers"
	"github.com/Esca6585dev/habarchy/backend/internal/app/webhooks"
	"github.com/Esca6585dev/habarchy/backend/internal/config"
	"github.com/Esca6585dev/habarchy/backend/internal/queue"
	"github.com/Esca6585dev/habarchy/backend/internal/worker"
	"github.com/Esca6585dev/habarchy/backend/pkg/crypto"
	"github.com/Esca6585dev/habarchy/backend/pkg/logger"
	"github.com/Esca6585dev/habarchy/backend/pkg/metrics"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, "worker:", err)
		os.Exit(1)
	}
}

func run() error {
	cfg, err := config.Load()
	if err != nil {
		return err
	}
	log := logger.New("worker", cfg.LogLevel, cfg.LogPretty)

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
	deliverySvc := delivery.New(db, rdb, providerSvc, webhookSvc, contactSvc, log, cfg.Queue.MaxRetry)
	deliverySvc.Events = events.NewPublisher(rdb.Raw())
	if cfg.Telemetry.MetricsEnabled && cfg.Telemetry.WorkerMetricsAddr != "" {
		msrv := metrics.Serve(cfg.Telemetry.WorkerMetricsAddr)
		defer func() { _ = msrv.Close() }()
		log.Info().Str("addr", cfg.Telemetry.WorkerMetricsAddr).Msg("worker metrics listening")
	}

	// SMPP sessions live only in the worker; they report receipts straight
	// into the delivery service.
	smppPool := worker.NewSMPPPool(log, deliverySvc)
	defer smppPool.Close()
	providerSvc.SMPPFactory = smppPool.Factory

	srv := asynq.NewServer(queue.RedisOpt(rdb.Options()), asynq.Config{
		Concurrency:    cfg.Queue.Concurrency,
		Queues:         queue.Weights(),
		Logger:         queue.AsynqLogger{Logger: log},
		LogLevel:       asynq.InfoLevel,
		RetryDelayFunc: queue.RetryDelay,
		IsFailure:      worker.IsFailure,
	})
	mux := worker.NewMux(deliverySvc, webhookSvc)

	if err := srv.Start(mux); err != nil {
		return err
	}
	log.Info().Int("concurrency", cfg.Queue.Concurrency).Msg("worker started")
	<-ctx.Done()
	log.Info().Msg("shutting down")
	srv.Shutdown()
	return nil
}
