package hdUser

import (
	"context"
	"github.com/gin-gonic/gin"
	"go-storage/internal/domain"
)

type UserHandlerInterface interface {
	GetMe(ctx *gin.Context)
	UpdateMe(ctx *gin.Context)
	UpdatePasswordMe(ctx *gin.Context)
	GetAllUsersOfYourCompany(ctx *gin.Context)
	RegistrationUser(ctx *gin.Context)
	GetUserByID(ctx *gin.Context)
	UpdateUser(ctx *gin.Context)
	ChangePassword(ctx *gin.Context)
	DeactivateUser(ctx *gin.Context)
	UpdateUserRole(ctx *gin.Context)
	ActivateUser(ctx *gin.Context)
	GetAllUsers(ctx *gin.Context)
	TransferUserToCompany(ctx *gin.Context)
}

type UseCaseUserInterface interface {
	RegisterUser(ctx context.Context, u *domain.User) (*domain.User, error)
	GetUserByID(ctx context.Context, id string) (*domain.User, error)
	GetUsersByCompany(ctx context.Context, companyID string) ([]*domain.User, error)
	UpdateUser(ctx context.Context, userID string, u *domain.User) (*domain.User, error)
	ChangePassword(ctx context.Context, userID, oldPassword, newPassword string) error
	DeactivateUser(ctx context.Context, userID string) error
	UpdateUserRole(ctx context.Context, userID string, roleID string) error
	ActivateUser(ctx context.Context, userID string) error
	GetAllUsers(ctx context.Context) ([]*domain.User, error)
	TransferUserToCompany(ctx context.Context, userID string, companyID string) error
}

type UseCaseAuthInterface interface {
	GetRoleByName(ctx context.Context, roleName string) (*domain.Role, error)
	GetDefaultRole(ctx context.Context) (domain.Role, error)
}
