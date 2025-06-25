package http

import (
	"context"
	"database/sql"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"go-storage/internal/config"
	"go-storage/internal/delivery/http/handlers/hdAuth"
	"go-storage/internal/delivery/http/handlers/hdCompany"
	"go-storage/internal/delivery/http/handlers/hdFileFolder"
	"go-storage/internal/delivery/http/handlers/hdUser"
	"go-storage/internal/delivery/http/middleware"
	"go-storage/internal/repository/minio"
	"go-storage/internal/repository/postgres/rpAuth"
	"go-storage/internal/repository/postgres/rpChunkedUpload"
	"go-storage/internal/repository/postgres/rpCompany"
	"go-storage/internal/repository/postgres/rpFiles"
	"go-storage/internal/repository/postgres/rpUser"
	"go-storage/internal/usecase/ucAuthUser"
	"go-storage/internal/usecase/ucCompany"
	"go-storage/internal/usecase/ucFileFolder"
	"go-storage/internal/usecase/ucUser"
	"go-storage/pkg/logger"
	"go-storage/pkg/storage"
)

func Router(log logger.Logger, db *sql.DB, cnf config.Config) *gin.Engine {
	r := gin.Default()
	r.Use(middleware.Logger(log), middleware.Config(cnf))

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	var CompanyRepo = rpCompany.NewRepository(db)
	var AuthRepo = rpAuth.NewRepositoryAuth(db)
	var UserRepo = rpUser.NewRepository(db)
	// Initialize MinIO client
	minioClient, err := storage.NewMinIOClient(cnf.Minio)
	if err != nil {
		panic("Failed to initialize MinIO client: " + err.Error())
	}

	// Ensure bucket exists
	if err := storage.EnsureBucket(context.Background(), minioClient, cnf.Minio.BucketName); err != nil {
		panic("Failed to ensure bucket exists: " + err.Error())
	}

	var FilesRepo = rpFiles.NewRepository(db)
	var ChunkedUploadRepo = rpChunkedUpload.NewRepository(db)
	var StorageRepo = minio.NewStorageRepository(minioClient, cnf.Minio.BucketName)

	var CompanyUseCase = ucCompany.NewUseCase(CompanyRepo)
	var AuthUseCase = ucAuthUser.NewUseCaseAuth(AuthRepo)
	var UserUseCase = ucUser.NewUseCaseUser(UserRepo, AuthRepo)
	// Initialize file system UseCase
	var FileFolderUseCase = ucFileFolder.NewUseCaseFileFolder(FilesRepo, StorageRepo, ChunkedUploadRepo, &cnf.FileServer)

	var CompanyHandler = hdCompany.NewHandlerCompany(CompanyUseCase)
	var AuthHandler = hdAuth.NewHandlerAuth(UserUseCase, AuthUseCase)
	var UserHandler = hdUser.NewHandlerUser(UserUseCase, AuthUseCase)
	var FileFolderHandler = hdFileFolder.NewHandlerFileFolder(FileFolderUseCase)

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

	// File and folder management endpoints
	files := protected.Group("/files")
	files.Use(authMiddleware.RequireAnyPermission([]string{"file:read", "file:write", "file:delete"}))
	{
		// File operations
		files.POST("/upload", FileFolderHandler.UploadFile)
		files.GET("/:id", FileFolderHandler.GetFileInfo)
		files.GET("/:id/download", FileFolderHandler.DownloadFile)
		files.PUT("/:id/rename", FileFolderHandler.RenameFile)
		files.PUT("/:id/move", FileFolderHandler.MoveFile)
		files.DELETE("/:id", FileFolderHandler.DeleteFile)

		// Upload strategy
		files.GET("/upload-strategy", FileFolderHandler.GetUploadStrategy)

		// Chunked upload
		chunked := files.Group("/chunked")
		{
			chunked.POST("/init", FileFolderHandler.InitChunkedUpload)
			chunked.POST("/:uploadId/chunk/:chunkIndex", FileFolderHandler.UploadChunk)
			chunked.GET("/:uploadId/status", FileFolderHandler.GetChunkedUploadStatus)
			chunked.POST("/:uploadId/complete", FileFolderHandler.CompleteChunkedUpload)
			chunked.DELETE("/:uploadId/abort", FileFolderHandler.AbortChunkedUpload)
		}

		// Resource monitoring
		files.GET("/stats", FileFolderHandler.GetResourceStats)
	}

	folders := protected.Group("/folders")
	folders.Use(authMiddleware.RequireAnyPermission([]string{"file:read", "file:write", "file:delete"}))
	{
		folders.POST("/", FileFolderHandler.CreateFolder)
		folders.POST("/contents", FileFolderHandler.GetFolderContents)
		folders.PUT("/:path/rename", FileFolderHandler.FolderRename)
		folders.PUT("/:path/move", FileFolderHandler.MoveFolder)
		folders.DELETE("/:path", FileFolderHandler.DeleteFolder)
	}

	return r
}
