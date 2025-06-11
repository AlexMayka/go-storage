package hdUser

import (
	"github.com/gin-gonic/gin"
	"go-storage/internal/utils/valid"
	"go-storage/pkg/errors"
	"go-storage/pkg/logger"
	"net/http"
)

type HandlerUser struct {
	uc     UseCaseUserInterface
	ucAuth UseCaseAuthInterface
}

func NewHandlerUser(uc UseCaseUserInterface, ucAuth UseCaseAuthInterface) *HandlerUser {
	return &HandlerUser{
		uc:     uc,
		ucAuth: ucAuth,
	}
}

func (h *HandlerUser) RegistrationUser(ctx *gin.Context) {
	log := logger.FromContext(ctx)
	var inputData RegistrationUserDto

	if err := ctx.ShouldBind(&inputData); err != nil {
		log.Error("Bind json err", err)
		errors.HandleError(ctx, errors.BadRequest("Invalid JSON"))
		return
	}

	if ok := valid.CheckEmail(inputData.Email); !ok {
		log.Error("Invalid email")
		errors.HandleError(ctx, errors.BadRequest("Invalid email"))
		return
	}

	if inputData.Phone != "" && !valid.CheckPhone(inputData.Phone) {
		log.Error("Invalid phone")
		errors.HandleError(ctx, errors.BadRequest("Invalid phone"))
		return
	}

	if ok := valid.CheckName(inputData.FirstName); !ok {
		log.Error("Invalid first_name")
		errors.HandleError(ctx, errors.BadRequest("Invalid first_name"))
		return
	}

	if ok := valid.CheckName(inputData.LastName); !ok {
		log.Error("Invalid last_name")
		errors.HandleError(ctx, errors.BadRequest("Invalid last_name"))
		return
	}

	if inputData.SecondName != "" && !valid.CheckName(inputData.SecondName) {
		log.Error("Invalid second_name")
		errors.HandleError(ctx, errors.BadRequest("Invalid second_name"))
		return
	}

	role, err := h.ucAuth.GetRoleByName(ctx, inputData.RoleName)
	if err != nil {
		log.Error("Get role err", err)
		errors.HandleError(ctx, err)
		return
	}

	inputData.RoleId = role.ID

	domainObj := ToDomainCreate(inputData)
	user, errUc := h.uc.RegisterUser(ctx, domainObj)
	if errUc != nil {
		log.Error("Register err", errUc)
		errors.HandleError(ctx, errUc)
		return
	}

	ctx.JSON(http.StatusOK, ToResponseCreate(user))
}
