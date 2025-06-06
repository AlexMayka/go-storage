package rpCompany

import (
	"context"
	"go-storage/internal/domain"
)

type RepositoryInterface interface {
	Create(ctx context.Context, c *domain.Company) (*domain.Company, error)
	GetCompanyById(ctx context.Context, id string) (*domain.Company, error)
	GetAllCompanies(ctx context.Context) ([]*domain.Company, error)
	DeleteCompany(ctx context.Context, id string) error
	UpdateIsActive(ctx context.Context, id string, on bool) error
	Update(ctx context.Context, c *domain.Company) (*domain.Company, error)
}
