-- +goose Up
-- +goose StatementBegin
INSERT INTO permissions (id, name)
VALUES
    ('00000000-0000-0000-0000-000000000001', 'company:read'),
    ('00000000-0000-0000-0000-000000000002', 'company:create'),
    ('00000000-0000-0000-0000-000000000003', 'company:update'),
    ('00000000-0000-0000-0000-000000000004', 'company:delete'),

    ('00000000-0000-0000-0000-000000000005', 'user:create'),
    ('00000000-0000-0000-0000-000000000006', 'user:read'),
    ('00000000-0000-0000-0000-000000000007', 'user:update'),
    ('00000000-0000-0000-0000-000000000008', 'user:delete'),
    ('00000000-0000-0000-0000-000000000009', 'user:read_company'),
    ('00000000-0000-0000-0000-000000000010', 'user:manage_company'),

    ('00000000-0000-0000-0000-000000000011', 'file:create'),
    ('00000000-0000-0000-0000-000000000012', 'file:read'),
    ('00000000-0000-0000-0000-000000000013', 'file:update'),
    ('00000000-0000-0000-0000-000000000014', 'file:delete');

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DELETE FROM permissions;
-- +goose StatementEnd
