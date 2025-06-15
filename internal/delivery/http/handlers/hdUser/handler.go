package hdUser

import (
	"github.com/gin-gonic/gin"
	"go-storage/internal/config"
	"go-storage/internal/domain"
	"go-storage/internal/utils/valid"
	"go-storage/pkg/errors"
	"go-storage/pkg/jwt"
	"go-storage/pkg/logger"
	"net/http"
)

type HandlerUser struct {
	userCase     UseCaseUserInterface
	userCaseAuth UseCaseAuthInterface
}

func NewHandlerUser(userCase UseCaseUserInterface, userCaseAuth UseCaseAuthInterface) *HandlerUser {
	return &HandlerUser{
		userCase:     userCase,
		userCaseAuth: userCaseAuth,
	}
}

// RegistrationUser
// @Summary      Register new user
// @Description  Creates a new user account with JWT authentication token
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        user  body         RequestRegistrationUserDto  true  "User registration payload"
// @Success      200   {object}     ResponseRegisterUserDto
// @Failure      400,500  {object}  errors.ErrorResponse
// @Router       /user/register [post]
func (h *HandlerUser) RegistrationUser(ctx *gin.Context) {
	log := logger.FromContext(ctx)
	var inputData RequestRegistrationUserDto

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

	role, err := h.userCaseAuth.GetRoleByName(ctx, inputData.RoleName)
	if err != nil {
		log.Error("Get role err", err)
		errors.HandleError(ctx, err)
		return
	}

	inputData.RoleId = role.ID

	domainObj := ToDomainCreate(inputData)
	user, errUc := h.userCase.RegisterUser(ctx, domainObj)
	if errUc != nil {
		log.Error("Register err", errUc)
		errors.HandleError(ctx, errUc)
		return
	}

	userAuth, err := h.generateUserAuth(ctx, user)
	if err != nil {
		log.Error("Generate token err", err)
		errors.HandleError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, ToResponseCreate(userAuth))
}

// GetUserByID
// @Summary      Get user by ID
// @Description  Returns a user by their UUID
// @Tags         users
// @Security     BearerAuth
// @Produce      json
// @Param        id   path      string  true  "User ID (UUID)"
// @Success      200  {object}  ResponseUserDto
// @Failure      400,404,500  {object}  errors.ErrorResponse
// @Failure      401  {object}  errors.ErrorResponse
// @Router       /users/{id} [get]
func (h *HandlerUser) GetUserByID(ctx *gin.Context) {
	log := logger.FromContext(ctx)
	userID := ctx.Param("id")

	if userID == "" {
		errors.HandleError(ctx, errors.BadRequest("User ID is required"))
		return
	}

	user, err := h.userCase.GetUserByID(ctx, userID)
	if err != nil {
		log.Error("Get user err", err)
		errors.HandleError(ctx, err)
		return
	}

	response := ToResponseUser(user)

	ctx.JSON(http.StatusOK, response)
}

// GetUsersByCompany
// @Summary      Get users by company
// @Description  Returns all users belonging to a specific company
// @Tags         users
// @Security     BearerAuth
// @Produce      json
// @Param        company_id   path      string  true  "Company ID (UUID)"
// @Success      200          {object}  ResponseUsersDto
// @Failure      400,404,500  {object}  errors.ErrorResponse
// @Failure      401  {object}  errors.ErrorResponse
// @Router       /users/company/{company_id} [get]
func (h *HandlerUser) GetUsersByCompany(ctx *gin.Context) {
	log := logger.FromContext(ctx)
	companyID := ctx.Param("company_id")

	if companyID == "" {
		errors.HandleError(ctx, errors.BadRequest("Company ID is required"))
		return
	}

	users, err := h.userCase.GetUsersByCompany(ctx, companyID)
	if err != nil {
		log.Error("Get users err", err)
		errors.HandleError(ctx, err)
		return
	}

	response := ToResponseUsers(users)

	ctx.JSON(http.StatusOK, response)
}

// UpdateUser
// @Summary      Update user
// @Description  Updates user profile information
// @Tags         users
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        id    path      string                true  "User ID (UUID)"
// @Param        user  body      RequestUpdateUserDto  true  "User update payload"
// @Success      200   {object}  ResponseUserDto
// @Failure      400,404,500  {object}  errors.ErrorResponse
// @Failure      401  {object}  errors.ErrorResponse
// @Router       /users/{id} [put]
func (h *HandlerUser) UpdateUser(ctx *gin.Context) {
	log := logger.FromContext(ctx)
	userID := ctx.Param("id")
	var inputData RequestUpdateUserDto

	if userID == "" {
		errors.HandleError(ctx, errors.BadRequest("User ID is required"))
		return
	}

	if err := ctx.ShouldBind(&inputData); err != nil {
		log.Error("Bind json err", err)
		errors.HandleError(ctx, errors.BadRequest("Invalid JSON"))
		return
	}

	if inputData.Email != "" && !valid.CheckEmail(inputData.Email) {
		errors.HandleError(ctx, errors.BadRequest("Invalid email"))
		return
	}

	if inputData.Phone != "" && !valid.CheckPhone(inputData.Phone) {
		errors.HandleError(ctx, errors.BadRequest("Invalid phone"))
		return
	}

	if inputData.FirstName != "" && !valid.CheckName(inputData.FirstName) {
		errors.HandleError(ctx, errors.BadRequest("Invalid first_name"))
		return
	}

	if inputData.LastName != "" && !valid.CheckName(inputData.LastName) {
		errors.HandleError(ctx, errors.BadRequest("Invalid last_name"))
		return
	}

	if inputData.SecondName != "" && !valid.CheckName(inputData.SecondName) {
		errors.HandleError(ctx, errors.BadRequest("Invalid second_name"))
		return
	}

	domainObj := ToDomainUpdate(inputData)
	user, err := h.userCase.UpdateUser(ctx, userID, domainObj)
	if err != nil {
		log.Error("Update user err", err)
		errors.HandleError(ctx, err)
		return
	}

	response := ToResponseUser(user)

	ctx.JSON(http.StatusOK, response)
}

// ChangePassword
// @Summary      Change user password
// @Description  Changes user password with old password verification
// @Tags         users
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        id        path      string                     true  "User ID (UUID)"
// @Param        password  body      RequestChangePasswordDto  true  "Password change payload"
// @Success      200       {object}  ResponseMessageDto
// @Failure      400,401,404,500  {object}  errors.ErrorResponse
// @Router       /users/{id}/password [put]
func (h *HandlerUser) ChangePassword(ctx *gin.Context) {
	log := logger.FromContext(ctx)
	userID := ctx.Param("id")
	var inputData RequestChangePasswordDto

	if userID == "" {
		errors.HandleError(ctx, errors.BadRequest("User ID is required"))
		return
	}

	if err := ctx.ShouldBind(&inputData); err != nil {
		log.Error("Bind json err", err)
		errors.HandleError(ctx, errors.BadRequest("Invalid JSON"))
		return
	}

	if len(inputData.NewPassword) < 6 {
		errors.HandleError(ctx, errors.BadRequest("Password must be at least 6 characters"))
		return
	}

	err := h.userCase.ChangePassword(ctx, userID, inputData.OldPassword, inputData.NewPassword)
	if err != nil {
		log.Error("Change password err", err)
		errors.HandleError(ctx, err)
		return
	}

	response := ToResponseMessage("Password changed successfully")
	ctx.JSON(http.StatusOK, response)
}

// DeactivateUser
// @Summary      Deactivate user
// @Description  Deactivates user account (soft delete)
// @Tags         users
// @Security     BearerAuth
// @Produce      json
// @Param        id   path      string  true  "User ID (UUID)"
// @Success      200  {object}  ResponseMessageDto
// @Failure      400,404,500  {object}  errors.ErrorResponse
// @Failure      401  {object}  errors.ErrorResponse
// @Router       /users/{id} [delete]
func (h *HandlerUser) DeactivateUser(ctx *gin.Context) {
	log := logger.FromContext(ctx)
	userID := ctx.Param("id")

	if userID == "" {
		errors.HandleError(ctx, errors.BadRequest("User ID is required"))
		return
	}

	err := h.userCase.DeactivateUser(ctx, userID)
	if err != nil {
		log.Error("Deactivate user err", err)
		errors.HandleError(ctx, err)
		return
	}

	response := ToResponseMessage("User deactivated successfully")
	ctx.JSON(http.StatusOK, response)
}

func (h *HandlerUser) generateUserAuth(ctx *gin.Context, user *domain.User) (*domain.UserWithAuth, error) {
	cfg := config.FromContext(ctx.Request.Context())
	if cfg == nil {
		return nil, errors.InternalServer("config not found in context")
	}

	tokenResp, err := jwt.CreateTokenWithExpiry(user.ID, user.RoleId, []byte(cfg.App.JwtSecret))
	if err != nil {
		return nil, errors.InternalServer("failed to create jwt token")
	}

	return &domain.UserWithAuth{
		User:      user,
		Token:     tokenResp.Token,
		ExpiresAt: tokenResp.ExpiresAt,
	}, nil
}
