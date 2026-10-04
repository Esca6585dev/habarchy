// Package stats builds the dashboard, usage and health read models.
package stats

import (
	"context"
	"time"

	"github.com/google/uuid"

	"github.com/Esca6585dev/habarchy/backend/internal/adapters/postgres"
	"github.com/Esca6585dev/habarchy/backend/internal/adapters/postgres/sqlcgen"
	"github.com/Esca6585dev/habarchy/backend/internal/domain"
	"github.com/Esca6585dev/habarchy/backend/internal/queue"
)

// Service holds stats use cases.
type Service struct {
	db    *postgres.DB
	queue *queue.Client // optional: queue depth for health
}

// New creates the service.
func New(db *postgres.DB, q *queue.Client) *Service { return &Service{db: db, queue: q} }

// Window is a half-open [From, To) time range.
type Window struct {
	From time.Time
	To   time.Time
}

// LastDays returns a window covering the last n days up to now (UTC days).
func LastDays(n int) Window {
	now := time.Now().UTC()
	to := now.Truncate(24 * time.Hour).Add(24 * time.Hour)
	return Window{From: to.AddDate(0, 0, -n), To: to}
}

// Dashboard is the admin home page payload.
type Dashboard struct {
	Window   Window            `json:"window"`
	Totals   Totals            `json:"totals"`
	Latency  Latency           `json:"latency"`
	Daily    []DailyPoint      `json:"daily"`
	Failures []sqlcgen.Message `json:"recent_failures"`
	Channels map[string]Totals `json:"by_channel"`
}

// Totals are outcome counters.
type Totals struct {
	Total      int64 `json:"total"`
	Sent       int64 `json:"sent"`
	Delivered  int64 `json:"delivered"`
	Failed     int64 `json:"failed"`
	Pending    int64 `json:"pending"`
	CostMicros int64 `json:"cost_micros"`
}

// Latency is creation-to-sent / delivered percentiles in seconds.
type Latency struct {
	P50SentSec      float64 `json:"p50_sent_sec"`
	P95SentSec      float64 `json:"p95_sent_sec"`
	P95DeliveredSec float64 `json:"p95_delivered_sec"`
	Samples         int64   `json:"samples"`
}

// DailyPoint is one bar of the chart.
type DailyPoint struct {
	Day     string `json:"day"`
	Channel string `json:"channel"`
	Totals
}

// Dashboard aggregates everything the home page shows.
func (s *Service) Dashboard(ctx context.Context, projectID uuid.UUID, w Window) (*Dashboard, error) {
	q := s.db.Queries
	tot, err := q.DashboardTotals(ctx, sqlcgen.DashboardTotalsParams{ProjectID: projectID, FromTs: w.From, ToTs: w.To})
	if err != nil {
		return nil, err
	}
	lat, err := q.DeliveryLatency(ctx, sqlcgen.DeliveryLatencyParams{ProjectID: projectID, FromTs: w.From, ToTs: w.To})
	if err != nil {
		return nil, err
	}
	daily, err := q.DashboardDaily(ctx, sqlcgen.DashboardDailyParams{ProjectID: projectID, FromTs: w.From, ToTs: w.To})
	if err != nil {
		return nil, err
	}
	failures, err := q.RecentFailures(ctx, sqlcgen.RecentFailuresParams{ProjectID: projectID, FromTs: w.From, RowLimit: 10})
	if err != nil {
		return nil, err
	}
	d := &Dashboard{
		Window:   w,
		Totals:   Totals{Total: tot.Total, Sent: tot.Sent, Delivered: tot.Delivered, Failed: tot.Failed, Pending: tot.Pending, CostMicros: tot.CostMicros},
		Latency:  Latency{P50SentSec: lat.P50SentSec, P95SentSec: lat.P95SentSec, P95DeliveredSec: lat.P95DeliveredSec, Samples: lat.Samples},
		Daily:    make([]DailyPoint, 0, len(daily)),
		Failures: failures,
		Channels: map[string]Totals{},
	}
	for _, row := range daily {
		t := Totals{Total: row.Total, Sent: row.Sent, Delivered: row.Delivered, Failed: row.Failed, Pending: row.Pending, CostMicros: row.CostMicros}
		d.Daily = append(d.Daily, DailyPoint{Day: row.Day.Format("2006-01-02"), Channel: string(row.Channel), Totals: t})
		c := d.Channels[string(row.Channel)]
		c.Total += t.Total
		c.Sent += t.Sent
		c.Delivered += t.Delivered
		c.Failed += t.Failed
		c.Pending += t.Pending
		c.CostMicros += t.CostMicros
		d.Channels[string(row.Channel)] = c
	}
	return d, nil
}

// UsageRow is one line of the usage report.
type UsageRow struct {
	Day        string `json:"day,omitempty"`
	Channel    string `json:"channel,omitempty"`
	Queued     int64  `json:"queued"`
	Sent       int64  `json:"sent"`
	Delivered  int64  `json:"delivered"`
	Failed     int64  `json:"failed"`
	CostMicros int64  `json:"cost_micros"`
	Currency   string `json:"currency"`
}

// Usage reads usage_daily for [from, to] (dates) grouped by day, channel or
// both. Today's row is aggregated by the scheduler every few minutes, so
// the report may lag by that interval.
func (s *Service) Usage(ctx context.Context, projectID uuid.UUID, from, to time.Time, groupBy string) ([]UsageRow, error) {
	rows, err := s.db.Queries.ListUsageDaily(ctx, sqlcgen.ListUsageDailyParams{ProjectID: projectID, FromDay: from, ToDay: to})
	if err != nil {
		return nil, err
	}
	acc := map[string]*UsageRow{}
	order := []string{}
	for _, r := range rows {
		var key string
		row := UsageRow{Currency: r.Currency}
		switch groupBy {
		case "channel":
			key = string(r.Channel)
			row.Channel = key
		case "day":
			key = r.Day.Format("2006-01-02")
			row.Day = key
		default:
			key = r.Day.Format("2006-01-02") + "|" + string(r.Channel)
			row.Day, row.Channel = r.Day.Format("2006-01-02"), string(r.Channel)
		}
		cur, ok := acc[key]
		if !ok {
			cur = &row
			acc[key] = cur
			order = append(order, key)
		}
		cur.Queued += r.Queued
		cur.Sent += r.Sent
		cur.Delivered += r.Delivered
		cur.Failed += r.Failed
		cur.CostMicros += r.CostMicros
	}
	out := make([]UsageRow, 0, len(order))
	for _, k := range order {
		out = append(out, *acc[k])
	}
	return out, nil
}

// AggregateUsage recomputes usage_daily for one UTC day from messages.
func (s *Service) AggregateUsage(ctx context.Context, day time.Time) (int, error) {
	start := time.Date(day.Year(), day.Month(), day.Day(), 0, 0, 0, 0, time.UTC)
	rows, err := s.db.Queries.AggregateUsageForDay(ctx, sqlcgen.AggregateUsageForDayParams{DayStart: start, DayEnd: start.Add(24 * time.Hour)})
	if err != nil {
		return 0, err
	}
	for _, r := range rows {
		if err := s.db.Queries.UpsertUsageDaily(ctx, sqlcgen.UpsertUsageDailyParams{
			ProjectID: r.ProjectID, Day: start, Channel: r.Channel, Queued: r.Queued, Sent: r.Sent, Delivered: r.Delivered,
			Failed: r.Failed, CostMicros: r.CostMicros, Currency: r.Currency,
		}); err != nil {
			return 0, err
		}
	}
	return len(rows), nil
}

// ProviderHealth is one row of the health page.
type ProviderHealth struct {
	ID         uuid.UUID  `json:"id"`
	Name       string     `json:"name"`
	Channel    string     `json:"channel"`
	Type       string     `json:"type"`
	IsActive   bool       `json:"is_active"`
	Priority   int32      `json:"priority"`
	OkCount    int64      `json:"ok_count"`
	FailCount  int64      `json:"failed_count"`
	LastSentAt *time.Time `json:"last_sent_at"`
	Status     string     `json:"status"` // healthy | degraded | failing | idle | disabled
}

// Health is the system health payload for one project.
type Health struct {
	Queues    []queue.QueueStats `json:"queues"`
	Providers []ProviderHealth   `json:"providers"`
	Devices   struct {
		Active int64 `json:"active"`
		Total  int64 `json:"total"`
	} `json:"devices"`
	Contacts int64 `json:"contacts"`
}

// Health gathers queue depth, provider outcomes (last 24 h) and counts.
func (s *Service) Health(ctx context.Context, projectID uuid.UUID) (*Health, error) {
	h := &Health{Queues: []queue.QueueStats{}}
	if s.queue != nil {
		if qs, err := s.queue.Stats(); err == nil {
			h.Queues = qs
		}
	}
	rows, err := s.db.Queries.ProviderStats(ctx, sqlcgen.ProviderStatsParams{ProjectID: projectID, FromTs: time.Now().Add(-24 * time.Hour)})
	if err != nil {
		return nil, err
	}
	h.Providers = make([]ProviderHealth, 0, len(rows))
	for _, r := range rows {
		p := ProviderHealth{ID: r.ID, Name: r.Name, Channel: string(r.Channel), Type: string(r.Type), IsActive: r.IsActive, Priority: r.Priority,
			OkCount: r.OkCount, FailCount: r.FailedCount}
		if !r.LastSentAt.IsZero() {
			ts := r.LastSentAt
			p.LastSentAt = &ts
		}
		switch {
		case !r.IsActive:
			p.Status = "disabled"
		case r.OkCount == 0 && r.FailedCount == 0:
			p.Status = "idle"
		case r.FailedCount == 0:
			p.Status = "healthy"
		case r.FailedCount*5 < r.OkCount:
			p.Status = "degraded"
		default:
			p.Status = "failing"
		}
		h.Providers = append(h.Providers, p)
	}
	dev, err := s.db.Queries.CountDevices(ctx, projectID)
	if err != nil {
		return nil, err
	}
	h.Devices.Active, h.Devices.Total = dev.Active, dev.Total
	if h.Contacts, err = s.db.Queries.CountContacts(ctx, projectID); err != nil {
		return nil, err
	}
	return h, nil
}

// ChannelsOf lists the channels present in a dashboard (for charts).
func ChannelsOf(d *Dashboard) []domain.Channel {
	out := make([]domain.Channel, 0, len(d.Channels))
	for _, c := range domain.AllChannels {
		if _, ok := d.Channels[string(c)]; ok {
			out = append(out, c)
		}
	}
	return out
}
