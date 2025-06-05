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

func (h *HandlerCompany) RegistrationCompany(ctx *gin.Context) {
	log := logger.FromContext(ctx)
	var inputData RequestRegisterCompanyDto

	if err := ctx.ShouldBindJSON(&inputData); err != nil {
		log.Error("Bind json err", err)
		errors.HandleError(ctx, errors.BadRequest("Invalid JSON"))
		return
	}

	domainObj := ToDomain(inputData)
	company, errUc := h.userCase.RegisterCompany(ctx, domainObj)
	if errUc != nil {
		log.Error("Register err", errUc)
		errors.HandleError(ctx, errUc)
		return
	}

	answer := ToResponseCompany(company)

	ctx.JSON(http.StatusOK, answer)
}

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

func (h *HandlerCompany) DeleteCompany(ctx *gin.Context) {
	log := logger.FromContext(ctx)
	id := ctx.Param("id")
	if _, errParse := uuid.Parse(id); errParse != nil {
		log.Error("Parse uuid err", errParse)
		errors.HandleError(ctx, errParse)
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
