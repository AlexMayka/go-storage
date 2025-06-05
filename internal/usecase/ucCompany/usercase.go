package ucCompany

import (
	"context"
	"github.com/google/uuid"
	"go-storage/internal/domain"
	"go-storage/internal/repository/postgres/rpCompany"
	"go-storage/internal/utils"
	"go-storage/pkg/errors"
)

type UseCaseCompany struct {
	repo rpCompany.RepositoryInterface
}

func NewUseCase(repo rpCompany.RepositoryInterface) *UseCaseCompany {
	return &UseCaseCompany{
		repo: repo,
	}
}

func (u *UseCaseCompany) RegisterCompany(ctx context.Context, c *domain.Company) (*domain.Company, error) {
	if c.Name == "" {
		return nil, errors.EmptyField("name")
	}

	if c.Description == "" {
		return nil, errors.EmptyField("description")
	}

	c.ID = uuid.NewString()
	c.Path = utils.NormalizationOfName(c.Name)

	return u.repo.Create(ctx, c)
}

func (u *UseCaseCompany) GetCompanyById(ctx context.Context, id string) (*domain.Company, error) {
	return u.repo.GetCompanyById(ctx, id)
}

func (u *UseCaseCompany) GetAllCompanies(ctx context.Context) ([]*domain.Company, error) {
	return u.repo.GetAllCompanies(ctx)
}

func (u *UseCaseCompany) DeleteCompany(ctx context.Context, id string) error {
	return u.repo.DeleteCompany(ctx, id)
}
