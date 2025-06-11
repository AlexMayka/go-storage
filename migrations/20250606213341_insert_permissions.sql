-- +goose Up
-- +goose StatementBegin
INSERT INTO permissions (id, name)
VALUES
    ('00000000-0000-0000-0000-000000000001', 'company:read'),
    ('00000000-0000-0000-0000-000000000002', 'company:create'),
    ('00000000-0000-0000-0000-000000000003', 'company:update'),
    ('00000000-0000-0000-0000-000000000004', 'company:delete'),
    ('00000000-0000-0000-0000-000000000005', 'user:create');

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DELETE FROM permissions;
-- +goose StatementEnd
