-- +goose Up
CREATE TABLE roles (
    role_id SERIAL PRIMARY KEY,
    name TEXT UNIQUE,
    description TEXT
);

INSERT INTO roles (name, description)
VALUES 
  ('user', 'just user'),
  ('manager', 'just manager'),
  ('admin', 'just admin');

CREATE TABLE permissions (
    permission_id SERIAL PRIMARY KEY,
    code TEXT UNIQUE,
    description TEXT
);

CREATE TABLE roles_permissions (
    permission_id INT REFERENCES roles(role_id) ON DELETE CASCADE,
    role_id INT REFERENCES roles(role_id) ON DELETE CASCADE
);

CREATE TABLE users (
    user_id UUID PRIMARY KEY,
    email TEXT NOT NULL UNIQUE,
    name TEXT,
    password_hash TEXT NOT NULL,
    role_id INT REFERENCES roles(role_id) ON DELETE SET NULL,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL
);

-- +goose Down
DROP TABLE users;
DROP TABLE roles_permissions;
DROP TABLE permissions;
DROP TABLE roles;
