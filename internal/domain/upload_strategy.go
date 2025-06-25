package domain

import (
	"context"
	"go-storage/internal/config"
	"io"
)

type FileUploadStrategy interface {
	CanHandle(size int64, config *config.FileServer) bool
	GetStrategy() UploadStrategy
	Upload(ctx context.Context, uploadCtx *FileUploadContext, reader io.Reader) (*File, error)
}

type MemoryUploadStrategy struct{}

func (s *MemoryUploadStrategy) CanHandle(size int64, config *config.FileServer) bool {
	return size <= config.SmallFileThreshold
}

func (s *MemoryUploadStrategy) GetStrategy() UploadStrategy {
	return UploadStrategyMemory
}

func (s *MemoryUploadStrategy) Upload(ctx context.Context, uploadCtx *FileUploadContext, reader io.Reader) (*File, error) {
	return uploadCtx.File, nil
}

type StreamUploadStrategy struct{}

func (s *StreamUploadStrategy) CanHandle(size int64, config *config.FileServer) bool {
	return size > config.SmallFileThreshold && size <= config.MediumFileThreshold
}

func (s *StreamUploadStrategy) GetStrategy() UploadStrategy {
	return UploadStrategyStream
}

func (s *StreamUploadStrategy) Upload(ctx context.Context, uploadCtx *FileUploadContext, reader io.Reader) (*File, error) {
	return uploadCtx.File, nil
}

type ChunkedUploadStrategy struct{}

func (s *ChunkedUploadStrategy) CanHandle(size int64, config *config.FileServer) bool {
	return size > config.MediumFileThreshold && size <= config.MaxFileSize
}

func (s *ChunkedUploadStrategy) GetStrategy() UploadStrategy {
	return UploadStrategyChunked
}

func (s *ChunkedUploadStrategy) Upload(ctx context.Context, uploadCtx *FileUploadContext, reader io.Reader) (*File, error) {
	return uploadCtx.File, nil
}

type UploadStrategySelector struct {
	strategies []FileUploadStrategy
	config     *config.FileServer
}

func NewUploadStrategySelector(config *config.FileServer) *UploadStrategySelector {
	return &UploadStrategySelector{
		strategies: []FileUploadStrategy{
			&MemoryUploadStrategy{},
			&StreamUploadStrategy{},
			&ChunkedUploadStrategy{},
		},
		config: config,
	}
}

func (sel *UploadStrategySelector) SelectStrategy(size int64) (FileUploadStrategy, error) {
	if size > sel.config.MaxFileSize {
		return nil, NewValidationError("file size exceeds maximum allowed size")
	}

	for _, strategy := range sel.strategies {
		if strategy.CanHandle(size, sel.config) {
			return strategy, nil
		}
	}

	return &ChunkedUploadStrategy{}, nil
}

func (sel *UploadStrategySelector) GetStrategyInfo(size int64) (*StrategyInfo, error) {
	strategy, err := sel.SelectStrategy(size)
	if err != nil {
		return nil, err
	}

	info := &StrategyInfo{
		Strategy:       strategy.GetStrategy(),
		FileSize:       size,
		RequiresChunks: strategy.GetStrategy() == UploadStrategyChunked,
	}

	if info.RequiresChunks {
		info.ChunkSize = sel.config.ChunkSize
		info.TotalChunks = calculateTotalChunks(size, sel.config.ChunkSize)
	}

	return info, nil
}

type StrategyInfo struct {
	Strategy       UploadStrategy `json:"strategy"`
	FileSize       int64          `json:"file_size"`
	RequiresChunks bool           `json:"requires_chunks"`
	ChunkSize      int64          `json:"chunk_size,omitempty"`
	TotalChunks    int            `json:"total_chunks,omitempty"`
}

func calculateTotalChunks(fileSize, chunkSize int64) int {
	if chunkSize <= 0 {
		return 0
	}
	chunks := fileSize / chunkSize
	if fileSize%chunkSize != 0 {
		chunks++
	}
	return int(chunks)
}
