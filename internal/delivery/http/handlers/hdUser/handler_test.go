package hdUser

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go-storage/internal/config"
	"go-storage/internal/domain"
	"go-storage/pkg/errors"
)

type MockUseCaseUser struct {
	mock.Mock
}

func (m *MockUseCaseUser) RegisterUser(ctx context.Context, u *domain.User) (*domain.User, error) {
	args := m.Called(ctx, u)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.User), args.Error(1)
}

func (m *MockUseCaseUser) Login(ctx context.Context, login, password string) (*domain.User, error) {
	args := m.Called(ctx, login, password)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.User), args.Error(1)
}

func (m *MockUseCaseUser) GetUserByID(ctx context.Context, id string) (*domain.User, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.User), args.Error(1)
}

func (m *MockUseCaseUser) GetUsersByCompany(ctx context.Context, companyID string) ([]*domain.User, error) {
	args := m.Called(ctx, companyID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.User), args.Error(1)
}

func (m *MockUseCaseUser) UpdateUser(ctx context.Context, userID string, u *domain.User) (*domain.User, error) {
	args := m.Called(ctx, userID, u)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.User), args.Error(1)
}

func (m *MockUseCaseUser) ChangePassword(ctx context.Context, userID, oldPassword, newPassword string) error {
	args := m.Called(ctx, userID, oldPassword, newPassword)
	return args.Error(0)
}

func (m *MockUseCaseUser) DeactivateUser(ctx context.Context, userID string) error {
	args := m.Called(ctx, userID)
	return args.Error(0)
}

func (m *MockUseCaseUser) RefreshToken(ctx context.Context, userID string) (*domain.User, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.User), args.Error(1)
}

type MockUseCaseAuth struct {
	mock.Mock
}

func (m *MockUseCaseAuth) GetRoleByName(ctx context.Context, roleName string) (*domain.Role, error) {
	args := m.Called(ctx, roleName)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Role), args.Error(1)
}

func setupTestRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	return gin.New()
}

func setupTestContext(ctx *gin.Context) {
	cfg := &config.Config{
		App: config.App{
			JwtSecret: "test-secret-key-for-testing-purposes-only",
		},
	}
	// Используем правильный способ передачи конфига через context
	newCtx := config.WithConfig(ctx.Request.Context(), *cfg)
	ctx.Request = ctx.Request.WithContext(newCtx)
}

func createTestUser() *domain.User {
	return &domain.User{
		ID:        "test-user-id",
		FirstName: "John",
		LastName:  "Doe",
		Username:  "johndoe",
		Email:     "john@example.com",
		Phone:     "+1234567890",
		CompanyId: "test-company-id",
		RoleId:    "test-role-id",
		IsActive:  true,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}

func createTestRole() *domain.Role {
	return &domain.Role{
		ID:   "test-role-id",
		Name: "user",
	}
}

func TestHandlerUser_RegistrationUser_Success(t *testing.T) {
	mockUserCase := new(MockUseCaseUser)
	mockAuthCase := new(MockUseCaseAuth)
	handler := NewHandlerUser(mockUserCase, mockAuthCase)

	router := setupTestRouter()
	router.Use(func(c *gin.Context) {
		setupTestContext(c)
		c.Next()
	})
	router.POST("/users/register", handler.RegistrationUser)

	testUser := createTestUser()
	testRole := createTestRole()

	requestData := RequestRegistrationUserDto{
		FirstName: "John",
		LastName:  "Doe",
		Username:  "johndoe",
		Email:     "john@example.com",
		Phone:     "+1234567890",
		Password:  "password123",
		CompanyId: "test-company-id",
		RoleName:  "user",
	}

	mockAuthCase.On("GetRoleByName", mock.Anything, "user").Return(testRole, nil)
	mockUserCase.On("RegisterUser", mock.Anything, mock.AnythingOfType("*domain.User")).Return(testUser, nil)

	jsonData, _ := json.Marshal(requestData)
	req, _ := http.NewRequest("POST", "/users/register", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response ResponseRegisterUserDto
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "success", response.Status)
	assert.Equal(t, testUser.ID, response.Answer.ID)
	assert.NotEmpty(t, response.Answer.Auth.Token)

	mockUserCase.AssertExpectations(t)
	mockAuthCase.AssertExpectations(t)
}

func TestHandlerUser_RegistrationUser_InvalidJSON(t *testing.T) {
	mockUserCase := new(MockUseCaseUser)
	mockAuthCase := new(MockUseCaseAuth)
	handler := NewHandlerUser(mockUserCase, mockAuthCase)

	router := setupTestRouter()
	router.POST("/users/register", handler.RegistrationUser)

	req, _ := http.NewRequest("POST", "/users/register", bytes.NewBuffer([]byte("invalid json")))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestHandlerUser_RegistrationUser_InvalidEmail(t *testing.T) {
	// Setup
	mockUserCase := new(MockUseCaseUser)
	mockAuthCase := new(MockUseCaseAuth)
	handler := NewHandlerUser(mockUserCase, mockAuthCase)

	router := setupTestRouter()
	router.POST("/users/register", handler.RegistrationUser)

	// Test data with invalid email
	requestData := RequestRegistrationUserDto{
		FirstName: "John",
		LastName:  "Doe",
		Username:  "johndoe",
		Email:     "invalid-email",
		Password:  "password123",
		CompanyId: "test-company-id",
		RoleName:  "user",
	}

	// Request
	jsonData, _ := json.Marshal(requestData)
	req, _ := http.NewRequest("POST", "/users/register", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Assertions
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestHandlerUser_GetUserByID_Success(t *testing.T) {
	mockUserCase := new(MockUseCaseUser)
	mockAuthCase := new(MockUseCaseAuth)
	handler := NewHandlerUser(mockUserCase, mockAuthCase)

	router := setupTestRouter()
	router.GET("/users/:id", handler.GetUserByID)

	testUser := createTestUser()

	mockUserCase.On("GetUserByID", mock.Anything, "test-user-id").Return(testUser, nil)

	req, _ := http.NewRequest("GET", "/users/test-user-id", nil)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response ResponseUserDto
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "success", response.Status)
	assert.Equal(t, testUser.ID, response.Answer.ID)

	mockUserCase.AssertExpectations(t)
}

func TestHandlerUser_GetUserByID_NotFound(t *testing.T) {
	mockUserCase := new(MockUseCaseUser)
	mockAuthCase := new(MockUseCaseAuth)
	handler := NewHandlerUser(mockUserCase, mockAuthCase)

	router := setupTestRouter()
	router.GET("/users/:id", handler.GetUserByID)

	mockUserCase.On("GetUserByID", mock.Anything, "non-existent-id").Return(nil, errors.NotFound("user not found"))

	req, _ := http.NewRequest("GET", "/users/non-existent-id", nil)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)

	mockUserCase.AssertExpectations(t)
}

func TestHandlerUser_UpdateUser_Success(t *testing.T) {
	mockUserCase := new(MockUseCaseUser)
	mockAuthCase := new(MockUseCaseAuth)
	handler := NewHandlerUser(mockUserCase, mockAuthCase)

	router := setupTestRouter()
	router.PUT("/users/:id", handler.UpdateUser)

	testUser := createTestUser()
	testUser.FirstName = "UpdatedJohn"

	requestData := RequestUpdateUserDto{
		FirstName: "UpdatedJohn",
	}

	mockUserCase.On("UpdateUser", mock.Anything, "test-user-id", mock.AnythingOfType("*domain.User")).Return(testUser, nil)

	jsonData, _ := json.Marshal(requestData)
	req, _ := http.NewRequest("PUT", "/users/test-user-id", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response ResponseUserDto
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "success", response.Status)
	assert.Equal(t, "UpdatedJohn", response.Answer.FirstName)

	mockUserCase.AssertExpectations(t)
}

func TestHandlerUser_ChangePassword_Success(t *testing.T) {
	mockUserCase := new(MockUseCaseUser)
	mockAuthCase := new(MockUseCaseAuth)
	handler := NewHandlerUser(mockUserCase, mockAuthCase)

	router := setupTestRouter()
	router.PUT("/users/:id/password", handler.ChangePassword)

	requestData := RequestChangePasswordDto{
		OldPassword: "oldpass123",
		NewPassword: "newpass123",
	}

	mockUserCase.On("ChangePassword", mock.Anything, "test-user-id", "oldpass123", "newpass123").Return(nil)

	jsonData, _ := json.Marshal(requestData)
	req, _ := http.NewRequest("PUT", "/users/test-user-id/password", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response ResponseMessageDto
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "success", response.Status)
	assert.Equal(t, "Password changed successfully", response.Message)

	mockUserCase.AssertExpectations(t)
}

func TestHandlerUser_ChangePassword_TooShort(t *testing.T) {
	mockUserCase := new(MockUseCaseUser)
	mockAuthCase := new(MockUseCaseAuth)
	handler := NewHandlerUser(mockUserCase, mockAuthCase)

	router := setupTestRouter()
	router.PUT("/users/:id/password", handler.ChangePassword)

	requestData := RequestChangePasswordDto{
		OldPassword: "oldpass123",
		NewPassword: "123",
	}

	jsonData, _ := json.Marshal(requestData)
	req, _ := http.NewRequest("PUT", "/users/test-user-id/password", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestHandlerUser_DeactivateUser_Success(t *testing.T) {
	mockUserCase := new(MockUseCaseUser)
	mockAuthCase := new(MockUseCaseAuth)
	handler := NewHandlerUser(mockUserCase, mockAuthCase)

	router := setupTestRouter()
	router.DELETE("/users/:id", handler.DeactivateUser)

	mockUserCase.On("DeactivateUser", mock.Anything, "test-user-id").Return(nil)

	req, _ := http.NewRequest("DELETE", "/users/test-user-id", nil)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response ResponseMessageDto
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "success", response.Status)
	assert.Equal(t, "User deactivated successfully", response.Message)

	mockUserCase.AssertExpectations(t)
}

func TestHandlerUser_GetUsersByCompany_Success(t *testing.T) {
	mockUserCase := new(MockUseCaseUser)
	mockAuthCase := new(MockUseCaseAuth)
	handler := NewHandlerUser(mockUserCase, mockAuthCase)

	router := setupTestRouter()
	router.GET("/companies/:company_id/users", handler.GetUsersByCompany)

	testUsers := []*domain.User{createTestUser()}

	mockUserCase.On("GetUsersByCompany", mock.Anything, "test-company-id").Return(testUsers, nil)

	req, _ := http.NewRequest("GET", "/companies/test-company-id/users", nil)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response ResponseUsersDto
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "success", response.Status)
	assert.Len(t, response.Answer, 1)
	assert.Equal(t, testUsers[0].ID, response.Answer[0].ID)

	mockUserCase.AssertExpectations(t)
}

func TestHandlerUser_UpdateUser_InvalidEmail(t *testing.T) {
	mockUserCase := new(MockUseCaseUser)
	mockAuthCase := new(MockUseCaseAuth)
	handler := NewHandlerUser(mockUserCase, mockAuthCase)

	router := setupTestRouter()
	router.PUT("/users/:id", handler.UpdateUser)

	requestData := RequestUpdateUserDto{
		Email: "invalid-email-format",
	}

	jsonData, _ := json.Marshal(requestData)
	req, _ := http.NewRequest("PUT", "/users/test-user-id", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}
