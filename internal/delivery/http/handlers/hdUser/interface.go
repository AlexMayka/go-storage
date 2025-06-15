package hdUser

import (
	"context"
	"github.com/gin-gonic/gin"
	"go-storage/internal/domain"
)

type UserHandlerInterface interface {
	GetUserByID(ctx *gin.Context)
	GetUsersByCompany(ctx *gin.Context)
	UpdateUser(ctx *gin.Context)
	ChangePassword(ctx *gin.Context)
	DeactivateUser(ctx *gin.Context)
}

type UseCaseUserInterface interface {
	RegisterUser(ctx context.Context, u *domain.User) (*domain.User, error)
	GetUserByID(ctx context.Context, id string) (*domain.User, error)
	GetUsersByCompany(ctx context.Context, companyID string) ([]*domain.User, error)
	UpdateUser(ctx context.Context, userID string, u *domain.User) (*domain.User, error)
	ChangePassword(ctx context.Context, userID, oldPassword, newPassword string) error
	DeactivateUser(ctx context.Context, userID string) error
}

type UseCaseAuthInterface interface {
	GetRoleByName(ctx context.Context, roleName string) (*domain.Role, error)
}
