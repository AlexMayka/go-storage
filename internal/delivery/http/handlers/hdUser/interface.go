package hdUser

import (
	"context"
	"github.com/gin-gonic/gin"
	"go-storage/internal/domain"
)

type UserHandlerInterface interface {
	RegistrationUser(ctx *gin.Context)
}

type UseCaseUserInterface interface {
	RegisterUser(ctx context.Context, u *domain.User) (*domain.User, error)
}

type UseCaseAuthInterface interface {
	GetRoleByName(ctx context.Context, roleName string) (*domain.Role, error)
}
