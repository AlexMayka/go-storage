package rpChunkedUpload

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"go-storage/internal/domain"
)

func setupMockDB(t *testing.T) (*sql.DB, sqlmock.Sqlmock, *RepositoryChunkedUpload) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create mock database: %v", err)
	}

	repo := NewRepository(db)
	return db, mock, repo
}

func TestNewRepository(t *testing.T) {
	db, _, _ := setupMockDB(t)
	defer db.Close()

	repo := NewRepository(db)
	assert.NotNil(t, repo)
	assert.Equal(t, db, repo.db)
}

func TestCreateChunkedUpload_Success(t *testing.T) {
	db, mock, repo := setupMockDB(t)
	defer db.Close()

	parentPath, _ := domain.NewPath("/")
	targetPath, _ := domain.NewPath("/test.zip")

	upload := &domain.ChunkedUpload{
		ID:             "upload-id",
		FileName:       "test.zip",
		TotalSize:      1024000,
		ChunkSize:      5242880,
		TotalChunks:    1,
		UploadedChunks: 0,
		UploadedSize:   0,
		Status:         "active",
		CompanyID:      "company-id",
		UserCreateID:   "user-id",
		ParentPath:     parentPath,
		TargetPath:     targetPath,
		MimeType:       "application/zip",
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
		ExpiresAt:      time.Now().Add(24 * time.Hour),
		Chunks:         make(map[int]*domain.ChunkInfo),
	}

	mock.ExpectExec(`INSERT INTO chunked_uploads`).
		WithArgs(
			upload.ID, upload.FileName, upload.TotalSize, upload.ChunkSize, upload.TotalChunks,
			upload.UploadedChunks, upload.UploadedSize, upload.Status, upload.CompanyID, upload.UserCreateID,
			upload.ParentPath.String(), upload.TargetPath.String(), upload.MimeType,
			upload.CreatedAt, upload.UpdatedAt, upload.ExpiresAt,
		).
		WillReturnResult(sqlmock.NewResult(1, 1))

	result, err := repo.CreateChunkedUpload(context.Background(), upload)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, upload.ID, result.ID)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestCreateChunkedUpload_DatabaseError(t *testing.T) {
	db, mock, repo := setupMockDB(t)
	defer db.Close()

	parentPath, _ := domain.NewPath("/")
	targetPath, _ := domain.NewPath("/test.zip")

	upload := &domain.ChunkedUpload{
		ID:           "upload-id",
		FileName:     "test.zip",
		CompanyID:    "company-id",
		UserCreateID: "user-id",
		ParentPath:   parentPath,
		TargetPath:   targetPath,
	}

	mock.ExpectExec(`INSERT INTO chunked_uploads`).
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(),
			sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(),
			sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(),
			sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnError(sql.ErrConnDone)

	result, err := repo.CreateChunkedUpload(context.Background(), upload)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "unable to create chunked upload session")
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetChunkedUpload_Success(t *testing.T) {
	db, mock, repo := setupMockDB(t)
	defer db.Close()

	uploadID := "upload-id"
	companyID := "company-id"

	rows := sqlmock.NewRows([]string{
		"id", "file_name", "total_size", "chunk_size", "total_chunks",
		"uploaded_chunks", "uploaded_size", "status", "company_id", "user_created",
		"parent_path", "target_path", "mime_type",
		"created_at", "updated_at", "expires_at",
	}).AddRow(
		uploadID, "test.zip", 1024000, 5242880, 1,
		0, 0, "active", companyID, "user-id",
		"/", "/test.zip", "application/zip",
		time.Now(), time.Now(), time.Now().Add(24*time.Hour),
	)

	mock.ExpectQuery(`SELECT .+ FROM chunked_uploads WHERE id = \$1 AND company_id = \$2`).
		WithArgs(uploadID, companyID).
		WillReturnRows(rows)

	result, err := repo.GetChunkedUpload(context.Background(), companyID, uploadID)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, uploadID, result.ID)
	assert.Equal(t, "test.zip", result.FileName)
	assert.NotNil(t, result.Chunks)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetChunkedUpload_NotFound(t *testing.T) {
	db, mock, repo := setupMockDB(t)
	defer db.Close()

	uploadID := "nonexistent-id"
	companyID := "company-id"

	mock.ExpectQuery(`SELECT .+ FROM chunked_uploads WHERE id = \$1 AND company_id = \$2`).
		WithArgs(uploadID, companyID).
		WillReturnError(sql.ErrNoRows)

	result, err := repo.GetChunkedUpload(context.Background(), companyID, uploadID)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "chunked upload session not found")
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUpdateChunkedUpload_Success(t *testing.T) {
	db, mock, repo := setupMockDB(t)
	defer db.Close()

	parentPath, _ := domain.NewPath("/")
	targetPath, _ := domain.NewPath("/test.zip")

	upload := &domain.ChunkedUpload{
		ID:             "upload-id",
		UploadedChunks: 1,
		UploadedSize:   5242880,
		Status:         "in_progress",
		CompanyID:      "company-id",
		ParentPath:     parentPath,
		TargetPath:     targetPath,
	}

	mock.ExpectExec(`UPDATE chunked_uploads SET`).
		WithArgs(
			upload.ID, upload.UploadedChunks, upload.UploadedSize, upload.Status, sqlmock.AnyArg(), upload.CompanyID,
		).
		WillReturnResult(sqlmock.NewResult(1, 1))

	result, err := repo.UpdateChunkedUpload(context.Background(), upload)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, upload.ID, result.ID)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestDeleteChunkedUpload_Success(t *testing.T) {
	db, mock, repo := setupMockDB(t)
	defer db.Close()

	uploadID := "upload-id"
	companyID := "company-id"

	mock.ExpectExec(`DELETE FROM chunked_uploads WHERE id = \$1 AND company_id = \$2`).
		WithArgs(uploadID, companyID).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err := repo.DeleteChunkedUpload(context.Background(), companyID, uploadID)

	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestAddChunk_Success(t *testing.T) {
	db, mock, repo := setupMockDB(t)
	defer db.Close()

	uploadID := "upload-id"
	chunkIndex := 0
	etag := "etag-123"
	size := int64(5242880)

	mock.ExpectExec(`INSERT INTO upload_chunks`).
		WithArgs(
			uploadID, chunkIndex, size, etag, true, sqlmock.AnyArg(), 0,
		).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err := repo.AddChunk(context.Background(), uploadID, chunkIndex, etag, size)

	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestAddChunk_DatabaseError(t *testing.T) {
	db, mock, repo := setupMockDB(t)
	defer db.Close()

	uploadID := "upload-id"
	chunkIndex := 0
	etag := "etag-123"
	size := int64(5242880)

	mock.ExpectExec(`INSERT INTO upload_chunks`).
		WithArgs(
			uploadID, chunkIndex, size, etag, true, sqlmock.AnyArg(), 0,
		).
		WillReturnError(sql.ErrConnDone)

	err := repo.AddChunk(context.Background(), uploadID, chunkIndex, etag, size)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "unable to add chunk info")
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetUploadProgress_Success(t *testing.T) {
	db, mock, repo := setupMockDB(t)
	defer db.Close()

	uploadID := "upload-id"
	chunksJSON := `[{"index":0,"size":5242880,"etag":"etag-123","uploaded":true,"uploaded_at":"2023-01-01T00:00:00Z","retries":0}]`

	rows := sqlmock.NewRows([]string{
		"id", "file_name", "total_size", "chunk_size", "total_chunks",
		"uploaded_chunks", "uploaded_size", "status", "company_id", "user_created",
		"parent_path", "target_path", "mime_type",
		"created_at", "updated_at", "expires_at", "chunks",
	}).AddRow(
		uploadID, "test.zip", 5242880, 5242880, 1,
		1, 5242880, "completed", "company-id", "user-id",
		"/", "/test.zip", "application/zip",
		time.Now(), time.Now(), time.Now().Add(24*time.Hour), chunksJSON,
	)

	mock.ExpectQuery(`SELECT cu\.id, cu\.file_name, cu\.total_size, cu\.chunk_size, cu\.total_chunks,`).
		WithArgs(uploadID).
		WillReturnRows(rows)

	result, err := repo.GetUploadProgress(context.Background(), uploadID)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, uploadID, result.ID)
	assert.NotNil(t, result.Chunks)
	assert.Len(t, result.Chunks, 1)
	assert.Contains(t, result.Chunks, 0)
	assert.Equal(t, "etag-123", result.Chunks[0].ETag)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetUploadProgress_EmptyChunks(t *testing.T) {
	db, mock, repo := setupMockDB(t)
	defer db.Close()

	uploadID := "upload-id"
	chunksJSON := "[]"

	rows := sqlmock.NewRows([]string{
		"id", "file_name", "total_size", "chunk_size", "total_chunks",
		"uploaded_chunks", "uploaded_size", "status", "company_id", "user_created",
		"parent_path", "target_path", "mime_type",
		"created_at", "updated_at", "expires_at", "chunks",
	}).AddRow(
		uploadID, "test.zip", 5242880, 5242880, 1,
		0, 0, "active", "company-id", "user-id",
		"/", "/test.zip", "application/zip",
		time.Now(), time.Now(), time.Now().Add(24*time.Hour), chunksJSON,
	)

	mock.ExpectQuery(`SELECT cu\.id, cu\.file_name, cu\.total_size, cu\.chunk_size, cu\.total_chunks,`).
		WithArgs(uploadID).
		WillReturnRows(rows)

	result, err := repo.GetUploadProgress(context.Background(), uploadID)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, uploadID, result.ID)
	assert.NotNil(t, result.Chunks)
	assert.Len(t, result.Chunks, 0)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetUploadProgress_NotFound(t *testing.T) {
	db, mock, repo := setupMockDB(t)
	defer db.Close()

	uploadID := "nonexistent-id"

	mock.ExpectQuery(`SELECT cu\.id, cu\.file_name, cu\.total_size, cu\.chunk_size, cu\.total_chunks,`).
		WithArgs(uploadID).
		WillReturnError(sql.ErrNoRows)

	result, err := repo.GetUploadProgress(context.Background(), uploadID)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "chunked upload session not found")
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestCleanupExpiredUploads_Success(t *testing.T) {
	db, mock, repo := setupMockDB(t)
	defer db.Close()

	// Mock getting expired uploads
	expiredRows := sqlmock.NewRows([]string{"id", "company_id", "file_name"}).
		AddRow("expired-1", "company-1", "expired1.zip").
		AddRow("expired-2", "company-1", "expired2.zip")

	mock.ExpectQuery(`SELECT id, company_id, file_name FROM chunked_uploads WHERE expires_at < NOW\(\)`).
		WillReturnRows(expiredRows)

	// Mock deletion
	mock.ExpectExec(`DELETE FROM chunked_uploads WHERE expires_at < NOW\(\)`).
		WillReturnResult(sqlmock.NewResult(2, 2))

	err := repo.CleanupExpiredUploads(context.Background())

	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestCleanupExpiredUploads_QueryError(t *testing.T) {
	db, mock, repo := setupMockDB(t)
	defer db.Close()

	mock.ExpectQuery(`SELECT id, company_id, file_name FROM chunked_uploads WHERE expires_at < NOW\(\)`).
		WillReturnError(sql.ErrConnDone)

	err := repo.CleanupExpiredUploads(context.Background())

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "unable to get expired uploads")
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestCleanupExpiredUploads_DeleteError(t *testing.T) {
	db, mock, repo := setupMockDB(t)
	defer db.Close()

	// Mock getting expired uploads
	expiredRows := sqlmock.NewRows([]string{"id", "company_id", "file_name"}).
		AddRow("expired-1", "company-1", "expired1.zip")

	mock.ExpectQuery(`SELECT id, company_id, file_name FROM chunked_uploads WHERE expires_at < NOW\(\)`).
		WillReturnRows(expiredRows)

	// Mock deletion error
	mock.ExpectExec(`DELETE FROM chunked_uploads WHERE expires_at < NOW\(\)`).
		WillReturnError(sql.ErrConnDone)

	err := repo.CleanupExpiredUploads(context.Background())

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "unable to cleanup expired uploads")
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestEnsureChunksTable_Success(t *testing.T) {
	db, mock, repo := setupMockDB(t)
	defer db.Close()

	mock.ExpectExec(`CREATE TABLE IF NOT EXISTS upload_chunks`).
		WillReturnResult(sqlmock.NewResult(0, 0))

	err := repo.EnsureChunksTable(context.Background())

	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestEnsureChunkedUploadsTable_Success(t *testing.T) {
	db, mock, repo := setupMockDB(t)
	defer db.Close()

	mock.ExpectExec(`CREATE TABLE IF NOT EXISTS chunked_uploads`).
		WillReturnResult(sqlmock.NewResult(0, 0))

	err := repo.EnsureChunkedUploadsTable(context.Background())

	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}
