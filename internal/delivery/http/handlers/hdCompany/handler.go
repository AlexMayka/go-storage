package hdCompany

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go-storage/pkg/errors"
	"go-storage/pkg/logger"
	"net/http"
)

type HandlerCompany struct {
	userCase UseCaseCompanyInterface
}

func NewHandlerCompany(userCase UseCaseCompanyInterface) *HandlerCompany {
	return &HandlerCompany{
		userCase: userCase,
	}
}

// RegistrationCompany
// @Summary      Register new company
// @Description  Creates a new company and assigns a unique storage path
// @Tags         companies
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        company  body      RequestRegisterCompanyDto  true  "Company payload"
// @Success      200      {object}  ResponseCompanyDto
// @Failure      400,500  {object}  errors.ErrorResponse
// @Failure      401,403  {object}  errors.ErrorResponse
// @Router       /companies/ [post]
func (h *HandlerCompany) RegistrationCompany(ctx *gin.Context) {
	log := logger.FromContext(ctx)
	var inputData RequestRegisterCompanyDto

	if err := ctx.ShouldBindJSON(&inputData); err != nil {
		log.Error("func registrationCompany: Error in parse input param", "func", "registrationCompany", "err", err)
		errors.HandleError(ctx, errors.BadRequest("Invalid JSON"))
		return
	}

	domainObj := ToDomainCreate(inputData)
	company, errUc := h.userCase.RegisterCompany(ctx, domainObj)
	if errUc != nil {
		log.Error("func registrationCompany: Error work UseCase/Repository", "func", "registrationCompany", "err", errUc)
		errors.HandleError(ctx, errUc)
		return
	}

	ctx.JSON(http.StatusOK, ToResponseCompany(company))
}

// GetCompanyById
// @Summary      Get company by ID
// @Description  Returns a company by its UUID
// @Tags         companies
// @Security     BearerAuth
// @Produce      json
// @Param        id   path      string  true  "Company ID (UUID)"
// @Success      200  {object}  ResponseCompanyDto
// @Failure      400,404,500  {object}  errors.ErrorResponse
// @Failure      401,403  {object}  errors.ErrorResponse
// @Router       /companies/{id} [get]
func (h *HandlerCompany) GetCompanyById(ctx *gin.Context) {
	log := logger.FromContext(ctx)
	id := ctx.Param("id")

	if _, errParse := uuid.Parse(id); errParse != nil {
		log.Error("func getCompanyById: Invalid UUID", "func", "getCompanyById", "err", errParse)
		errors.HandleError(ctx, errors.BadRequest("Invalid ID"))
		return
	}

	company, errUc := h.userCase.GetCompanyById(ctx, id)
	if errUc != nil {
		log.Error("func getCompanyById: Error work UseCase/Repository", "func", "getCompanyById", "err", errUc)
		errors.HandleError(ctx, errUc)
		return
	}

	ctx.JSON(http.StatusOK, ToResponseCompany(company))
}

// GetAllCompanies
// @Summary      Get all companies
// @Description  Returns a list of all active companies
// @Tags         companies
// @Security     BearerAuth
// @Produce      json
// @Success      200  {object}  ResponseCompaniesDto
// @Failure      500  {object}  errors.ErrorResponse
// @Failure      401,403  {object}  errors.ErrorResponse
// @Router       /companies/ [get]
func (h *HandlerCompany) GetAllCompanies(ctx *gin.Context) {
	log := logger.FromContext(ctx)
	companies, errUc := h.userCase.GetAllCompanies(ctx)
	if errUc != nil {
		log.Error("func getAllCompanies: Error work UseCase/Repository", "func", "getAllCompanies", "err", errUc)
		errors.HandleError(ctx, errUc)
		return
	}

	ctx.JSON(http.StatusOK, ToResponseCompanies(companies))
}

// DeleteCompany
// @Summary      Delete company
// @Description  Soft-deletes (deactivates) a company by ID
// @Tags         companies
// @Security     BearerAuth
// @Produce      json
// @Param        id   path      string  true  "Company ID (UUID)"
// @Success      200  {object}  ResponseDeleteCompanyDto
// @Failure      400,404,500  {object}  errors.ErrorResponse
// @Failure      401,403  {object}  errors.ErrorResponse
// @Router       /companies/{id} [delete]
func (h *HandlerCompany) DeleteCompany(ctx *gin.Context) {
	log := logger.FromContext(ctx)
	id := ctx.Param("id")
	if _, errParse := uuid.Parse(id); errParse != nil {
		log.Error("func deleteCompany: Invalid UUID", "func", "deleteCompany", "err", errParse)
		errors.HandleError(ctx, errors.BadRequest("Invalid ID"))
		return
	}

	errUc := h.userCase.DeleteCompany(ctx, id)
	if errUc != nil {
		log.Error("func deleteCompany: Error work UseCase/Repository", "func", "deleteCompany", "err", errUc)
		errors.HandleError(ctx, errUc)
		return
	}

	ctx.JSON(http.StatusOK, ToResponseDelete())
}

// UpdateCompany
// @Summary      Update company
// @Description  Updates name and/or description of the company
// @Tags         companies
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        id       path      string                   true  "Company ID (UUID)"
// @Param        company  body      RequestUpdateCompanyDto  true  "Fields to update"
// @Success      200      {object}  ResponseCompanyDto
// @Failure      400,404,500  {object}  errors.ErrorResponse
// @Failure      401,403  {object}  errors.ErrorResponse
// @Router       /companies/{id} [put]
func (h *HandlerCompany) UpdateCompany(ctx *gin.Context) {
	log := logger.FromContext(ctx)
	id := ctx.Param("id")

	if _, errParse := uuid.Parse(id); errParse != nil {
		log.Error("func updateCompany: Invalid UUID", "func", "updateCompany", "err", errParse)
		errors.HandleError(ctx, errors.BadRequest("Invalid ID"))
		return
	}

	var inputData RequestUpdateCompanyDto
	if err := ctx.ShouldBindJSON(&inputData); err != nil {
		log.Error("func updateCompany: Error in parse input param", "func", "updateCompany", "err", err)
		errors.HandleError(ctx, errors.BadRequest("Invalid JSON"))
		return
	}

	domainObj := ToDomainUpdate(inputData)
	company, errUs := h.userCase.UpdateCompany(ctx, id, domainObj)
	if errUs != nil {
		log.Error("func updateCompany: Error work UseCase/Repository", "func", "updateCompany", "err", errUs)
		errors.HandleError(ctx, errUs)
		return
	}

	ctx.JSON(http.StatusOK, ToResponseCompany(company))
}

// GetMyCompany
// @Summary      Get my company
// @Description  Returns your company information
// @Tags         companies
// @Security     BearerAuth
// @Produce      json
// @Success      200  {object}  ResponseCompanyDto
// @Failure      400,404,500  {object}  errors.ErrorResponse
// @Failure      401,403  {object}  errors.ErrorResponse
// @Router       /companies/me [get]
func (h *HandlerCompany) GetMyCompany(ctx *gin.Context) {
	log := logger.FromContext(ctx)
	companyID := ctx.GetString("company_id")

	if companyID == "" {
		log.Error("func getMyCompany: Company ID is required", "func", "getMyCompany", "err", "empty companyId from JWT")
		errors.HandleError(ctx, errors.BadRequest("Company ID is required"))
		return
	}

	company, errUc := h.userCase.GetCompanyById(ctx, companyID)
	if errUc != nil {
		log.Error("func getMyCompany: Error work UseCase/Repository", "func", "getMyCompany", "err", errUc)
		errors.HandleError(ctx, errUc)
		return
	}

	ctx.JSON(http.StatusOK, ToResponseCompany(company))
}

// UpdateMyCompany
// @Summary      Update my company
// @Description  Updates your company information (company admin only)
// @Tags         companies
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        company  body      RequestUpdateCompanyDto  true  "Fields to update"
// @Success      200      {object}  ResponseCompanyDto
// @Failure      400,404,500  {object}  errors.ErrorResponse
// @Failure      401,403  {object}  errors.ErrorResponse
// @Router       /companies/me [put]
func (h *HandlerCompany) UpdateMyCompany(ctx *gin.Context) {
	log := logger.FromContext(ctx)
	companyID := ctx.GetString("company_id")

	if companyID == "" {
		log.Error("func updateMyCompany: Company ID is required", "func", "updateMyCompany", "err", "empty companyId from JWT")
		errors.HandleError(ctx, errors.BadRequest("Company ID is required"))
		return
	}

	var inputData RequestUpdateCompanyDto
	if err := ctx.ShouldBindJSON(&inputData); err != nil {
		log.Error("func updateMyCompany: Error in parse input param", "func", "updateMyCompany", "err", err)
		errors.HandleError(ctx, errors.BadRequest("Invalid JSON"))
		return
	}

	domainObj := ToDomainUpdate(inputData)
	company, errUs := h.userCase.UpdateCompany(ctx, companyID, domainObj)
	if errUs != nil {
		log.Error("func updateMyCompany: Error work UseCase/Repository", "func", "updateMyCompany", "err", errUs)
		errors.HandleError(ctx, errUs)
		return
	}

	ctx.JSON(http.StatusOK, ToResponseCompany(company))
}
