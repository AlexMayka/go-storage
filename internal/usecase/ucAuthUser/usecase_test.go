package ucAuthUser

import (
	"context"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go-storage/internal/domain"
	"testing"
)

type MockRepository struct {
	mock.Mock
}

func (m *MockRepository) GetRoleById(ctx context.Context, roleId string) (*domain.Role, error) {
	args := m.Called(ctx, roleId)
	return args.Get(0).(*domain.Role), args.Error(1)
}

func (m *MockRepository) GetRoleByName(ctx context.Context, roleName string) (*domain.Role, error) {
	args := m.Called(ctx, roleName)
	return args.Get(0).(*domain.Role), args.Error(1)
}

func (m *MockRepository) GetPermissionById(ctx context.Context, permissionId string) (*domain.Permission, error) {
	args := m.Called(ctx, permissionId)
	return args.Get(0).(*domain.Permission), args.Error(1)
}

func (m *MockRepository) GetPermissionByIds(ctx context.Context, ids []string) (*[]domain.Permission, error) {
	args := m.Called(ctx, ids)
	return args.Get(0).(*[]domain.Permission), args.Error(1)
}

func (m *MockRepository) GetPermissionsIdByRoleId(ctx context.Context, roleId string) (*[]domain.RolePermissions, error) {
	args := m.Called(ctx, roleId)
	return args.Get(0).(*[]domain.RolePermissions), args.Error(1)
}

func TestGetPermissionByRoleName(t *testing.T) {
	ctx := context.Background()
	mockRepo := new(MockRepository)
	usecase := NewUseCaseAuth(mockRepo)

	roleName := "admin"
	role := &domain.Role{
		ID:   "role-1",
		Name: "admin",
	}

	rolePermissions := &[]domain.RolePermissions{
		{RoleID: "role-1", PermissionsID: "perm-1"},
		{RoleID: "role-1", PermissionsID: "perm-2"},
	}

	permissions := &[]domain.Permission{
		{ID: "perm-1", Name: "read"},
		{ID: "perm-2", Name: "write"},
	}

	mockRepo.On("GetRoleByName", ctx, roleName).Return(role, nil)
	mockRepo.On("GetPermissionsIdByRoleId", ctx, role.ID).Return(rolePermissions, nil)
	mockRepo.On("GetPermissionByIds", ctx, []string{"perm-1", "perm-2"}).Return(permissions, nil)

	result, err := usecase.GetPermissionByRoleName(ctx, roleName)

	assert.NoError(t, err)
	assert.Equal(t, permissions, result)

	mockRepo.AssertExpectations(t)
}
