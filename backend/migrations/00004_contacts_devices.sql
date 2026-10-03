-- +goose Up
CREATE TABLE contacts (
    id               uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    project_id       uuid NOT NULL REFERENCES projects (id) ON DELETE CASCADE,
    external_id      text NOT NULL DEFAULT '',
    phone            text NOT NULL DEFAULT '',    -- E.164
    email            text NOT NULL DEFAULT '',    -- lower-cased
    telegram_chat_id text NOT NULL DEFAULT '',
    locale           text NOT NULL DEFAULT 'tk',
    tags             text[] NOT NULL DEFAULT '{}',
    attributes       jsonb NOT NULL DEFAULT '{}',
    created_at       timestamptz NOT NULL DEFAULT now(),
    updated_at       timestamptz NOT NULL DEFAULT now()
);
CREATE UNIQUE INDEX contacts_external_id_idx ON contacts (project_id, external_id) WHERE external_id <> '';
CREATE INDEX contacts_phone_idx ON contacts (project_id, phone) WHERE phone <> '';
CREATE INDEX contacts_email_idx ON contacts (project_id, email) WHERE email <> '';
CREATE INDEX contacts_telegram_idx ON contacts (project_id, telegram_chat_id) WHERE telegram_chat_id <> '';
CREATE INDEX contacts_tags_idx ON contacts USING gin (tags);

CREATE TABLE devices (
    id           uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    project_id   uuid NOT NULL REFERENCES projects (id) ON DELETE CASCADE,
    contact_id   uuid REFERENCES contacts (id) ON DELETE SET NULL,
    platform     device_platform NOT NULL,
    fcm_token    text NOT NULL UNIQUE,
    app_version  text NOT NULL DEFAULT '',
    last_seen_at timestamptz NOT NULL DEFAULT now(),
    is_active    boolean NOT NULL DEFAULT true,
    created_at   timestamptz NOT NULL DEFAULT now()
);
CREATE INDEX devices_contact_idx ON devices (contact_id) WHERE is_active;
CREATE INDEX devices_project_idx ON devices (project_id);

-- +goose Down
DROP TABLE IF EXISTS devices;
DROP TABLE IF EXISTS contacts;
