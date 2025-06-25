package hdFileFolder

import (
	"fmt"
	"go-storage/internal/domain"
	"go-storage/internal/utils/valid"
	"time"
)

func ToDomainCreateFolder(dto *RequestCreateFolder) (*domain.File, error) {
	parentPath, err := domain.NewPath(dto.ParentPath)
	if err != nil {
		return nil, err
	}

	name := valid.NormalizationOfName(dto.Name)

	fullPath := parentPath.Join(name)

	return &domain.File{
		Name:     dto.Name,
		Type:     domain.FileTypeFolder,
		FullPath: fullPath,
	}, nil
}

func ToDomainGetFolder(dto *RequestGetFolder) (*domain.Path, *domain.FileType, error) {
	path, err := domain.NewPath(dto.Path)
	if err != nil {
		return nil, nil, err
	}

	if dto.Type != "" && !dto.Type.IsValid() {
		return nil, nil, fmt.Errorf("invalid file type")
	}

	var fileTypePtr *domain.FileType
	if dto.Type != "" {
		fileTypePtr = &dto.Type
	}

	return &path, fileTypePtr, nil
}

func ToResponseFolder(dto *domain.File) *ResponseFolder {
	return &ResponseFolder{
		Status: "success",
		Time:   time.Now(),
		Folder: DtoFileToFolder(dto),
	}
}

func ToResponseFolders(dto []*domain.File) *ResponseGetFolder {
	var answer = make([]*FolderFileDTO, len(dto))
	for index, value := range dto {
		answer[index] = DtoFileToFolder(value)
	}

	return &ResponseGetFolder{
		Status: "success",
		Time:   time.Now(),
		Files:  answer,
	}
}

func ToResponsePath(path *domain.Path) *ResponsePath {
	return &ResponsePath{
		Status: "success",
		Time:   time.Now(),
		Path:   path.String(),
	}
}

func DtoFileToFolder(dto *domain.File) *FolderFileDTO {
	return &FolderFileDTO{
		ID:           dto.ID,
		Name:         dto.Name,
		Type:         dto.Type,
		FullPath:     dto.FullPath.String(),
		ParentID:     dto.ParentID,
		CompanyID:    dto.CompanyId,
		UserCreateID: dto.UserCreateID,

		MimeType:    dto.MimeType,
		Size:        dto.Size,
		StoragePath: dto.StoragePath,

		CreatedAt: dto.CreatedAt,
		UpdatedAt: dto.UpdatedAt,
	}
}

func ToResponseFile(file *domain.File) *ResponseFile {
	return &ResponseFile{
		Status: "success",
		Time:   time.Now(),
		File:   DtoFileToFolder(file),
	}
}

func ToFileDTO(file *domain.File) *FolderFileDTO {
	return DtoFileToFolder(file)
}

func ToResponseUploadStrategy(strategy *domain.StrategyInfo) *ResponseUploadStrategy {
	return &ResponseUploadStrategy{
		Status:   "success",
		Time:     time.Now(),
		Strategy: strategy,
	}
}

func ToResponseInitChunkedUpload(chunkedUpload *domain.ChunkedUpload, strategy *domain.StrategyInfo) *ResponseInitChunkedUpload {
	return &ResponseInitChunkedUpload{
		Status:      "success",
		Time:        time.Now(),
		UploadID:    chunkedUpload.ID,
		ChunkSize:   chunkedUpload.ChunkSize,
		TotalChunks: chunkedUpload.TotalChunks,
		Strategy:    strategy,
	}
}

func ToResponseUploadChunk(chunkedUpload *domain.ChunkedUpload, chunkIndex int) *ResponseUploadChunk {
	return &ResponseUploadChunk{
		Status:         "success",
		Time:           time.Now(),
		ChunkIndex:     chunkIndex,
		Progress:       chunkedUpload.GetProgress(),
		UploadedChunks: chunkedUpload.UploadedChunks,
		TotalChunks:    chunkedUpload.TotalChunks,
	}
}

func ToResponseChunkedUploadStatus(chunkedUpload *domain.ChunkedUpload) *ResponseChunkedUploadStatus {
	return &ResponseChunkedUploadStatus{
		Status:         "success",
		Time:           time.Now(),
		UploadID:       chunkedUpload.ID,
		FileName:       chunkedUpload.FileName,
		Progress:       chunkedUpload.GetProgress(),
		UploadedChunks: chunkedUpload.UploadedChunks,
		TotalChunks:    chunkedUpload.TotalChunks,
		UploadedSize:   chunkedUpload.GetUploadedSize(),
		TotalSize:      chunkedUpload.TotalSize,
		UploadStatus:   chunkedUpload.Status,
		MissingChunks:  chunkedUpload.GetMissingChunks(),
		ExpiresAt:      chunkedUpload.ExpiresAt,
	}
}

func ToResponseCompleteChunkedUpload(file *domain.File) *ResponseCompleteChunkedUpload {
	return &ResponseCompleteChunkedUpload{
		Status: "success",
		Time:   time.Now(),
		File:   ToFileDTO(file),
	}
}

func ToResponseSuccess(message string) map[string]interface{} {
	return map[string]interface{}{
		"status":  "success",
		"time":    time.Now(),
		"message": message,
	}
}

func ToResponseResourceStats(stats *domain.ResourceStats) map[string]interface{} {
	return map[string]interface{}{
		"status": "success",
		"time":   time.Now(),
		"stats":  stats,
	}
}
