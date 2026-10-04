-- +goose NO TRANSACTION
-- +goose Up
-- Android SMS gateway: a phone with a SIM card runs the Habarchy Gateway
-- app, long-polls its outbox and sends the SMS itself.
ALTER TYPE provider_type ADD VALUE IF NOT EXISTS 'android_sms';

-- One row per android_sms provider: the pairing key hash and the last
-- heartbeat the phone sent (battery, signal, operator, app version).
CREATE TABLE IF NOT EXISTS gateway_devices (
    provider_id  uuid PRIMARY KEY REFERENCES providers(id) ON DELETE CASCADE,
    project_id   uuid NOT NULL REFERENCES projects(id) ON DELETE CASCADE,
    key_hash     bytea NOT NULL UNIQUE,
    last_seen_at timestamptz,
    device_info  jsonb NOT NULL DEFAULT '{}'::jsonb,
    created_at   timestamptz NOT NULL DEFAULT now(),
    updated_at   timestamptz NOT NULL DEFAULT now()
);

-- Work queue for the phone. The worker inserts a row and waits for the
-- phone to report sent/failed; expired rows are skipped by the phone.
CREATE TABLE IF NOT EXISTS gateway_outbox (
    id            uuid PRIMARY KEY,
    provider_id   uuid NOT NULL REFERENCES providers(id) ON DELETE CASCADE,
    message_id    uuid NOT NULL,
    to_address    text NOT NULL,
    text          text NOT NULL,
    sim_slot      int  NOT NULL DEFAULT -1,
    status        text NOT NULL DEFAULT 'pending', -- pending | leased | sent | failed | delivered | expired
    error_code    text NOT NULL DEFAULT '',
    error_message text NOT NULL DEFAULT '',
    parts         int  NOT NULL DEFAULT 0,
    expires_at    timestamptz NOT NULL,
    leased_at     timestamptz,
    sent_at       timestamptz,
    delivered_at  timestamptz,
    created_at    timestamptz NOT NULL DEFAULT now(),
    updated_at    timestamptz NOT NULL DEFAULT now()
);
CREATE INDEX IF NOT EXISTS gateway_outbox_pending_idx ON gateway_outbox (provider_id, created_at) WHERE status IN ('pending', 'leased');
CREATE INDEX IF NOT EXISTS gateway_outbox_message_idx ON gateway_outbox (message_id);

-- SMS received by the phone, forwarded to the project webhook (sms.inbound).
CREATE TABLE IF NOT EXISTS gateway_inbound (
    id           uuid PRIMARY KEY,
    provider_id  uuid NOT NULL REFERENCES providers(id) ON DELETE CASCADE,
    project_id   uuid NOT NULL REFERENCES projects(id) ON DELETE CASCADE,
    from_address text NOT NULL,
    text         text NOT NULL,
    received_at  timestamptz NOT NULL,
    created_at   timestamptz NOT NULL DEFAULT now()
);
CREATE INDEX IF NOT EXISTS gateway_inbound_project_idx ON gateway_inbound (project_id, received_at DESC);

-- +goose Down
DROP TABLE IF EXISTS gateway_inbound;
DROP TABLE IF EXISTS gateway_outbox;
DROP TABLE IF EXISTS gateway_devices;
-- PostgreSQL cannot drop an enum value; 'android_sms' stays in provider_type.
