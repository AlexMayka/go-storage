-- +goose Up
-- +goose StatementBegin
CREATE TABLE roles (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name varchar(255) UNIQUE NOT NULL,
    is_default bool default false
);

CREATE INDEX idx_roles_name ON roles(name);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS idx_roles_name;
DROP TABLE IF EXISTS roles;
-- +goose StatementEnd
