-- +goose Up
-- +goose StatementBegin
INSERT INTO roles (id, name)
VALUES
    ('00000000-0000-0000-0000-000000000001', 'super_admin'),
    ('00000000-0000-0000-0000-000000000002', 'company_admin'),
    ('00000000-0000-0000-0000-000000000003', 'user');
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DELETE FROM roles;
-- +goose StatementEnd
