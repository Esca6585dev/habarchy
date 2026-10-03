-- +goose Up
-- Enumerations shared by all tables. Keep in sync with internal/domain/enums.go.
CREATE TYPE channel AS ENUM ('sms', 'email', 'push', 'telegram');
CREATE TYPE message_status AS ENUM ('queued', 'processing', 'sent', 'delivered', 'failed', 'cancelled');
CREATE TYPE message_priority AS ENUM ('high', 'normal', 'low');
CREATE TYPE provider_type AS ENUM ('http_sms', 'smpp', 'smtp', 'fcm', 'telegram_bot');
CREATE TYPE member_role AS ENUM ('owner', 'admin', 'developer', 'viewer');
CREATE TYPE device_platform AS ENUM ('android', 'ios', 'web');
CREATE TYPE project_status AS ENUM ('active', 'suspended', 'archived');
CREATE TYPE event_type AS ENUM ('queued', 'attempt', 'provider_response', 'sent', 'delivered', 'failed', 'cancelled', 'webhook_sent');

-- +goose Down
DROP TYPE IF EXISTS event_type;
DROP TYPE IF EXISTS project_status;
DROP TYPE IF EXISTS device_platform;
DROP TYPE IF EXISTS member_role;
DROP TYPE IF EXISTS provider_type;
DROP TYPE IF EXISTS message_priority;
DROP TYPE IF EXISTS message_status;
DROP TYPE IF EXISTS channel;
