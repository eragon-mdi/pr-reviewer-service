-- Rollback: Drop pr_members table

DROP INDEX IF EXISTS idx_pr_members_role_id;
DROP INDEX IF EXISTS idx_pr_members_member_id;
DROP INDEX IF EXISTS idx_pr_members_pr_id;
DROP TABLE IF EXISTS pr_members;

