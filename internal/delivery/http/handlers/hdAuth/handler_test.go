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
			JwtSecret: "test-secret-key-for-testing-purposes-only",
		},
	}
	// Используем правильный способ передачи конфига через context
	newCtx := config.WithConfig(c.Request.Context(), *cfg)
	c.Request = c.Request.WithContext(newCtx)
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
		Password: "Hello123!",
	}

	// Mock expectations
	mockUserCase.On("Login", mock.Anything, "john@example.com", "Hello123!").Return(testUser, nil)

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
		Password: "WrongPass123!",
	}

	// Mock expectations
	mockUserCase.On("Login", mock.Anything, "john@example.com", "WrongPass123!").Return(nil, errors.Unauthorized("invalid credentials"))

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
		Password: "Hello123!",
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

func TestHandlerAuth_RefreshToken_InvalidToken(t *testing.T) {
	mockUserCase := new(MockUseCaseUser)
	mockAuthCase := new(MockUseCaseAuth)
	handler := NewHandlerAuth(mockUserCase, mockAuthCase)

	router := setupTestRouter()
	router.Use(func(c *gin.Context) {
		setupTestContext(c)
		c.Next()
	})
	router.POST("/auth/refresh-token", handler.RefreshToken)

	requestData := RequestRefreshTokenDto{
		Token: "invalid.jwt.token",
	}

	// Request
	jsonData, _ := json.Marshal(requestData)
	req, _ := http.NewRequest("POST", "/auth/refresh-token", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Should return 401 due to invalid token
	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

// Additional tests for missing fields and invalid JSON
func TestHandlerAuth_Login_MissingPassword(t *testing.T) {
	mockUserCase := new(MockUseCaseUser)
	mockAuthCase := new(MockUseCaseAuth)
	handler := NewHandlerAuth(mockUserCase, mockAuthCase)

	router := setupTestRouter()
	router.POST("/auth/login", handler.Login)

	requestData := RequestLoginDto{
		Login:    "john@example.com",
		Password: "",
	}

	jsonData, _ := json.Marshal(requestData)
	req, _ := http.NewRequest("POST", "/auth/login", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestHandlerAuth_Login_InvalidJSON(t *testing.T) {
	mockUserCase := new(MockUseCaseUser)
	mockAuthCase := new(MockUseCaseAuth)
	handler := NewHandlerAuth(mockUserCase, mockAuthCase)

	router := setupTestRouter()
	router.POST("/auth/login", handler.Login)

	req, _ := http.NewRequest("POST", "/auth/login", bytes.NewBuffer([]byte("invalid json")))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestHandlerAuth_Login_InvalidPassword(t *testing.T) {
	mockUserCase := new(MockUseCaseUser)
	mockAuthCase := new(MockUseCaseAuth)
	handler := NewHandlerAuth(mockUserCase, mockAuthCase)

	router := setupTestRouter()
	router.POST("/auth/login", handler.Login)

	requestData := RequestLoginDto{
		Login:    "john@example.com",
		Password: "123", // Too short
	}

	jsonData, _ := json.Marshal(requestData)
	req, _ := http.NewRequest("POST", "/auth/login", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestHandlerAuth_RefreshToken_InvalidJSON(t *testing.T) {
	mockUserCase := new(MockUseCaseUser)
	mockAuthCase := new(MockUseCaseAuth)
	handler := NewHandlerAuth(mockUserCase, mockAuthCase)

	router := setupTestRouter()
	router.POST("/auth/refresh-token", handler.RefreshToken)

	req, _ := http.NewRequest("POST", "/auth/refresh-token", bytes.NewBuffer([]byte("invalid json")))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}