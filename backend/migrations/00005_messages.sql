-- +goose Up
CREATE TABLE batches (
    id              uuid PRIMARY KEY,
    project_id      uuid NOT NULL REFERENCES projects (id) ON DELETE CASCADE,
    template_key    text NOT NULL DEFAULT '',
    channel         channel NOT NULL,
    total           integer NOT NULL DEFAULT 0,
    queued          integer NOT NULL DEFAULT 0,
    sent            integer NOT NULL DEFAULT 0,
    delivered       integer NOT NULL DEFAULT 0,
    failed          integer NOT NULL DEFAULT 0,
    status          text NOT NULL DEFAULT 'processing',   -- processing | completed
    idempotency_key text,
    created_at      timestamptz NOT NULL DEFAULT now(),
    completed_at    timestamptz
);
CREATE INDEX batches_project_idx ON batches (project_id, created_at DESC);
CREATE UNIQUE INDEX batches_idempotency_idx ON batches (project_id, idempotency_key) WHERE idempotency_key IS NOT NULL;

-- messages is range-partitioned by month on created_at. The primary key
-- therefore includes created_at; look-ups by id alone still work because
-- the id is a UUID v7 and the API always scopes by project.
CREATE TABLE messages (
    id                  uuid NOT NULL,
    project_id          uuid NOT NULL REFERENCES projects (id) ON DELETE CASCADE,
    batch_id            uuid REFERENCES batches (id) ON DELETE SET NULL,
    channel             channel NOT NULL,
    to_address          text NOT NULL,
    contact_id          uuid REFERENCES contacts (id) ON DELETE SET NULL,
    template_key        text NOT NULL DEFAULT '',
    template_version    integer,
    rendered_subject    text NOT NULL DEFAULT '',
    rendered_body       text NOT NULL,
    status              message_status NOT NULL DEFAULT 'queued',
    priority            message_priority NOT NULL DEFAULT 'normal',
    provider_id         uuid REFERENCES providers (id) ON DELETE SET NULL,
    provider_message_id text NOT NULL DEFAULT '',
    error_code          text NOT NULL DEFAULT '',
    error_message       text NOT NULL DEFAULT '',
    attempts            integer NOT NULL DEFAULT 0,
    scheduled_at        timestamptz,
    sent_at             timestamptz,
    delivered_at        timestamptz,
    cost_micros         bigint NOT NULL DEFAULT 0,
    currency            text NOT NULL DEFAULT 'TMT',
    metadata            jsonb NOT NULL DEFAULT '{}',
    idempotency_key     text,
    is_test             boolean NOT NULL DEFAULT false,
    created_at          timestamptz NOT NULL DEFAULT now(),
    updated_at          timestamptz NOT NULL DEFAULT now(),
    PRIMARY KEY (id, created_at)
) PARTITION BY RANGE (created_at);

CREATE INDEX messages_project_created_idx ON messages (project_id, created_at DESC);
CREATE INDEX messages_status_scheduled_idx ON messages (status, scheduled_at);
CREATE INDEX messages_batch_idx ON messages (batch_id) WHERE batch_id IS NOT NULL;
CREATE INDEX messages_provider_msg_idx ON messages (provider_message_id) WHERE provider_message_id <> '';
CREATE INDEX messages_idempotency_idx ON messages (project_id, idempotency_key) WHERE idempotency_key IS NOT NULL;

-- A unique index on a partitioned table must include the partition key,
-- which would not give a real per-project guarantee. Idempotency keys are
-- therefore enforced in this small side table (and in Redis for 24 h).
CREATE TABLE message_idempotency_keys (
    project_id      uuid NOT NULL REFERENCES projects (id) ON DELETE CASCADE,
    idempotency_key text NOT NULL,
    message_id      uuid NOT NULL,
    created_at      timestamptz NOT NULL DEFAULT now(),
    PRIMARY KEY (project_id, idempotency_key)
);
CREATE INDEX message_idempotency_keys_created_idx ON message_idempotency_keys (created_at);

-- Per-message timeline. No FK to messages: FKs to partitioned tables must
-- reference the full (id, created_at) key, which the API never has at hand.
CREATE TABLE message_events (
    id          uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    message_id  uuid NOT NULL,
    type        event_type NOT NULL,
    provider_id uuid,
    payload     jsonb NOT NULL DEFAULT '{}',
    created_at  timestamptz NOT NULL DEFAULT now()
);
CREATE INDEX message_events_message_idx ON message_events (message_id, created_at);

-- +goose StatementBegin
-- ensure_messages_partition creates the monthly partition that contains
-- p_day if it does not exist yet, and returns its name. The scheduler calls
-- it daily for the current and the next month.
CREATE OR REPLACE FUNCTION ensure_messages_partition(p_day date) RETURNS text
LANGUAGE plpgsql AS $$
DECLARE
    start_date date := date_trunc('month', p_day)::date;
    end_date   date := (date_trunc('month', p_day) + interval '1 month')::date;
    part_name  text := 'messages_' || to_char(start_date, 'YYYY_MM');
BEGIN
    IF to_regclass(part_name) IS NULL THEN
        EXECUTE format(
            'CREATE TABLE %I PARTITION OF messages FOR VALUES FROM (%L) TO (%L)',
            part_name, start_date, end_date
        );
    END IF;
    RETURN part_name;
END
$$;
-- +goose StatementEnd

-- +goose StatementBegin
DO $$
BEGIN
    PERFORM ensure_messages_partition(current_date);
    PERFORM ensure_messages_partition((current_date + interval '1 month')::date);
    -- Safety net for rows outside the pre-created ranges (e.g. clock skew).
    EXECUTE 'CREATE TABLE IF NOT EXISTS messages_default PARTITION OF messages DEFAULT';
END
$$;
-- +goose StatementEnd

-- +goose Down
DROP FUNCTION IF EXISTS ensure_messages_partition(date);
DROP TABLE IF EXISTS message_events;
DROP TABLE IF EXISTS message_idempotency_keys;
DROP TABLE IF EXISTS messages;   -- drops all partitions
DROP TABLE IF EXISTS batches;
