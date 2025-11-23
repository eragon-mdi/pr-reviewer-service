-- Migration: Create pr_members junction table
-- Many-to-many relationship between pull requests and members with roles
-- Tracks which members are assigned to which PRs and their roles (reviewer, approver, etc.)

CREATE TABLE IF NOT EXISTS pr_members (
    pr_id INT NOT NULL,
    member_id INT NOT NULL,
    role_id INT NOT NULL,
    assigned_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    PRIMARY KEY (pr_id, member_id),
    CONSTRAINT fk_pr_members_pr 
        FOREIGN KEY (pr_id) REFERENCES pull_requests(id) ON DELETE CASCADE,
    CONSTRAINT fk_pr_members_member 
        FOREIGN KEY (member_id) REFERENCES members(id) ON DELETE CASCADE,
    CONSTRAINT fk_pr_members_role 
        FOREIGN KEY (role_id) REFERENCES roles(id) ON DELETE RESTRICT
);

-- Index for fast lookup of PR reviewers
CREATE INDEX IF NOT EXISTS idx_pr_members_pr_id ON pr_members(pr_id);

-- Index for fast lookup of member PRs
CREATE INDEX IF NOT EXISTS idx_pr_members_member_id ON pr_members(member_id);

-- Index for filtering by role
CREATE INDEX IF NOT EXISTS idx_pr_members_role_id ON pr_members(role_id);

