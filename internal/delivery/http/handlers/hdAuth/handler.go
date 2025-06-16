package hdAuth

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

type HandlerAuth struct {
	userCase UseCaseUserInterface
	authCase UseCaseAuthInterface
}

func NewHandlerAuth(uc UseCaseUserInterface, ac UseCaseAuthInterface) *HandlerAuth {
	return &HandlerAuth{
		userCase: uc,
		authCase: ac,
	}
}

// Login
// @Summary      User login
// @Description  Authenticates user and returns JWT token
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        credentials  body     RequestLoginDto  true  "Login credentials"
// @Success      200         {object}  ResponseLoginDto
// @Failure      400,401,500 {object}  errors.ErrorResponse
// @Router       /auth/login [post]
func (h *HandlerAuth) Login(ctx *gin.Context) {
	log := logger.FromContext(ctx)
	var inputData RequestLoginDto

	if err := ctx.ShouldBind(&inputData); err != nil {
		log.Error("func login: Error in parse input param", "func", "login", "err", err)
		errors.HandleError(ctx, errors.BadRequest("Invalid JSON"))
		return
	}

	if inputData.Login == "" {
		log.Error("func login: Login is required", "func", "login", "err", "empty login from request")
		errors.HandleError(ctx, errors.BadRequest("Login is required"))
		return
	}

	if inputData.Password == "" {
		log.Error("func login: Password is required", "func", "login", "err", "empty password from request")
		errors.HandleError(ctx, errors.BadRequest("Password is required"))
		return
	}

	if _, err := valid.CheckPassword(inputData.Password); err != nil {
		log.Error("func login: Invalid password", "func", "login", "err", err)
		errors.HandleError(ctx, errors.BadRequest(err.Error()))
		return
	}

	user, err := h.userCase.Login(ctx, inputData.Login, inputData.Password)
	if err != nil {
		log.Error("func login: Error work UseCase/Repository", "func", "login", "err", err)
		errors.HandleError(ctx, err)
		return
	}

	userAuth, err := h.generateUserAuth(ctx, user)
	if err != nil {
		log.Error("func login: Error generate token", "func", "login", "err", err)
		errors.HandleError(ctx, err)
		return
	}

	response := ToResponseLogin(userAuth)

	ctx.JSON(http.StatusOK, *response)
}

// RefreshToken
// @Summary      Refresh JWT token
// @Description  Refreshes JWT authentication token
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        token  body      RequestRefreshTokenDto  true  "Refresh token payload"
// @Success      200    {object}  ResponseRefreshTokenDto
// @Failure      400,401,500  {object}  errors.ErrorResponse
// @Router       /auth/refresh-token [post]
func (h *HandlerAuth) RefreshToken(ctx *gin.Context) {
	log := logger.FromContext(ctx)
	var inputData RequestRefreshTokenDto

	if err := ctx.ShouldBind(&inputData); err != nil {
		log.Error("func refreshToken: Error in parse input param", "func", "refreshToken", "err", err)
		errors.HandleError(ctx, errors.BadRequest("Invalid JSON"))
		return
	}

	cfg := config.FromContext(ctx.Request.Context())
	if cfg == nil {
		log.Error("func refreshToken: Config not found in context", "func", "refreshToken", "err", "missing config")
		errors.HandleError(ctx, errors.InternalServer("config not found in context"))
		return
	}

	claims, err := jwt.ParseToken(inputData.Token, []byte(cfg.App.JwtSecret))
	if err != nil {
		log.Error("func refreshToken: Invalid token", "func", "refreshToken", "err", err)
		errors.HandleError(ctx, errors.Unauthorized("invalid token"))
		return
	}

	user, err := h.userCase.RefreshToken(ctx, claims.UserID)
	if err != nil {
		log.Error("func refreshToken: Error work UseCase/Repository", "func", "refreshToken", "err", err)
		errors.HandleError(ctx, err)
		return
	}

	userAuth, err := h.generateUserAuth(ctx, user)
	if err != nil {
		log.Error("func refreshToken: Error generate token", "func", "refreshToken", "err", err)
		errors.HandleError(ctx, err)
		return
	}

	response := ToResponseRefreshToken(userAuth)

	ctx.JSON(http.StatusOK, response)
}

func (h *HandlerAuth) generateUserAuth(ctx *gin.Context, user *domain.User) (*domain.UserWithAuth, error) {
	cfg := config.FromContext(ctx.Request.Context())
	if cfg == nil {
		return nil, errors.InternalServer("config not found in context")
	}

	tokenResp, err := jwt.CreateTokenWithExpiry(user.ID, user.RoleId, user.CompanyId, []byte(cfg.App.JwtSecret))
	if err != nil {
		return nil, errors.InternalServer("failed to create jwt token")
	}

	return &domain.UserWithAuth{
		User:      user,
		Token:     tokenResp.Token,
		ExpiresAt: tokenResp.ExpiresAt,
	}, nil
}
