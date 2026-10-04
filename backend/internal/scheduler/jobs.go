package scheduler

import (
	"context"
	"time"

	"github.com/Esca6585dev/habarchy/backend/internal/adapters/postgres"
	"github.com/Esca6585dev/habarchy/backend/internal/app/stats"
	"github.com/Esca6585dev/habarchy/backend/internal/config"
)

// DefaultJobs returns the maintenance jobs every deployment needs.
// Usage aggregation and quota handling are added in step 4.
func DefaultJobs(db *postgres.DB, cfg *config.Config) []Job {
	st := stats.New(db, nil)
	return []Job{
		{
			Name:     "aggregate_usage_daily",
			Interval: 5 * time.Minute,
			Timeout:  4 * time.Minute,
			Run: func(ctx context.Context) error {
				// Today and yesterday: late deliveries / receipts change
				// yesterday's counters for a while after midnight.
				now := time.Now().UTC()
				for _, day := range []time.Time{now, now.AddDate(0, 0, -1)} {
					if _, err := st.AggregateUsage(ctx, day); err != nil {
						return err
					}
				}
				return nil
			},
		},
		{
			Name:     "ensure_message_partitions",
			Interval: 6 * time.Hour,
			Timeout:  time.Minute,
			Run:      EnsurePartitions(db),
		},
		{
			Name:     "expire_idempotency_keys",
			Interval: time.Hour,
			Timeout:  5 * time.Minute,
			Run: func(ctx context.Context) error {
				_, err := db.Queries.DeleteExpiredIdempotencyKeys(ctx, time.Now().Add(-cfg.Security.IdempotencyTTL))
				return err
			},
		},
		{
			Name:     "expire_refresh_tokens",
			Interval: 12 * time.Hour,
			Timeout:  5 * time.Minute,
			Run: func(ctx context.Context) error {
				_, err := db.Queries.DeleteExpiredRefreshTokens(ctx)
				return err
			},
		},
	}
}

// EnsurePartitions creates the messages partitions for this month and the
// next two, so a scheduler outage of up to two months cannot route rows
// into the default partition.
func EnsurePartitions(db *postgres.DB) func(ctx context.Context) error {
	return func(ctx context.Context) error {
		now := time.Now().UTC()
		for i := 0; i < 3; i++ {
			day := time.Date(now.Year(), now.Month()+time.Month(i), 1, 0, 0, 0, 0, time.UTC)
			if _, err := db.Queries.EnsureMessagesPartition(ctx, day); err != nil {
				return err
			}
		}
		return nil
	}
}
