package rpAuth

import (
	"context"
	"github.com/stretchr/testify/require"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"go-storage/internal/domain"
)

func TestGetRoleById(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()
	repo := NewRepositoryAuth(db)

	expected := &domain.Role{
		ID:   "123",
		Name: "admin",
	}

	rows := sqlmock.NewRows([]string{"id", "name"}).AddRow(expected.ID, expected.Name)
	mock.ExpectQuery(`SELECT id, "name" FROM roles WHERE id = \$1`).
		WithArgs(expected.ID).
		WillReturnRows(rows)

	result, err := repo.GetRoleById(context.Background(), expected.ID)
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, expected, result)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetRoleByName(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()
	repo := NewRepositoryAuth(db)

	expected := &domain.Role{
		ID:   "123",
		Name: "admin",
	}

	rows := sqlmock.NewRows([]string{"id", "name"}).AddRow(expected.ID, expected.Name)
	mock.ExpectQuery(`SELECT id, "name" FROM roles WHERE name = \$1`).
		WithArgs(expected.Name).
		WillReturnRows(rows)

	result, err := repo.GetRoleByName(context.Background(), expected.Name)
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, expected, result)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetPermissionById(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()
	repo := NewRepositoryAuth(db)

	expected := &domain.Permission{
		ID:   "perm-1",
		Name: "company:read",
	}

	rows := sqlmock.NewRows([]string{"id", "name"}).AddRow(expected.ID, expected.Name)
	mock.ExpectQuery(`SELECT id, "name" FROM permissions WHERE id = \$1`).
		WithArgs(expected.ID).
		WillReturnRows(rows)

	result, err := repo.GetPermissionById(context.Background(), expected.ID)
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, expected, result)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetPermissionsIdByRoleId(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := NewRepositoryAuth(db)
	roleId := "role-123"

	rows := sqlmock.NewRows([]string{"role_id", "permission_id"}).
		AddRow(roleId, "perm-1").
		AddRow(roleId, "perm-2")

	mock.ExpectQuery(`SELECT role_id, permission_id FROM role_permissions WHERE role_id = \$1`).
		WithArgs(roleId).
		WillReturnRows(rows)

	result, err := repo.GetPermissionsIdByRoleId(context.Background(), roleId)
	require.NoError(t, err)
	require.NotNil(t, result)
	require.Len(t, *result, 2)

	assert.Equal(t, roleId, (*result)[0].RoleID)
	assert.Equal(t, "perm-1", (*result)[0].PermissionsID)
	assert.Equal(t, "perm-2", (*result)[1].PermissionsID)

	assert.NoError(t, mock.ExpectationsWereMet())
}
