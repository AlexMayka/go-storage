package ucUser

import (
	"context"
	"go-storage/internal/domain"
)

type RepositoryUserInterface interface {
	CreateUser(ctx context.Context, u *domain.User) (*domain.User, error)
	GetUserByID(ctx context.Context, id string) (*domain.User, error)
	GetUserByEmail(ctx context.Context, email string) (*domain.User, error)
	GetUserByUsername(ctx context.Context, username string) (*domain.User, error)
	GetUsersByCompanyID(ctx context.Context, companyID string) ([]*domain.User, error)
	UpdateUser(ctx context.Context, u *domain.User) (*domain.User, error)
	UpdatePassword(ctx context.Context, userID, hashedPassword string) error
	UpdateIsActive(ctx context.Context, userID string, isActive bool) error
	UpdateLastLogin(ctx context.Context, userID string) error
}

type UseCaseUserInterface interface {
	RegisterUser(ctx context.Context, u *domain.User) (*domain.User, error)
	Login(ctx context.Context, login, password string) (*domain.User, error)
	GetUserByID(ctx context.Context, id string) (*domain.User, error)
	GetUsersByCompany(ctx context.Context, companyID string) ([]*domain.User, error)
	UpdateUser(ctx context.Context, userID string, u *domain.User) (*domain.User, error)
	ChangePassword(ctx context.Context, userID, oldPassword, newPassword string) error
	DeactivateUser(ctx context.Context, userID string) error
	RefreshToken(ctx context.Context, userID string) (*domain.User, error)
}
