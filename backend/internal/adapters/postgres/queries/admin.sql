-- Admin-panel read models: message log, dashboard, usage, health.

-- name: AdminListMessages :many
SELECT * FROM messages
WHERE project_id = @project_id
  AND (sqlc.narg('status')::message_status IS NULL OR status = sqlc.narg('status')::message_status)
  AND (sqlc.narg('channel')::channel IS NULL OR channel = sqlc.narg('channel')::channel)
  AND (sqlc.narg('from_ts')::timestamptz IS NULL OR created_at >= sqlc.narg('from_ts')::timestamptz)
  AND (sqlc.narg('to_ts')::timestamptz IS NULL OR created_at < sqlc.narg('to_ts')::timestamptz)
  AND (sqlc.narg('cursor_id')::uuid IS NULL OR id < sqlc.narg('cursor_id')::uuid)
  AND (sqlc.narg('search')::text IS NULL
       OR to_address ILIKE '%' || sqlc.narg('search')::text || '%'
       OR template_key ILIKE '%' || sqlc.narg('search')::text || '%'
       OR provider_message_id = sqlc.narg('search')::text)
  AND (sqlc.narg('batch_id')::uuid IS NULL OR batch_id = sqlc.narg('batch_id')::uuid)
  AND (sqlc.narg('contact_id')::uuid IS NULL OR contact_id = sqlc.narg('contact_id')::uuid)
  AND (sqlc.narg('include_test')::boolean IS NULL OR sqlc.narg('include_test')::boolean OR NOT is_test)
ORDER BY id DESC
LIMIT @row_limit;

-- name: DashboardDaily :many
-- Counts per day per channel in [from, to), by message creation time.
SELECT
    (created_at AT TIME ZONE 'UTC')::date                      AS day,
    channel,
    count(*)::bigint                                           AS total,
    count(*) FILTER (WHERE status = 'sent')::bigint            AS sent,
    count(*) FILTER (WHERE status = 'delivered')::bigint       AS delivered,
    count(*) FILTER (WHERE status = 'failed')::bigint          AS failed,
    count(*) FILTER (WHERE status IN ('queued','processing'))::bigint AS pending,
    coalesce(sum(cost_micros), 0)::bigint                      AS cost_micros
FROM messages
WHERE project_id = @project_id AND created_at >= @from_ts AND created_at < @to_ts AND NOT is_test
GROUP BY 1, 2
ORDER BY 1, 2;

-- name: DashboardTotals :one
SELECT
    count(*)::bigint                                                  AS total,
    count(*) FILTER (WHERE status = 'sent')::bigint                   AS sent,
    count(*) FILTER (WHERE status = 'delivered')::bigint              AS delivered,
    count(*) FILTER (WHERE status = 'failed')::bigint                 AS failed,
    count(*) FILTER (WHERE status IN ('queued','processing'))::bigint AS pending,
    coalesce(sum(cost_micros), 0)::bigint                             AS cost_micros
FROM messages
WHERE project_id = @project_id AND created_at >= @from_ts AND created_at < @to_ts AND NOT is_test;

-- name: DeliveryLatency :one
-- Seconds from creation to sent / delivered; percentiles over the window.
SELECT
    coalesce(percentile_cont(0.5)  WITHIN GROUP (ORDER BY extract(epoch FROM (sent_at - created_at))), 0)::float8      AS p50_sent_sec,
    coalesce(percentile_cont(0.95) WITHIN GROUP (ORDER BY extract(epoch FROM (sent_at - created_at))), 0)::float8      AS p95_sent_sec,
    coalesce(percentile_cont(0.95) WITHIN GROUP (ORDER BY extract(epoch FROM (delivered_at - created_at))), 0)::float8 AS p95_delivered_sec,
    count(*) FILTER (WHERE sent_at IS NOT NULL)::bigint AS samples
FROM messages
WHERE project_id = @project_id AND created_at >= @from_ts AND created_at < @to_ts AND NOT is_test AND sent_at IS NOT NULL;

-- name: ProviderStats :many
-- Per-provider outcome counts in the window (for the health page).
SELECT
    p.id, p.name, p.channel, p.type, p.is_active, p.priority,
    count(m.id) FILTER (WHERE m.status IN ('sent','delivered'))::bigint AS ok_count,
    count(m.id) FILTER (WHERE m.status = 'failed')::bigint              AS failed_count,
    coalesce(max(m.sent_at), '0001-01-01 00:00:00+00'::timestamptz)::timestamptz AS last_sent_at
FROM providers p
LEFT JOIN messages m ON m.provider_id = p.id AND m.created_at >= @from_ts
WHERE p.project_id = @project_id
GROUP BY p.id
ORDER BY p.channel, p.priority;

-- name: RecentFailures :many
SELECT * FROM messages
WHERE project_id = @project_id AND status = 'failed' AND created_at >= @from_ts AND NOT is_test
ORDER BY id DESC LIMIT @row_limit;

-- name: GlobalCounts :one
SELECT
    (SELECT count(*) FROM projects)::bigint                                              AS projects,
    (SELECT count(*) FROM users WHERE is_active)::bigint                                 AS users,
    (SELECT count(*) FROM messages WHERE created_at >= now() - interval '24 hours')::bigint AS messages_24h,
    (SELECT count(*) FROM messages WHERE status IN ('queued','processing'))::bigint      AS pending,
    (SELECT count(*) FROM webhook_deliveries WHERE delivered_at IS NULL AND next_retry_at IS NOT NULL)::bigint AS webhooks_pending;

-- name: CountAuditLogs :one
SELECT count(*) FROM audit_logs
WHERE (sqlc.narg('project_id')::uuid IS NULL OR project_id = sqlc.narg('project_id')::uuid);

-- name: AdminListContacts :many
SELECT * FROM contacts
WHERE project_id = @project_id AND deleted_at IS NULL
  AND (sqlc.narg('search')::text IS NULL
       OR phone ILIKE '%' || sqlc.narg('search')::text || '%'
       OR email ILIKE '%' || sqlc.narg('search')::text || '%'
       OR external_id ILIKE '%' || sqlc.narg('search')::text || '%'
       OR name ILIKE '%' || sqlc.narg('search')::text || '%')
ORDER BY created_at DESC LIMIT @row_limit OFFSET @row_offset;

-- name: CountContacts :one
SELECT count(*) FROM contacts WHERE project_id = @project_id AND deleted_at IS NULL;

-- name: CountDevices :one
SELECT count(*) FILTER (WHERE is_active)::bigint AS active, count(*)::bigint AS total FROM devices WHERE project_id = @project_id;

-- name: ListWebhookDeliveriesFiltered :many
SELECT * FROM webhook_deliveries
WHERE project_id = @project_id
  AND (sqlc.narg('event')::text IS NULL OR event = sqlc.narg('event')::text)
  AND (sqlc.narg('only_failed')::boolean IS NULL OR NOT sqlc.narg('only_failed')::boolean OR delivered_at IS NULL)
  AND (sqlc.narg('message_id')::uuid IS NULL OR message_id = sqlc.narg('message_id')::uuid)
ORDER BY created_at DESC LIMIT @row_limit OFFSET @row_offset;
