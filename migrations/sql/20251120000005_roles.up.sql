CREATE TABLE IF NOT EXISTS roles (
    id SERIAL PRIMARY KEY,
    role VARCHAR(20) NOT NULL UNIQUE
);

INSERT INTO roles (role) VALUES ('author'), ('reviewer'), ('approver')
ON CONFLICT (role) DO NOTHING;

