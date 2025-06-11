package rpUser

import (
	"context"
	"database/sql"
	"go-storage/internal/domain"
	pkgErrors "go-storage/pkg/errors"
	"time"
)

type RepositoryUser struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) RepositoryUser {
	return RepositoryUser{db: db}
}

func (r *RepositoryUser) CreateUser(ctx context.Context, u *domain.User) (*domain.User, error) {
	const isActive = true
	createDate := time.Now()

	_, err := r.db.QueryContext(ctx, QueryCreateUser,
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
