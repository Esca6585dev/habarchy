-- name: CreateProvider :one
INSERT INTO providers (id, project_id, name, channel, type, priority, is_active, credentials_enc, rate_limit_per_sec)
VALUES (@id, @project_id, @name, @channel, @type, @priority, @is_active, @credentials_enc, @rate_limit_per_sec)
RETURNING *;

-- name: GetProvider :one
SELECT * FROM providers WHERE id = @id AND project_id = @project_id;

-- name: GetProviderByID :one
SELECT * FROM providers WHERE id = @id;

-- name: ListProviders :many
SELECT * FROM providers WHERE project_id = @project_id ORDER BY channel, priority, created_at;

-- name: ListActiveProvidersForChannel :many
SELECT * FROM providers
WHERE project_id = @project_id AND channel = @channel AND is_active
ORDER BY priority, created_at;

-- name: UpdateProvider :one
UPDATE providers SET
    name               = @name,
    priority           = @priority,
    is_active          = @is_active,
    credentials_enc    = @credentials_enc,
    rate_limit_per_sec = @rate_limit_per_sec,
    updated_at         = now()
WHERE id = @id AND project_id = @project_id
RETURNING *;

-- name: DeleteProvider :execrows
DELETE FROM providers WHERE id = @id AND project_id = @project_id;
