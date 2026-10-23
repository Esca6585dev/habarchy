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

-- name: ListActiveProvidersForChannelWithShared :many
-- The project's own active providers first (by priority), then any active
-- providers shared with it by other tenants.
SELECT p.* FROM providers p
WHERE p.channel = @channel AND p.is_active
  AND (p.project_id = @project_id
       OR p.id IN (SELECT provider_id FROM provider_shares WHERE grantee_project_id = @project_id))
ORDER BY (p.project_id = @project_id) DESC, p.priority, p.created_at;

-- name: ShareProvider :exec
INSERT INTO provider_shares (provider_id, owner_project_id, grantee_project_id)
VALUES (@provider_id, @owner_project_id, @grantee_project_id)
ON CONFLICT (provider_id, grantee_project_id) DO NOTHING;

-- name: RevokeProviderShare :exec
DELETE FROM provider_shares WHERE provider_id = @provider_id AND grantee_project_id = @grantee_project_id;

-- name: ListProviderShares :many
SELECT ps.grantee_project_id, pr.name AS grantee_name, pr.slug AS grantee_slug, ps.created_at
FROM provider_shares ps JOIN projects pr ON pr.id = ps.grantee_project_id
WHERE ps.provider_id = @provider_id
ORDER BY ps.created_at DESC;

-- name: ListSharedProvidersForProject :many
SELECT sqlc.embed(p), op.name AS owner_name, op.slug AS owner_slug
FROM provider_shares ps
JOIN providers p ON p.id = ps.provider_id
JOIN projects op ON op.id = p.project_id
WHERE ps.grantee_project_id = @grantee_project_id AND p.is_active
ORDER BY p.created_at DESC;
