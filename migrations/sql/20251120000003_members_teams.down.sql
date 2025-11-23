-- Rollback: Drop members_teams table

DROP INDEX IF EXISTS idx_members_teams_member_id;
DROP INDEX IF EXISTS idx_members_teams_team_id;
DROP TABLE IF EXISTS members_teams;

