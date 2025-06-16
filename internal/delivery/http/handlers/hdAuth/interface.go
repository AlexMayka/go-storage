package hdAuth

import (
	"context"
	"github.com/gin-gonic/gin"
	"go-storage/internal/domain"
)

type AuthHandlerInterface interface {
	Login(ctx *gin.Context)
	RefreshToken(ctx *gin.Context)
}

type UseCaseUserInterface interface {
	Login(ctx context.Context, login, password string) (*domain.User, error)
	RefreshToken(ctx context.Context, userID string) (*domain.User, error)
}

type UseCaseAuthInterface interface {
	GetRoleByName(ctx context.Context, roleName string) (*domain.Role, error)
}
