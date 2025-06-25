package hdFileFolder

import (
	"context"
	"go-storage/internal/domain"
	"io"
)

type UseCaseFileFolder interface {
	// Folder operations
	CreateFolder(ctx context.Context, folder *domain.File) (*domain.File, error)
	GetFolderContents(ctx context.Context, companyID string, path *domain.Path, fileType *domain.FileType) ([]*domain.File, error)
	MoveFolder(ctx context.Context, companyID string, folderPath *domain.Path, newPath *domain.Path) (*domain.Path, error)
	DeleteFolder(ctx context.Context, companyID string, folderPath *domain.Path) error

	// File operations
	UploadFile(ctx context.Context, companyID, userID string, parentPath *domain.Path, filename string, size int64, reader io.Reader) (*domain.File, error)
	DownloadFile(ctx context.Context, companyID, fileID string) (io.ReadCloser, *domain.File, error)
	GetFileInfo(ctx context.Context, companyID, fileID string) (*domain.File, error)
	RenameFile(ctx context.Context, companyID, fileID, newName string) (*domain.File, error)
	MoveFile(ctx context.Context, companyID, fileID string, newParentPath *domain.Path) (*domain.File, error)
	DeleteFile(ctx context.Context, companyID, fileID string) error

	// Upload strategy
	GetUploadStrategy(ctx context.Context, fileSize int64) (*domain.StrategyInfo, error)

	// Chunked upload operations
	InitChunkedUpload(ctx context.Context, companyID, userID, filename string, fileSize int64, parentPath *domain.Path, mimeType string) (*domain.ChunkedUpload, error)
	UploadChunk(ctx context.Context, companyID, uploadID string, chunkIndex int, chunkData io.Reader, chunkSize int64) (*domain.ChunkedUpload, error)
	GetChunkedUploadStatus(ctx context.Context, companyID, uploadID string) (*domain.ChunkedUpload, error)
	CompleteChunkedUpload(ctx context.Context, companyID, uploadID string) (*domain.File, error)
	AbortChunkedUpload(ctx context.Context, companyID, uploadID string) error

	// Resource monitoring
	GetResourceStats(ctx context.Context) (*domain.ResourceStats, error)
}
