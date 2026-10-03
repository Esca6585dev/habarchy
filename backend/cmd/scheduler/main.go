// Command scheduler runs periodic maintenance: creating message partitions
// ahead of time, aggregating usage, expiring idempotency keys and refresh
// tokens. Scheduled (delayed) messages are handled by asynq itself.
package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Esca6585dev/habarchy/backend/internal/adapters/postgres"
	"github.com/Esca6585dev/habarchy/backend/internal/config"
	"github.com/Esca6585dev/habarchy/backend/internal/scheduler"
	"github.com/Esca6585dev/habarchy/backend/pkg/logger"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, "scheduler:", err)
		os.Exit(1)
	}
}

func run() error {
	cfg, err := config.Load()
	if err != nil {
		return err
	}
	log := logger.New("scheduler", cfg.LogLevel, cfg.LogPretty)

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	db, err := postgres.Connect(ctx, cfg.Database)
	if err != nil {
		return err
	}
	defer db.Close()

	jobs := scheduler.DefaultJobs(db, cfg)
	runner := scheduler.NewRunner(log, jobs...)
	log.Info().Int("jobs", len(jobs)).Msg("scheduler started")
	runner.Run(ctx)
	log.Info().Msg("scheduler stopped")
	return nil
}

// keep time imported for future cron-style configuration.
var _ = time.Second
