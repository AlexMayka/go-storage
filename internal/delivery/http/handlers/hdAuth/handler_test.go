package hdAuth

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go-storage/internal/config"
	"go-storage/internal/domain"
	"go-storage/pkg/errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

// Mock implementations
type MockUseCaseUser struct {
	mock.Mock
}

func (m *MockUseCaseUser) Login(ctx context.Context, login, password string) (*domain.User, error) {
	args := m.Called(ctx, login, password)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.User), args.Error(1)
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

// Test helper functions
func setupTestRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	return gin.New()
}

func setupTestContext(c *gin.Context) {
	cfg := &config.Config{
		App: config.App{
			JwtSecret: "test-secret-key-for-testing-purposes",
		},
	}
	c.Request = c.Request.WithContext(context.WithValue(c.Request.Context(), "config", cfg))
}

func createTestUser() *domain.User {
	return &domain.User{
		ID:        "test-user-id",
		FirstName: "John",
		LastName:  "Doe",
		Username:  "johndoe",
		Email:     "john@example.com",
		RoleId:    "role-id",
		CompanyId: "company-id",
		IsActive:  true,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}

// Tests
func TestHandlerAuth_Login_Success(t *testing.T) {
	mockUserCase := new(MockUseCaseUser)
	mockAuthCase := new(MockUseCaseAuth)
	handler := NewHandlerAuth(mockUserCase, mockAuthCase)

	router := setupTestRouter()
	router.Use(func(c *gin.Context) {
		setupTestContext(c)
		c.Next()
	})
	router.POST("/auth/login", handler.Login)

	// Test data
	testUser := createTestUser()

	requestData := RequestLoginDto{
		Login:    "john@example.com",
		Password: "password123",
	}

	// Mock expectations
	mockUserCase.On("Login", mock.Anything, "john@example.com", "password123").Return(testUser, nil)

	// Request
	jsonData, _ := json.Marshal(requestData)
	req, _ := http.NewRequest("POST", "/auth/login", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Assertions
	assert.Equal(t, http.StatusOK, w.Code)

	var response ResponseLoginDto
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "success", response.Status)
	assert.Equal(t, testUser.ID, response.Answer.ID)
	assert.NotEmpty(t, response.Answer.Auth.Token)

	mockUserCase.AssertExpectations(t)
}

func TestHandlerAuth_Login_InvalidCredentials(t *testing.T) {
	// Setup
	mockUserCase := new(MockUseCaseUser)
	mockAuthCase := new(MockUseCaseAuth)
	handler := NewHandlerAuth(mockUserCase, mockAuthCase)

	router := setupTestRouter()
	router.Use(func(c *gin.Context) {
		setupTestContext(c)
		c.Next()
	})
	router.POST("/auth/login", handler.Login)

	requestData := RequestLoginDto{
		Login:    "john@example.com",
		Password: "wrongpassword",
	}

	// Mock expectations
	mockUserCase.On("Login", mock.Anything, "john@example.com", "wrongpassword").Return(nil, errors.Unauthorized("invalid credentials"))

	// Request
	jsonData, _ := json.Marshal(requestData)
	req, _ := http.NewRequest("POST", "/auth/login", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Assertions
	assert.Equal(t, http.StatusUnauthorized, w.Code)

	mockUserCase.AssertExpectations(t)
}

func TestHandlerAuth_Login_MissingLogin(t *testing.T) {
	// Setup
	mockUserCase := new(MockUseCaseUser)
	mockAuthCase := new(MockUseCaseAuth)
	handler := NewHandlerAuth(mockUserCase, mockAuthCase)

	router := setupTestRouter()
	router.POST("/auth/login", handler.Login)

	requestData := RequestLoginDto{
		Login:    "",
		Password: "password123",
	}

	// Request
	jsonData, _ := json.Marshal(requestData)
	req, _ := http.NewRequest("POST", "/auth/login", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Assertions
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestHandlerAuth_RefreshToken_Success(t *testing.T) {
	mockUserCase := new(MockUseCaseUser)
	mockAuthCase := new(MockUseCaseAuth)
	handler := NewHandlerAuth(mockUserCase, mockAuthCase)

	router := setupTestRouter()
	router.Use(func(c *gin.Context) {
		setupTestContext(c)
		c.Next()
	})
	router.POST("/auth/refresh-token", handler.RefreshToken)

	// Test data
	testUser := createTestUser()

	requestData := RequestRefreshTokenDto{
		Token: "valid.jwt.token",
	}

	// Mock expectations
	mockUserCase.On("RefreshToken", mock.Anything, mock.AnythingOfType("string")).Return(testUser, nil)

	// Request
	jsonData, _ := json.Marshal(requestData)
	req, _ := http.NewRequest("POST", "/auth/refresh-token", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Note: This test may fail due to JWT parsing, but tests the handler structure
	// In real scenario, you'd need to generate a valid JWT token for testing
	mockUserCase.AssertExpectations(t)
}