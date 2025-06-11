package ucUser

import (
	"context"
	"go-storage/internal/domain"
)

type RepositoryUserInterface interface {
	CreateUser(ctx context.Context, u *domain.User) (*domain.User, error)
}
