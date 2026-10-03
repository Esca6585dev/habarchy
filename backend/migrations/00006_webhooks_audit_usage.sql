-- +goose Up
CREATE TABLE webhook_deliveries (
    id            uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    project_id    uuid NOT NULL REFERENCES projects (id) ON DELETE CASCADE,
    message_id    uuid,
    batch_id      uuid,
    event         text NOT NULL,          -- message.sent | message.delivered | message.failed | batch.completed
    url           text NOT NULL,
    payload       jsonb NOT NULL,
    signature     text NOT NULL,
    response_code integer,
    response_body text NOT NULL DEFAULT '',
    attempts      integer NOT NULL DEFAULT 0,
    next_retry_at timestamptz,
    delivered_at  timestamptz,
    created_at    timestamptz NOT NULL DEFAULT now()
);
CREATE INDEX webhook_deliveries_project_idx ON webhook_deliveries (project_id, created_at DESC);
CREATE INDEX webhook_deliveries_message_idx ON webhook_deliveries (message_id) WHERE message_id IS NOT NULL;
CREATE INDEX webhook_deliveries_retry_idx ON webhook_deliveries (next_retry_at) WHERE delivered_at IS NULL;

CREATE TABLE audit_logs (
    id          uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    project_id  uuid REFERENCES projects (id) ON DELETE SET NULL,
    user_id     uuid REFERENCES users (id) ON DELETE SET NULL,
    action      text NOT NULL,            -- e.g. provider.create, api_key.revoke
    entity_type text NOT NULL DEFAULT '',
    entity_id   text NOT NULL DEFAULT '',
    changes     jsonb NOT NULL DEFAULT '{}',
    ip          text NOT NULL DEFAULT '',
    user_agent  text NOT NULL DEFAULT '',
    created_at  timestamptz NOT NULL DEFAULT now()
);
CREATE INDEX audit_logs_project_idx ON audit_logs (project_id, created_at DESC);
CREATE INDEX audit_logs_user_idx ON audit_logs (user_id, created_at DESC);

-- Aggregated by the scheduler from messages; read by usage reports and quotas.
CREATE TABLE usage_daily (
    project_id  uuid NOT NULL REFERENCES projects (id) ON DELETE CASCADE,
    day         date NOT NULL,
    channel     channel NOT NULL,
    queued      bigint NOT NULL DEFAULT 0,
    sent        bigint NOT NULL DEFAULT 0,
    delivered   bigint NOT NULL DEFAULT 0,
    failed      bigint NOT NULL DEFAULT 0,
    cost_micros bigint NOT NULL DEFAULT 0,
    currency    text NOT NULL DEFAULT 'TMT',
    PRIMARY KEY (project_id, day, channel)
);

-- +goose Down
DROP TABLE IF EXISTS usage_daily;
DROP TABLE IF EXISTS audit_logs;
DROP TABLE IF EXISTS webhook_deliveries;
