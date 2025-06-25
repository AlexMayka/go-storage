package ucAuthUser

import (
	"context"
	"go-storage/internal/domain"
)

type RepositoryAuthInterface interface {
	GetRoleById(ctx context.Context, roleId string) (*domain.Role, error)
	GetRoleByName(ctx context.Context, roleId string) (*domain.Role, error)
	GetPermissionById(ctx context.Context, permissionId string) (*domain.Permission, error)
	GetPermissionByIds(ctx context.Context, permissionsId []string) (*[]domain.Permission, error)
	GetPermissionsIdByRoleId(ctx context.Context, roleId string) (*[]domain.RolePermissions, error)
}
