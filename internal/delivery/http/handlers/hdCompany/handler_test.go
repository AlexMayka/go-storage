package hdCompany

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go-storage/internal/domain"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

type mockUseCase struct {
	mock.Mock
}

func (m *mockUseCase) GetAllCompanies(ctx context.Context) ([]*domain.Company, error) {
	args := m.Called(ctx)
	return args.Get(0).([]*domain.Company), args.Error(1)
}

func (m *mockUseCase) RegisterCompany(ctx context.Context, c *domain.Company) (*domain.Company, error) {
	args := m.Called(ctx, c)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Company), args.Error(1)
}

func (m *mockUseCase) GetCompanyById(ctx context.Context, id string) (*domain.Company, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Company), args.Error(1)
}

func (m *mockUseCase) DeleteCompany(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *mockUseCase) UpdateCompany(ctx context.Context, id string, c *domain.Company) (*domain.Company, error) {
	args := m.Called(ctx, id, c)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Company), args.Error(1)
}

type reqBuilder = func(body io.Reader) (*http.Request, error)

func setupTestHandler(body io.Reader, reqBuilder reqBuilder) (*mockUseCase, *HandlerCompany, *gin.Context, *httptest.ResponseRecorder) {
	gin.SetMode(gin.TestMode)
	mockUC := new(mockUseCase)
	handler := NewHandlerCompany(mockUC)

	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)
	ctx.Request, _ = reqBuilder(body)

	return mockUC, handler, ctx, w
}

// Func GetAllCompanies

var reqGetAllCompanies = func(body io.Reader) (*http.Request, error) {
	return http.NewRequest(http.MethodGet, "companies", nil)
}

func TestHandlerCompany_GetAllCompanies_Success(t *testing.T) {
	mockUC, handler, ctx, w := setupTestHandler(nil, reqGetAllCompanies)

	expected := []*domain.Company{
		{ID: "Id1", Name: "Test1", Description: "Description1", Path: "company1"},
		{ID: "Id2", Name: "Test2", Description: "Description2", Path: "company2"},
	}

	mockUC.On("GetAllCompanies", mock.Anything).Return(expected, nil)

	handler.GetAllCompanies(ctx)

	assert.Equal(t, http.StatusOK, w.Code)

	for _, value := range expected {
		assert.Contains(t, w.Body.String(), value.ID)
		assert.Contains(t, w.Body.String(), value.Name)
		assert.Contains(t, w.Body.String(), value.Description)
		assert.Contains(t, w.Body.String(), value.Path)
	}

	mockUC.AssertExpectations(t)
}

func TestHandlerCompany_GetAllCompanies_UseCaseError(t *testing.T) {
	mockUC, handler, ctx, w := setupTestHandler(nil, reqGetAllCompanies)
	mockUC.On("GetAllCompanies", mock.Anything).Return([]*domain.Company{}, errors.New("some error"))
	handler.GetAllCompanies(ctx)
	assert.Equal(t, http.StatusInternalServerError, w.Code)
	mockUC.AssertExpectations(t)
}

func TestHandlerCompany_GetAllCompanies_Empty(t *testing.T) {
	mockUC, handler, ctx, w := setupTestHandler(nil, reqGetAllCompanies)
	mockUC.On("GetAllCompanies", mock.Anything).Return([]*domain.Company{}, nil)
	handler.GetAllCompanies(ctx)
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), `"answer":[]`)
	mockUC.AssertExpectations(t)
}

// Func RegistrationCompany

var reqRegisterCompany = func(body io.Reader) (*http.Request, error) {
	reqPost, _ := http.NewRequest(http.MethodPost, "companies", body)
	reqPost.Header.Set("Content-Type", "application/json")
	return reqPost, nil
}

func TestHandlerCompany_RegisterCompany_Success(t *testing.T) {
	cases := []struct {
		name        string
		inputName   string
		inputDesc   string
		returnModel *domain.Company
	}{
		{
			name:      "Task 1",
			inputName: "Name1",
			inputDesc: "Desc1",
			returnModel: &domain.Company{
				ID:          "id1",
				Name:        "Name1",
				Description: "Desc1",
				Path:        "name1",
			},
		},
		{
			name:      "Task 2",
			inputName: "Name2",
			inputDesc: "Desc2",
			returnModel: &domain.Company{
				ID:          "id2",
				Name:        "Name2",
				Description: "Desc2",
				Path:        "name2",
			},
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			bodyMap := map[string]string{
				"name":        c.inputName,
				"description": c.inputDesc,
			}
			bodyJSON, _ := json.Marshal(bodyMap)

			mockUC, handler, ctx, w := setupTestHandler(bytes.NewBuffer(bodyJSON), reqRegisterCompany)

			mockUC.On("RegisterCompany", mock.Anything, mock.Anything).Return(c.returnModel, nil)

			handler.RegistrationCompany(ctx)

			assert.Equal(t, http.StatusOK, w.Code)
			assert.Contains(t, w.Body.String(), c.returnModel.ID)
			assert.Contains(t, w.Body.String(), c.returnModel.Name)
			assert.Contains(t, w.Body.String(), c.returnModel.Path)
			mockUC.AssertExpectations(t)
		})
	}
}

func TestHandlerCompany_RegisterCompany_InvalidJSON(t *testing.T) {
	body := []byte(`{"name": "Test", "description": "no-quote}`)

	_, handler, ctx, w := setupTestHandler(bytes.NewBuffer(body), reqRegisterCompany)

	handler.RegistrationCompany(ctx)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "Invalid JSON")
}

func TestHandlerCompany_RegisterCompany_UseCaseError(t *testing.T) {
	body := map[string]string{
		"name":        "Company X",
		"description": "Description X",
	}
	bodyJSON, _ := json.Marshal(body)

	mockUC, handler, ctx, w := setupTestHandler(bytes.NewBuffer(bodyJSON), reqRegisterCompany)

	mockUC.On("RegisterCompany", mock.Anything, mock.Anything).
		Return(nil, errors.New("internal error"))

	handler.RegistrationCompany(ctx)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Contains(t, w.Body.String(), "internal_error")
	mockUC.AssertExpectations(t)
}

func TestHandlerCompany_RegisterCompany_EmptyFields(t *testing.T) {
	body := map[string]string{
		"name":        "",
		"description": "Valid",
	}
	bodyJSON, _ := json.Marshal(body)

	_, handler, ctx, w := setupTestHandler(bytes.NewBuffer(bodyJSON), reqRegisterCompany)

	handler.RegistrationCompany(ctx)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "Invalid JSON")
}

// Func GetCompanyById

func TestHandlerCompany_GetCompanyById_Success(t *testing.T) {
	id := "a3b2c1d4-1111-2222-3333-444455556666"
	expected := &domain.Company{
		ID:          id,
		Name:        "Example Co",
		Description: "A test company",
		Path:        "example-co",
	}

	reqURL := "/companies/" + id
	req := httptest.NewRequest(http.MethodGet, reqURL, nil)

	rec := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(rec)
	ctx.Request = req
	ctx.Params = gin.Params{{Key: "id", Value: id}}

	mockUC := new(mockUseCase)
	handler := NewHandlerCompany(mockUC)

	mockUC.On("GetCompanyById", mock.Anything, id).Return(expected, nil)

	handler.GetCompanyById(ctx)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Contains(t, rec.Body.String(), expected.ID)
	assert.Contains(t, rec.Body.String(), expected.Name)
	mockUC.AssertExpectations(t)
}

func TestHandlerCompany_GetCompanyById_InvalidUUID(t *testing.T) {
	invalidID := "not-a-uuid"
	reqURL := "/companies/" + invalidID
	req := httptest.NewRequest(http.MethodGet, reqURL, nil)

	rec := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(rec)
	ctx.Request = req
	ctx.Params = gin.Params{{Key: "id", Value: invalidID}}

	handler := NewHandlerCompany(nil)

	handler.GetCompanyById(ctx)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
	assert.Contains(t, rec.Body.String(), "Invalid ID")
}

func TestHandlerCompany_GetCompanyById_UseCaseError(t *testing.T) {
	id := "a3b2c1d4-1111-2222-3333-444455556666"
	reqURL := "/companies/" + id
	req := httptest.NewRequest(http.MethodGet, reqURL, nil)

	rec := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(rec)
	ctx.Request = req
	ctx.Params = gin.Params{{Key: "id", Value: id}}

	mockUC := new(mockUseCase)
	handler := NewHandlerCompany(mockUC)

	mockUC.On("GetCompanyById", mock.Anything, id).Return(nil, errors.New("not found"))

	handler.GetCompanyById(ctx)

	assert.Equal(t, http.StatusInternalServerError, rec.Code)
	mockUC.AssertExpectations(t)
}

// Func DeleteCompany
func TestHandlerCompany_DeleteCompany_Success(t *testing.T) {
	id := "a3b2c1d4-1111-2222-3333-444455556666"
	req := httptest.NewRequest(http.MethodDelete, "/companies/"+id, nil)

	rec := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(rec)
	ctx.Request = req
	ctx.Params = gin.Params{{Key: "id", Value: id}}

	mockUC := new(mockUseCase)
	handler := NewHandlerCompany(mockUC)

	mockUC.On("DeleteCompany", mock.Anything, id).Return(nil)

	handler.DeleteCompany(ctx)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Contains(t, rec.Body.String(), `"status":"success"`)
	mockUC.AssertExpectations(t)
}

func TestHandlerCompany_DeleteCompany_InvalidUUID(t *testing.T) {
	invalidID := "not-a-uuid"
	req := httptest.NewRequest(http.MethodDelete, "/companies/"+invalidID, nil)

	rec := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(rec)
	ctx.Request = req
	ctx.Params = gin.Params{{Key: "id", Value: invalidID}}

	handler := NewHandlerCompany(nil)

	handler.DeleteCompany(ctx)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestHandlerCompany_DeleteCompany_UseCaseError(t *testing.T) {
	id := "a3b2c1d4-1111-2222-3333-444455556666"
	req := httptest.NewRequest(http.MethodDelete, "/companies/"+id, nil)

	rec := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(rec)
	ctx.Request = req
	ctx.Params = gin.Params{{Key: "id", Value: id}}

	mockUC := new(mockUseCase)
	handler := NewHandlerCompany(mockUC)

	mockUC.On("DeleteCompany", mock.Anything, id).Return(errors.New("delete failed"))

	handler.DeleteCompany(ctx)

	assert.Equal(t, http.StatusInternalServerError, rec.Code)
	mockUC.AssertExpectations(t)
}

// Func UpdateCompany
func TestHandlerCompany_UpdateCompany_Success(t *testing.T) {
	id := "a3b2c1d4-1111-2222-3333-444455556666"
	body := map[string]string{
		"name":        "UpdatedName",
		"description": "UpdatedDesc",
	}
	bodyJSON, _ := json.Marshal(body)

	req := httptest.NewRequest(http.MethodPut, "/companies/"+id, bytes.NewBuffer(bodyJSON))
	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(rec)
	ctx.Request = req
	ctx.Params = gin.Params{{Key: "id", Value: id}}

	mockUC := new(mockUseCase)
	handler := NewHandlerCompany(mockUC)

	expected := &domain.Company{
		ID:          id,
		Name:        "UpdatedName",
		Description: "UpdatedDesc",
		Path:        "updatedname",
	}
	mockUC.On("UpdateCompany", mock.Anything, id, mock.AnythingOfType("*domain.Company")).Return(expected, nil)

	handler.UpdateCompany(ctx)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Contains(t, rec.Body.String(), expected.ID)
	assert.Contains(t, rec.Body.String(), expected.Name)
	mockUC.AssertExpectations(t)
}

func TestHandlerCompany_UpdateCompany_InvalidUUID(t *testing.T) {
	invalidID := "invalid-uuid"
	body := map[string]string{
		"name":        "UpdatedName",
		"description": "UpdatedDesc",
	}
	bodyJSON, _ := json.Marshal(body)

	req := httptest.NewRequest(http.MethodPut, "/companies/"+invalidID, bytes.NewBuffer(bodyJSON))
	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(rec)
	ctx.Request = req
	ctx.Params = gin.Params{{Key: "id", Value: invalidID}}

	handler := NewHandlerCompany(nil)

	handler.UpdateCompany(ctx)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
	assert.Contains(t, rec.Body.String(), "Invalid ID")
}

func TestHandlerCompany_UpdateCompany_InvalidJSON(t *testing.T) {
	id := "a3b2c1d4-1111-2222-3333-444455556666"
	body := []byte(`{"name": "Name", "description": "MissingEndQuote}`)

	req := httptest.NewRequest(http.MethodPut, "/companies/"+id, bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(rec)
	ctx.Request = req
	ctx.Params = gin.Params{{Key: "id", Value: id}}

	handler := NewHandlerCompany(nil)

	handler.UpdateCompany(ctx)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
	assert.Contains(t, rec.Body.String(), "Invalid JSON")
}

func TestHandlerCompany_UpdateCompany_UseCaseError(t *testing.T) {
	id := "a3b2c1d4-1111-2222-3333-444455556666"
	body := map[string]string{
		"name":        "UpdatedName",
		"description": "UpdatedDesc",
	}
	bodyJSON, _ := json.Marshal(body)

	req := httptest.NewRequest(http.MethodPut, "/companies/"+id, bytes.NewBuffer(bodyJSON))
	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(rec)
	ctx.Request = req
	ctx.Params = gin.Params{{Key: "id", Value: id}}

	mockUC := new(mockUseCase)
	handler := NewHandlerCompany(mockUC)

	mockUC.On("UpdateCompany", mock.Anything, id, mock.AnythingOfType("*domain.Company")).Return(nil, errors.New("update failed"))

	handler.UpdateCompany(ctx)

	assert.Equal(t, http.StatusInternalServerError, rec.Code)
	mockUC.AssertExpectations(t)
}

// Tests for GetMyCompany
func TestHandlerCompany_GetMyCompany_Success(t *testing.T) {
	mockUseCase := new(mockUseCase)
	handler := NewHandlerCompany(mockUseCase)

	router := gin.New()
	router.Use(func(c *gin.Context) {
		c.Set("company_id", "test-company-id")
		c.Next()
	})
	router.GET("/companies/me", handler.GetMyCompany)

	testCompany := &domain.Company{
		ID:          "test-company-id",
		Name:        "Test Company",
		Description: "Test Description",
	}

	mockUseCase.On("GetCompanyById", mock.Anything, "test-company-id").Return(testCompany, nil)

	req, _ := http.NewRequest("GET", "/companies/me", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response ResponseCompanyDto
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "success", response.Status)
	assert.Equal(t, testCompany.ID, response.Answer.Id)

	mockUseCase.AssertExpectations(t)
}

func TestHandlerCompany_GetMyCompany_MissingCompanyID(t *testing.T) {
	mockUseCase := new(mockUseCase)
	handler := NewHandlerCompany(mockUseCase)

	router := gin.New()
	router.GET("/companies/me", handler.GetMyCompany)

	req, _ := http.NewRequest("GET", "/companies/me", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

// Tests for UpdateMyCompany
func TestHandlerCompany_UpdateMyCompany_Success(t *testing.T) {
	mockUseCase := new(mockUseCase)
	handler := NewHandlerCompany(mockUseCase)

	router := gin.New()
	router.Use(func(c *gin.Context) {
		c.Set("company_id", "test-company-id")
		c.Next()
	})
	router.PUT("/companies/me", handler.UpdateMyCompany)

	testCompany := &domain.Company{
		ID:          "test-company-id",
		Name:        "Updated Company",
		Description: "Updated Description",
	}

	requestData := RequestUpdateCompanyDto{
		Name:        "Updated Company",
		Description: "Updated Description",
	}

	mockUseCase.On("UpdateCompany", mock.Anything, "test-company-id", mock.AnythingOfType("*domain.Company")).Return(testCompany, nil)

	jsonData, _ := json.Marshal(requestData)
	req, _ := http.NewRequest("PUT", "/companies/me", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response ResponseCompanyDto
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "success", response.Status)
	assert.Equal(t, testCompany.Name, response.Answer.Name)

	mockUseCase.AssertExpectations(t)
}

func TestHandlerCompany_UpdateMyCompany_MissingCompanyID(t *testing.T) {
	mockUseCase := new(mockUseCase)
	handler := NewHandlerCompany(mockUseCase)

	router := gin.New()
	router.PUT("/companies/me", handler.UpdateMyCompany)

	requestData := RequestUpdateCompanyDto{
		Name: "Updated Company",
	}

	jsonData, _ := json.Marshal(requestData)
	req, _ := http.NewRequest("PUT", "/companies/me", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestHandlerCompany_UpdateMyCompany_InvalidJSON(t *testing.T) {
	mockUseCase := new(mockUseCase)
	handler := NewHandlerCompany(mockUseCase)

	router := gin.New()
	router.Use(func(c *gin.Context) {
		c.Set("company_id", "test-company-id")
		c.Next()
	})
	router.PUT("/companies/me", handler.UpdateMyCompany)

	req, _ := http.NewRequest("PUT", "/companies/me", bytes.NewBuffer([]byte("invalid json")))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}
