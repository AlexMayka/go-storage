package rpUser

import (
	"context"
	"database/sql"
	"errors"
	"go-storage/internal/domain"
	pkgErrors "go-storage/pkg/errors"
	"time"
)

type RepositoryUser struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *RepositoryUser {
	return &RepositoryUser{db: db}
}

func (r *RepositoryUser) CreateUser(ctx context.Context, u *domain.User) (*domain.User, error) {
	const isActive = true
	createDate := time.Now()

	_, err := r.db.ExecContext(ctx, QueryCreateUser,
		u.ID, u.FirstName, u.SecondName, u.LastName, u.Username, u.Email, u.Phone, u.Password, u.CompanyId, u.RoleId,
		createDate, createDate, createDate, isActive)

	if err != nil {
		return nil, pkgErrors.Database("unable to create user")
	}

	return &domain.User{
		ID:         u.ID,
		FirstName:  u.FirstName,
		SecondName: u.SecondName,
		LastName:   u.LastName,
		Username:   u.Username,
		Email:      u.Email,
		Phone:      u.Phone,
		Password:   u.Password,
		CompanyId:  u.CompanyId,
		RoleId:     u.RoleId,
		LastLogin:  createDate,
		CreatedAt:  createDate,
		UpdatedAt:  createDate,
		IsActive:   isActive,
	}, nil
}

func (r *RepositoryUser) GetUserByID(ctx context.Context, id string) (*domain.User, error) {
	var user domain.User
	row := r.db.QueryRowContext(ctx, QueryGetUserByID, id)

	if err := row.Scan(&user.ID, &user.FirstName, &user.SecondName, &user.LastName, &user.Username, &user.Email, &user.Phone, &user.Password, &user.CompanyId, &user.RoleId, &user.LastLogin, &user.CreatedAt, &user.UpdatedAt, &user.IsActive); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, pkgErrors.NotFound("user not found")
		}
		return nil, pkgErrors.Database("unable to get user")
	}

	return &user, nil
}

func (r *RepositoryUser) GetUserByIDWithCompany(ctx context.Context, id string, companyId string) (*domain.User, error) {
	var user domain.User
	row := r.db.QueryRowContext(ctx, QueryGetUserByIDWithCompany, id, companyId)

	if err := row.Scan(&user.ID, &user.FirstName, &user.SecondName, &user.LastName, &user.Username, &user.Email, &user.Phone, &user.Password, &user.CompanyId, &user.RoleId, &user.LastLogin, &user.CreatedAt, &user.UpdatedAt, &user.IsActive); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, pkgErrors.NotFound("user not found")
		}
		return nil, pkgErrors.Database("unable to get user")
	}

	return &user, nil
}

func (r *RepositoryUser) GetUserByEmail(ctx context.Context, email string) (*domain.User, error) {
	var user domain.User
	row := r.db.QueryRowContext(ctx, QueryGetUserByEmail, email)

	if err := row.Scan(&user.ID, &user.FirstName, &user.SecondName, &user.LastName, &user.Username, &user.Email, &user.Phone, &user.Password, &user.CompanyId, &user.RoleId, &user.LastLogin, &user.CreatedAt, &user.UpdatedAt, &user.IsActive); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, pkgErrors.NotFound("user not found")
		}

		return nil, pkgErrors.Database("unable to get user")
	}

	return &user, nil
}

func (r *RepositoryUser) GetUserByUsername(ctx context.Context, username string) (*domain.User, error) {
	var user domain.User
	row := r.db.QueryRowContext(ctx, QueryGetUserByUsername, username)

	if err := row.Scan(&user.ID, &user.FirstName, &user.SecondName, &user.LastName, &user.Username, &user.Email, &user.Phone, &user.Password, &user.CompanyId, &user.RoleId, &user.LastLogin, &user.CreatedAt, &user.UpdatedAt, &user.IsActive); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, pkgErrors.NotFound("user not found")
		}
		return nil, pkgErrors.Database("unable to get user")
	}

	return &user, nil
}

func (r *RepositoryUser) GetUsersByCompanyID(ctx context.Context, companyID string) ([]*domain.User, error) {
	var users []*domain.User
	rows, err := r.db.QueryContext(ctx, QueryGetUsersByCompanyID, companyID)

	if err != nil {
		return nil, pkgErrors.Database("unable to query users by company")
	}
	defer func(rows *sql.Rows) {
		_ = rows.Close()
	}(rows)

	for rows.Next() {
		var user domain.User
		if err := rows.Scan(&user.ID, &user.FirstName, &user.SecondName, &user.LastName, &user.Username, &user.Email, &user.Phone, &user.Password, &user.CompanyId, &user.RoleId, &user.LastLogin, &user.CreatedAt, &user.UpdatedAt, &user.IsActive); err != nil {
			return nil, pkgErrors.Database("unable to query users by company")
		}
		users = append(users, &user)
	}

	if err := rows.Err(); err != nil {
		return nil, pkgErrors.Database("unable to query users by company")
	}

	return users, nil
}

func (r *RepositoryUser) UpdateUser(ctx context.Context, u *domain.User) (*domain.User, error) {
	updateDate := time.Now()
	_, err := r.db.ExecContext(ctx, QueryUpdateUser, u.ID, u.FirstName, u.SecondName, u.LastName, u.Username, u.Email, u.Phone, updateDate)
	if err != nil {
		return nil, pkgErrors.Database("unable to update user")
	}

	u.UpdatedAt = updateDate
	return u, nil
}

func (r *RepositoryUser) UpdateUserWithCompany(ctx context.Context, u *domain.User, companyId string) (*domain.User, error) {
	updateDate := time.Now()
	_, err := r.db.ExecContext(ctx, QueryUpdateUserWithCompany, u.ID, u.FirstName, u.SecondName, u.LastName, u.Username, u.Email, u.Phone, updateDate, companyId)
	if err != nil {
		return nil, pkgErrors.Database("unable to update user")
	}

	u.UpdatedAt = updateDate
	return u, nil
}

func (r *RepositoryUser) UpdatePassword(ctx context.Context, userID, hashedPassword string) error {
	updateDate := time.Now()
	_, err := r.db.ExecContext(ctx, QueryUpdatePassword, userID, hashedPassword, updateDate)
	if err != nil {
		return pkgErrors.Database("unable to update password")
	}
	return nil
}

func (r *RepositoryUser) UpdatePasswordWithCompany(ctx context.Context, userID, hashedPassword string, companyId string) error {
	updateDate := time.Now()
	_, err := r.db.ExecContext(ctx, QueryUpdatePasswordWithCompany, userID, hashedPassword, updateDate, companyId)
	if err != nil {
		return pkgErrors.Database("unable to update password")
	}
	return nil
}

func (r *RepositoryUser) UpdateIsActive(ctx context.Context, userID string, isActive bool) error {
	updateDate := time.Now()
	_, err := r.db.ExecContext(ctx, QueryUpdateIsActive, userID, isActive, updateDate)
	if err != nil {
		return pkgErrors.Database("unable to update user status")
	}
	return nil
}

func (r *RepositoryUser) UpdateIsActiveWithCompany(ctx context.Context, userID string, isActive bool, companyId string) error {
	updateDate := time.Now()
	_, err := r.db.ExecContext(ctx, QueryUpdateIsActiveWithCompany, userID, isActive, updateDate, companyId)
	if err != nil {
		return pkgErrors.Database("unable to update user status")
	}
	return nil
}

func (r *RepositoryUser) UpdateLastLogin(ctx context.Context, userID string) error {
	updateDate := time.Now()
	_, err := r.db.ExecContext(ctx, QueryUpdateLastLogin, userID, updateDate, updateDate)
	if err != nil {
		return pkgErrors.Database("unable to update last login")
	}
	return nil
}
