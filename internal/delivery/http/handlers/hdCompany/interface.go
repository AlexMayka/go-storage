package hdCompany

import (
	"context"
	"github.com/gin-gonic/gin"
	"go-storage/internal/domain"
)

type CompanyHandlerInterface interface {
	RegistrationCompany(c *gin.Context)
	GetCompanyById(ctx *gin.Context)
	GetAllCompanies(ctx *gin.Context)
	DeleteCompany(ctx *gin.Context)
	UpdateCompany(ctx *gin.Context)

	GetMyCompany(ctx *gin.Context)
	UpdateMyCompany(ctx *gin.Context)
}

type UseCaseCompanyInterface interface {
	RegisterCompany(ctx context.Context, c *domain.Company) (*domain.Company, error)
	GetCompanyById(ctx context.Context, id string) (*domain.Company, error)
	GetAllCompanies(ctx context.Context) ([]*domain.Company, error)
	DeleteCompany(ctx context.Context, id string) error
	UpdateCompany(ctx context.Context, id string, c *domain.Company) (*domain.Company, error)
}
