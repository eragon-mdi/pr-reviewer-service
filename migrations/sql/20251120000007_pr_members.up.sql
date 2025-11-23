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

CREATE INDEX IF NOT EXISTS idx_pr_members_pr_id ON pr_members(pr_id);

CREATE INDEX IF NOT EXISTS idx_pr_members_member_id ON pr_members(member_id);

CREATE INDEX IF NOT EXISTS idx_pr_members_role_id ON pr_members(role_id);

