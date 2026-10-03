package scheduler

import (
	"context"
	"errors"
	"sync/atomic"
	"testing"
	"time"

	"github.com/rs/zerolog"
)

func TestRunnerRunsImmediatelyAndRepeats(t *testing.T) {
	var runs atomic.Int32
	r := NewRunner(zerolog.Nop(), Job{
		Name:     "tick",
		Interval: 10 * time.Millisecond,
		Run: func(context.Context) error {
			runs.Add(1)
			return nil
		},
	}, Job{
		Name:     "fails",
		Interval: 10 * time.Millisecond,
		Run:      func(context.Context) error { return errors.New("boom") },
	}, Job{
		Name:     "panics",
		Interval: 10 * time.Millisecond,
		Run:      func(context.Context) error { panic("oops") },
	})
	ctx, cancel := context.WithTimeout(context.Background(), 75*time.Millisecond)
	defer cancel()
	r.Run(ctx)
	if n := runs.Load(); n < 3 {
		t.Fatalf("expected at least 3 runs, got %d", n)
	}
}
