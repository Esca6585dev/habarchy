-- +goose Up
CREATE TABLE users (
    id              uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    email           text NOT NULL,
    password_hash   text NOT NULL,
    full_name       text NOT NULL DEFAULT '',
    is_active       boolean NOT NULL DEFAULT true,
    totp_enabled    boolean NOT NULL DEFAULT false,
    totp_secret_enc bytea,
    last_login_at   timestamptz,
    created_at      timestamptz NOT NULL DEFAULT now(),
    updated_at      timestamptz NOT NULL DEFAULT now()
);
CREATE UNIQUE INDEX users_email_lower_idx ON users (lower(email));

-- Refresh tokens are rotated on every use; only the SHA-256 hash is stored.
CREATE TABLE refresh_tokens (
    id          uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id     uuid NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    token_hash  bytea NOT NULL UNIQUE,
    user_agent  text NOT NULL DEFAULT '',
    ip          text NOT NULL DEFAULT '',
    expires_at  timestamptz NOT NULL,
    revoked_at  timestamptz,
    replaced_by uuid,
    created_at  timestamptz NOT NULL DEFAULT now()
);
CREATE INDEX refresh_tokens_user_idx ON refresh_tokens (user_id);
CREATE INDEX refresh_tokens_expires_idx ON refresh_tokens (expires_at);

CREATE TABLE projects (
    id                 uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    name               text NOT NULL,
    slug               text NOT NULL UNIQUE,
    status             project_status NOT NULL DEFAULT 'active',
    daily_quota        bigint NOT NULL DEFAULT 0,   -- 0 = unlimited
    monthly_quota      bigint NOT NULL DEFAULT 0,   -- 0 = unlimited
    webhook_url        text NOT NULL DEFAULT '',
    webhook_secret_enc bytea,
    default_locale     text NOT NULL DEFAULT 'tk',
    allowed_ips        text[] NOT NULL DEFAULT '{}',
    auto_channel_order channel[] NOT NULL DEFAULT '{push,sms,email}',
    created_at         timestamptz NOT NULL DEFAULT now(),
    updated_at         timestamptz NOT NULL DEFAULT now()
);

CREATE TABLE project_members (
    project_id uuid NOT NULL REFERENCES projects (id) ON DELETE CASCADE,
    user_id    uuid NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    role       member_role NOT NULL DEFAULT 'viewer',
    created_at timestamptz NOT NULL DEFAULT now(),
    PRIMARY KEY (project_id, user_id)
);
CREATE INDEX project_members_user_idx ON project_members (user_id);

CREATE TABLE api_keys (
    id           uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    project_id   uuid NOT NULL REFERENCES projects (id) ON DELETE CASCADE,
    name         text NOT NULL,
    prefix       text NOT NULL,            -- hb_live_ | hb_test_
    hint         text NOT NULL,            -- last 4 characters for display
    key_hash     bytea NOT NULL UNIQUE,    -- sha256(plaintext)
    scopes       text[] NOT NULL DEFAULT '{}',
    ip_allowlist text[] NOT NULL DEFAULT '{}',
    last_used_at timestamptz,
    expires_at   timestamptz,
    revoked_at   timestamptz,
    created_at   timestamptz NOT NULL DEFAULT now()
);
CREATE INDEX api_keys_project_idx ON api_keys (project_id);

-- +goose Down
DROP TABLE IF EXISTS api_keys;
DROP TABLE IF EXISTS project_members;
DROP TABLE IF EXISTS projects;
DROP TABLE IF EXISTS refresh_tokens;
DROP TABLE IF EXISTS users;
