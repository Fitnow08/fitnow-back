-- +goose Up
CREATE TABLE IF NOT EXISTS roles (
    id          BIGINT PRIMARY KEY,
    name        TEXT NOT NULL UNIQUE,
    created_at  TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at  TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    version     BIGINT DEFAULT  1 NOT NULL
);

CREATE TRIGGER update_roles_updated_at BEFORE UPDATE
    ON roles FOR EACH ROW
EXECUTE PROCEDURE update_updated_at_column();

INSERT INTO roles (id, name) VALUES
    (1, 'user'),
    (2, 'admin'),
    (3, 'trainer')
ON CONFLICT (id) DO NOTHING;

ALTER TABLE users
    ADD COLUMN IF NOT EXISTS role_id BIGINT NOT NULL DEFAULT 1
        REFERENCES roles(id);

-- +goose Down
ALTER TABLE users DROP COLUMN IF EXISTS role_id;
DROP TRIGGER IF EXISTS update_roles_updated_at ON roles;
DROP TABLE IF EXISTS roles;
