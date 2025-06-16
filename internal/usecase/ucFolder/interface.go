package ucFolder

import (
	"context"
	"go-storage/internal/domain"
)

type RepositoryFolderInterface interface {
	GetFolders(ctx context.Context) ([]*domain.Folder, error)
	GetFolderByName(ctx context.Context, name string, companyId string) (*domain.Folder, error)
	GetFolderById(ctx context.Context, id string) (*domain.Folder, error)
	GetFolderByPath(ctx context.Context, path string) (*domain.Folder, error)
	GetFoldersByParentId(ctx context.Context, parentId string) ([]*domain.Folder, error)
	GetFoldersByCompanyId(ctx context.Context, companyId string) ([]*domain.Folder, error)
	GetFoldersByUserCreateId(ctx context.Context, userCreateId string) ([]*domain.Folder, error)
	UpdateFolder(ctx context.Context, f *domain.Folder) (*domain.Folder, error)
	InsertFolder(ctx context.Context, f *domain.Folder) (*domain.Folder, error)
}
