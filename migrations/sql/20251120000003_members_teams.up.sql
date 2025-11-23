-- Migration: Create members_teams junction table
-- Many-to-many relationship between members and teams
-- A member can belong to multiple teams, a team can have multiple members

CREATE TABLE IF NOT EXISTS members_teams (
    team_id INT NOT NULL,
    member_id INT NOT NULL,
    PRIMARY KEY (team_id, member_id),
    CONSTRAINT fk_members_teams_team 
        FOREIGN KEY (team_id) REFERENCES teams(id) ON DELETE CASCADE,
    CONSTRAINT fk_members_teams_member 
        FOREIGN KEY (member_id) REFERENCES members(id) ON DELETE CASCADE
);

-- Index for fast lookup of team members
CREATE INDEX IF NOT EXISTS idx_members_teams_team_id ON members_teams(team_id);

-- Index for fast lookup of member teams
CREATE INDEX IF NOT EXISTS idx_members_teams_member_id ON members_teams(member_id);

