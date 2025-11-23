-- Rollback: Drop members table

DROP INDEX IF EXISTS idx_members_is_active;
DROP INDEX IF EXISTS idx_members_uuid;
DROP TABLE IF EXISTS members;

