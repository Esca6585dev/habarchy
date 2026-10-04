-- +goose Up
-- Soft delete for contacts: a deleted contact is kept (deleted_at set) and
-- hidden from every normal query, so it can be restored or kept as a record
-- within the project. An admin purges it explicitly; nothing deletes it
-- automatically.
ALTER TABLE contacts ADD COLUMN IF NOT EXISTS deleted_at timestamptz;

-- The external_id uniqueness must ignore trashed rows, so re-importing the
-- same external_id after a delete creates a fresh active contact.
DROP INDEX IF EXISTS contacts_external_id_idx;
CREATE UNIQUE INDEX contacts_external_id_idx ON contacts (project_id, external_id) WHERE external_id <> '' AND deleted_at IS NULL;
CREATE INDEX IF NOT EXISTS contacts_deleted_idx ON contacts (project_id, deleted_at) WHERE deleted_at IS NOT NULL;

-- +goose Down
DROP INDEX IF EXISTS contacts_deleted_idx;
DROP INDEX IF EXISTS contacts_external_id_idx;
CREATE UNIQUE INDEX contacts_external_id_idx ON contacts (project_id, external_id) WHERE external_id <> '';
ALTER TABLE contacts DROP COLUMN IF EXISTS deleted_at;
