package hdCompany

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go-storage/internal/usecase/ucCompany"
	"go-storage/pkg/errors"
	"go-storage/pkg/logger"
	"net/http"
)

type HandlerCompany struct {
	userCase ucCompany.UseCaseCompanyInterface
}

func NewHandlerCompany(userCase ucCompany.UseCaseCompanyInterface) *HandlerCompany {
	return &HandlerCompany{
		userCase: userCase,
	}
}

// RegistrationCompany
// @Summary      Register new company
// @Description  Creates a new company and assigns a unique storage path
// @Tags         companies
// @Accept       json
// @Produce      json
// @Param        company  body      RequestRegisterCompanyDto  true  "Company payload"
// @Success      200      {object}  ResponseCompanyDto
// @Failure      400,500  {object}  errors.ErrorResponse
// @Router       /companies/ [post]
func (h *HandlerCompany) RegistrationCompany(ctx *gin.Context) {
	log := logger.FromContext(ctx)
	var inputData RequestRegisterCompanyDto

	if err := ctx.ShouldBindJSON(&inputData); err != nil {
		log.Error("Bind json err", err)
		errors.HandleError(ctx, errors.BadRequest("Invalid JSON"))
		return
	}

	domainObj := ToDomainCreate(inputData)
	company, errUc := h.userCase.RegisterCompany(ctx, domainObj)
	if errUc != nil {
		log.Error("Register err", errUc)
		errors.HandleError(ctx, errUc)
		return
	}

	answer := ToResponseCompany(company)

	ctx.JSON(http.StatusOK, answer)
}

// GetCompanyById
// @Summary      Get company by ID
// @Description  Returns a company by its UUID
// @Tags         companies
// @Produce      json
// @Param        id   path      string  true  "Company ID (UUID)"
// @Success      200  {object}  ResponseCompanyDto
// @Failure      400,404,500  {object}  errors.ErrorResponse
// @Router       /companies/{id} [get]
func (h *HandlerCompany) GetCompanyById(ctx *gin.Context) {
	log := logger.FromContext(ctx)
	id := ctx.Param("id")

	if _, errParse := uuid.Parse(id); errParse != nil {
		log.Error("Parse uuid err", errParse)
		errors.HandleError(ctx, errors.BadRequest("Invalid ID"))
		return
	}

	company, errUc := h.userCase.GetCompanyById(ctx, id)
	if errUc != nil {
		log.Error("Error by GetCompanyById", errUc)
		errors.HandleError(ctx, errUc)
		return
	}

	ctx.JSON(http.StatusOK, ToResponseCompany(company))
}

// GetAllCompanies
// @Summary      Get all companies
// @Description  Returns a list of all active companies
// @Tags         companies
// @Produce      json
// @Success      200  {object}  ResponseCompaniesDto
// @Failure      500  {object}  errors.ErrorResponse
// @Router       /companies/ [get]
func (h *HandlerCompany) GetAllCompanies(ctx *gin.Context) {
	log := logger.FromContext(ctx)
	companies, errUc := h.userCase.GetAllCompanies(ctx)
	if errUc != nil {
		log.Error("Error by GetAllCompanies", errUc)
		errors.HandleError(ctx, errUc)
		return
	}

	ctx.JSON(http.StatusOK, ToResponseCompanies(companies))
}

// DeleteCompany
// @Summary      Delete company
// @Description  Soft-deletes (deactivates) a company by ID
// @Tags         companies
// @Produce      json
// @Param        id   path      string  true  "Company ID (UUID)"
// @Success      200  {object}  ResponseDeleteCompanyDto
// @Failure      400,404,500  {object}  errors.ErrorResponse
// @Router       /companies/{id} [delete]
func (h *HandlerCompany) DeleteCompany(ctx *gin.Context) {
	log := logger.FromContext(ctx)
	id := ctx.Param("id")
	if _, errParse := uuid.Parse(id); errParse != nil {
		log.Error("Parse uuid err", errParse)
		errors.HandleError(ctx, errors.BadRequest("Invalid ID"))
		return
	}

	errUc := h.userCase.DeleteCompany(ctx, id)
	if errUc != nil {
		log.Error("Error by DeleteCompany", errUc)
		errors.HandleError(ctx, errUc)
		return
	}

	ctx.JSON(http.StatusOK, ToResponseDelete())
}

// UpdateCompany
// @Summary      Update company
// @Description  Updates name and/or description of the company
// @Tags         companies
// @Accept       json
// @Produce      json
// @Param        id       path      string                   true  "Company ID (UUID)"
// @Param        company  body      RequestUpdateCompanyDto  true  "Fields to update"
// @Success      200      {object}  ResponseCompanyDto
// @Failure      400,404,500  {object}  errors.ErrorResponse
// @Router       /companies/{id} [put]
func (h *HandlerCompany) UpdateCompany(ctx *gin.Context) {
	log := logger.FromContext(ctx)
	id := ctx.Param("id")

	if _, errParse := uuid.Parse(id); errParse != nil {
		log.Error("Parse uuid err", errParse)
		errors.HandleError(ctx, errors.BadRequest("Invalid ID"))
		return
	}

	var inputData RequestUpdateCompanyDto
	if err := ctx.ShouldBindJSON(&inputData); err != nil {
		log.Error("Bind json err", err)
		errors.HandleError(ctx, errors.BadRequest("Invalid JSON"))
		return
	}

	domainObj := ToDomainUpdate(inputData)
	company, errUs := h.userCase.UpdateCompany(ctx, id, domainObj)
	if errUs != nil {
		log.Error("Error by UpdateCompany", errUs)
		errors.HandleError(ctx, errUs)
		return
	}

	ctx.JSON(http.StatusOK, ToResponseCompany(company))
}
