-- Migration: Create teams table
-- Teams represent groups of members with unique names

CREATE TABLE IF NOT EXISTS teams (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL UNIQUE
);

-- Index for fast lookup by name (most common query pattern)
CREATE INDEX IF NOT EXISTS idx_teams_name ON teams(name);

