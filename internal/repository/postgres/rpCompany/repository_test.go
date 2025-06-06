package rpCompany

import (
	"context"
	"errors"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go-storage/internal/domain"
	"testing"
	"time"
)

func TestRepositoryCompany_Create(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := NewRepository(db)

	c := &domain.Company{
		ID:          "123",
		Name:        "Test Company",
		Path:        "test-company",
		Description: "Just a test",
	}

	mock.ExpectExec("INSERT INTO companies").
		WithArgs(c.ID, c.Name, c.Path, c.Description, sqlmock.AnyArg(), sqlmock.AnyArg(), true).
		WillReturnResult(sqlmock.NewResult(1, 1))

	result, err := repo.Create(context.Background(), c)
	assert.NoError(t, err)
	assert.Equal(t, c.ID, result.ID)
	assert.Equal(t, c.Name, result.Name)
	assert.Equal(t, true, result.IsActive)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestRepositoryCompany_GetCompanyById(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := NewRepository(db)

	createdAt := time.Now()
	updatedAt := createdAt.Add(1 * time.Hour)

	expected := &domain.Company{
		ID:          "123",
		Name:        "Test Company",
		Path:        "test-company",
		Description: "Description",
		CreatedAt:   createdAt,
		UpdatedAt:   updatedAt,
		IsActive:    true,
	}

	rows := sqlmock.NewRows([]string{
		"id", "name", "path", "description", "created_at", "updated_at", "is_active",
	}).AddRow(
		expected.ID,
		expected.Name,
		expected.Path,
		expected.Description,
		expected.CreatedAt,
		expected.UpdatedAt,
		expected.IsActive,
	)

	mock.ExpectQuery("SELECT (.+) FROM companies WHERE id = ?").
		WithArgs(expected.ID).
		WillReturnRows(rows)

	result, err := repo.GetCompanyById(context.Background(), expected.ID)
	assert.NoError(t, err)
	assert.Equal(t, expected.ID, result.ID)
	assert.Equal(t, expected.Name, result.Name)
	assert.Equal(t, expected.Path, result.Path)
	assert.Equal(t, expected.Description, result.Description)
	assert.WithinDuration(t, expected.CreatedAt, result.CreatedAt, time.Second)
	assert.WithinDuration(t, expected.UpdatedAt, result.UpdatedAt, time.Second)
	assert.Equal(t, expected.IsActive, result.IsActive)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestRepositoryCompany_GetAllCompanies(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := NewRepository(db)

	now := time.Now()

	expectedCompanies := []*domain.Company{
		{
			ID:          "1",
			Name:        "Company One",
			Path:        "company-one",
			Description: "Desc One",
			CreatedAt:   now,
			UpdatedAt:   now,
			IsActive:    true,
		},
		{
			ID:          "2",
			Name:        "Company Two",
			Path:        "company-two",
			Description: "Desc Two",
			CreatedAt:   now,
			UpdatedAt:   now,
			IsActive:    false,
		},
	}

	rows := sqlmock.NewRows([]string{
		"id", "name", "path", "description", "created_at", "updated_at", "is_active",
	}).AddRow(
		expectedCompanies[0].ID,
		expectedCompanies[0].Name,
		expectedCompanies[0].Path,
		expectedCompanies[0].Description,
		expectedCompanies[0].CreatedAt,
		expectedCompanies[0].UpdatedAt,
		expectedCompanies[0].IsActive,
	).AddRow(
		expectedCompanies[1].ID,
		expectedCompanies[1].Name,
		expectedCompanies[1].Path,
		expectedCompanies[1].Description,
		expectedCompanies[1].CreatedAt,
		expectedCompanies[1].UpdatedAt,
		expectedCompanies[1].IsActive,
	)

	mock.ExpectQuery(`(?i)FROM companies WHERE is_active = true`).
		WillReturnRows(rows)

	result, err := repo.GetAllCompanies(context.Background())
	assert.NoError(t, err)
	assert.Len(t, result, len(expectedCompanies))

	for i := range result {
		assert.Equal(t, expectedCompanies[i].ID, result[i].ID)
		assert.Equal(t, expectedCompanies[i].Name, result[i].Name)
		assert.Equal(t, expectedCompanies[i].Path, result[i].Path)
		assert.Equal(t, expectedCompanies[i].Description, result[i].Description)
		assert.WithinDuration(t, expectedCompanies[i].CreatedAt, result[i].CreatedAt, time.Second)
		assert.WithinDuration(t, expectedCompanies[i].UpdatedAt, result[i].UpdatedAt, time.Second)
		assert.Equal(t, expectedCompanies[i].IsActive, result[i].IsActive)
	}

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestRepositoryCompany_DeleteCompany(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := NewRepository(db)
	companyID := "123"

	mock.ExpectExec(`DELETE FROM companies`).
		WithArgs(companyID).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err = repo.DeleteCompany(context.Background(), companyID)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestRepositoryCompany_UpdateIsActive(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := NewRepository(db)
	companyID := "123"
	isActive := true

	t.Run("success", func(t *testing.T) {
		mock.ExpectExec(`(?i)UPDATE companies SET is_active = \$1 WHERE id = \$2`).
			WithArgs(isActive, companyID).
			WillReturnResult(sqlmock.NewResult(1, 1))

		err := repo.UpdateIsActive(context.Background(), companyID, isActive)
		assert.NoError(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("database error", func(t *testing.T) {
		mock.ExpectExec(`(?i)UPDATE companies SET is_active = \$1 WHERE id = \$2`).
			WithArgs(isActive, companyID).
			WillReturnError(errors.New("db failure"))

		err := repo.UpdateIsActive(context.Background(), companyID, isActive)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "unable to delete company")
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestRepositoryCompany_Update(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := NewRepository(db)

	createdAt := time.Now().Add(-24 * time.Hour)
	c := &domain.Company{
		ID:          "123",
		Name:        "New Name",
		Path:        "new-name",
		Description: "Updated description",
		CreatedAt:   createdAt,
		IsActive:    true,
	}

	queryRegex := `(?i)UPDATE companies SET id = \$1, name = \$2, storage_path = \$3, description = \$4, created_at = \$5, update_at = \$6, is_active = \$7 WHERE id = \$1 AND is_active = true`

	t.Run("success", func(t *testing.T) {
		mock.ExpectExec(queryRegex).
			WithArgs(
				c.ID,
				c.Name,
				c.Path,
				c.Description,
				c.CreatedAt,
				sqlmock.AnyArg(), // updateDate
				c.IsActive,
			).
			WillReturnResult(sqlmock.NewResult(0, 1))

		result, err := repo.Update(context.Background(), c)
		require.NoError(t, err)
		require.NotNil(t, result)

		assert.Equal(t, c.ID, result.ID)
		assert.Equal(t, c.Name, result.Name)
		assert.Equal(t, c.Path, result.Path)
		assert.Equal(t, c.Description, result.Description)
		assert.Equal(t, c.CreatedAt, result.CreatedAt)
		assert.Equal(t, c.IsActive, result.IsActive)
		assert.False(t, result.UpdatedAt.IsZero())

		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("database error", func(t *testing.T) {
		mock.ExpectExec(queryRegex).
			WithArgs(
				c.ID,
				c.Name,
				c.Path,
				c.Description,
				c.CreatedAt,
				sqlmock.AnyArg(), // updateDate
				c.IsActive,
			).
			WillReturnError(errors.New("update failed"))

		_, err := repo.Update(context.Background(), c)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "unable to insert company")
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}
