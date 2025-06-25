package rpFiles

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"go-storage/internal/domain"
)

func setupMockDB(t *testing.T) (*sql.DB, sqlmock.Sqlmock, *RepositoryFiles) {
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

func TestCreateFile_Success(t *testing.T) {
	db, mock, repo := setupMockDB(t)
	defer db.Close()

	path, _ := domain.NewPath("/test.txt")
	file := &domain.File{
		ID:           "file-id",
		Name:         "test.txt",
		Type:         domain.FileTypeFile,
		FullPath:     path,
		ParentID:     nil,
		CompanyId:    "company-id",
		UserCreateID: "user-id",
		MimeType:     stringPtr("text/plain"),
		Size:         int64Ptr(1024),
		Hash:         stringPtr("hash123"),
		StoragePath:  stringPtr("storage/path"),
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
		IsActive:     true,
	}

	mock.ExpectExec(`INSERT INTO files`).
		WithArgs(
			file.ID, file.Name, file.Type, file.FullPath.String(), file.ParentID, file.CompanyId, file.UserCreateID,
			file.MimeType, file.Size, file.Hash, file.StoragePath,
			file.CreatedAt, file.UpdatedAt, file.IsActive,
		).
		WillReturnResult(sqlmock.NewResult(1, 1))

	result, err := repo.CreateFile(context.Background(), file)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, file.ID, result.ID)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestCreateFile_UniqueConstraintViolation(t *testing.T) {
	db, mock, repo := setupMockDB(t)
	defer db.Close()

	path, _ := domain.NewPath("/test.txt")
	file := &domain.File{
		ID:           "file-id",
		Name:         "test.txt",
		Type:         domain.FileTypeFile,
		FullPath:     path,
		CompanyId:    "company-id",
		UserCreateID: "user-id",
		IsActive:     true,
	}

	mock.ExpectExec(`INSERT INTO files`).
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(),
			sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(),
			sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnError(&mockError{message: "idx_unique_name_in_folder constraint violated"})

	result, err := repo.CreateFile(context.Background(), file)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "file with this name already exists")
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetFile_Success(t *testing.T) {
	db, mock, repo := setupMockDB(t)
	defer db.Close()

	fileID := "file-id"
	companyID := "company-id"
	path := "/test.txt"

	rows := sqlmock.NewRows([]string{
		"id", "name", "type", "full_path", "parent_id", "company_id", "user_created",
		"mime_type", "size", "hash", "storage_path",
		"created_at", "updated_at", "is_active",
	}).AddRow(
		fileID, "test.txt", domain.FileTypeFile, path, nil, companyID, "user-id",
		"text/plain", 1024, "hash123", "storage/path",
		time.Now(), time.Now(), true,
	)

	mock.ExpectQuery(`SELECT .+ FROM files WHERE id = \$1 AND company_id = \$2 AND is_active = true`).
		WithArgs(fileID, companyID).
		WillReturnRows(rows)

	result, err := repo.GetFile(context.Background(), companyID, fileID)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, fileID, result.ID)
	assert.Equal(t, "test.txt", result.Name)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetFile_NotFound(t *testing.T) {
	db, mock, repo := setupMockDB(t)
	defer db.Close()

	fileID := "nonexistent-id"
	companyID := "company-id"

	mock.ExpectQuery(`SELECT .+ FROM files WHERE id = \$1 AND company_id = \$2 AND is_active = true`).
		WithArgs(fileID, companyID).
		WillReturnError(sql.ErrNoRows)

	result, err := repo.GetFile(context.Background(), companyID, fileID)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "file not found")
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetFileByPath_Success(t *testing.T) {
	db, mock, repo := setupMockDB(t)
	defer db.Close()

	path, _ := domain.NewPath("/test.txt")
	companyID := "company-id"

	rows := sqlmock.NewRows([]string{
		"id", "name", "type", "full_path", "parent_id", "company_id", "user_created",
		"mime_type", "size", "hash", "storage_path",
		"created_at", "updated_at", "is_active",
	}).AddRow(
		"file-id", "test.txt", domain.FileTypeFile, "/test.txt", nil, companyID, "user-id",
		"text/plain", 1024, "hash123", "storage/path",
		time.Now(), time.Now(), true,
	)

	mock.ExpectQuery(`SELECT .+ FROM files WHERE full_path = \$1 AND company_id = \$2 AND is_active = true`).
		WithArgs(path.String(), companyID).
		WillReturnRows(rows)

	result, err := repo.GetFileByPath(context.Background(), companyID, &path)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "test.txt", result.Name)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetFolderContents_Success(t *testing.T) {
	db, mock, repo := setupMockDB(t)
	defer db.Close()

	path, _ := domain.NewPath("/")
	companyID := "company-id"

	rows := sqlmock.NewRows([]string{
		"id", "name", "type", "full_path", "parent_id", "company_id", "user_created",
		"mime_type", "size", "hash", "storage_path",
		"created_at", "updated_at", "is_active",
	}).AddRow(
		"file1", "file1.txt", domain.FileTypeFile, "/file1.txt", "parent-id", companyID, "user-id",
		"text/plain", 1024, "hash1", "storage/path1",
		time.Now(), time.Now(), true,
	).AddRow(
		"folder1", "folder1", domain.FileTypeFolder, "/folder1", "parent-id", companyID, "user-id",
		nil, nil, nil, nil,
		time.Now(), time.Now(), true,
	)

	parentRows := sqlmock.NewRows([]string{
		"id", "name", "type", "full_path", "parent_id", "company_id", "user_created",
		"mime_type", "size", "hash", "storage_path",
		"created_at", "updated_at", "is_active",
	}).AddRow(
		"parent-id", "root", domain.FileTypeFolder, "/", nil, companyID, "user-id",
		nil, nil, nil, nil,
		time.Now(), time.Now(), true,
	)

	mock.ExpectQuery(`SELECT .+ FROM files WHERE full_path = \$1 AND company_id = \$2 AND is_active = true`).
		WithArgs(path.String(), companyID).
		WillReturnRows(parentRows)

	mock.ExpectQuery(`SELECT .+ FROM files WHERE parent_id = \$1 AND company_id = \$2 AND type = \$3 AND is_active = true ORDER BY name ASC`).
		WithArgs("parent-id", companyID, domain.FileTypeFile).
		WillReturnRows(rows)

	fileType := domain.FileTypeFile
	result, err := repo.GetFolderContents(context.Background(), companyID, &path, &fileType)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Len(t, result, 2)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUpdateFile_Success(t *testing.T) {
	db, mock, repo := setupMockDB(t)
	defer db.Close()

	path, _ := domain.NewPath("/updated.txt")
	file := &domain.File{
		ID:          "file-id",
		Name:        "updated.txt",
		FullPath:    path,
		ParentID:    stringPtr("parent-id"),
		MimeType:    stringPtr("text/plain"),
		Size:        int64Ptr(2048),
		Hash:        stringPtr("newhash"),
		StoragePath: stringPtr("new/storage/path"),
		CompanyId:   "company-id",
	}

	mock.ExpectExec(`UPDATE files SET`).
		WithArgs(
			file.ID, file.Name, file.FullPath.String(), file.ParentID, file.MimeType,
			file.Size, file.Hash, file.StoragePath, sqlmock.AnyArg(), file.CompanyId,
		).
		WillReturnResult(sqlmock.NewResult(1, 1))

	result, err := repo.UpdateFile(context.Background(), file)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, file.ID, result.ID)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestRenameFile_Success(t *testing.T) {
	db, mock, repo := setupMockDB(t)
	defer db.Close()

	fileID := "file-id"
	companyID := "company-id"
	newName := "renamed.txt"

	rows := sqlmock.NewRows([]string{
		"id", "name", "type", "full_path", "parent_id", "company_id", "user_created",
		"mime_type", "size", "hash", "storage_path",
		"created_at", "updated_at", "is_active",
	}).AddRow(
		fileID, "old.txt", domain.FileTypeFile, "/old.txt", nil, companyID, "user-id",
		"text/plain", 1024, "hash123", "storage/path",
		time.Now(), time.Now(), true,
	)

	mock.ExpectQuery(`SELECT .+ FROM files WHERE id = \$1 AND company_id = \$2 AND is_active = true`).
		WithArgs(fileID, companyID).
		WillReturnRows(rows)

	mock.ExpectExec(`UPDATE files SET name = \$2, full_path = \$3, updated_at = \$4 WHERE id = \$1 AND company_id = \$5 AND is_active = true`).
		WithArgs(fileID, newName, "/renamed.txt", sqlmock.AnyArg(), companyID).
		WillReturnResult(sqlmock.NewResult(1, 1))

	result, err := repo.RenameFile(context.Background(), companyID, fileID, newName)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, newName, result.Name)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestDeleteFile_Success(t *testing.T) {
	db, mock, repo := setupMockDB(t)
	defer db.Close()

	fileID := "file-id"
	companyID := "company-id"

	mock.ExpectExec(`UPDATE files SET is_active = false, updated_at = \$3 WHERE id = \$1 AND company_id = \$2 AND is_active = true`).
		WithArgs(fileID, companyID, sqlmock.AnyArg()).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err := repo.DeleteFile(context.Background(), companyID, fileID)

	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestCreateFolder_Success(t *testing.T) {
	db, mock, repo := setupMockDB(t)
	defer db.Close()

	path, _ := domain.NewPath("/test-folder")
	folder := &domain.File{
		ID:           "folder-id",
		Name:         "test-folder",
		Type:         domain.FileTypeFolder,
		FullPath:     path,
		ParentID:     nil,
		CompanyId:    "company-id",
		UserCreateID: "user-id",
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
		IsActive:     true,
	}

	mock.ExpectExec(`INSERT INTO files`).
		WithArgs(
			folder.ID, folder.Name, folder.Type, folder.FullPath.String(), folder.ParentID, folder.CompanyId, folder.UserCreateID,
			folder.CreatedAt, folder.UpdatedAt, folder.IsActive,
		).
		WillReturnResult(sqlmock.NewResult(1, 1))

	result, err := repo.CreateFolder(context.Background(), folder)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, folder.ID, result.ID)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetFolder_Success(t *testing.T) {
	db, mock, repo := setupMockDB(t)
	defer db.Close()

	path, _ := domain.NewPath("/test-folder")
	companyID := "company-id"

	rows := sqlmock.NewRows([]string{
		"id", "name", "type", "full_path", "parent_id", "company_id", "user_created",
		"mime_type", "size", "hash", "storage_path",
		"created_at", "updated_at", "is_active",
	}).AddRow(
		"folder-id", "test-folder", domain.FileTypeFolder, "/test-folder", nil, companyID, "user-id",
		nil, nil, nil, nil,
		time.Now(), time.Now(), true,
	)

	mock.ExpectQuery(`SELECT .+ FROM files WHERE full_path = \$1 AND company_id = \$2 AND type = 'folder' AND is_active = true`).
		WithArgs(path.String(), companyID).
		WillReturnRows(rows)

	result, err := repo.GetFolder(context.Background(), companyID, &path)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "test-folder", result.Name)
	assert.Equal(t, domain.FileTypeFolder, result.Type)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestMoveFolder_Success(t *testing.T) {
	db, mock, repo := setupMockDB(t)
	defer db.Close()

	oldPath, _ := domain.NewPath("/old-folder")
	newPath, _ := domain.NewPath("/new-folder")
	companyID := "company-id"

	rows := sqlmock.NewRows([]string{
		"id", "name", "type", "full_path", "parent_id", "company_id", "user_created",
		"mime_type", "size", "hash", "storage_path",
		"created_at", "updated_at", "is_active",
	}).AddRow(
		"folder-id", "old-folder", domain.FileTypeFolder, "/old-folder", nil, companyID, "user-id",
		nil, nil, nil, nil,
		time.Now(), time.Now(), true,
	)

	mock.ExpectQuery(`SELECT .+ FROM files WHERE full_path = \$1 AND company_id = \$2 AND type = 'folder' AND is_active = true`).
		WithArgs(oldPath.String(), companyID).
		WillReturnRows(rows)

	mock.ExpectExec(`UPDATE files SET full_path = REPLACE\(full_path, \$1, \$2\), updated_at = \$5 WHERE full_path LIKE \$3 AND company_id = \$4 AND is_active = true`).
		WithArgs(oldPath.String(), newPath.String(), "/old-folder%", companyID, sqlmock.AnyArg()).
		WillReturnResult(sqlmock.NewResult(1, 1))

	result, err := repo.MoveFolder(context.Background(), companyID, &oldPath, &newPath)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, newPath.String(), result.String())
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestDeleteFolder_Success(t *testing.T) {
	db, mock, repo := setupMockDB(t)
	defer db.Close()

	folderPath, _ := domain.NewPath("/test-folder")
	companyID := "company-id"

	mock.ExpectExec(`UPDATE files SET is_active = false, updated_at = \$3 WHERE full_path = \$1 AND company_id = \$2 AND is_active = true`).
		WithArgs(folderPath.String(), companyID, sqlmock.AnyArg()).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err := repo.DeleteFolder(context.Background(), companyID, &folderPath)

	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func stringPtr(s string) *string {
	return &s
}

func int64Ptr(i int64) *int64 {
	return &i
}

type mockError struct {
	message string
}

func (e *mockError) Error() string {
	return e.message
}
