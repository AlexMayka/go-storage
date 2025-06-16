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
	{
		companyAdmin := companies.Group("/")
		companyAdmin.Use(authMiddleware.RequireAnyPermission([]string{"company:read:all", "company:create", "company:update:all", "company:delete"}))
		{
			companyAdmin.GET("/", CompanyHandler.GetAllCompanies)
			companyAdmin.POST("/", CompanyHandler.RegistrationCompany)
			companyAdmin.GET("/:id", CompanyHandler.GetCompanyById)
			companyAdmin.DELETE("/:id", CompanyHandler.DeleteCompany)
			companyAdmin.PUT("/:id", CompanyHandler.UpdateCompany)
		}

		companyOwn := companies.Group("/")
		companyOwn.Use(authMiddleware.RequireAnyPermission([]string{"company:read:own", "company:update:own"}))
		{
			companyOwn.GET("/me", CompanyHandler.GetMyCompany)
			companyOwn.PUT("/me", CompanyHandler.UpdateMyCompany)
		}
	}

	users := protected.Group("/users")
	{
		userOwn := users.Group("/")
		userOwn.Use(authMiddleware.RequireAnyPermission([]string{"user:read", "user:update"}))
		{
			userOwn.GET("/me", UserHandler.GetMe)
			userOwn.PUT("/me", UserHandler.UpdateMe)
			userOwn.PUT("/me/password", UserHandler.UpdatePasswordMe)
		}

		// Company user management endpoints (for company_admin and super_admin)
		userCompany := users.Group("/")
		userCompany.Use(authMiddleware.RequireAnyPermission([]string{"user:read_company", "user:manage_company"}))
		{
			userCompany.GET("/company", UserHandler.GetAllUsersOfYourCompany)
		}

		// Super admin user management endpoints
		userAdmin := users.Group("/")
		userAdmin.Use(authMiddleware.RequireAnyPermission([]string{"user:create", "user:read", "user:update", "user:delete"}))
		{
			userAdmin.POST("/", UserHandler.RegistrationUser)
			userAdmin.GET("/", UserHandler.GetAllUsers)
			userAdmin.GET("/:id", UserHandler.GetUserByID)
			userAdmin.PUT("/:id", UserHandler.UpdateUser)
			userAdmin.PUT("/:id/password", UserHandler.ChangePassword)
			userAdmin.PUT("/:id/role", UserHandler.UpdateUserRole)
			userAdmin.PUT("/:id/activate", UserHandler.ActivateUser)
			userAdmin.DELETE("/:id", UserHandler.DeactivateUser)
			userAdmin.PUT("/:id/transfer", UserHandler.TransferUserToCompany)
		}
	}

	return r
}
