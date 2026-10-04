-- name: UpsertGatewayDevice :one
INSERT INTO gateway_devices (provider_id, project_id, key_hash)
VALUES (@provider_id, @project_id, @key_hash)
ON CONFLICT (provider_id) DO UPDATE SET key_hash = EXCLUDED.key_hash, updated_at = now()
RETURNING *;

-- name: GetGatewayDeviceByKeyHash :one
SELECT * FROM gateway_devices WHERE key_hash = @key_hash;

-- name: GetGatewayDevice :one
SELECT * FROM gateway_devices WHERE provider_id = @provider_id;

-- name: ListGatewayDevicesForProject :many
SELECT * FROM gateway_devices WHERE project_id = @project_id;

-- name: TouchGatewayDevice :exec
UPDATE gateway_devices SET last_seen_at = now(), device_info = @device_info, updated_at = now()
WHERE provider_id = @provider_id;

-- name: EnqueueGatewayOutbox :one
INSERT INTO gateway_outbox (id, provider_id, message_id, to_address, text, sim_slot, expires_at)
VALUES (@id, @provider_id, @message_id, @to_address, @text, @sim_slot, @expires_at)
RETURNING *;

-- name: LeaseGatewayOutbox :many
UPDATE gateway_outbox o SET status = 'leased', leased_at = now(), updated_at = now()
WHERE o.id IN (
    SELECT g.id FROM gateway_outbox g
    WHERE g.provider_id = @provider_id AND g.status = 'pending' AND g.expires_at > now()
    ORDER BY g.created_at
    LIMIT @row_limit
    FOR UPDATE SKIP LOCKED
)
RETURNING o.*;

-- name: GetGatewayOutbox :one
SELECT * FROM gateway_outbox WHERE id = @id;

-- name: GetGatewayOutboxForProvider :one
SELECT * FROM gateway_outbox WHERE id = @id AND provider_id = @provider_id;

-- name: SetGatewayOutboxResult :one
UPDATE gateway_outbox SET
    status        = @status,
    error_code    = @error_code,
    error_message = @error_message,
    parts         = CASE WHEN @parts::int > 0 THEN @parts::int ELSE parts END,
    sent_at       = CASE WHEN @status::text IN ('sent', 'delivered') AND sent_at IS NULL THEN now() ELSE sent_at END,
    delivered_at  = CASE WHEN @status::text = 'delivered' THEN now() ELSE delivered_at END,
    updated_at    = now()
WHERE id = @id
RETURNING *;

-- name: ExpireGatewayOutbox :execrows
UPDATE gateway_outbox SET status = 'expired', updated_at = now()
WHERE id = @id AND status IN ('pending', 'leased');

-- name: CountGatewayOutboxPending :one
SELECT count(*)::bigint FROM gateway_outbox
WHERE provider_id = @provider_id AND status IN ('pending', 'leased') AND expires_at > now();

-- name: PurgeGatewayOutbox :execrows
DELETE FROM gateway_outbox WHERE created_at < @before;

-- name: InsertGatewayInbound :one
INSERT INTO gateway_inbound (id, provider_id, project_id, from_address, text, received_at)
VALUES (@id, @provider_id, @project_id, @from_address, @text, @received_at)
RETURNING *;

-- name: ListGatewayInbound :many
SELECT * FROM gateway_inbound WHERE project_id = @project_id
ORDER BY received_at DESC LIMIT @row_limit OFFSET @row_offset;

-- name: CountGatewayInbound :one
SELECT count(*)::bigint FROM gateway_inbound WHERE project_id = @project_id;
