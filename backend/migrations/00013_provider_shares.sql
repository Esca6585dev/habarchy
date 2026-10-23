-- +goose Up
-- Provider sharing: a tenant can let another tenant send through one of its
-- phones (android_sms). The 0796 SMPP line is admin-only and never shared;
-- that is enforced in the service.
CREATE TABLE provider_shares (
    provider_id        uuid NOT NULL REFERENCES providers (id) ON DELETE CASCADE,
    owner_project_id   uuid NOT NULL REFERENCES projects (id) ON DELETE CASCADE,
    grantee_project_id uuid NOT NULL REFERENCES projects (id) ON DELETE CASCADE,
    created_at         timestamptz NOT NULL DEFAULT now(),
    PRIMARY KEY (provider_id, grantee_project_id),
    CHECK (owner_project_id <> grantee_project_id)
);
CREATE INDEX provider_shares_grantee_idx ON provider_shares (grantee_project_id);

-- +goose Down
DROP TABLE provider_shares;
