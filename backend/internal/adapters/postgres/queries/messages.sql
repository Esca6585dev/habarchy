-- name: CreateMessage :one
INSERT INTO messages (
    id, project_id, batch_id, channel, to_address, contact_id, template_key, template_version,
    rendered_subject, rendered_body, status, priority, scheduled_at, metadata, idempotency_key, is_test
) VALUES (
    @id, @project_id, @batch_id, @channel, @to_address, @contact_id, @template_key, @template_version,
    @rendered_subject, @rendered_body, @status, @priority, @scheduled_at, @metadata, @idempotency_key, @is_test
)
RETURNING *;

-- name: GetMessage :one
SELECT * FROM messages WHERE id = @id AND project_id = @project_id;

-- name: GetMessageByID :one
SELECT * FROM messages WHERE id = @id;

-- name: GetMessageByProviderMessageID :one
SELECT * FROM messages WHERE provider_id = @provider_id AND provider_message_id = @provider_message_id
ORDER BY created_at DESC LIMIT 1;

-- name: ListMessages :many
SELECT * FROM messages
WHERE project_id = @project_id
  AND (sqlc.narg('status')::message_status IS NULL OR status = sqlc.narg('status')::message_status)
  AND (sqlc.narg('channel')::channel IS NULL OR channel = sqlc.narg('channel')::channel)
  AND (sqlc.narg('from_ts')::timestamptz IS NULL OR created_at >= sqlc.narg('from_ts')::timestamptz)
  AND (sqlc.narg('to_ts')::timestamptz IS NULL OR created_at < sqlc.narg('to_ts')::timestamptz)
  AND (sqlc.narg('cursor_id')::uuid IS NULL OR id < sqlc.narg('cursor_id')::uuid)
ORDER BY id DESC
LIMIT @row_limit;

-- name: MarkMessageProcessing :one
UPDATE messages SET status = 'processing', attempts = attempts + 1, updated_at = now()
WHERE id = @id AND status IN ('queued', 'processing')
RETURNING *;

-- name: MarkMessageSent :exec
UPDATE messages SET
    status              = 'sent',
    provider_id         = @provider_id,
    provider_message_id = @provider_message_id,
    cost_micros         = @cost_micros,
    currency            = @currency,
    error_code          = '',
    error_message       = '',
    sent_at             = now(),
    updated_at          = now()
WHERE id = @id;

-- name: MarkMessageDelivered :execrows
UPDATE messages SET status = 'delivered', delivered_at = now(), updated_at = now()
WHERE id = @id AND status = 'sent';

-- name: MarkMessageFailed :exec
UPDATE messages SET
    status        = 'failed',
    provider_id   = COALESCE(@provider_id, provider_id),
    error_code    = @error_code,
    error_message = @error_message,
    updated_at    = now()
WHERE id = @id;

-- name: RequeueMessage :exec
UPDATE messages SET status = 'queued', error_code = @error_code, error_message = @error_message, updated_at = now()
WHERE id = @id;

-- name: CancelMessage :execrows
UPDATE messages SET status = 'cancelled', updated_at = now()
WHERE id = @id AND project_id = @project_id AND status = 'queued';

-- name: ResetMessageForResend :one
UPDATE messages SET
    status = 'queued', attempts = 0, provider_id = NULL, provider_message_id = '',
    error_code = '', error_message = '', sent_at = NULL, delivered_at = NULL, updated_at = now()
WHERE id = @id AND project_id = @project_id AND status IN ('failed', 'cancelled')
RETURNING *;

-- name: CountMessagesSince :one
SELECT count(*) FROM messages
WHERE project_id = @project_id AND created_at >= @since AND status <> 'cancelled' AND NOT is_test;

-- name: CreateMessageEvent :one
INSERT INTO message_events (message_id, type, provider_id, payload)
VALUES (@message_id, @type, @provider_id, @payload)
RETURNING *;

-- name: ListMessageEvents :many
SELECT * FROM message_events WHERE message_id = @message_id ORDER BY created_at, id;

-- name: ReserveIdempotencyKey :one
INSERT INTO message_idempotency_keys (project_id, idempotency_key, message_id)
VALUES (@project_id, @idempotency_key, @message_id)
ON CONFLICT (project_id, idempotency_key) DO NOTHING
RETURNING message_id;

-- name: GetIdempotencyKey :one
SELECT message_id FROM message_idempotency_keys WHERE project_id = @project_id AND idempotency_key = @idempotency_key;

-- name: DeleteExpiredIdempotencyKeys :execrows
DELETE FROM message_idempotency_keys WHERE created_at < @before;

-- name: EnsureMessagesPartition :one
SELECT ensure_messages_partition(@day::date) AS partition_name;

-- name: CreateBatch :one
INSERT INTO batches (id, project_id, template_key, channel, total, queued, idempotency_key)
VALUES (@id, @project_id, @template_key, @channel, @total, @total, @idempotency_key)
RETURNING *;

-- name: GetBatch :one
SELECT * FROM batches WHERE id = @id AND project_id = @project_id;

-- name: GetBatchByIdempotencyKey :one
SELECT * FROM batches WHERE project_id = @project_id AND idempotency_key = @idempotency_key;

-- name: RefreshBatchCounters :one
-- Counters only; the processing -> completed transition is done once by
-- CompleteBatchIfDone so the batch.completed webhook fires exactly once.
UPDATE batches b SET
    queued    = s.queued,
    sent      = s.sent,
    delivered = s.delivered,
    failed    = s.failed
FROM (
    SELECT
        count(*) FILTER (WHERE status IN ('queued', 'processing'))::int AS queued,
        count(*) FILTER (WHERE status = 'sent')::int                    AS sent,
        count(*) FILTER (WHERE status = 'delivered')::int               AS delivered,
        count(*) FILTER (WHERE status IN ('failed', 'cancelled'))::int  AS failed
    FROM messages WHERE batch_id = @id
) s
WHERE b.id = @id
RETURNING b.*;

-- name: CompleteBatchIfDone :one
UPDATE batches b SET status = 'completed', completed_at = now()
WHERE b.id = @id AND b.status = 'processing'
  AND NOT EXISTS (SELECT 1 FROM messages m WHERE m.batch_id = b.id AND m.status IN ('queued', 'processing'))
RETURNING b.*;

-- name: ListMessagesByBatch :many
SELECT * FROM messages WHERE batch_id = @batch_id ORDER BY id LIMIT @row_limit OFFSET @row_offset;

-- name: ListBatches :many
SELECT * FROM batches WHERE project_id = @project_id ORDER BY created_at DESC LIMIT @row_limit OFFSET @row_offset;

-- name: ListDueWebhookDeliveries :many
SELECT * FROM webhook_deliveries
WHERE delivered_at IS NULL AND next_retry_at IS NOT NULL AND next_retry_at <= now()
ORDER BY next_retry_at LIMIT @row_limit;
