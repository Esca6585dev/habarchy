-- name: CreateAttachment :one
INSERT INTO attachments (id, uploaded_by, content_type, size_bytes, filename, data)
VALUES (@id, @uploaded_by, @content_type, @size_bytes, @filename, @data)
RETURNING id, uploaded_by, content_type, size_bytes, filename, created_at;

-- name: GetAttachment :one
SELECT * FROM attachments WHERE id = @id;
