package ucFileFolder

import (
	"context"
	"go-storage/internal/domain"
	"io"
)

type RepositoryFileFolder interface {
	// File operations
	CreateFile(ctx context.Context, file *domain.File) (*domain.File, error)
	GetFile(ctx context.Context, companyID, fileID string) (*domain.File, error)
	GetFolderContents(ctx context.Context, companyID string, path *domain.Path, fileType *domain.FileType) ([]*domain.File, error)
	UpdateFile(ctx context.Context, file *domain.File) (*domain.File, error)
	DeleteFile(ctx context.Context, companyID, fileID string) error

	// File path operations
	GetFileByPath(ctx context.Context, companyID string, path *domain.Path) (*domain.File, error)
	MoveFile(ctx context.Context, companyID, fileID string, newPath *domain.Path) (*domain.File, error)
	RenameFile(ctx context.Context, companyID, fileID, newName string) (*domain.File, error)

	// Folder operations
	CreateFolder(ctx context.Context, folder *domain.File) (*domain.File, error)
	GetFolder(ctx context.Context, companyID string, path *domain.Path) (*domain.File, error)
	MoveFolder(ctx context.Context, companyID string, oldPath, newPath *domain.Path) (*domain.Path, error)
	DeleteFolder(ctx context.Context, companyID string, path *domain.Path) error
}

type StorageRepository interface {
	// File storage operations
	StoreFile(ctx context.Context, key string, reader io.Reader, size int64, mimeType string) (string, error)
	GetFile(ctx context.Context, key string) (io.ReadCloser, error)
	DeleteFile(ctx context.Context, key string) error

	// Chunked upload operations
	InitChunkedUpload(ctx context.Context, key string, mimeType string) (string, error)
	UploadChunk(ctx context.Context, uploadID, key string, chunkIndex int, reader io.Reader, size int64) (string, error)
	CompleteChunkedUpload(ctx context.Context, uploadID, key string, parts []string) error
	AbortChunkedUpload(ctx context.Context, uploadID, key string) error

	// File info operations
	GetFileInfo(ctx context.Context, key string) (*domain.StorageFileInfo, error)
}

type ChunkedUploadRepository interface {
	// Chunked upload session management
	CreateChunkedUpload(ctx context.Context, upload *domain.ChunkedUpload) (*domain.ChunkedUpload, error)
	GetChunkedUpload(ctx context.Context, companyID, uploadID string) (*domain.ChunkedUpload, error)
	UpdateChunkedUpload(ctx context.Context, upload *domain.ChunkedUpload) (*domain.ChunkedUpload, error)
	DeleteChunkedUpload(ctx context.Context, companyID, uploadID string) error

	// Chunk tracking
	AddChunk(ctx context.Context, uploadID string, chunkIndex int, etag string, size int64) error
	GetUploadProgress(ctx context.Context, uploadID string) (*domain.ChunkedUpload, error)

	// Cleanup operations
	CleanupExpiredUploads(ctx context.Context) error
}
