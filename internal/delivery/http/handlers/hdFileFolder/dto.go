package hdFileFolder

import (
	"go-storage/internal/domain"
	"time"
)

type FolderFileDTO struct {
	ID           string          `json:"id"`
	Name         string          `json:"name"`
	Type         domain.FileType `json:"type"`
	FullPath     string          `json:"full_path"`
	ParentID     *string         `json:"parent_id,omitempty"`
	CompanyID    string          `json:"company_id"`
	UserCreateID string          `json:"user_created_id"`

	MimeType    *string `json:"mime_type,omitempty"`
	Size        *int64  `json:"size,omitempty"`
	StoragePath *string `json:"storage_path,omitempty"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type RequestCreateFolder struct {
	Name       string `json:"name" binding:"required"`
	ParentPath string `json:"parentPath" binding:"required"`
}

type RequestGetFolder struct {
	Path string          `json:"path" binding:"required"`
	Type domain.FileType `json:"type,omitempty"`
}

type RequestRenameFolder struct {
	Name string `json:"name" binding:"required"`
}

type RequestMoveFolder struct {
	ParentPath string `json:"parentPath" binding:"required"`
}

type ResponseFolder struct {
	Status string         `json:"status"`
	Time   time.Time      `json:"time"`
	Folder *FolderFileDTO `json:"folder"`
}

type ResponseGetFolder struct {
	Status string           `json:"status"`
	Time   time.Time        `json:"time"`
	Files  []*FolderFileDTO `json:"files"`
}

type ResponsePath struct {
	Status string    `json:"status"`
	Time   time.Time `json:"time"`
	Path   string    `json:"path"`
}

type RequestUploadFile struct {
	ParentPath string `form:"parentPath" binding:"required"`
}

type RequestDownloadFile struct {
	ID     string `uri:"id" binding:"required,uuid"`
	Inline bool   `form:"inline"`
}

type RequestGetFileInfo struct {
	ID string `uri:"id" binding:"required,uuid"`
}

type RequestRenameFile struct {
	Name string `json:"name" binding:"required"`
}

type RequestMoveFile struct {
	ParentPath string `json:"parentPath" binding:"required"`
}

type RequestInitChunkedUpload struct {
	FileName   string `json:"fileName" binding:"required"`
	FileSize   int64  `json:"fileSize" binding:"required,min=1"`
	ParentPath string `json:"parentPath" binding:"required"`
	MimeType   string `json:"mimeType,omitempty"`
}

type RequestUploadChunk struct {
	UploadID   string `uri:"uploadId" binding:"required,uuid"`
	ChunkIndex string `uri:"chunkIndex" binding:"required"`
}

type RequestCompleteChunkedUpload struct {
	UploadID string `uri:"uploadId" binding:"required,uuid"`
}

type RequestGetChunkedUploadStatus struct {
	UploadID string `uri:"uploadId" binding:"required,uuid"`
}

type ResponseFile struct {
	Status string         `json:"status"`
	Time   time.Time      `json:"time"`
	File   *FolderFileDTO `json:"file"`
}

type ResponseInitChunkedUpload struct {
	Status      string               `json:"status"`
	Time        time.Time            `json:"time"`
	UploadID    string               `json:"upload_id"`
	ChunkSize   int64                `json:"chunk_size"`
	TotalChunks int                  `json:"total_chunks"`
	Strategy    *domain.StrategyInfo `json:"strategy"`
}

type ResponseUploadChunk struct {
	Status         string    `json:"status"`
	Time           time.Time `json:"time"`
	ChunkIndex     int       `json:"chunk_index"`
	Progress       float64   `json:"progress"`
	UploadedChunks int       `json:"uploaded_chunks"`
	TotalChunks    int       `json:"total_chunks"`
}

type ResponseChunkedUploadStatus struct {
	Status         string                     `json:"status"`
	Time           time.Time                  `json:"time"`
	UploadID       string                     `json:"upload_id"`
	FileName       string                     `json:"file_name"`
	Progress       float64                    `json:"progress"`
	UploadedChunks int                        `json:"uploaded_chunks"`
	TotalChunks    int                        `json:"total_chunks"`
	UploadedSize   int64                      `json:"uploaded_size"`
	TotalSize      int64                      `json:"total_size"`
	UploadStatus   domain.ChunkedUploadStatus `json:"upload_status"`
	MissingChunks  []int                      `json:"missing_chunks,omitempty"`
	ExpiresAt      time.Time                  `json:"expires_at"`
}

type ResponseCompleteChunkedUpload struct {
	Status string         `json:"status"`
	Time   time.Time      `json:"time"`
	File   *FolderFileDTO `json:"file"`
}

type ResponseUploadStrategy struct {
	Status   string               `json:"status"`
	Time     time.Time            `json:"time"`
	Strategy *domain.StrategyInfo `json:"strategy"`
}

type ResponseSuccess struct {
	Status  string    `json:"status"`
	Time    time.Time `json:"time"`
	Message string    `json:"message"`
}

type ResponseResourceStats struct {
	Status string                `json:"status"`
	Time   time.Time             `json:"time"`
	Stats  *domain.ResourceStats `json:"stats"`
}
