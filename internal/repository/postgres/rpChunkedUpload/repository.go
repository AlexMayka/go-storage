package rpChunkedUpload

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"time"

	"go-storage/internal/domain"
	pkgErrors "go-storage/pkg/errors"
)

type RepositoryChunkedUpload struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *RepositoryChunkedUpload {
	return &RepositoryChunkedUpload{db: db}
}

func (r *RepositoryChunkedUpload) CreateChunkedUpload(ctx context.Context, upload *domain.ChunkedUpload) (*domain.ChunkedUpload, error) {
	_, err := r.db.ExecContext(ctx, QueryCreateChunkedUpload,
		upload.ID, upload.FileName, upload.TotalSize, upload.ChunkSize, upload.TotalChunks,
		upload.UploadedChunks, upload.UploadedSize, upload.Status, upload.CompanyID, upload.UserCreateID,
		upload.ParentPath.String(), upload.TargetPath.String(), upload.MimeType,
		upload.CreatedAt, upload.UpdatedAt, upload.ExpiresAt,
	)
	if err != nil {
		return nil, pkgErrors.Database("unable to create chunked upload session")
	}

	return upload, nil
}

func (r *RepositoryChunkedUpload) GetChunkedUpload(ctx context.Context, companyID, uploadID string) (*domain.ChunkedUpload, error) {
	var upload domain.ChunkedUpload
	var parentPathStr, targetPathStr string

	row := r.db.QueryRowContext(ctx, QueryGetChunkedUpload, uploadID, companyID)

	err := row.Scan(
		&upload.ID, &upload.FileName, &upload.TotalSize, &upload.ChunkSize, &upload.TotalChunks,
		&upload.UploadedChunks, &upload.UploadedSize, &upload.Status, &upload.CompanyID, &upload.UserCreateID,
		&parentPathStr, &targetPathStr, &upload.MimeType,
		&upload.CreatedAt, &upload.UpdatedAt, &upload.ExpiresAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, pkgErrors.NotFound("chunked upload session not found")
		}
		return nil, pkgErrors.Database("unable to get chunked upload session")
	}

	// Parse paths
	parentPath, err := domain.NewPath(parentPathStr)
	if err != nil {
		return nil, pkgErrors.Database("invalid parent path")
	}
	upload.ParentPath = parentPath

	targetPath, err := domain.NewPath(targetPathStr)
	if err != nil {
		return nil, pkgErrors.Database("invalid target path")
	}
	upload.TargetPath = targetPath

	// Initialize chunks map
	upload.Chunks = make(map[int]*domain.ChunkInfo)

	return &upload, nil
}

func (r *RepositoryChunkedUpload) UpdateChunkedUpload(ctx context.Context, upload *domain.ChunkedUpload) (*domain.ChunkedUpload, error) {
	upload.UpdatedAt = time.Now()

	_, err := r.db.ExecContext(ctx, QueryUpdateChunkedUpload,
		upload.ID, upload.UploadedChunks, upload.UploadedSize, upload.Status, upload.UpdatedAt, upload.CompanyID,
	)
	if err != nil {
		return nil, pkgErrors.Database("unable to update chunked upload session")
	}

	return upload, nil
}

func (r *RepositoryChunkedUpload) DeleteChunkedUpload(ctx context.Context, companyID, uploadID string) error {
	_, err := r.db.ExecContext(ctx, QueryDeleteChunkedUpload, uploadID, companyID)
	if err != nil {
		return pkgErrors.Database("unable to delete chunked upload session")
	}

	return nil
}

func (r *RepositoryChunkedUpload) AddChunk(ctx context.Context, uploadID string, chunkIndex int, etag string, size int64) error {
	_, err := r.db.ExecContext(ctx, QueryAddChunk,
		uploadID, chunkIndex, size, etag, true, time.Now(), 0,
	)
	if err != nil {
		return pkgErrors.Database("unable to add chunk info")
	}

	return nil
}

func (r *RepositoryChunkedUpload) GetUploadProgress(ctx context.Context, uploadID string) (*domain.ChunkedUpload, error) {
	var upload domain.ChunkedUpload
	var parentPathStr, targetPathStr string
	var chunksJSON string

	row := r.db.QueryRowContext(ctx, QueryGetUploadProgress, uploadID)

	err := row.Scan(
		&upload.ID, &upload.FileName, &upload.TotalSize, &upload.ChunkSize, &upload.TotalChunks,
		&upload.UploadedChunks, &upload.UploadedSize, &upload.Status, &upload.CompanyID, &upload.UserCreateID,
		&parentPathStr, &targetPathStr, &upload.MimeType,
		&upload.CreatedAt, &upload.UpdatedAt, &upload.ExpiresAt,
		&chunksJSON,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, pkgErrors.NotFound("chunked upload session not found")
		}
		return nil, pkgErrors.Database("unable to get upload progress")
	}

	// Parse paths
	parentPath, err := domain.NewPath(parentPathStr)
	if err != nil {
		return nil, pkgErrors.Database("invalid parent path")
	}
	upload.ParentPath = parentPath

	targetPath, err := domain.NewPath(targetPathStr)
	if err != nil {
		return nil, pkgErrors.Database("invalid target path")
	}
	upload.TargetPath = targetPath

	// Parse chunks JSON
	upload.Chunks = make(map[int]*domain.ChunkInfo)

	if chunksJSON != "[]" {
		var chunks []struct {
			Index      int       `json:"index"`
			Size       int64     `json:"size"`
			ETag       string    `json:"etag"`
			Uploaded   bool      `json:"uploaded"`
			UploadedAt time.Time `json:"uploaded_at"`
			Retries    int       `json:"retries"`
		}

		if err := json.Unmarshal([]byte(chunksJSON), &chunks); err != nil {
			return nil, pkgErrors.Database("unable to parse chunks data")
		}

		for _, chunk := range chunks {
			upload.Chunks[chunk.Index] = &domain.ChunkInfo{
				Index:      chunk.Index,
				Size:       chunk.Size,
				ETag:       chunk.ETag,
				Uploaded:   chunk.Uploaded,
				UploadedAt: chunk.UploadedAt,
				Retries:    chunk.Retries,
			}
		}
	}

	return &upload, nil
}

func (r *RepositoryChunkedUpload) CleanupExpiredUploads(ctx context.Context) error {
	// First get expired uploads to clean up storage
	rows, err := r.db.QueryContext(ctx, QueryGetExpiredUploads)
	if err != nil {
		return pkgErrors.Database("unable to get expired uploads")
	}
	defer rows.Close()

	var expiredUploads []struct {
		ID        string
		CompanyID string
		FileName  string
	}

	for rows.Next() {
		var upload struct {
			ID        string
			CompanyID string
			FileName  string
		}
		if err := rows.Scan(&upload.ID, &upload.CompanyID, &upload.FileName); err != nil {
			continue // Skip this one, continue cleanup
		}
		expiredUploads = append(expiredUploads, upload)
	}

	// Delete from database
	_, err = r.db.ExecContext(ctx, QueryCleanupExpiredUploads)
	if err != nil {
		return pkgErrors.Database("unable to cleanup expired uploads")
	}

	return nil
}

// Helper method to check if we need to create the chunks table
func (r *RepositoryChunkedUpload) EnsureChunksTable(ctx context.Context) error {
	createTableSQL := `
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
	`

	_, err := r.db.ExecContext(ctx, createTableSQL)
	return err
}

func (r *RepositoryChunkedUpload) EnsureChunkedUploadsTable(ctx context.Context) error {
	createTableSQL := `
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
		FOREIGN KEY (company_id) REFERENCES companies(id) ON DELETE CASCADE
	);
	
	CREATE INDEX IF NOT EXISTS idx_chunked_uploads_company ON chunked_uploads(company_id);
	CREATE INDEX IF NOT EXISTS idx_chunked_uploads_status ON chunked_uploads(status);
	CREATE INDEX IF NOT EXISTS idx_chunked_uploads_expires ON chunked_uploads(expires_at);
	CREATE INDEX IF NOT EXISTS idx_chunked_uploads_user ON chunked_uploads(user_created);
	`

	_, err := r.db.ExecContext(ctx, createTableSQL)
	return err
}
