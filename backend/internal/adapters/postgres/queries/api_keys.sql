-- name: CreateAPIKey :one
INSERT INTO api_keys (project_id, name, prefix, hint, key_hash, scopes, ip_allowlist, expires_at, require_signature)
VALUES (@project_id, @name, @prefix, @hint, @key_hash, @scopes, @ip_allowlist, @expires_at, @require_signature)
RETURNING *;

-- name: GetAPIKeyByHash :one
SELECT * FROM api_keys WHERE key_hash = @key_hash;

-- name: GetAPIKey :one
SELECT * FROM api_keys WHERE id = @id AND project_id = @project_id;

-- name: ListAPIKeys :many
SELECT * FROM api_keys WHERE project_id = @project_id ORDER BY created_at DESC;

-- name: RevokeAPIKey :execrows
UPDATE api_keys SET revoked_at = now() WHERE id = @id AND project_id = @project_id AND revoked_at IS NULL;

-- name: TouchAPIKey :exec
UPDATE api_keys SET last_used_at = now() WHERE id = @id;

-- name: UpdateAPIKey :one
UPDATE api_keys SET
    name              = @name,
    scopes            = @scopes,
    ip_allowlist      = @ip_allowlist,
    require_signature = @require_signature
WHERE id = @id AND project_id = @project_id AND revoked_at IS NULL
RETURNING *;
