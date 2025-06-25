package domain

import (
	"fmt"
	"time"
)

type ChunkedUploadStatus string

const (
	ChunkedUploadStatusActive    ChunkedUploadStatus = "active"
	ChunkedUploadStatusCompleted ChunkedUploadStatus = "completed"
	ChunkedUploadStatusFailed    ChunkedUploadStatus = "failed"
	ChunkedUploadStatusExpired   ChunkedUploadStatus = "expired"
)

type ChunkedUpload struct {
	ID             string
	FileName       string
	TotalSize      int64
	ChunkSize      int64
	TotalChunks    int
	UploadedChunks int
	UploadedSize   int64
	Status         ChunkedUploadStatus
	CompanyID      string
	UserCreateID   string
	ParentPath     Path
	TargetPath     Path

	MimeType string

	CreatedAt time.Time
	UpdatedAt time.Time
	ExpiresAt time.Time

	Chunks map[int]*ChunkInfo
}

type ChunkInfo struct {
	Index      int
	Size       int64
	ETag       string
	Uploaded   bool
	UploadedAt time.Time
	Retries    int
}

type UploadStrategy string

const (
	UploadStrategyMemory  UploadStrategy = "memory"
	UploadStrategyStream  UploadStrategy = "stream"
	UploadStrategyChunked UploadStrategy = "chunked"
)

type FileUploadContext struct {
	File         *File
	Strategy     UploadStrategy
	ChunkSession *ChunkedUpload
	TempPath     string
	CompanyID    string
	UserID       string
}

func (cu *ChunkedUpload) IsComplete() bool {
	return cu.UploadedChunks == cu.TotalChunks && cu.Status == ChunkedUploadStatusActive
}

func (cu *ChunkedUpload) IsExpired() bool {
	return time.Now().After(cu.ExpiresAt) || cu.Status == ChunkedUploadStatusExpired
}

func (cu *ChunkedUpload) GetProgress() float64 {
	if cu.TotalChunks == 0 {
		return 0.0
	}
	return float64(cu.UploadedChunks) / float64(cu.TotalChunks) * 100.0
}

func (cu *ChunkedUpload) GetUploadedSize() int64 {
	return cu.UploadedSize
}

func (cu *ChunkedUpload) AddChunk(chunkIndex int, size int64, etag string) {
	if cu.Chunks == nil {
		cu.Chunks = make(map[int]*ChunkInfo)
	}

	chunk := &ChunkInfo{
		Index:      chunkIndex,
		Size:       size,
		ETag:       etag,
		Uploaded:   true,
		UploadedAt: time.Now(),
	}

	if existingChunk, exists := cu.Chunks[chunkIndex]; !exists || !existingChunk.Uploaded {
		cu.UploadedChunks++
		cu.UploadedSize += size
	}

	cu.Chunks[chunkIndex] = chunk
	cu.UpdatedAt = time.Now()
}

func (cu *ChunkedUpload) GetMissingChunks() []int {
	missing := make([]int, 0)

	for i := 0; i < cu.TotalChunks; i++ {
		if chunk, exists := cu.Chunks[i]; !exists || !chunk.Uploaded {
			missing = append(missing, i)
		}
	}

	return missing
}

func (cu *ChunkedUpload) MarkAsCompleted() {
	cu.Status = ChunkedUploadStatusCompleted
	cu.UpdatedAt = time.Now()
}

func (cu *ChunkedUpload) MarkAsFailed() {
	cu.Status = ChunkedUploadStatusFailed
	cu.UpdatedAt = time.Now()
}

func (cu *ChunkedUpload) Validate() error {
	if cu.FileName == "" {
		return NewValidationError("filename is required")
	}

	if cu.TotalSize <= 0 {
		return NewValidationError("total size must be positive")
	}

	if cu.ChunkSize <= 0 {
		return NewValidationError("chunk size must be positive")
	}

	if cu.CompanyID == "" {
		return NewValidationError("company_id is required")
	}

	if cu.UserCreateID == "" {
		return NewValidationError("user_created_id is required")
	}

	return nil
}

func NewValidationError(message string) error {
	return fmt.Errorf("validation error: %s", message)
}
