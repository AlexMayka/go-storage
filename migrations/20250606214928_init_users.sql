-- +goose Up
-- +goose StatementBegin
CREATE TABLE users (
   id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
   first_name VARCHAR(255) NOT NULL,
   second_name VARCHAR(255),
   last_name VARCHAR(255) NOT NULL,
   username VARCHAR(255) NOT NULL,
   email VARCHAR(255) UNIQUE NOT NULL,
   phone VARCHAR(255) NOT NULL,
   password VARCHAR(255) NOT NULL,
   company_id UUID,
   role_id UUID NOT NULL,
   last_login TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
   created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
   updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
   is_active BOOL NOT NULL DEFAULT true,
   FOREIGN KEY (company_id) REFERENCES companies(id) ON DELETE CASCADE,
   FOREIGN KEY (role_id) REFERENCES roles(id) ON DELETE RESTRICT ON UPDATE CASCADE
);

CREATE UNIQUE INDEX idx_users_email ON users(email);
CREATE INDEX idx_users_company_id ON users(company_id);
CREATE INDEX idx_users_role_id ON users(role_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS idx_users_email;
DROP INDEX IF EXISTS idx_users_company_id;
DROP INDEX IF EXISTS idx_users_role_id;
DROP TABLE IF EXISTS users;
-- +goose StatementEnd
