-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS chunked_uploads (
    id VARCHAR(255) PRIMARY KEY,
    file_name VARCHAR(255) NOT NULL,
    total_size BIGINT NOT NULL,
    chunk_size BIGINT NOT NULL,
    total_chunks INTEGER NOT NULL,
    uploaded_chunks INTEGER DEFAULT 0,
    uploaded_size BIGINT DEFAULT 0,
    status VARCHAR(50) NOT NULL DEFAULT 'active',
    company_id UUID NOT NULL,
    user_created UUID NOT NULL,
    parent_path VARCHAR(1000) NOT NULL,
    target_path VARCHAR(1000) NOT NULL,
    mime_type VARCHAR(255),
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    expires_at TIMESTAMP NOT NULL,
    FOREIGN KEY (company_id) REFERENCES companies(id) ON DELETE CASCADE,
    FOREIGN KEY (user_created) REFERENCES users(id) ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS idx_chunked_uploads_company ON chunked_uploads(company_id);
CREATE INDEX IF NOT EXISTS idx_chunked_uploads_status ON chunked_uploads(status);
CREATE INDEX IF NOT EXISTS idx_chunked_uploads_expires ON chunked_uploads(expires_at);
CREATE INDEX IF NOT EXISTS idx_chunked_uploads_user ON chunked_uploads(user_created);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS chunked_uploads;
-- +goose StatementEnd