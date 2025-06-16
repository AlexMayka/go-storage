-- +goose Up
-- +goose StatementBegin
INSERT INTO roles (id, name, is_default)
VALUES
    ('00000000-0000-0000-0000-000000000001', 'super_admin', false),
    ('00000000-0000-0000-0000-000000000002', 'company_admin', false),
    ('00000000-0000-0000-0000-000000000003', 'user', true);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DELETE FROM roles;
-- +goose StatementEnd
