-- +goose Up
-- +goose StatementBegin
CREATE TABLE companies (
   id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
   storage_path VARCHAR(255) NOT NULL,
   name VARCHAR(255) NOT NULL,
   description TEXT NOT NULL,
   created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
   update_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
   is_active BOOL NOT NULL DEFAULT true
);

CREATE INDEX idx_companies_name ON companies (name);
CREATE UNIQUE INDEX idx_companies_storage_path ON companies (storage_path);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS idx_companies_name;
DROP INDEX IF EXISTS idx_companies_storage_path;
DROP TABLE IF EXISTS companies;
-- +goose StatementEnd
