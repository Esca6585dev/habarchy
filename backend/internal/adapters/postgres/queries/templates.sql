-- name: CreateTemplate :one
INSERT INTO templates (project_id, key, channel, locale, subject, body, required_vars)
VALUES (@project_id, @key, @channel, @locale, @subject, @body, @required_vars)
RETURNING *;

-- name: GetTemplate :one
SELECT * FROM templates WHERE id = @id AND project_id = @project_id;

-- name: FindTemplate :one
SELECT * FROM templates
WHERE project_id = @project_id AND key = @key AND channel = @channel AND locale = @locale AND is_active;

-- name: ListTemplates :many
SELECT * FROM templates WHERE project_id = @project_id ORDER BY key, channel, locale;

-- name: ListTemplatesByKey :many
SELECT * FROM templates WHERE project_id = @project_id AND key = @key ORDER BY channel, locale;

-- name: UpdateTemplate :one
UPDATE templates SET
    subject       = @subject,
    body          = @body,
    required_vars = @required_vars,
    is_active     = @is_active,
    version       = version + 1,
    updated_at    = now()
WHERE id = @id AND project_id = @project_id
RETURNING *;

-- name: DeleteTemplate :execrows
DELETE FROM templates WHERE id = @id AND project_id = @project_id;

-- name: CreateTemplateVersion :one
INSERT INTO template_versions (template_id, version, subject, body, required_vars, created_by)
VALUES (@template_id, @version, @subject, @body, @required_vars, @created_by)
RETURNING *;

-- name: ListTemplateVersions :many
SELECT * FROM template_versions WHERE template_id = @template_id ORDER BY version DESC;

-- name: GetTemplateVersion :one
SELECT * FROM template_versions WHERE template_id = @template_id AND version = @version;
