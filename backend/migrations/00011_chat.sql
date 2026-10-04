-- +goose Up
-- Internal team chat: user profiles (avatar + bio), attachments (avatars
-- and chat images), channels (public / private groups / direct), members
-- and messages.
ALTER TABLE users ADD COLUMN IF NOT EXISTS bio text NOT NULL DEFAULT '';
ALTER TABLE users ADD COLUMN IF NOT EXISTS avatar_id uuid;

CREATE TABLE IF NOT EXISTS attachments (
    id           uuid PRIMARY KEY,
    uploaded_by  uuid REFERENCES users (id) ON DELETE SET NULL,
    content_type text NOT NULL,
    size_bytes   integer NOT NULL,
    filename     text NOT NULL DEFAULT '',
    data         bytea NOT NULL,
    created_at   timestamptz NOT NULL DEFAULT now()
);

CREATE TYPE chat_channel_kind AS ENUM ('public', 'private', 'direct');

CREATE TABLE IF NOT EXISTS chat_channels (
    id         uuid PRIMARY KEY,
    kind       chat_channel_kind NOT NULL,
    name       text NOT NULL DEFAULT '',
    topic      text NOT NULL DEFAULT '',
    created_by uuid REFERENCES users (id) ON DELETE SET NULL,
    dm_key     text,                       -- sorted "uidA:uidB" for direct channels
    created_at timestamptz NOT NULL DEFAULT now(),
    updated_at timestamptz NOT NULL DEFAULT now()
);
CREATE UNIQUE INDEX chat_channels_dm_key_idx ON chat_channels (dm_key) WHERE dm_key IS NOT NULL;
CREATE INDEX chat_channels_kind_idx ON chat_channels (kind);

CREATE TABLE IF NOT EXISTS chat_members (
    channel_id   uuid NOT NULL REFERENCES chat_channels (id) ON DELETE CASCADE,
    user_id      uuid NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    role         text NOT NULL DEFAULT 'member',  -- owner | member
    last_read_at timestamptz NOT NULL DEFAULT now(),
    joined_at    timestamptz NOT NULL DEFAULT now(),
    PRIMARY KEY (channel_id, user_id)
);
CREATE INDEX chat_members_user_idx ON chat_members (user_id);

CREATE TABLE IF NOT EXISTS chat_messages (
    id            uuid PRIMARY KEY,         -- uuid v7, sorts by time
    channel_id    uuid NOT NULL REFERENCES chat_channels (id) ON DELETE CASCADE,
    user_id       uuid REFERENCES users (id) ON DELETE SET NULL,
    body          text NOT NULL DEFAULT '',
    attachment_id uuid REFERENCES attachments (id) ON DELETE SET NULL,
    created_at    timestamptz NOT NULL DEFAULT now(),
    edited_at     timestamptz,
    deleted_at    timestamptz
);
CREATE INDEX chat_messages_channel_idx ON chat_messages (channel_id, id DESC) WHERE deleted_at IS NULL;

-- +goose Down
DROP TABLE IF EXISTS chat_messages;
DROP TABLE IF EXISTS chat_members;
DROP TABLE IF EXISTS chat_channels;
DROP TYPE IF EXISTS chat_channel_kind;
ALTER TABLE users DROP COLUMN IF EXISTS avatar_id;
ALTER TABLE users DROP COLUMN IF EXISTS bio;
DROP TABLE IF EXISTS attachments;
