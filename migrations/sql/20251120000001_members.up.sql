CREATE TABLE IF NOT EXISTS members (
    id SERIAL PRIMARY KEY,
    uuid UUID NOT NULL UNIQUE,
    name VARCHAR(100) NOT NULL,
    is_active BOOLEAN NOT NULL DEFAULT true
);

CREATE INDEX IF NOT EXISTS idx_members_uuid ON members(uuid);

CREATE INDEX IF NOT EXISTS idx_members_is_active ON members(is_active);

