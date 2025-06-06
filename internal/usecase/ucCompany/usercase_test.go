package ucCompany

import (
	"context"
	"errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go-storage/internal/domain"
	customErrors "go-storage/pkg/errors"
	"testing"
)

type rpCompanyMock struct {
	mock.Mock
}

func (m *rpCompanyMock) Create(ctx context.Context, c *domain.Company) (*domain.Company, error) {
	args := m.Called(ctx, c)
	var company *domain.Company
	if args.Get(0) != nil {
		company = args.Get(0).(*domain.Company)
	}
	return company, args.Error(1)
}

func (m *rpCompanyMock) GetCompanyById(ctx context.Context, id string) (*domain.Company, error) {
	args := m.Called(ctx, id)
	var company *domain.Company
	if args.Get(0) != nil {
		company = args.Get(0).(*domain.Company)
	}
	return company, args.Error(1)
}

func (m *rpCompanyMock) Update(ctx context.Context, c *domain.Company) (*domain.Company, error) {
	args := m.Called(ctx, c)
	var company *domain.Company
	if args.Get(0) != nil {
		company = args.Get(0).(*domain.Company)
	}
	return company, args.Error(1)
}

func (m *rpCompanyMock) GetAllCompanies(ctx context.Context) ([]*domain.Company, error) {
	args := m.Called(ctx)
	var companies []*domain.Company
	if args.Get(0) != nil {
		companies = args.Get(0).([]*domain.Company)
	}
	return companies, args.Error(1)
}

func (m *rpCompanyMock) UpdateIsActive(ctx context.Context, id string, active bool) error {
	args := m.Called(ctx, id, active)
	return args.Error(0)
}

func (m *rpCompanyMock) DeleteCompany(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func TestUseCaseCompany_RegisterCompany(t *testing.T) {
	mockRepo := new(rpCompanyMock)
	uc := NewUseCase(mockRepo)

	t.Run("valid input", func(t *testing.T) {
		input := &domain.Company{Name: "TestCo", Description: "Desc"}
		expected := &domain.Company{
			ID:          "generated-id",
			Name:        "TestCo",
			Description: "Desc",
			Path:        "testco",
		}

		mockRepo.On("Create", mock.Anything, mock.MatchedBy(func(c *domain.Company) bool {
			return c.Name == input.Name && c.Description == input.Description && c.Path == "testco"
		})).Return(expected, nil)

		result, err := uc.RegisterCompany(context.Background(), input)
		assert.NoError(t, err)
		assert.Equal(t, expected, result)
		mockRepo.AssertExpectations(t)
	})

	t.Run("empty name", func(t *testing.T) {
		_, err := uc.RegisterCompany(context.Background(), &domain.Company{Name: "", Description: "Desc"})
		appErr, ok := err.(*customErrors.AppError)
		assert.True(t, ok)
		assert.Equal(t, 400, appErr.Code)
		assert.Equal(t, "Field 'name' is required", appErr.Message)
	})

	t.Run("empty description", func(t *testing.T) {
		_, err := uc.RegisterCompany(context.Background(), &domain.Company{Name: "Valid", Description: ""})
		appErr, ok := err.(*customErrors.AppError)
		assert.True(t, ok)
		assert.Equal(t, 400, appErr.Code)
		assert.Equal(t, "Field 'description' is required", appErr.Message)
	})
}

func TestUseCaseCompany_UpdateCompany(t *testing.T) {
	t.Run("empty update fields", func(t *testing.T) {
		mockRepo := new(rpCompanyMock)
		uc := NewUseCase(mockRepo)

		_, err := uc.UpdateCompany(context.Background(), "id123", &domain.Company{})
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "Field 'name and description' is required")
	})

	t.Run("repo.GetCompanyById error", func(t *testing.T) {
		mockRepo := new(rpCompanyMock)
		uc := NewUseCase(mockRepo)

		mockRepo.On("GetCompanyById", mock.Anything, "id123").Return(nil, errors.New("not found"))

		_, err := uc.UpdateCompany(context.Background(), "id123", &domain.Company{Name: "X"})
		assert.Error(t, err)
	})

	t.Run("no changes", func(t *testing.T) {
		mockRepo := new(rpCompanyMock)
		uc := NewUseCase(mockRepo)

		original := &domain.Company{ID: "id123", Name: "OldName", Description: "OldDesc"}
		mockRepo.On("GetCompanyById", mock.Anything, "id123").Return(original, nil)

		result, err := uc.UpdateCompany(context.Background(), "id123", &domain.Company{
			Name:        "OldName",
			Description: "OldDesc",
		})

		assert.NoError(t, err)
		assert.Equal(t, original, result)
	})

	t.Run("successful update", func(t *testing.T) {
		mockRepo := new(rpCompanyMock)
		uc := NewUseCase(mockRepo)

		original := &domain.Company{ID: "id123", Name: "OldName", Description: "OldDesc"}
		updated := &domain.Company{ID: "id123", Name: "NewName", Description: "OldDesc"}

		mockRepo.On("GetCompanyById", mock.Anything, "id123").Return(original, nil)
		mockRepo.On("Update", mock.Anything, mock.AnythingOfType("*domain.Company")).Return(updated, nil)

		result, err := uc.UpdateCompany(context.Background(), "id123", &domain.Company{
			Name: "NewName",
		})

		assert.NoError(t, err)
		assert.Equal(t, updated, result)
	})
}
