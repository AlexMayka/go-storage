package ucCompany

import (
	"context"
	"github.com/google/uuid"
	"go-storage/internal/domain"
	"go-storage/internal/utils/valid"
	"go-storage/pkg/errors"
)

type UseCaseCompany struct {
	repo RepositoryCompanyInterface
}

func NewUseCase(repo RepositoryCompanyInterface) *UseCaseCompany {
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
	c.Path = valid.NormalizationOfName(c.Name)

	return u.repo.Create(ctx, c)
}

func (u *UseCaseCompany) UpdateCompany(ctx context.Context, id string, c *domain.Company) (*domain.Company, error) {
	if c.Name == "" && c.Description == "" {
		return nil, errors.EmptyField("name and description")
	}

	company, err := u.repo.GetCompanyById(ctx, id)
	if err != nil {
		return nil, err
	}

	if c.Name == company.Name && c.Description == company.Description {
		return company, nil
	}

	if c.Name != company.Name && c.Name != "" {
		company.Name = c.Name
	}

	if c.Description != company.Description && c.Description != "" {
		company.Description = c.Description
	}

	return u.repo.Update(ctx, company)
}

func (u *UseCaseCompany) GetCompanyById(ctx context.Context, id string) (*domain.Company, error) {
	return u.repo.GetCompanyById(ctx, id)
}

func (u *UseCaseCompany) GetAllCompanies(ctx context.Context) ([]*domain.Company, error) {
	return u.repo.GetAllCompanies(ctx)
}

func (u *UseCaseCompany) DeleteCompany(ctx context.Context, id string) error {
	return u.repo.UpdateIsActive(ctx, id, false)
}
