package ucUser

import (
	"context"
	"github.com/google/uuid"
	"go-storage/internal/domain"
	"go-storage/internal/utils/valid"
	"go-storage/pkg/auth"
	"go-storage/pkg/errors"
)

type UseCaseUser struct {
	repo     RepositoryUserInterface
	authRepo RepositoryAuthInterface
}

func NewUseCaseUser(repo RepositoryUserInterface, authRepo RepositoryAuthInterface) *UseCaseUser {
	return &UseCaseUser{
		repo:     repo,
		authRepo: authRepo,
	}
}

func (u *UseCaseUser) RegisterUser(ctx context.Context, c *domain.User) (*domain.User, error) {
	var errPs error
	if c.Password, errPs = auth.Hash(c.Password); errPs != nil {
		return nil, errors.InternalServer("error password")
	}

	c.ID = uuid.NewString()
	c.IsActive = true

	return u.repo.CreateUser(ctx, c)
}

func (u *UseCaseUser) Login(ctx context.Context, login, password string) (*domain.User, error) {
	var user *domain.User
	var err error

	if valid.CheckEmail(login) {
		user, err = u.repo.GetUserByEmail(ctx, login)
	} else {
		user, err = u.repo.GetUserByUsername(ctx, login)
	}

	if err != nil {
		return nil, err
	}

	if !user.IsActive {
		return nil, errors.Forbidden("user is deactivated")
	}

	if !auth.CheckPasswordHash(password, user.Password) {
		return nil, errors.Unauthorized("invalid credentials")
	}

	_ = u.repo.UpdateLastLogin(ctx, user.ID)

	return user, nil
}

func (u *UseCaseUser) GetUserByID(ctx context.Context, id string) (*domain.User, error) {
	return u.repo.GetUserByID(ctx, id)
}

func (u *UseCaseUser) GetUsersByCompany(ctx context.Context, companyID string) ([]*domain.User, error) {
	return u.repo.GetUsersByCompanyID(ctx, companyID)
}

func (u *UseCaseUser) UpdateUser(ctx context.Context, userID string, user *domain.User) (*domain.User, error) {
	existingUser, err := u.repo.GetUserByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	if user.FirstName != "" {
		existingUser.FirstName = user.FirstName
	}
	if user.SecondName != "" {
		existingUser.SecondName = user.SecondName
	}
	if user.LastName != "" {
		existingUser.LastName = user.LastName
	}
	if user.Email != "" && user.Email != existingUser.Email {
		existingUser.Email = user.Email
	}
	if user.Phone != "" && user.Phone != existingUser.Phone {
		existingUser.Phone = user.Phone
	}
	if user.Username != "" && user.Username != existingUser.Username {
		existingUser.Username = user.Username
	}

	return u.repo.UpdateUser(ctx, existingUser)
}

func (u *UseCaseUser) ChangePassword(ctx context.Context, userID, oldPassword, newPassword string) error {
	user, err := u.repo.GetUserByID(ctx, userID)
	if err != nil {
		return err
	}

	if !auth.CheckPasswordHash(oldPassword, user.Password) {
		return errors.Unauthorized("invalid old password")
	}

	hashedPassword, err := auth.Hash(newPassword)
	if err != nil {
		return errors.InternalServer("failed to hash password")
	}

	return u.repo.UpdatePassword(ctx, userID, hashedPassword)
}

func (u *UseCaseUser) DeactivateUser(ctx context.Context, userID string) error {
	_, err := u.repo.GetUserByID(ctx, userID)
	if err != nil {
		return err
	}

	return u.repo.UpdateIsActive(ctx, userID, false)
}

func (u *UseCaseUser) RefreshToken(ctx context.Context, userID string) (*domain.User, error) {
	user, err := u.repo.GetUserByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	if !user.IsActive {
		return nil, errors.Forbidden("user is deactivated")
	}

	return user, nil
}

func (u *UseCaseUser) UpdateUserRole(ctx context.Context, userID string, roleID string) error {
	_, err := u.repo.GetUserByID(ctx, userID)
	if err != nil {
		return err
	}

	_, err = u.authRepo.GetRoleById(ctx, roleID)
	if err != nil {
		return errors.BadRequest("invalid role ID")
	}

	return u.repo.UpdateUserRole(ctx, userID, roleID)
}

func (u *UseCaseUser) ActivateUser(ctx context.Context, userID string) error {
	_, err := u.repo.GetUserByID(ctx, userID)
	if err != nil {
		return err
	}

	return u.repo.UpdateIsActive(ctx, userID, true)
}

func (u *UseCaseUser) GetAllUsers(ctx context.Context) ([]*domain.User, error) {
	return u.repo.GetAllUsers(ctx)
}

func (u *UseCaseUser) TransferUserToCompany(ctx context.Context, userID string, companyID string) error {
	_, err := u.repo.GetUserByID(ctx, userID)
	if err != nil {
		return err
	}

	return u.repo.UpdateUserCompany(ctx, userID, companyID)
}
