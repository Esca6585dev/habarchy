-- +goose NO TRANSACTION
-- +goose Up
-- New channels (WhatsApp Cloud API, Slack) and contact groups for
-- broadcasting to colleagues / classmates / customers.
ALTER TYPE channel ADD VALUE IF NOT EXISTS 'whatsapp';
ALTER TYPE channel ADD VALUE IF NOT EXISTS 'slack';
ALTER TYPE provider_type ADD VALUE IF NOT EXISTS 'whatsapp_cloud';
ALTER TYPE provider_type ADD VALUE IF NOT EXISTS 'slack';

ALTER TABLE contacts ADD COLUMN IF NOT EXISTS name     text NOT NULL DEFAULT '';
ALTER TABLE contacts ADD COLUMN IF NOT EXISTS whatsapp text NOT NULL DEFAULT '';   -- E.164, defaults to phone when empty
ALTER TABLE contacts ADD COLUMN IF NOT EXISTS slack_id text NOT NULL DEFAULT '';   -- Slack user (U…) or channel (C…) id
CREATE INDEX IF NOT EXISTS contacts_name_idx ON contacts (project_id, lower(name));

CREATE TABLE IF NOT EXISTS contact_groups (
    id          uuid PRIMARY KEY,
    project_id  uuid NOT NULL REFERENCES projects (id) ON DELETE CASCADE,
    name        text NOT NULL,
    description text NOT NULL DEFAULT '',
    created_at  timestamptz NOT NULL DEFAULT now(),
    updated_at  timestamptz NOT NULL DEFAULT now(),
    UNIQUE (project_id, name)
);

CREATE TABLE IF NOT EXISTS contact_group_members (
    group_id   uuid NOT NULL REFERENCES contact_groups (id) ON DELETE CASCADE,
    contact_id uuid NOT NULL REFERENCES contacts (id) ON DELETE CASCADE,
    added_at   timestamptz NOT NULL DEFAULT now(),
    PRIMARY KEY (group_id, contact_id)
);
CREATE INDEX IF NOT EXISTS contact_group_members_contact_idx ON contact_group_members (contact_id);

-- +goose Down
DROP TABLE IF EXISTS contact_group_members;
DROP TABLE IF EXISTS contact_groups;
ALTER TABLE contacts DROP COLUMN IF EXISTS slack_id;
ALTER TABLE contacts DROP COLUMN IF EXISTS whatsapp;
ALTER TABLE contacts DROP COLUMN IF EXISTS name;
-- enum values cannot be dropped in PostgreSQL
