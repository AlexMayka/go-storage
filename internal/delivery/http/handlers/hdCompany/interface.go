package hdCompany

import "github.com/gin-gonic/gin"

type CompanyHandlerInterface interface {
	RegistrationCompany(c *gin.Context)
	GetCompanyById(ctx *gin.Context)
	GetAllCompanies(ctx *gin.Context)
	DeleteCompany(ctx *gin.Context)
}
