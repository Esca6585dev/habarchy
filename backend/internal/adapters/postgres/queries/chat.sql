-- name: CreateChannel :one
INSERT INTO chat_channels (id, kind, name, topic, created_by, dm_key)
VALUES (@id, @kind, @name, @topic, @created_by, sqlc.narg('dm_key'))
RETURNING *;

-- name: GetChannel :one
SELECT * FROM chat_channels WHERE id = @id;

-- name: GetDirectChannel :one
SELECT * FROM chat_channels WHERE dm_key = @dm_key;

-- name: AddChatMember :exec
INSERT INTO chat_members (channel_id, user_id, role, last_read_at)
VALUES (@channel_id, @user_id, @role, now())
ON CONFLICT (channel_id, user_id) DO NOTHING;

-- name: AutoJoinPublicChannels :exec
INSERT INTO chat_members (channel_id, user_id)
SELECT c.id, @user_id FROM chat_channels c WHERE c.kind = 'public'
ON CONFLICT (channel_id, user_id) DO NOTHING;

-- name: RemoveChatMember :execrows
DELETE FROM chat_members WHERE channel_id = @channel_id AND user_id = @user_id;

-- name: GetChatMember :one
SELECT * FROM chat_members WHERE channel_id = @channel_id AND user_id = @user_id;

-- name: ListChannelMembers :many
SELECT u.id, u.full_name, u.email, u.avatar_id, u.bio, m.role, m.joined_at
FROM chat_members m JOIN users u ON u.id = m.user_id
WHERE m.channel_id = @channel_id
ORDER BY u.full_name;

-- name: ListMyChannelRows :many
SELECT c.*, mem.last_read_at FROM chat_members mem
JOIN chat_channels c ON c.id = mem.channel_id
WHERE mem.user_id = @user_id;

-- name: ChannelUnreadCounts :many
SELECT mem.channel_id, count(m.id)::bigint AS unread
FROM chat_members mem
LEFT JOIN chat_messages m ON m.channel_id = mem.channel_id AND m.deleted_at IS NULL
  AND m.created_at > mem.last_read_at AND (m.user_id IS NULL OR m.user_id <> mem.user_id)
WHERE mem.user_id = @user_id
GROUP BY mem.channel_id;

-- name: ChannelLastMessages :many
SELECT DISTINCT ON (m.channel_id)
  m.channel_id, m.body, m.created_at, m.user_id, (m.attachment_id IS NOT NULL)::boolean AS has_attachment
FROM chat_messages m
JOIN chat_members mem ON mem.channel_id = m.channel_id AND mem.user_id = @user_id
WHERE m.deleted_at IS NULL
ORDER BY m.channel_id, m.id DESC;

-- name: DirectPeers :many
SELECT m2.channel_id, u.id AS user_id, u.full_name, u.avatar_id
FROM chat_members m2
JOIN chat_channels c ON c.id = m2.channel_id AND c.kind = 'direct'
JOIN users u ON u.id = m2.user_id
JOIN chat_members me ON me.channel_id = m2.channel_id AND me.user_id = @user_id
WHERE m2.user_id <> @user_id;

-- name: CreateChatMessage :one
INSERT INTO chat_messages (id, channel_id, user_id, body, attachment_id)
VALUES (@id, @channel_id, @user_id, @body, sqlc.narg('attachment_id'))
RETURNING *;

-- name: TouchChannel :exec
UPDATE chat_channels SET updated_at = now() WHERE id = @id;

-- name: ListChatMessages :many
SELECT m.id, m.channel_id, m.user_id, m.body, m.attachment_id, m.created_at, m.edited_at,
  u.full_name AS author_name, u.avatar_id AS author_avatar,
  a.content_type AS attachment_type, a.filename AS attachment_name, a.size_bytes AS attachment_size
FROM chat_messages m
LEFT JOIN users u ON u.id = m.user_id
LEFT JOIN attachments a ON a.id = m.attachment_id
WHERE m.channel_id = @channel_id AND m.deleted_at IS NULL
  AND (sqlc.narg('before')::uuid IS NULL OR m.id < sqlc.narg('before')::uuid)
ORDER BY m.id DESC
LIMIT @row_limit;

-- name: GetChatMessage :one
SELECT * FROM chat_messages WHERE id = @id AND channel_id = @channel_id;

-- name: SoftDeleteChatMessage :execrows
UPDATE chat_messages SET deleted_at = now() WHERE id = @id AND user_id = @user_id AND deleted_at IS NULL;

-- name: MarkChannelRead :exec
UPDATE chat_members SET last_read_at = now() WHERE channel_id = @channel_id AND user_id = @user_id;
