-- +goose Up
-- +goose StatementBegin
CREATE TABLE permissions (
     id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
     name VARCHAR(255) UNIQUE NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS permissions;
-- +goose StatementEnd
