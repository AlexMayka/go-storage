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
	_, err := r.db.ExecContext(ctx, QueryCreateUser,
		u.ID, u.FirstName, u.SecondName, u.LastName, u.Username,
		u.Email, u.Phone, u.Password,
		u.CompanyId, u.RoleId,
		u.LastLogin, u.CreatedAt, u.UpdatedAt, u.IsActive,
	)
	if err != nil {
		return nil, pkgErrors.Database("unable to create user")
	}
	return u, nil
}

func (r *RepositoryUser) GetUserByID(ctx context.Context, id string) (*domain.User, error) {
	var user domain.User
	row := r.db.QueryRowContext(ctx, QueryGetUserByID, id)

	if err := scanUser(row, &user); err != nil {
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

	if err := scanUser(row, &user); err != nil {
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

	if err := scanUser(row, &user); err != nil {
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

	if err := scanUser(row, &user); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, pkgErrors.NotFound("user not found")
		}
		return nil, pkgErrors.Database("unable to get user")
	}

	return &user, nil
}

func (r *RepositoryUser) GetUsersByCompanyID(ctx context.Context, companyID string) ([]*domain.User, error) {
	rows, err := r.db.QueryContext(ctx, QueryGetUsersByCompanyID, companyID)
	if err != nil {
		return nil, pkgErrors.Database("unable to get users")
	}
	defer rows.Close()

	var users []*domain.User
	for rows.Next() {
		var user domain.User
		if err := scanUserFromRows(rows, &user); err != nil {
			return nil, pkgErrors.Database("unable to scan user")
		}
		users = append(users, &user)
	}

	if err := rows.Err(); err != nil {
		return nil, pkgErrors.Database("unable to get users")
	}

	return users, nil
}

func (r *RepositoryUser) UpdateUser(ctx context.Context, u *domain.User) (*domain.User, error) {
	u.UpdatedAt = time.Now()
	_, err := r.db.ExecContext(ctx, QueryUpdateUser,
		u.ID, u.FirstName, u.SecondName, u.LastName, u.Username,
		u.Email, u.Phone, u.UpdatedAt,
	)
	if err != nil {
		return nil, pkgErrors.Database("unable to update user")
	}
	return u, nil
}

func (r *RepositoryUser) UpdateUserWithCompany(ctx context.Context, u *domain.User, companyId string) (*domain.User, error) {
	u.UpdatedAt = time.Now()
	_, err := r.db.ExecContext(ctx, QueryUpdateUserWithCompany,
		u.ID, u.FirstName, u.SecondName, u.LastName, u.Username,
		u.Email, u.Phone, u.UpdatedAt, companyId,
	)
	if err != nil {
		return nil, pkgErrors.Database("unable to update user")
	}
	return u, nil
}

func (r *RepositoryUser) UpdatePassword(ctx context.Context, userID, hashedPassword string) error {
	_, err := r.db.ExecContext(ctx, QueryUpdatePassword, userID, hashedPassword, time.Now())
	if err != nil {
		return pkgErrors.Database("unable to update password")
	}
	return nil
}

func (r *RepositoryUser) UpdatePasswordWithCompany(ctx context.Context, userID, hashedPassword string, companyId string) error {
	_, err := r.db.ExecContext(ctx, QueryUpdatePasswordWithCompany, userID, hashedPassword, time.Now(), companyId)
	if err != nil {
		return pkgErrors.Database("unable to update password")
	}
	return nil
}

func (r *RepositoryUser) UpdateIsActive(ctx context.Context, userID string, isActive bool) error {
	_, err := r.db.ExecContext(ctx, QueryUpdateIsActive, userID, isActive, time.Now())
	if err != nil {
		return pkgErrors.Database("unable to update user status")
	}
	return nil
}

func (r *RepositoryUser) UpdateIsActiveWithCompany(ctx context.Context, userID string, isActive bool, companyId string) error {
	_, err := r.db.ExecContext(ctx, QueryUpdateIsActiveWithCompany, userID, isActive, time.Now(), companyId)
	if err != nil {
		return pkgErrors.Database("unable to update user status")
	}
	return nil
}

func (r *RepositoryUser) UpdateLastLogin(ctx context.Context, userID string) error {
	_, err := r.db.ExecContext(ctx, QueryUpdateLastLogin, userID, time.Now(), time.Now())
	if err != nil {
		return pkgErrors.Database("unable to update last login")
	}
	return nil
}

func (r *RepositoryUser) UpdateUserRole(ctx context.Context, userID string, roleID string) error {
	_, err := r.db.ExecContext(ctx, QueryUpdateUserRole, userID, roleID, time.Now())
	if err != nil {
		return pkgErrors.Database("unable to update user role")
	}
	return nil
}

func (r *RepositoryUser) GetAllUsers(ctx context.Context) ([]*domain.User, error) {
	rows, err := r.db.QueryContext(ctx, QueryGetAllUsers)
	if err != nil {
		return nil, pkgErrors.Database("unable to get all users")
	}
	defer rows.Close()

	var users []*domain.User
	for rows.Next() {
		var user domain.User
		if err := scanUserFromRows(rows, &user); err != nil {
			return nil, pkgErrors.Database("unable to scan user")
		}
		users = append(users, &user)
	}

	if err := rows.Err(); err != nil {
		return nil, pkgErrors.Database("unable to get all users")
	}

	return users, nil
}

func (r *RepositoryUser) UpdateUserCompany(ctx context.Context, userID string, companyID string) error {
	_, err := r.db.ExecContext(ctx, QueryUpdateUserCompany, userID, companyID, time.Now())
	if err != nil {
		return pkgErrors.Database("unable to update user company")
	}
	return nil
}

func scanUser(row *sql.Row, user *domain.User) error {
	return row.Scan(
		&user.ID, &user.FirstName, &user.SecondName, &user.LastName, &user.Username,
		&user.Email, &user.Phone, &user.Password,
		&user.CompanyId, &user.RoleId,
		&user.LastLogin, &user.CreatedAt, &user.UpdatedAt, &user.IsActive,
	)
}

func scanUserFromRows(rows *sql.Rows, user *domain.User) error {
	return rows.Scan(
		&user.ID, &user.FirstName, &user.SecondName, &user.LastName, &user.Username,
		&user.Email, &user.Phone, &user.Password,
		&user.CompanyId, &user.RoleId,
		&user.LastLogin, &user.CreatedAt, &user.UpdatedAt, &user.IsActive,
	)
}
