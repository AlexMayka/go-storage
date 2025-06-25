package ucAuthUser

import (
	"context"
	"go-storage/internal/domain"
)

type UseCaseAuth struct {
	repo RepositoryAuthInterface
}

func NewUseCaseAuth(repo RepositoryAuthInterface) *UseCaseAuth {
	return &UseCaseAuth{
		repo: repo,
	}
}

func (u *UseCaseAuth) GetRoleById(ctx context.Context, roleId string) (*domain.Role, error) {
	return u.repo.GetRoleById(ctx, roleId)
}

func (u *UseCaseAuth) GetRoleByName(ctx context.Context, roleName string) (*domain.Role, error) {
	return u.repo.GetRoleByName(ctx, roleName)
}

func (u *UseCaseAuth) GetPermissionById(ctx context.Context, permissionId string) (*domain.Permission, error) {
	return u.repo.GetPermissionById(ctx, permissionId)
}

func (u *UseCaseAuth) GetPermissionByRoleName(ctx context.Context, roleName string) (*[]domain.Permission, error) {
	role, err := u.GetRoleByName(ctx, roleName)
	if err != nil {
		return nil, err
	}

	rolePermissions, err := u.repo.GetPermissionsIdByRoleId(ctx, role.ID)
	if err != nil {
		return nil, err
	}

	collectOfPermissions := make([]string, 0, len(*rolePermissions))
	for _, permission := range *rolePermissions {
		collectOfPermissions = append(collectOfPermissions, permission.PermissionsID)
	}

	permissions, err := u.repo.GetPermissionByIds(ctx, collectOfPermissions)
	if err != nil {
		return nil, err
	}

	return permissions, nil
}

func (u *UseCaseAuth) GetRolePermissionsByRoleId(ctx context.Context, roleId string) (*[]domain.Permission, error) {
	rolePermissions, err := u.repo.GetPermissionsIdByRoleId(ctx, roleId)
	if err != nil {
		return nil, err
	}

	permissionsIds := make([]string, len(*rolePermissions))
	for index, permission := range *rolePermissions {
		permissionsIds[index] = permission.PermissionsID
	}

	permissions, err := u.repo.GetPermissionByIds(ctx, permissionsIds)
	if err != nil {
		return nil, err
	}

	return permissions, nil
}

func (u *UseCaseAuth) GetDefaultRole(ctx context.Context) (domain.Role, error) {
	role, err := u.GetRoleByName(ctx, "user")
	if err != nil {
		return domain.Role{}, err
	}
	return *role, nil
}
