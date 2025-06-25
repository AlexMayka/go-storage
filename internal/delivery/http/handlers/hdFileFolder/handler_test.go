package hdFileFolder

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go-storage/internal/domain"
)

type mockUseCaseFileFolder struct {
	mock.Mock
}

func (m *mockUseCaseFileFolder) CreateFolder(ctx context.Context, folder *domain.File) (*domain.File, error) {
	args := m.Called(ctx, folder)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.File), args.Error(1)
}

func (m *mockUseCaseFileFolder) GetFolderContents(ctx context.Context, companyID string, path *domain.Path, fileType *domain.FileType) ([]*domain.File, error) {
	args := m.Called(ctx, companyID, path, fileType)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.File), args.Error(1)
}

func (m *mockUseCaseFileFolder) MoveFolder(ctx context.Context, companyID string, folderPath *domain.Path, newPath *domain.Path) (*domain.Path, error) {
	args := m.Called(ctx, companyID, folderPath, newPath)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Path), args.Error(1)
}

func (m *mockUseCaseFileFolder) DeleteFolder(ctx context.Context, companyID string, folderPath *domain.Path) error {
	args := m.Called(ctx, companyID, folderPath)
	return args.Error(0)
}

func (m *mockUseCaseFileFolder) UploadFile(ctx context.Context, companyID, userID string, parentPath *domain.Path, filename string, size int64, reader io.Reader) (*domain.File, error) {
	args := m.Called(ctx, companyID, userID, parentPath, filename, size, reader)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.File), args.Error(1)
}

func (m *mockUseCaseFileFolder) DownloadFile(ctx context.Context, companyID, fileID string) (io.ReadCloser, *domain.File, error) {
	args := m.Called(ctx, companyID, fileID)
	if args.Get(0) == nil {
		return nil, nil, args.Error(2)
	}
	return args.Get(0).(io.ReadCloser), args.Get(1).(*domain.File), args.Error(2)
}

func (m *mockUseCaseFileFolder) GetFileInfo(ctx context.Context, companyID, fileID string) (*domain.File, error) {
	args := m.Called(ctx, companyID, fileID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.File), args.Error(1)
}

func (m *mockUseCaseFileFolder) RenameFile(ctx context.Context, companyID, fileID, newName string) (*domain.File, error) {
	args := m.Called(ctx, companyID, fileID, newName)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.File), args.Error(1)
}

func (m *mockUseCaseFileFolder) MoveFile(ctx context.Context, companyID, fileID string, newParentPath *domain.Path) (*domain.File, error) {
	args := m.Called(ctx, companyID, fileID, newParentPath)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.File), args.Error(1)
}

func (m *mockUseCaseFileFolder) DeleteFile(ctx context.Context, companyID, fileID string) error {
	args := m.Called(ctx, companyID, fileID)
	return args.Error(0)
}

func (m *mockUseCaseFileFolder) GetUploadStrategy(ctx context.Context, fileSize int64) (*domain.StrategyInfo, error) {
	args := m.Called(ctx, fileSize)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.StrategyInfo), args.Error(1)
}

func (m *mockUseCaseFileFolder) InitChunkedUpload(ctx context.Context, companyID, userID, filename string, fileSize int64, parentPath *domain.Path, mimeType string) (*domain.ChunkedUpload, error) {
	args := m.Called(ctx, companyID, userID, filename, fileSize, parentPath, mimeType)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.ChunkedUpload), args.Error(1)
}

func (m *mockUseCaseFileFolder) UploadChunk(ctx context.Context, companyID, uploadID string, chunkIndex int, chunkData io.Reader, chunkSize int64) (*domain.ChunkedUpload, error) {
	args := m.Called(ctx, companyID, uploadID, chunkIndex, chunkData, chunkSize)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.ChunkedUpload), args.Error(1)
}

func (m *mockUseCaseFileFolder) GetChunkedUploadStatus(ctx context.Context, companyID, uploadID string) (*domain.ChunkedUpload, error) {
	args := m.Called(ctx, companyID, uploadID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.ChunkedUpload), args.Error(1)
}

func (m *mockUseCaseFileFolder) CompleteChunkedUpload(ctx context.Context, companyID, uploadID string) (*domain.File, error) {
	args := m.Called(ctx, companyID, uploadID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.File), args.Error(1)
}

func (m *mockUseCaseFileFolder) AbortChunkedUpload(ctx context.Context, companyID, uploadID string) error {
	args := m.Called(ctx, companyID, uploadID)
	return args.Error(0)
}

func (m *mockUseCaseFileFolder) GetResourceStats(ctx context.Context) (*domain.ResourceStats, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.ResourceStats), args.Error(1)
}

func setupTestRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	return gin.New()
}

func createTestContext(companyID, userID string) *gin.Context {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Set("company_id", companyID)
	c.Set("user_id", userID)
	return c
}

func createMultipartRequest(fieldName, fileName, content string) (*http.Request, error) {
	var buf bytes.Buffer
	writer := multipart.NewWriter(&buf)

	// Add parentPath form field
	err := writer.WriteField("parentPath", "/test")
	if err != nil {
		return nil, err
	}

	part, err := writer.CreateFormFile(fieldName, fileName)
	if err != nil {
		return nil, err
	}

	_, err = part.Write([]byte(content))
	if err != nil {
		return nil, err
	}

	err = writer.Close()
	if err != nil {
		return nil, err
	}

	req := httptest.NewRequest("POST", "/", &buf)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	return req, nil
}

func createTestFile() *domain.File {
	path, _ := domain.NewPath("/test/folder")
	mimeType := "text/plain"
	size := int64(1024)

	return &domain.File{
		ID:           "test-id",
		Name:         "test.txt",
		Type:         domain.FileTypeFile,
		FullPath:     path,
		CompanyId:    "company-123",
		UserCreateID: "user-123",
		MimeType:     &mimeType,
		Size:         &size,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}
}

func createTestFolder() *domain.File {
	path, _ := domain.NewPath("/test/folder")

	return &domain.File{
		ID:           "folder-id",
		Name:         "test-folder",
		Type:         domain.FileTypeFolder,
		FullPath:     path,
		CompanyId:    "company-123",
		UserCreateID: "user-123",
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}
}

func createTestChunkedUpload() *domain.ChunkedUpload {
	parentPath, _ := domain.NewPath("/test")
	targetPath, _ := domain.NewPath("/test/large-file.bin")

	return &domain.ChunkedUpload{
		ID:             "upload-123",
		FileName:       "large-file.bin",
		TotalSize:      1073741824, // 1GB
		ChunkSize:      5242880,    // 5MB
		TotalChunks:    205,
		UploadedChunks: 0,
		Status:         domain.ChunkedUploadStatusActive,
		CompanyID:      "company-123",
		UserCreateID:   "user-123",
		ParentPath:     parentPath,
		TargetPath:     targetPath,
		MimeType:       "application/octet-stream",
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
		ExpiresAt:      time.Now().Add(24 * time.Hour),
		Chunks:         make(map[int]*domain.ChunkInfo),
	}
}

func createTestStrategyInfo() *domain.StrategyInfo {
	return &domain.StrategyInfo{
		Strategy:       domain.UploadStrategyChunked,
		FileSize:       1073741824, // 1GB
		RequiresChunks: true,
		ChunkSize:      5242880, // 5MB
		TotalChunks:    205,
	}
}

func createTestResourceStats() *domain.ResourceStats {
	return &domain.ResourceStats{
		MemoryUsage: domain.ResourceUsage{
			Current:     52428800,   // 50MB
			Limit:       1073741824, // 1GB
			SystemUsed:  104857600,  // 100MB
			SystemTotal: 2147483648, // 2GB
		},
		ActiveUploads: 3,
		MaxUploads:    10,
		CircuitState:  domain.CircuitClosed,
		Failures:      0,
	}
}

func TestCreateFolder_Success(t *testing.T) {
	mockUC := new(mockUseCaseFileFolder)
	handler := NewHandlerFileFolder(mockUC)

	expectedFolder := createTestFolder()
	mockUC.On("CreateFolder", mock.Anything, mock.AnythingOfType("*domain.File")).Return(expectedFolder, nil)

	reqBody := `{"name":"test-folder","parentPath":"/test"}`
	req := httptest.NewRequest("POST", "/folders", strings.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req
	c.Set("company_id", "company-123")

	handler.CreateFolder(c)

	assert.Equal(t, http.StatusCreated, w.Code)
	mockUC.AssertExpectations(t)

	var response ResponseFolder
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "success", response.Status)
	assert.Equal(t, expectedFolder.ID, response.Folder.ID)
}

func TestCreateFolder_MissingCompanyID(t *testing.T) {
	mockUC := new(mockUseCaseFileFolder)
	handler := NewHandlerFileFolder(mockUC)

	reqBody := `{"name":"test-folder","parentPath":"/test"}`
	req := httptest.NewRequest("POST", "/folders", strings.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req
	// Не устанавливаем company_id

	handler.CreateFolder(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	mockUC.AssertNotCalled(t, "CreateFolder")
}

func TestCreateFolder_InvalidJSON(t *testing.T) {
	mockUC := new(mockUseCaseFileFolder)
	handler := NewHandlerFileFolder(mockUC)

	reqBody := `{"name":"test-folder","parentPath":"/test"` // Invalid JSON
	req := httptest.NewRequest("POST", "/folders", strings.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req
	c.Set("company_id", "company-123")

	handler.CreateFolder(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	mockUC.AssertNotCalled(t, "CreateFolder")
}

func TestGetFolderContents_Success(t *testing.T) {
	mockUC := new(mockUseCaseFileFolder)
	handler := NewHandlerFileFolder(mockUC)

	expectedFiles := []*domain.File{createTestFolder(), createTestFile()}
	mockUC.On("GetFolderContents", mock.Anything, "company-123", mock.AnythingOfType("*domain.Path"), mock.AnythingOfType("*domain.FileType")).Return(expectedFiles, nil)

	reqBody := `{"path":"/test","type":"folder"}`
	req := httptest.NewRequest("POST", "/folders/contents", strings.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req
	c.Set("company_id", "company-123")

	handler.GetFolderContents(c)

	assert.Equal(t, http.StatusOK, w.Code)
	mockUC.AssertExpectations(t)

	var response ResponseGetFolder
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "success", response.Status)
	assert.Len(t, response.Files, 2)
}

func TestUploadFile_Success(t *testing.T) {
	mockUC := new(mockUseCaseFileFolder)
	handler := NewHandlerFileFolder(mockUC)

	expectedFile := createTestFile()
	mockUC.On("UploadFile", mock.Anything, "company-123", "user-123", mock.AnythingOfType("*domain.Path"), "test.txt", int64(12), mock.AnythingOfType("multipart.sectionReadCloser")).Return(expectedFile, nil)

	req, err := createMultipartRequest("file", "test.txt", "test content")
	assert.NoError(t, err)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req
	c.Set("company_id", "company-123")
	c.Set("user_id", "user-123")

	handler.UploadFile(c)

	assert.Equal(t, http.StatusCreated, w.Code)
	mockUC.AssertExpectations(t)
}

func TestUploadFile_MissingFile(t *testing.T) {
	mockUC := new(mockUseCaseFileFolder)
	handler := NewHandlerFileFolder(mockUC)

	reqBody := `parentPath=/test`
	req := httptest.NewRequest("POST", "/files/upload", strings.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req
	c.Set("company_id", "company-123")
	c.Set("user_id", "user-123")

	handler.UploadFile(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	mockUC.AssertNotCalled(t, "UploadFile")
}

func TestDownloadFile_Success(t *testing.T) {
	mockUC := new(mockUseCaseFileFolder)
	handler := NewHandlerFileFolder(mockUC)

	testFile := createTestFile()
	mockReader := io.NopCloser(strings.NewReader("test file content"))
	mockUC.On("DownloadFile", mock.Anything, "company-123", "123e4567-e89b-12d3-a456-426614174000").Return(mockReader, testFile, nil)

	req := httptest.NewRequest("GET", "/files/123e4567-e89b-12d3-a456-426614174000/download", nil)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req
	c.Set("company_id", "company-123")
	c.Params = []gin.Param{{Key: "id", Value: "123e4567-e89b-12d3-a456-426614174000"}}

	handler.DownloadFile(c)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "text/plain", w.Header().Get("Content-Type"))
	assert.Contains(t, w.Header().Get("Content-Disposition"), "test.txt")
	mockUC.AssertExpectations(t)
}

func TestGetUploadStrategy_Success(t *testing.T) {
	mockUC := new(mockUseCaseFileFolder)
	handler := NewHandlerFileFolder(mockUC)

	expectedStrategy := createTestStrategyInfo()
	mockUC.On("GetUploadStrategy", mock.Anything, int64(1073741824)).Return(expectedStrategy, nil)

	req := httptest.NewRequest("GET", "/files/upload-strategy?fileSize=1073741824", nil)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	handler.GetUploadStrategy(c)

	assert.Equal(t, http.StatusOK, w.Code)
	mockUC.AssertExpectations(t)

	var response ResponseUploadStrategy
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "success", response.Status)
	assert.Equal(t, domain.UploadStrategyChunked, response.Strategy.Strategy)
}

func TestGetUploadStrategy_MissingFileSize(t *testing.T) {
	mockUC := new(mockUseCaseFileFolder)
	handler := NewHandlerFileFolder(mockUC)

	req := httptest.NewRequest("GET", "/files/upload-strategy", nil)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	handler.GetUploadStrategy(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	mockUC.AssertNotCalled(t, "GetUploadStrategy")
}

func TestInitChunkedUpload_Success(t *testing.T) {
	mockUC := new(mockUseCaseFileFolder)
	handler := NewHandlerFileFolder(mockUC)

	expectedUpload := createTestChunkedUpload()
	expectedStrategy := createTestStrategyInfo()

	mockUC.On("InitChunkedUpload", mock.Anything, "company-123", "user-123", "large-file.zip", int64(1073741824), mock.AnythingOfType("*domain.Path"), "application/zip").Return(expectedUpload, nil)
	mockUC.On("GetUploadStrategy", mock.Anything, int64(1073741824)).Return(expectedStrategy, nil)

	reqBody := `{"fileName":"large-file.zip","fileSize":1073741824,"parentPath":"/test","mimeType":"application/zip"}`
	req := httptest.NewRequest("POST", "/files/chunked/init", strings.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req
	c.Set("company_id", "company-123")
	c.Set("user_id", "user-123")

	handler.InitChunkedUpload(c)

	assert.Equal(t, http.StatusCreated, w.Code)
	mockUC.AssertExpectations(t)

	var response ResponseInitChunkedUpload
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "success", response.Status)
	assert.Equal(t, expectedUpload.ID, response.UploadID)
}

func TestUploadChunk_Success(t *testing.T) {
	mockUC := new(mockUseCaseFileFolder)
	handler := NewHandlerFileFolder(mockUC)

	expectedUpload := createTestChunkedUpload()
	expectedUpload.UploadedChunks = 1

	mockUC.On("UploadChunk", mock.Anything, "company-123", "123e4567-e89b-12d3-a456-426614174001", 0, mock.AnythingOfType("multipart.sectionReadCloser"), int64(5242880)).Return(expectedUpload, nil)

	req, err := createMultipartRequest("chunk", "chunk-0", strings.Repeat("a", 5242880))
	assert.NoError(t, err)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req
	c.Set("company_id", "company-123")
	c.Params = []gin.Param{
		{Key: "uploadId", Value: "123e4567-e89b-12d3-a456-426614174001"},
		{Key: "chunkIndex", Value: "0"},
	}

	handler.UploadChunk(c)

	assert.Equal(t, http.StatusOK, w.Code)
	mockUC.AssertExpectations(t)
}

func TestGetResourceStats_Success(t *testing.T) {
	mockUC := new(mockUseCaseFileFolder)
	handler := NewHandlerFileFolder(mockUC)

	expectedStats := createTestResourceStats()
	mockUC.On("GetResourceStats", mock.Anything).Return(expectedStats, nil)

	req := httptest.NewRequest("GET", "/files/stats", nil)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	handler.GetResourceStats(c)

	assert.Equal(t, http.StatusOK, w.Code)
	mockUC.AssertExpectations(t)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "success", response["status"])
	assert.NotNil(t, response["stats"])
}

func TestDeleteFile_UseCaseError(t *testing.T) {
	mockUC := new(mockUseCaseFileFolder)
	handler := NewHandlerFileFolder(mockUC)

	mockUC.On("DeleteFile", mock.Anything, "company-123", "123e4567-e89b-12d3-a456-426614174000").Return(errors.New("file not found"))

	req := httptest.NewRequest("DELETE", "/files/123e4567-e89b-12d3-a456-426614174000", nil)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req
	c.Set("company_id", "company-123")
	c.Params = []gin.Param{{Key: "id", Value: "123e4567-e89b-12d3-a456-426614174000"}}

	handler.DeleteFile(c)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	mockUC.AssertExpectations(t)
}

func TestCompleteChunkedUpload_Success(t *testing.T) {
	mockUC := new(mockUseCaseFileFolder)
	handler := NewHandlerFileFolder(mockUC)

	expectedFile := createTestFile()
	mockUC.On("CompleteChunkedUpload", mock.Anything, "company-123", "123e4567-e89b-12d3-a456-426614174001").Return(expectedFile, nil)

	req := httptest.NewRequest("POST", "/files/chunked/123e4567-e89b-12d3-a456-426614174001/complete", nil)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req
	c.Set("company_id", "company-123")
	c.Params = []gin.Param{{Key: "uploadId", Value: "123e4567-e89b-12d3-a456-426614174001"}}

	handler.CompleteChunkedUpload(c)

	assert.Equal(t, http.StatusOK, w.Code)
	mockUC.AssertExpectations(t)

	var response ResponseCompleteChunkedUpload
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "success", response.Status)
	assert.Equal(t, expectedFile.ID, response.File.ID)
}
