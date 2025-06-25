package domain

import (
	"errors"
	"fmt"
	"time"
)

type File struct {
	ID           string
	Name         string
	Type         FileType
	FullPath     Path
	ParentID     *string
	CompanyId    string
	UserCreateID string

	MimeType    *string
	Size        *int64
	Hash        *string
	StoragePath *string

	CreatedAt time.Time
	UpdatedAt time.Time
	IsActive  bool
}

type StorageFileInfo struct {
	Key          string
	Size         int64
	MimeType     string
	ETag         string
	LastModified time.Time
}

func (f *File) Validate() error {
	if f.Name == "" {
		return errors.New("file name cannot be empty")
	}
	if f.CompanyId == "" {
		return errors.New("company ID cannot be empty")
	}
	if f.UserCreateID == "" {
		return errors.New("user ID cannot be empty")
	}
	if f.Type == FileTypeFile && (f.Size == nil || *f.Size < 0) {
		return errors.New("file size must be non-negative")
	}
	return nil
}

func (f *File) IsFolder() bool {
	return f.Type == FileTypeFolder
}

func (f *File) IsFile() bool {
	return f.Type == FileTypeFile
}

func (f *File) GetDisplaySize() string {
	if f.Size == nil {
		return "0 B"
	}

	size := *f.Size
	if size < 1024 {
		return fmt.Sprintf("%d B", size)
	} else if size < 1024*1024 {
		return fmt.Sprintf("%.1f KB", float64(size)/1024)
	} else if size < 1024*1024*1024 {
		return fmt.Sprintf("%.1f MB", float64(size)/(1024*1024))
	} else {
		return fmt.Sprintf("%.1f GB", float64(size)/(1024*1024*1024))
	}
}
