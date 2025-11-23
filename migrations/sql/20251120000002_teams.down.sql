-- Rollback: Drop teams table

DROP INDEX IF EXISTS idx_teams_name;
DROP TABLE IF EXISTS teams;

