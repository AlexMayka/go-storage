package hdUser

import (
	"github.com/gin-gonic/gin"
	"go-storage/internal/utils/valid"
	"go-storage/pkg/errors"
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

// AVAILABLE TO ALL USERS

// GetMe
// @Summary      Get info about yourself
// @Description  Getting information about yourself by JWT
// @Tags         users
// @Accept       json
// @Produce      json
// @Success      200   {object}     ResponseUserDto
// @Failure      400,500  {object}  errors.ErrorResponse
// @Router       /users/me [get]
func (h *HandlerUser) GetMe(ctx *gin.Context) {
	log := logger.FromContext(ctx)
	userID := ctx.GetString("user_id")

	if userID == "" {
		log.Error("func getMe: User ID is required", "func", "getMe", "err", "empty userId from JWT")
		errors.HandleError(ctx, errors.BadRequest("User ID is required"))
		return
	}

	user, err := h.userCase.GetUserByID(ctx, userID)
	if err != nil {
		log.Error("func getMe: Error work UseCase/Repository", "func", "getMe", "err", err)
		errors.HandleError(ctx, err)
		return
	}

	response := ToResponseUser(user)
	ctx.JSON(http.StatusOK, response)
}

// UpdateMe
// @Summary      Update your profile
// @Description  Update your own profile information
// @Tags         users
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        user  body      RequestUpdateUserDto  true  "User update payload"
// @Success      200   {object}  ResponseUserDto
// @Failure      400,500  {object}  errors.ErrorResponse
// @Failure      401  {object}  errors.ErrorResponse
// @Router       /users/me [put]
func (h *HandlerUser) UpdateMe(ctx *gin.Context) {
	log := logger.FromContext(ctx)
	userID := ctx.GetString("user_id")

	if userID == "" {
		log.Error("func updateMe: User ID is required", "func", "UpdateMe", "err", "empty userId from JWT")
		errors.HandleError(ctx, errors.BadRequest("User ID is required"))
		return
	}

	var inputData RequestUpdateUserDto
	if err := ctx.ShouldBind(&inputData); err != nil {
		log.Error("func updateMe: Error in parse input param", "func", "UpdateMe", "err", err)
		errors.HandleError(ctx, errors.BadRequest("Invalid JSON"))
		return
	}

	updateUser := ToDomainUpdate(inputData)

	user, err := h.userCase.UpdateUser(ctx, userID, updateUser)
	if err != nil {
		log.Error("func updateMe: Error work UseCase/Repository", "func", "updateMe", "err", err)
		errors.HandleError(ctx, err)
		return
	}

	answer := ToResponseUser(user)
	ctx.JSON(http.StatusOK, answer)
}

// UpdatePasswordMe
// @Summary      Change your password
// @Description  Change your own password with old password verification
// @Tags         users
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        password  body      RequestChangePasswordDto  true  "Password change payload"
// @Success      200       {object}  ResponseMessageDto
// @Failure      400,401,500  {object}  errors.ErrorResponse
// @Router       /users/me/password [put]
func (h *HandlerUser) UpdatePasswordMe(ctx *gin.Context) {
	log := logger.FromContext(ctx)
	userID := ctx.GetString("user_id")

	if userID == "" {
		log.Error("func updatePasswordMe: User ID is required", "func", "updatePasswordMe", "err", "empty userId from JWT")
		errors.HandleError(ctx, errors.BadRequest("User ID is required"))
		return
	}

	var inputData RequestChangePasswordDto
	if err := ctx.ShouldBind(&inputData); err != nil {
		log.Error("func updatePasswordMe: Error in parse input param", "func", "updatePasswordMe", "err", err)
		errors.HandleError(ctx, errors.BadRequest("Invalid JSON"))
		return
	}

	if _, err := valid.CheckPassword(inputData.NewPassword); err != nil {
		log.Error("func updatePasswordMe: Invalid password", "func", "updatePasswordMe", "err", err)
		errors.HandleError(ctx, errors.BadRequest(err.Error()))
		return
	}

	err := h.userCase.ChangePassword(ctx, userID, inputData.OldPassword, inputData.NewPassword)
	if err != nil {
		log.Error("func updatePasswordMe: Error work UseCase/Repository", "func", "updatePasswordMe", "err", err)
		errors.HandleError(ctx, err)
		return
	}

	response := ToResponseMessage("Password changed successfully")
	ctx.JSON(http.StatusOK, response)
}

// GetAllUsersOfYourCompany
// @Summary      Get all users in your company
// @Description  Returns all users from your company (company admin only)
// @Tags         users
// @Security     BearerAuth
// @Produce      json
// @Success      200  {object}  ResponseUsersDto
// @Failure      400,500  {object}  errors.ErrorResponse
// @Failure      401  {object}  errors.ErrorResponse
// @Router       /users/company [get]
func (h *HandlerUser) GetAllUsersOfYourCompany(ctx *gin.Context) {
	log := logger.FromContext(ctx)
	companyID := ctx.GetString("company_id")

	if companyID == "" {
		log.Error("func getAllUsersOfYourCompany: Company ID is required", "func", "getAllUsersOfYourCompany", "err", "empty companyId from JWT")
		errors.HandleError(ctx, errors.BadRequest("Company ID is required"))
		return
	}

	users, err := h.userCase.GetUsersByCompany(ctx, companyID)
	if err != nil {
		log.Error("func getAllUsersOfYourCompany: Error work UseCase/Repository", "func", "getAllUsersOfYourCompany", "err", err)
		errors.HandleError(ctx, err)
		return
	}

	response := ToResponseUsers(users)

	ctx.JSON(http.StatusOK, response)
}

// AVAILABLE FOR ADMINS

// RegistrationUser
// @Summary      Register new user
// @Description  Register new user in your company (company admin only)
// @Tags         users
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        user  body      RequestRegistrationUserDto  true  "User registration payload"
// @Success      200   {object}  ResponseUserDto
// @Failure      400,500  {object}  errors.ErrorResponse
// @Failure      401  {object}  errors.ErrorResponse
// @Router       /users [post]
func (h *HandlerUser) RegistrationUser(ctx *gin.Context) {
	log := logger.FromContext(ctx)
	var inputData RequestRegistrationUserDto
	var companyId = ctx.GetString("company_id")

	if err := ctx.ShouldBind(&inputData); err != nil {
		log.Error("func registrationUser: Error in parse input param", "func", "registrationUser", "err", err)
		errors.HandleError(ctx, errors.BadRequest("Invalid JSON"))
		return
	}

	if ok := valid.CheckEmail(inputData.Email); !ok {
		log.Error("func registrationUser: Invalid email", "func", "registrationUser", "email", inputData.Email)
		errors.HandleError(ctx, errors.BadRequest("Invalid email"))
		return
	}

	if inputData.Phone != "" && !valid.CheckPhone(inputData.Phone) {
		log.Error("func registrationUser: Invalid phone", "func", "registrationUser", "phone", inputData.Phone)
		errors.HandleError(ctx, errors.BadRequest("Invalid phone"))
		return
	}

	if ok := valid.CheckName(inputData.FirstName); !ok {
		log.Error("func registrationUser: Invalid first_name", "func", "registrationUser", "firstName", inputData.FirstName)
		errors.HandleError(ctx, errors.BadRequest("Invalid first_name"))
		return
	}

	if ok := valid.CheckName(inputData.LastName); !ok {
		log.Error("func registrationUser: Invalid last_name", "func", "registrationUser", "lastName", inputData.LastName)
		errors.HandleError(ctx, errors.BadRequest("Invalid last_name"))
		return
	}

	if inputData.SecondName != "" && !valid.CheckName(inputData.SecondName) {
		log.Error("func registrationUser: Invalid second_name", "func", "registrationUser", "secondName", inputData.SecondName)
		errors.HandleError(ctx, errors.BadRequest("Invalid second_name"))
		return
	}

	role, err := h.userCaseAuth.GetDefaultRole(ctx)
	if err != nil {
		log.Error("func registrationUser: Error get default role", "func", "registrationUser", "err", err)
		errors.HandleError(ctx, err)
		return
	}

	inputData.RoleId = role.ID
	inputData.CompanyId = companyId

	domainObj := ToDomainCreate(inputData)
	user, errUc := h.userCase.RegisterUser(ctx, domainObj)
	if errUc != nil {
		log.Error("func registrationUser: Error work UseCase/Repository", "func", "registrationUser", "err", errUc)
		errors.HandleError(ctx, errUc)
		return
	}

	response := ToResponseUser(user)

	ctx.JSON(http.StatusOK, response)
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
		log.Error("func getUserByID: User ID is required", "func", "getUserByID", "err", "empty userId from param")
		errors.HandleError(ctx, errors.BadRequest("User ID is required"))
		return
	}

	user, err := h.userCase.GetUserByID(ctx, userID)
	if err != nil {
		log.Error("func getUserByID: Error work UseCase/Repository", "func", "getUserByID", "err", err)
		errors.HandleError(ctx, err)
		return
	}

	response := ToResponseUser(user)

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
		log.Error("func updateUser: User ID is required", "func", "updateUser", "err", "empty userId from param")
		errors.HandleError(ctx, errors.BadRequest("User ID is required"))
		return
	}

	if err := ctx.ShouldBind(&inputData); err != nil {
		log.Error("func updateUser: Error in parse input param", "func", "updateUser", "err", err)
		errors.HandleError(ctx, errors.BadRequest("Invalid JSON"))
		return
	}

	if inputData.Email != "" && !valid.CheckEmail(inputData.Email) {
		log.Error("func updateUser: Invalid email", "func", "updateUser", "email", inputData.Email)
		errors.HandleError(ctx, errors.BadRequest("Invalid email"))
		return
	}

	if inputData.Phone != "" && !valid.CheckPhone(inputData.Phone) {
		log.Error("func updateUser: Invalid phone", "func", "updateUser", "phone", inputData.Phone)
		errors.HandleError(ctx, errors.BadRequest("Invalid phone"))
		return
	}

	if inputData.FirstName != "" && !valid.CheckName(inputData.FirstName) {
		log.Error("func updateUser: Invalid first_name", "func", "updateUser", "firstName", inputData.FirstName)
		errors.HandleError(ctx, errors.BadRequest("Invalid first_name"))
		return
	}

	if inputData.LastName != "" && !valid.CheckName(inputData.LastName) {
		log.Error("func updateUser: Invalid last_name", "func", "updateUser", "lastName", inputData.LastName)
		errors.HandleError(ctx, errors.BadRequest("Invalid last_name"))
		return
	}

	if inputData.SecondName != "" && !valid.CheckName(inputData.SecondName) {
		log.Error("func updateUser: Invalid second_name", "func", "updateUser", "secondName", inputData.SecondName)
		errors.HandleError(ctx, errors.BadRequest("Invalid second_name"))
		return
	}

	domainObj := ToDomainUpdate(inputData)
	user, err := h.userCase.UpdateUser(ctx, userID, domainObj)
	if err != nil {
		log.Error("func updateUser: Error work UseCase/Repository", "func", "updateUser", "err", err)
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
		log.Error("func changePassword: User ID is required", "func", "changePassword", "err", "empty userId from param")
		errors.HandleError(ctx, errors.BadRequest("User ID is required"))
		return
	}

	if err := ctx.ShouldBind(&inputData); err != nil {
		log.Error("func changePassword: Error in parse input param", "func", "changePassword", "err", err)
		errors.HandleError(ctx, errors.BadRequest("Invalid JSON"))
		return
	}

	if _, err := valid.CheckPassword(inputData.NewPassword); err != nil {
		log.Error("func changePassword: Invalid password", "func", "changePassword", "err", err)
		errors.HandleError(ctx, errors.BadRequest(err.Error()))
		return
	}

	err := h.userCase.ChangePassword(ctx, userID, inputData.OldPassword, inputData.NewPassword)
	if err != nil {
		log.Error("func changePassword: Error work UseCase/Repository", "func", "changePassword", "err", err)
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
		log.Error("func deactivateUser: User ID is required", "func", "deactivateUser", "err", "empty userId from param")
		errors.HandleError(ctx, errors.BadRequest("User ID is required"))
		return
	}

	err := h.userCase.DeactivateUser(ctx, userID)
	if err != nil {
		log.Error("func deactivateUser: Error work UseCase/Repository", "func", "deactivateUser", "err", err)
		errors.HandleError(ctx, err)
		return
	}

	response := ToResponseMessage("User deactivated successfully")
	ctx.JSON(http.StatusOK, response)
}

// UpdateUserRole
// @Summary      Update user role
// @Description  Updates user role (company_admin only)
// @Tags         users
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        id    path      string                true  "User ID (UUID)"
// @Param        role  body      RequestUpdateRoleDto  true  "Role update payload"
// @Success      200   {object}  ResponseMessageDto
// @Failure      400,404,500  {object}  errors.ErrorResponse
// @Failure      401  {object}  errors.ErrorResponse
// @Router       /users/{id}/role [put]
func (h *HandlerUser) UpdateUserRole(ctx *gin.Context) {
	log := logger.FromContext(ctx)
	userID := ctx.Param("id")
	var inputData RequestUpdateRoleDto

	if userID == "" {
		log.Error("func updateUserRole: User ID is required", "func", "updateUserRole", "err", "empty userId from param")
		errors.HandleError(ctx, errors.BadRequest("User ID is required"))
		return
	}

	if err := ctx.ShouldBind(&inputData); err != nil {
		log.Error("func updateUserRole: Error in parse input param", "func", "updateUserRole", "err", err)
		errors.HandleError(ctx, errors.BadRequest("Invalid JSON"))
		return
	}

	if inputData.RoleId == "" {
		log.Error("func updateUserRole: Role ID is required", "func", "updateUserRole", "err", "empty roleId from request")
		errors.HandleError(ctx, errors.BadRequest("Role ID is required"))
		return
	}

	err := h.userCase.UpdateUserRole(ctx, userID, inputData.RoleId)
	if err != nil {
		log.Error("func updateUserRole: Error work UseCase/Repository", "func", "updateUserRole", "err", err)
		errors.HandleError(ctx, err)
		return
	}

	response := ToResponseMessage("User role updated successfully")
	ctx.JSON(http.StatusOK, response)
}

// ActivateUser
// @Summary      Activate user
// @Description  Activates user account
// @Tags         users
// @Security     BearerAuth
// @Produce      json
// @Param        id   path      string  true  "User ID (UUID)"
// @Success      200  {object}  ResponseMessageDto
// @Failure      400,404,500  {object}  errors.ErrorResponse
// @Failure      401  {object}  errors.ErrorResponse
// @Router       /users/{id}/activate [put]
func (h *HandlerUser) ActivateUser(ctx *gin.Context) {
	log := logger.FromContext(ctx)
	userID := ctx.Param("id")

	if userID == "" {
		log.Error("func activateUser: User ID is required", "func", "activateUser", "err", "empty userId from param")
		errors.HandleError(ctx, errors.BadRequest("User ID is required"))
		return
	}

	err := h.userCase.ActivateUser(ctx, userID)
	if err != nil {
		log.Error("func activateUser: Error work UseCase/Repository", "func", "activateUser", "err", err)
		errors.HandleError(ctx, err)
		return
	}

	response := ToResponseMessage("User activated successfully")
	ctx.JSON(http.StatusOK, response)
}

// GetAllUsers
// @Summary      Get all users (super admin only)
// @Description  Returns all users from all companies (super admin only)
// @Tags         users
// @Security     BearerAuth
// @Produce      json
// @Success      200  {object}  ResponseUsersDto
// @Failure      400,500  {object}  errors.ErrorResponse
// @Failure      401  {object}  errors.ErrorResponse
// @Router       /users/all [get]
func (h *HandlerUser) GetAllUsers(ctx *gin.Context) {
	log := logger.FromContext(ctx)

	users, err := h.userCase.GetAllUsers(ctx)
	if err != nil {
		log.Error("func getAllUsers: Error work UseCase/Repository", "func", "getAllUsers", "err", err)
		errors.HandleError(ctx, err)
		return
	}

	response := ToResponseUsers(users)
	ctx.JSON(http.StatusOK, response)
}

// TransferUserToCompany
// @Summary      Transfer user to another company
// @Description  Transfers user to different company (super admin only)
// @Tags         users
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        id       path      string                       true  "User ID (UUID)"
// @Param        company  body      RequestTransferCompanyDto    true  "Company transfer payload"
// @Success      200      {object}  ResponseMessageDto
// @Failure      400,404,500  {object}  errors.ErrorResponse
// @Failure      401  {object}  errors.ErrorResponse
// @Router       /users/{id}/company [put]
func (h *HandlerUser) TransferUserToCompany(ctx *gin.Context) {
	log := logger.FromContext(ctx)
	userID := ctx.Param("id")
	var inputData RequestTransferCompanyDto

	if userID == "" {
		log.Error("func transferUserToCompany: User ID is required", "func", "transferUserToCompany", "err", "empty userId from param")
		errors.HandleError(ctx, errors.BadRequest("User ID is required"))
		return
	}

	if err := ctx.ShouldBind(&inputData); err != nil {
		log.Error("func transferUserToCompany: Error in parse input param", "func", "transferUserToCompany", "err", err)
		errors.HandleError(ctx, errors.BadRequest("Invalid JSON"))
		return
	}

	if inputData.CompanyId == "" {
		log.Error("func transferUserToCompany: Company ID is required", "func", "transferUserToCompany", "err", "empty companyId from request")
		errors.HandleError(ctx, errors.BadRequest("Company ID is required"))
		return
	}

	err := h.userCase.TransferUserToCompany(ctx, userID, inputData.CompanyId)
	if err != nil {
		log.Error("func transferUserToCompany: Error work UseCase/Repository", "func", "transferUserToCompany", "err", err)
		errors.HandleError(ctx, err)
		return
	}

	response := ToResponseMessage("User transferred to company successfully")
	ctx.JSON(http.StatusOK, response)
}
