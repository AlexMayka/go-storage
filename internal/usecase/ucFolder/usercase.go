package ucFolder

import (
	"context"
	"go-storage/internal/domain"
)

type UseCaseFolder struct {
	repo RepositoryFolderInterface
}

func NewUseCaseFolder(repo RepositoryFolderInterface) *UseCaseFolder {
	return &UseCaseFolder{repo: repo}
}

func (u *UseCaseFolder) GetFolders(ctx context.Context) ([]*domain.Folder, error) {
	return u.repo.GetFolders(ctx)
}

func (u *UseCaseFolder) GetFolderById(ctx context.Context, id string) (*domain.Folder, error) {
	return u.repo.GetFolderById(ctx, id)
}

func (u *UseCaseFolder) GetFolderByName(ctx context.Context, name, companyId string) (*domain.Folder, error) {
	return u.repo.GetFolderByName(ctx, name, companyId)
}

func (u *UseCaseFolder) GetFolderByPath(ctx context.Context, path string) (*domain.Folder, error) {
	return u.repo.GetFolderByPath(ctx, path)
}

func (u *UseCaseFolder) GetFolderByParentId(ctx context.Context, parentId string) (*domain.Folder, error) {
	return u.repo.GetFolderByPath(ctx, parentId)
}

func (u *UseCaseFolder) GetFoldersByCompanyId(ctx context.Context, companyId string) ([]*domain.Folder, error) {
	return u.repo.GetFoldersByCompanyId(ctx, companyId)
}
