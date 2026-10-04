-- name: CreateUser :one
INSERT INTO users (email, password_hash, full_name)
VALUES (lower(@email), @password_hash, @full_name)
RETURNING *;

-- name: GetUserByID :one
SELECT * FROM users WHERE id = @id;

-- name: GetUserByEmail :one
SELECT * FROM users WHERE lower(email) = lower(@email);

-- name: ListUsers :many
SELECT * FROM users ORDER BY created_at DESC LIMIT @row_limit OFFSET @row_offset;

-- name: CountUsers :one
SELECT count(*) FROM users;

-- name: UpdateUserPassword :exec
UPDATE users SET password_hash = @password_hash, updated_at = now() WHERE id = @id;

-- name: UpdateUserProfile :one
UPDATE users SET full_name = @full_name, is_active = @is_active, updated_at = now()
WHERE id = @id RETURNING *;

-- name: UpdateUserSelfProfile :one
UPDATE users SET full_name = @full_name, bio = @bio, avatar_id = sqlc.narg('avatar_id'), updated_at = now()
WHERE id = @id RETURNING *;

-- name: SetUserTOTP :exec
UPDATE users SET totp_enabled = @totp_enabled, totp_secret_enc = @totp_secret_enc, updated_at = now()
WHERE id = @id;

-- name: TouchUserLogin :exec
UPDATE users SET last_login_at = now() WHERE id = @id;

-- name: CreateRefreshToken :one
INSERT INTO refresh_tokens (user_id, token_hash, user_agent, ip, expires_at)
VALUES (@user_id, @token_hash, @user_agent, @ip, @expires_at)
RETURNING *;

-- name: GetRefreshTokenByHash :one
SELECT * FROM refresh_tokens WHERE token_hash = @token_hash;

-- name: RevokeRefreshToken :exec
UPDATE refresh_tokens SET revoked_at = now(), replaced_by = @replaced_by WHERE id = @id;

-- name: RevokeAllUserRefreshTokens :exec
UPDATE refresh_tokens SET revoked_at = now() WHERE user_id = @user_id AND revoked_at IS NULL;

-- name: DeleteExpiredRefreshTokens :execrows
DELETE FROM refresh_tokens WHERE expires_at < now() - interval '7 days';
