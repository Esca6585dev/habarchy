-- +goose Up
-- Phone + OTP identity for end-user clients, and a platform-wide owner flag.
-- Phone-only users have an empty email, so the email uniqueness index becomes
-- partial (it already was case-insensitive) and a partial phone index is added.
ALTER TABLE users ADD COLUMN phone          text    NOT NULL DEFAULT '';
ALTER TABLE users ADD COLUMN phone_verified boolean NOT NULL DEFAULT false;
ALTER TABLE users ADD COLUMN platform_admin boolean NOT NULL DEFAULT false;

DROP INDEX users_email_lower_idx;
CREATE UNIQUE INDEX users_email_lower_idx ON users (lower(email)) WHERE email <> '';
CREATE UNIQUE INDEX users_phone_idx       ON users (phone)        WHERE phone <> '';

-- +goose Down
DROP INDEX users_phone_idx;
DROP INDEX users_email_lower_idx;
CREATE UNIQUE INDEX users_email_lower_idx ON users (lower(email));
ALTER TABLE users DROP COLUMN platform_admin;
ALTER TABLE users DROP COLUMN phone_verified;
ALTER TABLE users DROP COLUMN phone;
