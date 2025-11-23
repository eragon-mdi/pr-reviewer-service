-- Migration: Create pull_requests table
-- Pull requests represent code review requests with status tracking and optimistic locking

CREATE TABLE IF NOT EXISTS pull_requests (
    id SERIAL PRIMARY KEY,
    uuid UUID NOT NULL UNIQUE,
    title VARCHAR(255) NOT NULL,
    author_id INT NOT NULL,
    status_id INT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    merged_at TIMESTAMPTZ,
    version INT NOT NULL DEFAULT 1,
    CONSTRAINT fk_pull_requests_author 
        FOREIGN KEY (author_id) REFERENCES members(id) ON DELETE RESTRICT,
    CONSTRAINT fk_pull_requests_status 
        FOREIGN KEY (status_id) REFERENCES statuses(id) ON DELETE RESTRICT
);

-- Index for fast lookup by UUID (most common query pattern)
CREATE INDEX IF NOT EXISTS idx_pull_requests_uuid ON pull_requests(uuid);

-- Index for filtering by author
CREATE INDEX IF NOT EXISTS idx_pull_requests_author_id ON pull_requests(author_id);

-- Index for filtering by status
CREATE INDEX IF NOT EXISTS idx_pull_requests_status_id ON pull_requests(status_id);

-- Index for filtering by created_at (for sorting)
CREATE INDEX IF NOT EXISTS idx_pull_requests_created_at ON pull_requests(created_at);

