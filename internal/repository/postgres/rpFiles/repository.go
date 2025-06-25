package rpFiles

import (
	"context"
	"database/sql"
	"errors"
	"strings"
	"time"

	"go-storage/internal/domain"
	pkgErrors "go-storage/pkg/errors"
)

type RepositoryFiles struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *RepositoryFiles {
	return &RepositoryFiles{db: db}
}

func (r *RepositoryFiles) CreateFile(ctx context.Context, file *domain.File) (*domain.File, error) {
	_, err := r.db.ExecContext(ctx, QueryCreateFile,
		file.ID, file.Name, file.Type, file.FullPath.String(), file.ParentID, file.CompanyId, file.UserCreateID,
		file.MimeType, file.Size, file.Hash, file.StoragePath,
		file.CreatedAt, file.UpdatedAt, file.IsActive,
	)
	if err != nil {
		if strings.Contains(err.Error(), "idx_unique_name_in_folder") {
			return nil, pkgErrors.BadRequest("file with this name already exists in the folder")
		}
		return nil, pkgErrors.Database("unable to create file")
	}
	return file, nil
}

func (r *RepositoryFiles) GetFile(ctx context.Context, companyID, fileID string) (*domain.File, error) {
	var file domain.File
	var fullPathStr string

	row := r.db.QueryRowContext(ctx, QueryGetFile, fileID, companyID)

	err := row.Scan(
		&file.ID, &file.Name, &file.Type, &fullPathStr, &file.ParentID, &file.CompanyId, &file.UserCreateID,
		&file.MimeType, &file.Size, &file.Hash, &file.StoragePath,
		&file.CreatedAt, &file.UpdatedAt, &file.IsActive,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, pkgErrors.NotFound("file not found")
		}
		return nil, pkgErrors.Database("unable to get file")
	}

	fullPath, err := domain.NewPath(fullPathStr)
	if err != nil {
		return nil, pkgErrors.Database("invalid file path")
	}
	file.FullPath = fullPath

	return &file, nil
}

func (r *RepositoryFiles) GetFileByPath(ctx context.Context, companyID string, path *domain.Path) (*domain.File, error) {
	var file domain.File
	var fullPathStr string

	row := r.db.QueryRowContext(ctx, QueryGetFileByPath, path.String(), companyID)

	err := row.Scan(
		&file.ID, &file.Name, &file.Type, &fullPathStr, &file.ParentID, &file.CompanyId, &file.UserCreateID,
		&file.MimeType, &file.Size, &file.Hash, &file.StoragePath,
		&file.CreatedAt, &file.UpdatedAt, &file.IsActive,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, pkgErrors.NotFound("file not found")
		}
		return nil, pkgErrors.Database("unable to get file")
	}

	fullPath, err := domain.NewPath(fullPathStr)
	if err != nil {
		return nil, pkgErrors.Database("invalid file path")
	}
	file.FullPath = fullPath

	return &file, nil
}

func (r *RepositoryFiles) GetFolderContents(ctx context.Context, companyID string, path *domain.Path, fileType *domain.FileType) ([]*domain.File, error) {
	var rows *sql.Rows
	var err error

	if fileType != nil {
		parentFolder, err := r.GetFileByPath(ctx, companyID, path)
		if err != nil {
			return nil, err
		}
		rows, err = r.db.QueryContext(ctx, QueryGetFolderContentsByType, parentFolder.ID, companyID, *fileType)
	} else {
		pathPattern := path.String()
		if pathPattern == "/" {
			pathPattern = ""
		}
		pathPattern += "/%"

		rows, err = r.db.QueryContext(ctx, QueryGetFolderContentsByPath, pathPattern, companyID, path.String())
	}

	if err != nil {
		return nil, pkgErrors.Database("unable to get folder contents")
	}
	defer rows.Close()

	var files []*domain.File
	for rows.Next() {
		file, err := scanFile(rows)
		if err != nil {
			return nil, pkgErrors.Database("unable to scan file")
		}
		files = append(files, file)
	}

	if err := rows.Err(); err != nil {
		return nil, pkgErrors.Database("unable to get folder contents")
	}

	return files, nil
}

func (r *RepositoryFiles) UpdateFile(ctx context.Context, file *domain.File) (*domain.File, error) {
	file.UpdatedAt = time.Now()

	_, err := r.db.ExecContext(ctx, QueryUpdateFile,
		file.ID, file.Name, file.FullPath.String(), file.ParentID, file.MimeType,
		file.Size, file.Hash, file.StoragePath, file.UpdatedAt, file.CompanyId,
	)
	if err != nil {
		return nil, pkgErrors.Database("unable to update file")
	}

	return file, nil
}

func (r *RepositoryFiles) RenameFile(ctx context.Context, companyID, fileID, newName string) (*domain.File, error) {
	file, err := r.GetFile(ctx, companyID, fileID)
	if err != nil {
		return nil, err
	}

	newPath := file.FullPath.GetParent().Join(newName)
	_, err = r.db.ExecContext(ctx, QueryUpdateFileName,
		fileID, newName, newPath.String(), time.Now(), companyID,
	)
	if err != nil {
		if strings.Contains(err.Error(), "idx_unique_name_in_folder") {
			return nil, pkgErrors.BadRequest("file with this name already exists")
		}
		return nil, pkgErrors.Database("unable to rename file")
	}

	file.Name = newName
	file.FullPath = newPath
	file.UpdatedAt = time.Now()

	return file, nil
}

func (r *RepositoryFiles) MoveFile(ctx context.Context, companyID, fileID string, newParentPath *domain.Path) (*domain.File, error) {
	file, err := r.GetFile(ctx, companyID, fileID)
	if err != nil {
		return nil, err
	}

	newPath := newParentPath.Join(file.Name)
	var newParentID *string
	if newParentPath.String() != "/" {
		parent, err := r.GetFileByPath(ctx, companyID, newParentPath)
		if err != nil {
			return nil, pkgErrors.BadRequest("destination folder not found")
		}
		if parent.Type != domain.FileTypeFolder {
			return nil, pkgErrors.BadRequest("destination must be a folder")
		}
		newParentID = &parent.ID
	}

	_, err = r.db.ExecContext(ctx, QueryUpdateFileParent,
		fileID, newParentID, newPath.String(), time.Now(), companyID,
	)
	if err != nil {
		if strings.Contains(err.Error(), "idx_unique_name_in_folder") {
			return nil, pkgErrors.BadRequest("file with this name already exists in destination")
		}
		return nil, pkgErrors.Database("unable to move file")
	}

	file.ParentID = newParentID
	file.FullPath = newPath
	file.UpdatedAt = time.Now()

	return file, nil
}

func (r *RepositoryFiles) DeleteFile(ctx context.Context, companyID, fileID string) error {
	_, err := r.db.ExecContext(ctx, QueryDeleteFile, fileID, companyID, time.Now())
	if err != nil {
		return pkgErrors.Database("unable to delete file")
	}
	return nil
}

func (r *RepositoryFiles) CreateFolder(ctx context.Context, folder *domain.File) (*domain.File, error) {
	_, err := r.db.ExecContext(ctx, QueryCreateFolder,
		folder.ID, folder.Name, folder.Type, folder.FullPath.String(), folder.ParentID, folder.CompanyId, folder.UserCreateID,
		folder.CreatedAt, folder.UpdatedAt, folder.IsActive,
	)
	if err != nil {
		if strings.Contains(err.Error(), "idx_unique_name_in_folder") {
			return nil, pkgErrors.BadRequest("folder with this name already exists")
		}
		return nil, pkgErrors.Database("unable to create folder")
	}
	return folder, nil
}

func (r *RepositoryFiles) GetFolder(ctx context.Context, companyID string, path *domain.Path) (*domain.File, error) {
	var folder domain.File
	var fullPathStr string

	row := r.db.QueryRowContext(ctx, QueryGetFolder, path.String(), companyID)

	err := row.Scan(
		&folder.ID, &folder.Name, &folder.Type, &fullPathStr, &folder.ParentID, &folder.CompanyId, &folder.UserCreateID,
		&folder.MimeType, &folder.Size, &folder.Hash, &folder.StoragePath,
		&folder.CreatedAt, &folder.UpdatedAt, &folder.IsActive,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, pkgErrors.NotFound("folder not found")
		}
		return nil, pkgErrors.Database("unable to get folder")
	}

	fullPath, err := domain.NewPath(fullPathStr)
	if err != nil {
		return nil, pkgErrors.Database("invalid folder path")
	}
	folder.FullPath = fullPath

	return &folder, nil
}

func (r *RepositoryFiles) MoveFolder(ctx context.Context, companyID string, oldPath, newPath *domain.Path) (*domain.Path, error) {
	_, err := r.GetFolder(ctx, companyID, oldPath)
	if err != nil {
		return nil, err
	}

	oldPathStr := oldPath.String()
	newPathStr := newPath.String()
	pathPattern := oldPathStr + "%"

	_, err = r.db.ExecContext(ctx, QueryMoveFolderAndContents,
		oldPathStr, newPathStr, pathPattern, companyID, time.Now(),
	)
	if err != nil {
		if strings.Contains(err.Error(), "idx_unique_name_in_folder") {
			return nil, pkgErrors.BadRequest("folder with this name already exists at destination")
		}
		return nil, pkgErrors.Database("unable to move folder")
	}

	return newPath, nil
}

func (r *RepositoryFiles) DeleteFolder(ctx context.Context, companyID string, folderPath *domain.Path) error {
	_, err := r.db.ExecContext(ctx, QueryDeleteFolder, folderPath.String(), companyID, time.Now())
	if err != nil {
		return pkgErrors.Database("unable to delete folder")
	}
	return nil
}

func scanFile(rows *sql.Rows) (*domain.File, error) {
	var file domain.File
	var fullPathStr string

	err := rows.Scan(
		&file.ID, &file.Name, &file.Type, &fullPathStr, &file.ParentID, &file.CompanyId, &file.UserCreateID,
		&file.MimeType, &file.Size, &file.Hash, &file.StoragePath,
		&file.CreatedAt, &file.UpdatedAt, &file.IsActive,
	)
	if err != nil {
		return nil, err
	}

	fullPath, err := domain.NewPath(fullPathStr)
	if err != nil {
		return nil, err
	}
	file.FullPath = fullPath

	return &file, nil
}
