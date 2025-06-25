package ucFileFolder

import (
	"context"
	"fmt"
	"io"
	"path/filepath"
	"strings"
	"time"

	"github.com/google/uuid"
	"go-storage/internal/config"
	"go-storage/internal/domain"
	"go-storage/pkg/errors"
)

type UseCaseFileFolder struct {
	fileRepo         RepositoryFileFolder
	storageRepo      StorageRepository
	chunkedRepo      ChunkedUploadRepository
	resourceMonitor  *domain.ResourceMonitor
	strategySelector *domain.UploadStrategySelector
	config           *config.FileServer
}

func NewUseCaseFileFolder(
	fileRepo RepositoryFileFolder,
	storageRepo StorageRepository,
	chunkedRepo ChunkedUploadRepository,
	config *config.FileServer,
) *UseCaseFileFolder {
	resourceMonitor := domain.NewResourceMonitor(config)
	strategySelector := domain.NewUploadStrategySelector(config)

	return &UseCaseFileFolder{
		fileRepo:         fileRepo,
		storageRepo:      storageRepo,
		chunkedRepo:      chunkedRepo,
		resourceMonitor:  resourceMonitor,
		strategySelector: strategySelector,
		config:           config,
	}
}

func (uc *UseCaseFileFolder) CreateFolder(ctx context.Context, folder *domain.File) (*domain.File, error) {
	if folder.Name == "" {
		return nil, errors.BadRequest("folder name is required")
	}

	if folder.CompanyId == "" {
		return nil, errors.BadRequest("company ID is required")
	}

	if folder.ParentID != nil {
		parent, err := uc.fileRepo.GetFile(ctx, folder.CompanyId, *folder.ParentID)
		if err != nil {
			return nil, errors.BadRequest("parent folder not found")
		}
		if parent.Type != domain.FileTypeFolder {
			return nil, errors.BadRequest("parent must be a folder")
		}
	}

	// Check for naming conflicts
	existingFile, err := uc.fileRepo.GetFileByPath(ctx, folder.CompanyId, &folder.FullPath)
	if err == nil && existingFile != nil {
		return nil, errors.BadRequest("folder with this name already exists")
	}

	folder.ID = uuid.NewString()
	folder.Type = domain.FileTypeFolder
	folder.CreatedAt = time.Now()
	folder.UpdatedAt = time.Now()
	folder.IsActive = true

	return uc.fileRepo.CreateFolder(ctx, folder)
}

func (uc *UseCaseFileFolder) GetFolderContents(ctx context.Context, companyID string, path *domain.Path, fileType *domain.FileType) ([]*domain.File, error) {
	if companyID == "" {
		return nil, errors.BadRequest("company ID is required")
	}

	return uc.fileRepo.GetFolderContents(ctx, companyID, path, fileType)
}

func (uc *UseCaseFileFolder) MoveFolder(ctx context.Context, companyID string, folderPath *domain.Path, newPath *domain.Path) (*domain.Path, error) {
	if companyID == "" {
		return nil, errors.BadRequest("company ID is required")
	}

	folder, err := uc.fileRepo.GetFileByPath(ctx, companyID, folderPath)
	if err != nil {
		return nil, errors.NotFound("folder not found")
	}

	if folder.Type != domain.FileTypeFolder {
		return nil, errors.BadRequest("specified path is not a folder")
	}

	// Check if destination already exists
	_, err = uc.fileRepo.GetFileByPath(ctx, companyID, newPath)
	if err == nil {
		return nil, errors.BadRequest("destination already exists")
	}

	// Update folder path
	result, err := uc.fileRepo.MoveFolder(ctx, companyID, folderPath, newPath)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (uc *UseCaseFileFolder) DeleteFolder(ctx context.Context, companyID string, folderPath *domain.Path) error {
	if companyID == "" {
		return errors.BadRequest("company ID is required")
	}

	// Check if folder exists
	folder, err := uc.fileRepo.GetFileByPath(ctx, companyID, folderPath)
	if err != nil {
		return errors.NotFound("folder not found")
	}

	if folder.Type != domain.FileTypeFolder {
		return errors.BadRequest("specified path is not a folder")
	}

	// Check if folder is empty (optional - could allow recursive delete)
	contents, err := uc.fileRepo.GetFolderContents(ctx, companyID, folderPath, nil)
	if err != nil {
		return err
	}

	if len(contents) > 0 {
		return errors.BadRequest("folder is not empty")
	}

	return uc.fileRepo.DeleteFolder(ctx, companyID, folderPath)
}

func (uc *UseCaseFileFolder) UploadFile(ctx context.Context, companyID, userID string, parentPath *domain.Path, filename string, size int64, reader io.Reader) (*domain.File, error) {
	if companyID == "" {
		return nil, errors.BadRequest("company ID is required")
	}

	if userID == "" {
		return nil, errors.BadRequest("user ID is required")
	}

	if filename == "" {
		return nil, errors.BadRequest("filename is required")
	}

	// Check file size limits
	if size > uc.config.MaxFileSize {
		return nil, errors.BadRequest("file size exceeds maximum allowed size")
	}

	// Check resource availability
	if err := uc.resourceMonitor.AcquireUploadSlot(ctx); err != nil {
		return nil, errors.TooManyRequests("too many concurrent uploads")
	}
	defer uc.resourceMonitor.ReleaseUploadSlot()

	// Get upload strategy
	strategy, err := uc.strategySelector.SelectStrategy(size)
	if err != nil {
		return nil, err
	}

	// Create file record
	file := &domain.File{
		ID:           uuid.NewString(),
		Name:         filename,
		Type:         domain.FileTypeFile,
		FullPath:     parentPath.Join(filename),
		CompanyId:    companyID,
		UserCreateID: userID,
		Size:         &size,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
		IsActive:     true,
	}

	// Determine MIME type
	mimeType := determineMimeType(filename)
	file.MimeType = &mimeType

	// Generate storage key
	storageKey := generateStorageKey(companyID, file.ID, filename)

	// Upload file using strategy
	uploadCtx := &domain.FileUploadContext{
		File:      file,
		Strategy:  strategy.GetStrategy(),
		CompanyID: companyID,
		UserID:    userID,
	}

	storagePath, err := uc.uploadWithStrategy(ctx, uploadCtx, reader, size, mimeType, storageKey)
	if err != nil {
		uc.resourceMonitor.RecordFailure()
		return nil, err
	}

	file.StoragePath = &storagePath
	uc.resourceMonitor.RecordSuccess()

	// Save file record to database
	return uc.fileRepo.CreateFile(ctx, file)
}

func (uc *UseCaseFileFolder) uploadWithStrategy(ctx context.Context, uploadCtx *domain.FileUploadContext, reader io.Reader, size int64, mimeType, storageKey string) (string, error) {
	switch uploadCtx.Strategy {
	case domain.UploadStrategyMemory:
		return uc.uploadMemoryStrategy(ctx, reader, size, mimeType, storageKey)
	case domain.UploadStrategyStream:
		return uc.uploadStreamStrategy(ctx, reader, size, mimeType, storageKey)
	case domain.UploadStrategyChunked:
		return "", errors.BadRequest("chunked upload should use separate endpoint")
	default:
		return "", errors.InternalServer("unknown upload strategy")
	}
}

func (uc *UseCaseFileFolder) uploadMemoryStrategy(ctx context.Context, reader io.Reader, size int64, mimeType, storageKey string) (string, error) {
	// For small files, read into memory and upload
	if !uc.resourceMonitor.AllocateMemory(size) {
		return "", errors.TooManyRequests("insufficient memory for upload")
	}
	defer uc.resourceMonitor.ReleaseMemory(size)

	return uc.storageRepo.StoreFile(ctx, storageKey, reader, size, mimeType)
}

func (uc *UseCaseFileFolder) uploadStreamStrategy(ctx context.Context, reader io.Reader, size int64, mimeType, storageKey string) (string, error) {
	// For medium files, stream directly to storage
	return uc.storageRepo.StoreFile(ctx, storageKey, reader, size, mimeType)
}

func (uc *UseCaseFileFolder) DownloadFile(ctx context.Context, companyID, fileID string) (io.ReadCloser, *domain.File, error) {
	if companyID == "" {
		return nil, nil, errors.BadRequest("company ID is required")
	}

	// Get file record
	file, err := uc.fileRepo.GetFile(ctx, companyID, fileID)
	if err != nil {
		return nil, nil, err
	}

	if file.Type != domain.FileTypeFile {
		return nil, nil, errors.BadRequest("specified ID is not a file")
	}

	if file.StoragePath == nil {
		return nil, nil, errors.InternalServer("file storage path not found")
	}

	// Get file from storage
	reader, err := uc.storageRepo.GetFile(ctx, *file.StoragePath)
	if err != nil {
		return nil, nil, errors.InternalServer("failed to retrieve file from storage")
	}

	return reader, file, nil
}

func (uc *UseCaseFileFolder) GetFileInfo(ctx context.Context, companyID, fileID string) (*domain.File, error) {
	if companyID == "" {
		return nil, errors.BadRequest("company ID is required")
	}

	return uc.fileRepo.GetFile(ctx, companyID, fileID)
}

func (uc *UseCaseFileFolder) RenameFile(ctx context.Context, companyID, fileID, newName string) (*domain.File, error) {
	if companyID == "" {
		return nil, errors.BadRequest("company ID is required")
	}

	if newName == "" {
		return nil, errors.BadRequest("new name is required")
	}

	return uc.fileRepo.RenameFile(ctx, companyID, fileID, newName)
}

func (uc *UseCaseFileFolder) MoveFile(ctx context.Context, companyID, fileID string, newParentPath *domain.Path) (*domain.File, error) {
	if companyID == "" {
		return nil, errors.BadRequest("company ID is required")
	}

	return uc.fileRepo.MoveFile(ctx, companyID, fileID, newParentPath)
}

func (uc *UseCaseFileFolder) DeleteFile(ctx context.Context, companyID, fileID string) error {
	if companyID == "" {
		return errors.BadRequest("company ID is required")
	}

	// Get file record
	file, err := uc.fileRepo.GetFile(ctx, companyID, fileID)
	if err != nil {
		return err
	}

	// Delete from storage if it exists
	if file.StoragePath != nil {
		if err := uc.storageRepo.DeleteFile(ctx, *file.StoragePath); err != nil {
			// Log error but don't fail the operation
			// Storage cleanup can be done later
		}
	}

	// Delete from database
	return uc.fileRepo.DeleteFile(ctx, companyID, fileID)
}

func (uc *UseCaseFileFolder) GetUploadStrategy(ctx context.Context, fileSize int64) (*domain.StrategyInfo, error) {
	return uc.strategySelector.GetStrategyInfo(fileSize)
}

func (uc *UseCaseFileFolder) InitChunkedUpload(ctx context.Context, companyID, userID, filename string, fileSize int64, parentPath *domain.Path, mimeType string) (*domain.ChunkedUpload, error) {
	if companyID == "" {
		return nil, errors.BadRequest("company ID is required")
	}

	if userID == "" {
		return nil, errors.BadRequest("user ID is required")
	}

	if filename == "" {
		return nil, errors.BadRequest("filename is required")
	}

	// Validate file size
	if fileSize > uc.config.MaxFileSize {
		return nil, errors.BadRequest("file size exceeds maximum allowed size")
	}

	if fileSize <= uc.config.MediumFileThreshold {
		return nil, errors.BadRequest("file is too small for chunked upload")
	}

	// Create chunked upload session
	targetPath := parentPath.Join(filename)
	upload := &domain.ChunkedUpload{
		ID:             uuid.NewString(),
		FileName:       filename,
		TotalSize:      fileSize,
		ChunkSize:      uc.config.ChunkSize,
		TotalChunks:    int((fileSize + uc.config.ChunkSize - 1) / uc.config.ChunkSize),
		UploadedChunks: 0,
		UploadedSize:   0,
		Status:         domain.ChunkedUploadStatusActive,
		CompanyID:      companyID,
		UserCreateID:   userID,
		ParentPath:     *parentPath,
		TargetPath:     targetPath,
		MimeType:       mimeType,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
		ExpiresAt:      time.Now().Add(uc.config.ChunkedSessionTTL),
		Chunks:         make(map[int]*domain.ChunkInfo),
	}

	// Initialize storage upload session
	storageKey := generateStorageKey(companyID, upload.ID, filename)
	storageUploadID, err := uc.storageRepo.InitChunkedUpload(ctx, storageKey, mimeType)
	if err != nil {
		return nil, errors.InternalServer("failed to initialize storage upload")
	}

	// Store the storage upload ID (could extend domain model)
	upload.ID = storageUploadID

	return uc.chunkedRepo.CreateChunkedUpload(ctx, upload)
}

func (uc *UseCaseFileFolder) UploadChunk(ctx context.Context, companyID, uploadID string, chunkIndex int, chunkData io.Reader, chunkSize int64) (*domain.ChunkedUpload, error) {
	if companyID == "" {
		return nil, errors.BadRequest("company ID is required")
	}

	// Get upload session
	upload, err := uc.chunkedRepo.GetChunkedUpload(ctx, companyID, uploadID)
	if err != nil {
		return nil, err
	}

	// Validate upload session
	if upload.Status != domain.ChunkedUploadStatusActive {
		return nil, errors.BadRequest("upload session is not active")
	}

	if upload.IsExpired() {
		return nil, errors.BadRequest("upload session has expired")
	}

	if chunkIndex < 0 || chunkIndex >= upload.TotalChunks {
		return nil, errors.BadRequest("invalid chunk index")
	}

	// Check if chunk already uploaded
	if chunk, exists := upload.Chunks[chunkIndex]; exists && chunk.Uploaded {
		return upload, nil // Already uploaded
	}

	// Upload chunk to storage
	storageKey := generateStorageKey(companyID, upload.ID, upload.FileName)
	etag, err := uc.storageRepo.UploadChunk(ctx, uploadID, storageKey, chunkIndex, chunkData, chunkSize)
	if err != nil {
		return nil, errors.InternalServer("failed to upload chunk to storage")
	}

	// Update chunk info
	upload.AddChunk(chunkIndex, chunkSize, etag)

	// Save progress
	return uc.chunkedRepo.UpdateChunkedUpload(ctx, upload)
}

func (uc *UseCaseFileFolder) GetChunkedUploadStatus(ctx context.Context, companyID, uploadID string) (*domain.ChunkedUpload, error) {
	if companyID == "" {
		return nil, errors.BadRequest("company ID is required")
	}

	return uc.chunkedRepo.GetChunkedUpload(ctx, companyID, uploadID)
}

func (uc *UseCaseFileFolder) CompleteChunkedUpload(ctx context.Context, companyID, uploadID string) (*domain.File, error) {
	if companyID == "" {
		return nil, errors.BadRequest("company ID is required")
	}

	// Get upload session
	upload, err := uc.chunkedRepo.GetChunkedUpload(ctx, companyID, uploadID)
	if err != nil {
		return nil, err
	}

	// Validate completion
	if !upload.IsComplete() {
		return nil, errors.BadRequest("upload is not complete")
	}

	// Collect chunk ETags for storage completion
	parts := make([]string, upload.TotalChunks)
	for i := 0; i < upload.TotalChunks; i++ {
		chunk, exists := upload.Chunks[i]
		if !exists || !chunk.Uploaded {
			return nil, errors.BadRequest("missing chunk")
		}
		parts[i] = chunk.ETag
	}

	// Complete storage upload
	storageKey := generateStorageKey(companyID, upload.ID, upload.FileName)
	err = uc.storageRepo.CompleteChunkedUpload(ctx, uploadID, storageKey, parts)
	if err != nil {
		return nil, errors.InternalServer("failed to complete storage upload")
	}

	// Create file record
	file := &domain.File{
		ID:           uuid.NewString(),
		Name:         upload.FileName,
		Type:         domain.FileTypeFile,
		FullPath:     upload.TargetPath,
		CompanyId:    upload.CompanyID,
		UserCreateID: upload.UserCreateID,
		MimeType:     &upload.MimeType,
		Size:         &upload.TotalSize,
		StoragePath:  &storageKey,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
		IsActive:     true,
	}

	// Mark upload as completed
	upload.MarkAsCompleted()
	_, err = uc.chunkedRepo.UpdateChunkedUpload(ctx, upload)
	if err != nil {
		// Log error but continue
	}

	// Save file record
	return uc.fileRepo.CreateFile(ctx, file)
}

func (uc *UseCaseFileFolder) AbortChunkedUpload(ctx context.Context, companyID, uploadID string) error {
	if companyID == "" {
		return errors.BadRequest("company ID is required")
	}

	// Get upload session
	upload, err := uc.chunkedRepo.GetChunkedUpload(ctx, companyID, uploadID)
	if err != nil {
		return err
	}

	// Abort storage upload
	storageKey := generateStorageKey(companyID, upload.ID, upload.FileName)
	if err := uc.storageRepo.AbortChunkedUpload(ctx, uploadID, storageKey); err != nil {
		// Log error but continue cleanup
	}

	// Delete upload session
	return uc.chunkedRepo.DeleteChunkedUpload(ctx, companyID, uploadID)
}

func (uc *UseCaseFileFolder) GetResourceStats(ctx context.Context) (*domain.ResourceStats, error) {
	return uc.resourceMonitor.GetResourceStats(), nil
}

// Helper functions

func determineMimeType(filename string) string {
	ext := strings.ToLower(filepath.Ext(filename))

	mimeTypes := map[string]string{
		".txt":  "text/plain",
		".pdf":  "application/pdf",
		".jpg":  "image/jpeg",
		".jpeg": "image/jpeg",
		".png":  "image/png",
		".gif":  "image/gif",
		".mp4":  "video/mp4",
		".zip":  "application/zip",
		".json": "application/json",
		".xml":  "application/xml",
	}

	if mimeType, exists := mimeTypes[ext]; exists {
		return mimeType
	}

	return "application/octet-stream"
}

func generateStorageKey(companyID, fileID, filename string) string {
	return fmt.Sprintf("companies/%s/files/%s/%s", companyID, fileID, filename)
}
