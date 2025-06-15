package http

import (
	"database/sql"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"go-storage/internal/config"
	"go-storage/internal/delivery/http/handlers/hdAuth"
	"go-storage/internal/delivery/http/handlers/hdCompany"
	"go-storage/internal/delivery/http/handlers/hdUser"
	"go-storage/internal/delivery/http/middleware"
	"go-storage/internal/repository/postgres/rpAuth"
	"go-storage/internal/repository/postgres/rpCompany"
	"go-storage/internal/repository/postgres/rpUser"
	"go-storage/internal/usecase/ucAuthUser"
	"go-storage/internal/usecase/ucCompany"
	"go-storage/internal/usecase/ucUser"
	"go-storage/pkg/logger"
)

func Router(log logger.Logger, db *sql.DB, cnf config.Config) *gin.Engine {
	r := gin.Default()
	r.Use(middleware.Logger(log), middleware.Config(cnf))

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	var CompanyRepo = rpCompany.NewRepository(db)
	var AuthRepo = rpAuth.NewRepositoryAuth(db)
	var UserRepo = rpUser.NewRepository(db)

	var CompanyUseCase = ucCompany.NewUseCase(CompanyRepo)
	var AuthUseCase = ucAuthUser.NewUseCaseAuth(AuthRepo)
	var UserUseCase = ucUser.NewUseCaseUser(UserRepo)

	var CompanyHandler = hdCompany.NewHandlerCompany(CompanyUseCase)
	var AuthHandler = hdAuth.NewHandlerAuth(UserUseCase, AuthRepo)
	var UserHandler = hdUser.NewHandlerUser(UserUseCase, AuthUseCase)

	authMiddleware := middleware.NewAuthMiddleware(AuthUseCase)

	api := r.Group("/api/v1/")

	auth := api.Group("/auth")
	{
		auth.POST("/login", AuthHandler.Login)
		auth.POST("/refresh-token", AuthHandler.RefreshToken)
	}

	api.POST("/users/register", UserHandler.RegistrationUser)

	protected := api.Group("/")
	protected.Use(authMiddleware.RequireAuth())

	companies := protected.Group("/companies")
	companies.Use(authMiddleware.RequireAnyPermission([]string{"company:read", "company:create", "company:update", "company:delete"}))
	{
		companies.GET("/", CompanyHandler.GetAllCompanies)
		companies.POST("/", CompanyHandler.RegistrationCompany)
		companies.GET("/:id", CompanyHandler.GetCompanyById)
		companies.DELETE("/:id", CompanyHandler.DeleteCompany)
		companies.PUT("/:id", CompanyHandler.UpdateCompany)
	}

	users := protected.Group("/users")
	{
		users.GET("/:id", UserHandler.GetUserByID)
		users.PUT("/:id", UserHandler.UpdateUser)
		users.PUT("/:id/password", UserHandler.ChangePassword)
		users.DELETE("/:id", UserHandler.DeactivateUser)
		users.GET("/company/:company_id", UserHandler.GetUsersByCompany)
	}

	return r
}
