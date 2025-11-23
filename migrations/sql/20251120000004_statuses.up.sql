-- Migration: Create statuses table
-- Statuses represent possible states of pull requests (OPEN, MERGED, etc.)

CREATE TABLE IF NOT EXISTS statuses (
    id SERIAL PRIMARY KEY,
    status VARCHAR(20) NOT NULL UNIQUE
);

-- Insert default statuses
INSERT INTO statuses (status) VALUES ('OPEN'), ('MERGED')
ON CONFLICT (status) DO NOTHING;

