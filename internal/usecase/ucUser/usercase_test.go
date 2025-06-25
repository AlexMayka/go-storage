package ucUser

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go-storage/internal/domain"
	"go-storage/pkg/auth"
)

type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) CreateUser(ctx context.Context, u *domain.User) (*domain.User, error) {
	args := m.Called(ctx, u)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.User), args.Error(1)
}

func (m *MockUserRepository) GetUserByID(ctx context.Context, id string) (*domain.User, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.User), args.Error(1)
}

func (m *MockUserRepository) GetUserByIDWithCompany(ctx context.Context, id string, companyId string) (*domain.User, error) {
	args := m.Called(ctx, id, companyId)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.User), args.Error(1)
}

func (m *MockUserRepository) GetUserByEmail(ctx context.Context, email string) (*domain.User, error) {
	args := m.Called(ctx, email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.User), args.Error(1)
}

func (m *MockUserRepository) GetUserByUsername(ctx context.Context, username string) (*domain.User, error) {
	args := m.Called(ctx, username)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.User), args.Error(1)
}

func (m *MockUserRepository) GetUsersByCompanyID(ctx context.Context, companyID string) ([]*domain.User, error) {
	args := m.Called(ctx, companyID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.User), args.Error(1)
}

func (m *MockUserRepository) UpdateUser(ctx context.Context, u *domain.User) (*domain.User, error) {
	args := m.Called(ctx, u)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.User), args.Error(1)
}

func (m *MockUserRepository) UpdateUserWithCompany(ctx context.Context, u *domain.User, companyId string) (*domain.User, error) {
	args := m.Called(ctx, u, companyId)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.User), args.Error(1)
}

func (m *MockUserRepository) UpdatePassword(ctx context.Context, userID, hashedPassword string) error {
	args := m.Called(ctx, userID, hashedPassword)
	return args.Error(0)
}

func (m *MockUserRepository) UpdatePasswordWithCompany(ctx context.Context, userID, hashedPassword string, companyId string) error {
	args := m.Called(ctx, userID, hashedPassword, companyId)
	return args.Error(0)
}

func (m *MockUserRepository) UpdateIsActive(ctx context.Context, userID string, isActive bool) error {
	args := m.Called(ctx, userID, isActive)
	return args.Error(0)
}

func (m *MockUserRepository) UpdateIsActiveWithCompany(ctx context.Context, userID string, isActive bool, companyId string) error {
	args := m.Called(ctx, userID, isActive, companyId)
	return args.Error(0)
}

func (m *MockUserRepository) UpdateLastLogin(ctx context.Context, userID string) error {
	args := m.Called(ctx, userID)
	return args.Error(0)
}

func (m *MockUserRepository) UpdateUserRole(ctx context.Context, userID string, roleID string) error {
	args := m.Called(ctx, userID, roleID)
	return args.Error(0)
}

func (m *MockUserRepository) GetAllUsers(ctx context.Context) ([]*domain.User, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.User), args.Error(1)
}

func (m *MockUserRepository) UpdateUserCompany(ctx context.Context, userID string, companyID string) error {
	args := m.Called(ctx, userID, companyID)
	return args.Error(0)
}

type MockAuthRepository struct {
	mock.Mock
}

func (m *MockAuthRepository) GetRoleById(ctx context.Context, roleId string) (*domain.Role, error) {
	args := m.Called(ctx, roleId)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Role), args.Error(1)
}

func (m *MockAuthRepository) GetRoleByName(ctx context.Context, roleId string) (*domain.Role, error) {
	args := m.Called(ctx, roleId)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Role), args.Error(1)
}

func TestNewUseCaseUser(t *testing.T) {
	mockUserRepo := &MockUserRepository{}
	mockAuthRepo := &MockAuthRepository{}

	useCase := NewUseCaseUser(mockUserRepo, mockAuthRepo)

	assert.NotNil(t, useCase)
	assert.Equal(t, mockUserRepo, useCase.repo)
	assert.Equal(t, mockAuthRepo, useCase.authRepo)
}

func TestRegisterUser_Success(t *testing.T) {
	mockUserRepo := &MockUserRepository{}
	mockAuthRepo := &MockAuthRepository{}
	useCase := NewUseCaseUser(mockUserRepo, mockAuthRepo)

	user := &domain.User{
		Username: "testuser",
		Email:    "test@example.com",
		Password: "password123",
	}

	expectedUser := &domain.User{
		ID:       "test-id",
		Username: "testuser",
		Email:    "test@example.com",
		IsActive: true,
	}

	mockUserRepo.On("CreateUser", mock.Anything, mock.MatchedBy(func(u *domain.User) bool {
		return u.Username == "testuser" && u.IsActive == true && u.ID != ""
	})).Return(expectedUser, nil)

	result, err := useCase.RegisterUser(context.Background(), user)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, expectedUser.Username, result.Username)
	mockUserRepo.AssertExpectations(t)
}

func TestRegisterUser_RepositoryError(t *testing.T) {
	mockUserRepo := &MockUserRepository{}
	mockAuthRepo := &MockAuthRepository{}
	useCase := NewUseCaseUser(mockUserRepo, mockAuthRepo)

	user := &domain.User{
		Username: "testuser",
		Email:    "test@example.com",
		Password: "password123",
	}

	mockUserRepo.On("CreateUser", mock.Anything, mock.MatchedBy(func(u *domain.User) bool {
		return u.Username == "testuser" && u.IsActive == true && u.ID != ""
	})).Return(nil, errors.New("database error"))

	result, err := useCase.RegisterUser(context.Background(), user)

	assert.Error(t, err)
	assert.Nil(t, result)
	mockUserRepo.AssertExpectations(t)
}

func TestLogin_ByEmail_Success(t *testing.T) {
	mockUserRepo := &MockUserRepository{}
	mockAuthRepo := &MockAuthRepository{}
	useCase := NewUseCaseUser(mockUserRepo, mockAuthRepo)

	hashedPassword, _ := auth.Hash("password123")
	user := &domain.User{
		ID:       "test-id",
		Email:    "test@example.com",
		Password: hashedPassword,
		IsActive: true,
	}

	mockUserRepo.On("GetUserByEmail", mock.Anything, "test@example.com").Return(user, nil)
	mockUserRepo.On("UpdateLastLogin", mock.Anything, "test-id").Return(nil)

	result, err := useCase.Login(context.Background(), "test@example.com", "password123")

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, user.ID, result.ID)
	mockUserRepo.AssertExpectations(t)
}

func TestLogin_ByUsername_Success(t *testing.T) {
	mockUserRepo := &MockUserRepository{}
	mockAuthRepo := &MockAuthRepository{}
	useCase := NewUseCaseUser(mockUserRepo, mockAuthRepo)

	hashedPassword, _ := auth.Hash("password123")
	user := &domain.User{
		ID:       "test-id",
		Username: "user_name", // underscore makes it not a valid email
		Password: hashedPassword,
		IsActive: true,
	}

	mockUserRepo.On("GetUserByUsername", mock.Anything, "user_name").Return(user, nil)
	mockUserRepo.On("UpdateLastLogin", mock.Anything, "test-id").Return(nil)

	result, err := useCase.Login(context.Background(), "user_name", "password123")

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, user.ID, result.ID)
	mockUserRepo.AssertExpectations(t)
}

func TestLogin_UserNotActive(t *testing.T) {
	mockUserRepo := &MockUserRepository{}
	mockAuthRepo := &MockAuthRepository{}
	useCase := NewUseCaseUser(mockUserRepo, mockAuthRepo)

	user := &domain.User{
		ID:       "test-id",
		Email:    "test@example.com",
		IsActive: false,
	}

	mockUserRepo.On("GetUserByEmail", mock.Anything, "test@example.com").Return(user, nil)

	result, err := useCase.Login(context.Background(), "test@example.com", "password123")

	assert.Error(t, err)
	assert.Nil(t, result)
	mockUserRepo.AssertNotCalled(t, "UpdateLastLogin")
}

func TestLogin_WrongPassword(t *testing.T) {
	mockUserRepo := &MockUserRepository{}
	mockAuthRepo := &MockAuthRepository{}
	useCase := NewUseCaseUser(mockUserRepo, mockAuthRepo)

	hashedPassword, _ := auth.Hash("correctpassword")
	user := &domain.User{
		ID:       "test-id",
		Email:    "test@example.com",
		Password: hashedPassword,
		IsActive: true,
	}

	mockUserRepo.On("GetUserByEmail", mock.Anything, "test@example.com").Return(user, nil)

	result, err := useCase.Login(context.Background(), "test@example.com", "wrongpassword")

	assert.Error(t, err)
	assert.Nil(t, result)
	mockUserRepo.AssertNotCalled(t, "UpdateLastLogin")
}

func TestGetUserByID_Success(t *testing.T) {
	mockUserRepo := &MockUserRepository{}
	mockAuthRepo := &MockAuthRepository{}
	useCase := NewUseCaseUser(mockUserRepo, mockAuthRepo)

	user := &domain.User{
		ID:       "test-id",
		Username: "testuser",
	}

	mockUserRepo.On("GetUserByID", mock.Anything, "test-id").Return(user, nil)

	result, err := useCase.GetUserByID(context.Background(), "test-id")

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, user.ID, result.ID)
	mockUserRepo.AssertExpectations(t)
}

func TestGetUsersByCompany_Success(t *testing.T) {
	mockUserRepo := &MockUserRepository{}
	mockAuthRepo := &MockAuthRepository{}
	useCase := NewUseCaseUser(mockUserRepo, mockAuthRepo)

	users := []*domain.User{
		{ID: "user1", CompanyId: "company1"},
		{ID: "user2", CompanyId: "company1"},
	}

	mockUserRepo.On("GetUsersByCompanyID", mock.Anything, "company1").Return(users, nil)

	result, err := useCase.GetUsersByCompany(context.Background(), "company1")

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Len(t, result, 2)
	mockUserRepo.AssertExpectations(t)
}

func TestUpdateUser_Success(t *testing.T) {
	mockUserRepo := &MockUserRepository{}
	mockAuthRepo := &MockAuthRepository{}
	useCase := NewUseCaseUser(mockUserRepo, mockAuthRepo)

	existingUser := &domain.User{
		ID:       "test-id",
		Username: "oldname",
		Email:    "old@example.com",
	}

	updateData := &domain.User{
		Username: "newname",
		Email:    "new@example.com",
	}

	expectedUser := &domain.User{
		ID:       "test-id",
		Username: "newname",
		Email:    "new@example.com",
	}

	mockUserRepo.On("GetUserByID", mock.Anything, "test-id").Return(existingUser, nil)
	mockUserRepo.On("UpdateUser", mock.Anything, mock.MatchedBy(func(u *domain.User) bool {
		return u.ID == "test-id" && u.Username == "newname"
	})).Return(expectedUser, nil)

	result, err := useCase.UpdateUser(context.Background(), "test-id", updateData)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "newname", result.Username)
	mockUserRepo.AssertExpectations(t)
}

func TestChangePassword_Success(t *testing.T) {
	mockUserRepo := &MockUserRepository{}
	mockAuthRepo := &MockAuthRepository{}
	useCase := NewUseCaseUser(mockUserRepo, mockAuthRepo)

	hashedOldPassword, _ := auth.Hash("oldpassword")
	user := &domain.User{
		ID:       "test-id",
		Password: hashedOldPassword,
	}

	mockUserRepo.On("GetUserByID", mock.Anything, "test-id").Return(user, nil)
	mockUserRepo.On("UpdatePassword", mock.Anything, "test-id", mock.AnythingOfType("string")).Return(nil)

	err := useCase.ChangePassword(context.Background(), "test-id", "oldpassword", "newpassword")

	assert.NoError(t, err)
	mockUserRepo.AssertExpectations(t)
}

func TestChangePassword_WrongOldPassword(t *testing.T) {
	mockUserRepo := &MockUserRepository{}
	mockAuthRepo := &MockAuthRepository{}
	useCase := NewUseCaseUser(mockUserRepo, mockAuthRepo)

	hashedOldPassword, _ := auth.Hash("correctoldpassword")
	user := &domain.User{
		ID:       "test-id",
		Password: hashedOldPassword,
	}

	mockUserRepo.On("GetUserByID", mock.Anything, "test-id").Return(user, nil)

	err := useCase.ChangePassword(context.Background(), "test-id", "wrongoldpassword", "newpassword")

	assert.Error(t, err)
	mockUserRepo.AssertNotCalled(t, "UpdatePassword")
}

func TestDeactivateUser_Success(t *testing.T) {
	mockUserRepo := &MockUserRepository{}
	mockAuthRepo := &MockAuthRepository{}
	useCase := NewUseCaseUser(mockUserRepo, mockAuthRepo)

	user := &domain.User{ID: "test-id"}
	mockUserRepo.On("GetUserByID", mock.Anything, "test-id").Return(user, nil)
	mockUserRepo.On("UpdateIsActive", mock.Anything, "test-id", false).Return(nil)

	err := useCase.DeactivateUser(context.Background(), "test-id")

	assert.NoError(t, err)
	mockUserRepo.AssertExpectations(t)
}

func TestActivateUser_Success(t *testing.T) {
	mockUserRepo := &MockUserRepository{}
	mockAuthRepo := &MockAuthRepository{}
	useCase := NewUseCaseUser(mockUserRepo, mockAuthRepo)

	user := &domain.User{ID: "test-id"}
	mockUserRepo.On("GetUserByID", mock.Anything, "test-id").Return(user, nil)
	mockUserRepo.On("UpdateIsActive", mock.Anything, "test-id", true).Return(nil)

	err := useCase.ActivateUser(context.Background(), "test-id")

	assert.NoError(t, err)
	mockUserRepo.AssertExpectations(t)
}

func TestRefreshToken_Success(t *testing.T) {
	mockUserRepo := &MockUserRepository{}
	mockAuthRepo := &MockAuthRepository{}
	useCase := NewUseCaseUser(mockUserRepo, mockAuthRepo)

	user := &domain.User{
		ID:       "test-id",
		Username: "testuser",
		IsActive: true,
	}

	mockUserRepo.On("GetUserByID", mock.Anything, "test-id").Return(user, nil)

	result, err := useCase.RefreshToken(context.Background(), "test-id")

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, user.ID, result.ID)
	mockUserRepo.AssertExpectations(t)
}

func TestRefreshToken_UserNotActive(t *testing.T) {
	mockUserRepo := &MockUserRepository{}
	mockAuthRepo := &MockAuthRepository{}
	useCase := NewUseCaseUser(mockUserRepo, mockAuthRepo)

	user := &domain.User{
		ID:       "test-id",
		Username: "testuser",
		IsActive: false,
	}

	mockUserRepo.On("GetUserByID", mock.Anything, "test-id").Return(user, nil)

	result, err := useCase.RefreshToken(context.Background(), "test-id")

	assert.Error(t, err)
	assert.Nil(t, result)
	mockUserRepo.AssertNotCalled(t, "UpdateLastLogin")
}

func TestUpdateUserRole_Success(t *testing.T) {
	mockUserRepo := &MockUserRepository{}
	mockAuthRepo := &MockAuthRepository{}
	useCase := NewUseCaseUser(mockUserRepo, mockAuthRepo)

	user := &domain.User{ID: "user-id"}
	role := &domain.Role{ID: "role-id", Name: "admin"}

	mockUserRepo.On("GetUserByID", mock.Anything, "user-id").Return(user, nil)
	mockAuthRepo.On("GetRoleById", mock.Anything, "role-id").Return(role, nil)
	mockUserRepo.On("UpdateUserRole", mock.Anything, "user-id", "role-id").Return(nil)

	err := useCase.UpdateUserRole(context.Background(), "user-id", "role-id")

	assert.NoError(t, err)
	mockUserRepo.AssertExpectations(t)
	mockAuthRepo.AssertExpectations(t)
}

func TestUpdateUserRole_InvalidRole(t *testing.T) {
	mockUserRepo := &MockUserRepository{}
	mockAuthRepo := &MockAuthRepository{}
	useCase := NewUseCaseUser(mockUserRepo, mockAuthRepo)

	user := &domain.User{ID: "user-id"}

	mockUserRepo.On("GetUserByID", mock.Anything, "user-id").Return(user, nil)
	mockAuthRepo.On("GetRoleById", mock.Anything, "invalid-role").Return(nil, errors.New("role not found"))

	err := useCase.UpdateUserRole(context.Background(), "user-id", "invalid-role")

	assert.Error(t, err)
	mockUserRepo.AssertNotCalled(t, "UpdateUserRole")
}

func TestGetAllUsers_Success(t *testing.T) {
	mockUserRepo := &MockUserRepository{}
	mockAuthRepo := &MockAuthRepository{}
	useCase := NewUseCaseUser(mockUserRepo, mockAuthRepo)

	users := []*domain.User{
		{ID: "user1", Username: "user1"},
		{ID: "user2", Username: "user2"},
	}

	mockUserRepo.On("GetAllUsers", mock.Anything).Return(users, nil)

	result, err := useCase.GetAllUsers(context.Background())

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Len(t, result, 2)
	mockUserRepo.AssertExpectations(t)
}

func TestTransferUserToCompany_Success(t *testing.T) {
	mockUserRepo := &MockUserRepository{}
	mockAuthRepo := &MockAuthRepository{}
	useCase := NewUseCaseUser(mockUserRepo, mockAuthRepo)

	user := &domain.User{ID: "user-id"}
	mockUserRepo.On("GetUserByID", mock.Anything, "user-id").Return(user, nil)
	mockUserRepo.On("UpdateUserCompany", mock.Anything, "user-id", "company-id").Return(nil)

	err := useCase.TransferUserToCompany(context.Background(), "user-id", "company-id")

	assert.NoError(t, err)
	mockUserRepo.AssertExpectations(t)
}
