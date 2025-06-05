package ucCompany

import (
	"context"
	"go-storage/internal/domain"
)

type UseCaseCompanyInterface interface {
	RegisterCompany(ctx context.Context, c *domain.Company) (*domain.Company, error)
	GetCompanyById(ctx context.Context, id string) (*domain.Company, error)
	GetAllCompanies(ctx context.Context) ([]*domain.Company, error)
	DeleteCompany(ctx context.Context, id string) error
}
