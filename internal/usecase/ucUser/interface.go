package ucUser

import (
	"context"
	"go-storage/internal/domain"
)

type RepositoryUserInterface interface {
	CreateUser(ctx context.Context, u *domain.User) (*domain.User, error)
	GetUserByID(ctx context.Context, id string) (*domain.User, error)
	GetUserByIDWithCompany(ctx context.Context, id string, companyId string) (*domain.User, error)
	GetUserByEmail(ctx context.Context, email string) (*domain.User, error)
	GetUserByUsername(ctx context.Context, username string) (*domain.User, error)
	GetUsersByCompanyID(ctx context.Context, companyID string) ([]*domain.User, error)
	UpdateUser(ctx context.Context, u *domain.User) (*domain.User, error)
	UpdateUserWithCompany(ctx context.Context, u *domain.User, companyId string) (*domain.User, error)
	UpdatePassword(ctx context.Context, userID, hashedPassword string) error
	UpdatePasswordWithCompany(ctx context.Context, userID, hashedPassword string, companyId string) error
	UpdateIsActive(ctx context.Context, userID string, isActive bool) error
	UpdateIsActiveWithCompany(ctx context.Context, userID string, isActive bool, companyId string) error
	UpdateLastLogin(ctx context.Context, userID string) error
	UpdateUserRole(ctx context.Context, userID string, roleID string) error
	GetAllUsers(ctx context.Context) ([]*domain.User, error)
	UpdateUserCompany(ctx context.Context, userID string, companyID string) error
}

type RepositoryAuthInterface interface {
	GetRoleById(ctx context.Context, roleId string) (*domain.Role, error)
	GetRoleByName(ctx context.Context, roleId string) (*domain.Role, error)
}
