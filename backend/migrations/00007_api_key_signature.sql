-- +goose Up
-- When true, every request made with the key must carry a valid
-- X-Signature / X-Timestamp pair (HMAC-SHA256 keyed with the API key).
ALTER TABLE api_keys ADD COLUMN require_signature boolean NOT NULL DEFAULT false;

-- +goose Down
ALTER TABLE api_keys DROP COLUMN IF EXISTS require_signature;
