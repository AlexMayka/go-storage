-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS upload_chunks (
    upload_id VARCHAR(255) NOT NULL,
    chunk_index INTEGER NOT NULL,
    size BIGINT NOT NULL,
    etag VARCHAR(255) NOT NULL,
    uploaded BOOLEAN DEFAULT false,
    uploaded_at TIMESTAMP,
    retries INTEGER DEFAULT 0,
    PRIMARY KEY (upload_id, chunk_index),
    FOREIGN KEY (upload_id) REFERENCES chunked_uploads(id) ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS idx_upload_chunks_upload_id ON upload_chunks(upload_id);
CREATE INDEX IF NOT EXISTS idx_upload_chunks_uploaded ON upload_chunks(uploaded);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS upload_chunks;
-- +goose StatementEnd