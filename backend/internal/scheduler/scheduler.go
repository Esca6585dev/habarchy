// Package scheduler is a small interval-based job runner. Jobs are plain
// functions; each runs on its own ticker, immediately on start, and never
// overlaps with itself.
package scheduler

import (
	"context"
	"sync"
	"time"

	"github.com/rs/zerolog"
)

// Job is a named periodic task.
type Job struct {
	Name     string
	Interval time.Duration
	// Timeout bounds a single run; zero means Interval.
	Timeout time.Duration
	Run     func(ctx context.Context) error
}

// Runner executes jobs until the context is cancelled.
type Runner struct {
	log  zerolog.Logger
	jobs []Job
}

// NewRunner creates a runner for jobs.
func NewRunner(log zerolog.Logger, jobs ...Job) *Runner {
	return &Runner{log: log, jobs: jobs}
}

// Run blocks until ctx is done and all jobs have returned.
func (r *Runner) Run(ctx context.Context) {
	var wg sync.WaitGroup
	for _, job := range r.jobs {
		wg.Add(1)
		go func(j Job) {
			defer wg.Done()
			r.loop(ctx, j)
		}(job)
	}
	wg.Wait()
}

func (r *Runner) loop(ctx context.Context, j Job) {
	ticker := time.NewTicker(j.Interval)
	defer ticker.Stop()
	for {
		r.runOnce(ctx, j)
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
		}
	}
}

func (r *Runner) runOnce(ctx context.Context, j Job) {
	timeout := j.Timeout
	if timeout <= 0 {
		timeout = j.Interval
	}
	runCtx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	start := time.Now()
	log := r.log.With().Str("job", j.Name).Logger()
	defer func() {
		if rec := recover(); rec != nil {
			log.Error().Interface("panic", rec).Msg("job panicked")
		}
	}()
	if err := j.Run(runCtx); err != nil {
		if ctx.Err() != nil {
			return // shutting down; not a job failure
		}
		log.Error().Err(err).Dur("took", time.Since(start)).Msg("job failed")
		return
	}
	log.Debug().Dur("took", time.Since(start)).Msg("job ok")
}
