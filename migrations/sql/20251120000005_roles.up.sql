-- Migration: Create roles table
-- Roles represent member roles in pull requests (author, reviewer, approver, etc.)

CREATE TABLE IF NOT EXISTS roles (
    id SERIAL PRIMARY KEY,
    role VARCHAR(20) NOT NULL UNIQUE
);

-- Insert default roles
INSERT INTO roles (role) VALUES ('author'), ('reviewer'), ('approver')
ON CONFLICT (role) DO NOTHING;

