package rpFolder

import (
	"context"
	"database/sql"
	"errors"
	"go-storage/internal/domain"
	pkgErrors "go-storage/pkg/errors"
	"time"
)

type RepositoryFolder struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *RepositoryFolder { return &RepositoryFolder{db: db} }

func scanFolder(scanner interface{ Scan(...interface{}) error }) (*domain.Folder, error) {
	var folder domain.Folder
	err := scanner.Scan(&folder.ID, &folder.Name, &folder.Path, &folder.ParentId, &folder.CompanyId, &folder.UserCreateId, &folder.CreatedAt, &folder.UpdatedAt, &folder.IsActive)
	return &folder, err
}

func (r *RepositoryFolder) GetFolderByName(ctx context.Context, name string, companyId string) (*domain.Folder, error) {
	row := r.db.QueryRowContext(ctx, QueryGetFolderByName, name, companyId)

	folder, err := scanFolder(row)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, pkgErrors.NotFound("folder not found")
		}
		return nil, pkgErrors.NotFound("folder not found")
	}

	return folder, nil
}

func (r *RepositoryFolder) GetFolderById(ctx context.Context, id string, companyId string) (*domain.Folder, error) {
	row := r.db.QueryRowContext(ctx, QueryGetFolderById, id, companyId)

	folder, err := scanFolder(row)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, pkgErrors.NotFound("folder not found")
		}
		return nil, pkgErrors.NotFound("folder not found")
	}

	return folder, nil
}

func (r *RepositoryFolder) GetFolderByPath(ctx context.Context, path string, companyId string) (*domain.Folder, error) {
	row := r.db.QueryRowContext(ctx, QueryGetFolderByPath, path, companyId)

	folder, err := scanFolder(row)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, pkgErrors.NotFound("folder not found")
		}
		return nil, pkgErrors.NotFound("folder not found")
	}
	return folder, nil
}

func (r *RepositoryFolder) GetFoldersByParentId(ctx context.Context, parentId string, companyId string) ([]*domain.Folder, error) {
	var folders []*domain.Folder
	rows, err := r.db.QueryContext(ctx, QueryGetFoldersByParentId, parentId, companyId)
	if err != nil {
		return nil, pkgErrors.Database("unable to query all folders by parentId")
	}

	defer func() {
		_ = rows.Close()
	}()

	for rows.Next() {
		folder, err := scanFolder(rows)
		if err != nil {
			return nil, pkgErrors.Database("unable to query all folders by parentId")
		}
		folders = append(folders, folder)
	}

	if err := rows.Err(); err != nil {
		return nil, pkgErrors.Database("unable to query all folders by parentId")
	}

	return folders, nil
}

func (r *RepositoryFolder) GetFoldersByCompanyId(ctx context.Context, companyId string) ([]*domain.Folder, error) {
	var folders []*domain.Folder
	rows, err := r.db.QueryContext(ctx, QueryGetFoldersByCompanyId, companyId)
	if err != nil {
		return nil, pkgErrors.Database("unable to query all folders by companyId")
	}

	defer func() {
		_ = rows.Close()
	}()

	for rows.Next() {
		folder, err := scanFolder(rows)
		if err != nil {
			return nil, pkgErrors.Database("unable to query all folders by companyId")
		}
		folders = append(folders, folder)
	}

	if err := rows.Err(); err != nil {
		return nil, pkgErrors.Database("unable to query all folders by companyId")
	}

	return folders, nil
}


func (r *RepositoryFolder) UpdateFolder(ctx context.Context, f *domain.Folder) (*domain.Folder, error) {
	updateDate := time.Now()

	_, err := r.db.ExecContext(ctx, QueryUpdateFolder, f.Name, f.Path, f.ParentId, updateDate, f.IsActive, f.ID)
	if err != nil {
		return nil, pkgErrors.Database("unable to update folder")
	}
	return &domain.Folder{
		ID:           f.ID,
		Name:         f.Name,
		Path:         f.Path,
		ParentId:     f.ParentId,
		CompanyId:    f.CompanyId,
		UserCreateId: f.UserCreateId,
		CreatedAt:    f.CreatedAt,
		UpdatedAt:    updateDate,
		IsActive:     f.IsActive,
	}, nil
}

func (r *RepositoryFolder) InsertFolder(ctx context.Context, f *domain.Folder) (*domain.Folder, error) {
	const isActive = true
	createDate := time.Now()

	_, err := r.db.ExecContext(ctx, QueryInsertFolder, f.ID, f.Name, f.Path, f.ParentId, f.CompanyId, f.UserCreateId, createDate, createDate, isActive)
	if err != nil {
		return nil, pkgErrors.Database("unable to insert folder")
	}
	return &domain.Folder{
		ID:           f.ID,
		Name:         f.Name,
		Path:         f.Path,
		ParentId:     f.ParentId,
		CompanyId:    f.CompanyId,
		UserCreateId: f.UserCreateId,
		CreatedAt:    createDate,
		UpdatedAt:    createDate,
		IsActive:     isActive,
	}, nil
}
