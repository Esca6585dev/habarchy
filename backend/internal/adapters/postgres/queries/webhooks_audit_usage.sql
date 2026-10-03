-- name: CreateWebhookDelivery :one
INSERT INTO webhook_deliveries (project_id, message_id, batch_id, event, url, payload, signature, next_retry_at)
VALUES (@project_id, @message_id, @batch_id, @event, @url, @payload, @signature, now())
RETURNING *;

-- name: GetWebhookDelivery :one
SELECT * FROM webhook_deliveries WHERE id = @id AND project_id = @project_id;

-- name: ListWebhookDeliveries :many
SELECT * FROM webhook_deliveries WHERE project_id = @project_id
ORDER BY created_at DESC LIMIT @row_limit OFFSET @row_offset;

-- name: RecordWebhookAttempt :exec
UPDATE webhook_deliveries SET
    attempts      = attempts + 1,
    response_code = @response_code,
    response_body = @response_body,
    next_retry_at = @next_retry_at,
    delivered_at  = @delivered_at
WHERE id = @id;

-- name: ResetWebhookDelivery :one
UPDATE webhook_deliveries SET attempts = 0, next_retry_at = now(), delivered_at = NULL
WHERE id = @id AND project_id = @project_id
RETURNING *;

-- name: CreateAuditLog :one
INSERT INTO audit_logs (project_id, user_id, action, entity_type, entity_id, changes, ip, user_agent)
VALUES (@project_id, @user_id, @action, @entity_type, @entity_id, @changes, @ip, @user_agent)
RETURNING *;

-- name: ListAuditLogs :many
SELECT * FROM audit_logs
WHERE (sqlc.narg('project_id')::uuid IS NULL OR project_id = sqlc.narg('project_id')::uuid)
ORDER BY created_at DESC LIMIT @row_limit OFFSET @row_offset;

-- name: UpsertUsageDaily :exec
INSERT INTO usage_daily (project_id, day, channel, queued, sent, delivered, failed, cost_micros, currency)
VALUES (@project_id, @day, @channel, @queued, @sent, @delivered, @failed, @cost_micros, @currency)
ON CONFLICT (project_id, day, channel) DO UPDATE SET
    queued      = EXCLUDED.queued,
    sent        = EXCLUDED.sent,
    delivered   = EXCLUDED.delivered,
    failed      = EXCLUDED.failed,
    cost_micros = EXCLUDED.cost_micros,
    currency    = EXCLUDED.currency;

-- name: AggregateUsageForDay :many
SELECT
    project_id,
    channel,
    count(*) FILTER (WHERE status IN ('queued', 'processing'))::bigint AS queued,
    count(*) FILTER (WHERE status = 'sent')::bigint                    AS sent,
    count(*) FILTER (WHERE status = 'delivered')::bigint               AS delivered,
    count(*) FILTER (WHERE status = 'failed')::bigint                  AS failed,
    coalesce(sum(cost_micros), 0)::bigint                              AS cost_micros,
    coalesce(max(currency), 'TMT')::text                               AS currency
FROM messages
WHERE created_at >= @day_start AND created_at < @day_end AND NOT is_test
GROUP BY project_id, channel;

-- name: ListUsageDaily :many
SELECT * FROM usage_daily
WHERE project_id = @project_id AND day >= @from_day AND day <= @to_day
ORDER BY day, channel;
