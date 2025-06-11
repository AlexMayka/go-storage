package http

import (
	"database/sql"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"go-storage/internal/delivery/http/handlers/hdCompany"
	"go-storage/internal/delivery/http/middleware"
	"go-storage/internal/repository/postgres/rpCompany"
	"go-storage/internal/usecase/ucCompany"
	"go-storage/pkg/logger"
)

func Router(log logger.Logger, db *sql.DB) *gin.Engine {
	r := gin.Default()
	r.Use(middleware.Logger(log))

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	var rpCmp ucCompany.RepositoryCompanyInterface = rpCompany.NewRepository(db)
	var ucCmp hdCompany.UseCaseCompanyInterface = ucCompany.NewUseCase(rpCmp)
	var hdCmp hdCompany.CompanyHandlerInterface = hdCompany.NewHandlerCompany(ucCmp)

	api := r.Group("/api/v1/")
	company := api.Group("/companies")
	{
		company.GET("/", hdCmp.GetAllCompanies)
		company.POST("/", hdCmp.RegistrationCompany)
		company.GET("/:id", hdCmp.GetCompanyById)
		company.DELETE("/:id", hdCmp.DeleteCompany)
		company.PUT("/:id", hdCmp.UpdateCompany)
	}

	return r
}
