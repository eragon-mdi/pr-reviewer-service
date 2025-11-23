-- Migration: Create members table
-- Members represent users in the system with unique UUID identifiers

CREATE TABLE IF NOT EXISTS members (
    id SERIAL PRIMARY KEY,
    uuid UUID NOT NULL UNIQUE,
    name VARCHAR(100) NOT NULL,
    is_active BOOLEAN NOT NULL DEFAULT true
);

-- Index for fast lookup by UUID (most common query pattern)
CREATE INDEX IF NOT EXISTS idx_members_uuid ON members(uuid);

-- Index for filtering active members
CREATE INDEX IF NOT EXISTS idx_members_is_active ON members(is_active);

