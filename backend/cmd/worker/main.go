// Command worker consumes the asynq queues and delivers messages and
// webhooks. Task handlers are registered as the build progresses.
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
	"github.com/Esca6585dev/habarchy/backend/internal/config"
	"github.com/Esca6585dev/habarchy/backend/internal/queue"
	"github.com/Esca6585dev/habarchy/backend/pkg/logger"
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

	srv := asynq.NewServer(queue.RedisOpt(rdb.Options()), asynq.Config{
		Concurrency: cfg.Queue.Concurrency,
		Queues:      queue.Weights(),
		Logger:      queue.AsynqLogger{Logger: log},
		LogLevel:    asynq.InfoLevel,
	})
	mux := asynq.NewServeMux()
	// Handlers (send:sms, send:email, send:push, send:telegram, webhook)
	// are registered here in step 3.

	if err := srv.Start(mux); err != nil {
		return err
	}
	log.Info().Int("concurrency", cfg.Queue.Concurrency).Msg("worker started")
	<-ctx.Done()
	log.Info().Msg("shutting down")
	srv.Shutdown()
	return nil
}
