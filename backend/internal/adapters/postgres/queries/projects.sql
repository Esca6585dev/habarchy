-- name: CreateProject :one
INSERT INTO projects (name, slug, daily_quota, monthly_quota, default_locale)
VALUES (@name, @slug, @daily_quota, @monthly_quota, @default_locale)
RETURNING *;

-- name: GetProject :one
SELECT * FROM projects WHERE id = @id;

-- name: GetProjectBySlug :one
SELECT * FROM projects WHERE slug = @slug;

-- name: ListProjects :many
SELECT * FROM projects ORDER BY created_at DESC LIMIT @row_limit OFFSET @row_offset;

-- name: ListProjectsForUser :many
SELECT sqlc.embed(p), m.role
FROM projects p
JOIN project_members m ON m.project_id = p.id
WHERE m.user_id = @user_id
ORDER BY p.created_at DESC;

-- name: UpdateProject :one
UPDATE projects SET
    name               = @name,
    status             = @status,
    daily_quota        = @daily_quota,
    monthly_quota      = @monthly_quota,
    webhook_url        = @webhook_url,
    default_locale     = @default_locale,
    allowed_ips        = @allowed_ips,
    auto_channel_order = @auto_channel_order,
    updated_at         = now()
WHERE id = @id
RETURNING *;

-- name: SetProjectWebhookSecret :exec
UPDATE projects SET webhook_secret_enc = @webhook_secret_enc, updated_at = now() WHERE id = @id;

-- name: DeleteProject :exec
DELETE FROM projects WHERE id = @id;

-- name: UpsertProjectMember :one
INSERT INTO project_members (project_id, user_id, role)
VALUES (@project_id, @user_id, @role)
ON CONFLICT (project_id, user_id) DO UPDATE SET role = EXCLUDED.role
RETURNING *;

-- name: GetProjectMember :one
SELECT * FROM project_members WHERE project_id = @project_id AND user_id = @user_id;

-- name: ListProjectMembers :many
SELECT sqlc.embed(m), u.email, u.full_name
FROM project_members m
JOIN users u ON u.id = m.user_id
WHERE m.project_id = @project_id
ORDER BY m.created_at;

-- name: RemoveProjectMember :exec
DELETE FROM project_members WHERE project_id = @project_id AND user_id = @user_id;

-- name: CountProjectOwners :one
SELECT count(*) FROM project_members WHERE project_id = @project_id AND role = 'owner';
