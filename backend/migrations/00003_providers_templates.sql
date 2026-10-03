-- +goose Up
CREATE TABLE providers (
    id                 uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    project_id         uuid NOT NULL REFERENCES projects (id) ON DELETE CASCADE,
    name               text NOT NULL,
    channel            channel NOT NULL,
    type               provider_type NOT NULL,
    priority           integer NOT NULL DEFAULT 100,   -- lower is tried first
    is_active          boolean NOT NULL DEFAULT true,
    credentials_enc    bytea NOT NULL,                 -- AES-256-GCM encrypted JSON
    rate_limit_per_sec integer NOT NULL DEFAULT 0,     -- 0 = unlimited
    created_at         timestamptz NOT NULL DEFAULT now(),
    updated_at         timestamptz NOT NULL DEFAULT now()
);
CREATE INDEX providers_project_channel_idx ON providers (project_id, channel, priority) WHERE is_active;

CREATE TABLE templates (
    id            uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    project_id    uuid NOT NULL REFERENCES projects (id) ON DELETE CASCADE,
    key           text NOT NULL,
    channel       channel NOT NULL,
    locale        text NOT NULL DEFAULT 'tk',
    subject       text NOT NULL DEFAULT '',
    body          text NOT NULL,
    required_vars text[] NOT NULL DEFAULT '{}',
    version       integer NOT NULL DEFAULT 1,
    is_active     boolean NOT NULL DEFAULT true,
    created_at    timestamptz NOT NULL DEFAULT now(),
    updated_at    timestamptz NOT NULL DEFAULT now(),
    UNIQUE (project_id, key, channel, locale)
);

-- Every save of a template appends an immutable version row.
CREATE TABLE template_versions (
    id            uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    template_id   uuid NOT NULL REFERENCES templates (id) ON DELETE CASCADE,
    version       integer NOT NULL,
    subject       text NOT NULL DEFAULT '',
    body          text NOT NULL,
    required_vars text[] NOT NULL DEFAULT '{}',
    created_by    uuid REFERENCES users (id) ON DELETE SET NULL,
    created_at    timestamptz NOT NULL DEFAULT now(),
    UNIQUE (template_id, version)
);

-- +goose Down
DROP TABLE IF EXISTS template_versions;
DROP TABLE IF EXISTS templates;
DROP TABLE IF EXISTS providers;
