package rpAuth

import (
	"context"
	"database/sql"
	"errors"
	"github.com/lib/pq"
	"go-storage/internal/domain"
	pkgErrors "go-storage/pkg/errors"
)

type RepositoryAuth struct {
	db *sql.DB
}

func NewRepositoryAuth(db *sql.DB) *RepositoryAuth { return &RepositoryAuth{db: db} }

func (r *RepositoryAuth) GetRoleById(ctx context.Context, roleId string) (*domain.Role, error) {
	var role domain.Role
	row := r.db.QueryRowContext(ctx, QueryGetRoleById, roleId)

	if err := row.Scan(&role.ID, &role.Name); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, pkgErrors.NotFound("role not found")
		}
		return nil, pkgErrors.Database("unable get to role")
	}

	return &role, nil
}

func (r *RepositoryAuth) GetRoleByName(ctx context.Context, roleId string) (*domain.Role, error) {
	var role domain.Role
	row := r.db.QueryRowContext(ctx, QueryGetRoleByName, roleId)

	if err := row.Scan(&role.ID, &role.Name); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, pkgErrors.NotFound("role not found")
		}
		return nil, pkgErrors.Database("unable get to role")
	}

	return &role, nil
}

func (r *RepositoryAuth) GetPermissionById(ctx context.Context, permissionId string) (*domain.Permission, error) {
	var permission domain.Permission
	row := r.db.QueryRowContext(ctx, QueryGetPermissionById, permissionId)

	if err := row.Scan(&permission.ID, &permission.Name); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, pkgErrors.NotFound("permission not found")
		}
		return nil, pkgErrors.Database("unable get to permission")
	}

	return &permission, nil
}

func (r *RepositoryAuth) GetPermissionByIds(ctx context.Context, permissionsId []string) (*[]domain.Permission, error) {
	if len(permissionsId) == 0 {
		return &[]domain.Permission{}, nil
	}

	var permissions []domain.Permission
	rows, err := r.db.QueryContext(ctx, GetPermissionByIds, pq.Array(permissionsId))

	if err != nil {
		return nil, pkgErrors.Database("unable get to permissions")
	}
	defer func(rows *sql.Rows) {
		_ = rows.Close()
	}(rows)

	for rows.Next() {
		var permission domain.Permission
		if err := rows.Scan(&permission.ID, &permission.Name); err != nil {
			return nil, pkgErrors.Database("unable get to permissions")
		}
		permissions = append(permissions, permission)
	}

	if err := rows.Err(); err != nil {
		return nil, pkgErrors.Database("unable get to permissions")
	}

	return &permissions, nil
}

func (r *RepositoryAuth) GetPermissionsIdByRoleId(ctx context.Context, roleId string) (*[]domain.RolePermissions, error) {
	var permissions []domain.RolePermissions
	rows, err := r.db.QueryContext(ctx, QueryGetPermissionsIdByRoleId, roleId)

	if err != nil {
		return nil, pkgErrors.Database("unable get to permissions")
	}
	defer func(rows *sql.Rows) {
		_ = rows.Close()
	}(rows)

	for rows.Next() {
		var permission domain.RolePermissions
		if err := rows.Scan(&permission.RoleID, &permission.PermissionsID); err != nil {
			return nil, pkgErrors.Database("unable get to permissions")
		}
		permissions = append(permissions, permission)
	}

	if err := rows.Err(); err != nil {
		return nil, pkgErrors.Database("unable get to permissions")
	}

	if len(permissions) == 0 {
		return nil, pkgErrors.NotFound("permissions for role not found")
	}

	return &permissions, nil
}
