package ucUser

import (
	"context"
	"github.com/google/uuid"
	"go-storage/internal/domain"
	"go-storage/pkg/auth"
	"go-storage/pkg/errors"
)

type UseCaseUser struct {
	repo RepositoryUserInterface
}

func NewUseCaseUser(repo RepositoryUserInterface) *UseCaseUser {
	return &UseCaseUser{
		repo: repo,
	}
}

func (u *UseCaseUser) RegisterUser(ctx context.Context, c *domain.User) (*domain.User, error) {
	if c.FirstName == "" {
		return nil, errors.EmptyField("first_name")
	}

	if c.LastName == "" {
		return nil, errors.EmptyField("last_name")
	}

	if c.Username == "" {
		return nil, errors.EmptyField("username")
	}

	if c.Email == "" {
		return nil, errors.EmptyField("email")
	}

	if c.Password == "" {
		return nil, errors.EmptyField("password")
	}

	if c.CompanyId == "" {
		return nil, errors.EmptyField("company_id")
	}

	if c.RoleId == "" {
		return nil, errors.EmptyField("role_id")
	}

	var errPs error
	if c.Password, errPs = auth.Hash(c.Password); errPs != nil {
		return nil, errors.InternalServer("error password")
	}

	c.ID = uuid.NewString()
	c.IsActive = true

	user, errDb := u.repo.CreateUser(ctx, c)
	if errDb != nil {
		return nil, errDb
	}

	return user, nil
}
