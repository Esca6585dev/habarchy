-- name: CreateContact :one
INSERT INTO contacts (project_id, external_id, name, phone, email, whatsapp, telegram_chat_id, slack_id, locale, tags, attributes)
VALUES (@project_id, @external_id, @name, @phone, @email, @whatsapp, @telegram_chat_id, @slack_id, @locale, @tags, @attributes)
RETURNING *;

-- name: GetContact :one
SELECT * FROM contacts WHERE id = @id AND project_id = @project_id AND deleted_at IS NULL;

-- name: GetContactAny :one
-- Fetches a contact regardless of soft-delete state (for restore / purge).
SELECT * FROM contacts WHERE id = @id AND project_id = @project_id;

-- name: GetContactByExternalID :one
SELECT * FROM contacts WHERE project_id = @project_id AND external_id = @external_id AND deleted_at IS NULL;

-- name: GetContactByPhone :one
SELECT * FROM contacts WHERE project_id = @project_id AND phone = @phone AND deleted_at IS NULL LIMIT 1;

-- name: GetContactByEmail :one
SELECT * FROM contacts WHERE project_id = @project_id AND email = lower(@email) AND deleted_at IS NULL LIMIT 1;

-- name: ListContacts :many
SELECT * FROM contacts
WHERE project_id = @project_id AND deleted_at IS NULL
  AND (sqlc.narg('tag')::text IS NULL OR sqlc.narg('tag')::text = ANY (tags))
  AND (sqlc.narg('search')::text IS NULL
       OR phone ILIKE '%' || sqlc.narg('search')::text || '%'
       OR email ILIKE '%' || sqlc.narg('search')::text || '%'
       OR external_id ILIKE '%' || sqlc.narg('search')::text || '%'
       OR name ILIKE '%' || sqlc.narg('search')::text || '%')
ORDER BY created_at DESC
LIMIT @row_limit OFFSET @row_offset;

-- name: UpdateContact :one
UPDATE contacts SET
    external_id      = @external_id,
    name             = @name,
    phone            = @phone,
    email            = @email,
    whatsapp         = @whatsapp,
    telegram_chat_id = @telegram_chat_id,
    slack_id         = @slack_id,
    locale           = @locale,
    tags             = @tags,
    attributes       = @attributes,
    updated_at       = now()
WHERE id = @id AND project_id = @project_id AND deleted_at IS NULL
RETURNING *;

-- name: DeleteContact :execrows
-- Soft delete: keep the row, hide it everywhere.
UPDATE contacts SET deleted_at = now(), updated_at = now() WHERE id = @id AND project_id = @project_id AND deleted_at IS NULL;

-- name: RestoreContact :one
UPDATE contacts SET deleted_at = NULL, updated_at = now() WHERE id = @id AND project_id = @project_id AND deleted_at IS NOT NULL
RETURNING *;

-- name: ListDeletedContacts :many
SELECT * FROM contacts WHERE project_id = @project_id AND deleted_at IS NOT NULL
  AND (sqlc.narg('search')::text IS NULL
       OR phone ILIKE '%' || sqlc.narg('search')::text || '%'
       OR email ILIKE '%' || sqlc.narg('search')::text || '%'
       OR external_id ILIKE '%' || sqlc.narg('search')::text || '%'
       OR name ILIKE '%' || sqlc.narg('search')::text || '%')
ORDER BY deleted_at DESC
LIMIT @row_limit OFFSET @row_offset;

-- name: CountDeletedContacts :one
SELECT count(*) FROM contacts WHERE project_id = @project_id AND deleted_at IS NOT NULL;

-- name: PurgeContact :execrows
-- Permanent delete; only a row already in the trash can be purged.
DELETE FROM contacts WHERE id = @id AND project_id = @project_id AND deleted_at IS NOT NULL;

-- name: PurgeDeletedContacts :execrows
DELETE FROM contacts WHERE project_id = @project_id AND deleted_at IS NOT NULL;

-- name: UpsertDevice :one
INSERT INTO devices (project_id, contact_id, platform, fcm_token, app_version, last_seen_at, is_active)
VALUES (@project_id, @contact_id, @platform, @fcm_token, @app_version, now(), true)
ON CONFLICT (fcm_token) DO UPDATE SET
    project_id   = EXCLUDED.project_id,
    contact_id   = COALESCE(EXCLUDED.contact_id, devices.contact_id),
    platform     = EXCLUDED.platform,
    app_version  = EXCLUDED.app_version,
    last_seen_at = now(),
    is_active    = true
RETURNING *;

-- name: GetDeviceByToken :one
SELECT * FROM devices WHERE fcm_token = @fcm_token;

-- name: ListActiveDevicesForContact :many
SELECT * FROM devices WHERE contact_id = @contact_id AND is_active ORDER BY last_seen_at DESC;

-- name: ListDevices :many
SELECT * FROM devices WHERE project_id = @project_id ORDER BY last_seen_at DESC LIMIT @row_limit OFFSET @row_offset;

-- name: DisableDeviceByToken :execrows
UPDATE devices SET is_active = false WHERE fcm_token = @fcm_token AND is_active;

-- name: DeleteDeviceByToken :execrows
DELETE FROM devices WHERE fcm_token = @fcm_token AND project_id = @project_id;
