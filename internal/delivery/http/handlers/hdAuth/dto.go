package hdAuth

import (
	"go-storage/internal/domain"
	"time"
)

type JwtAuth struct {
	Token      string    `json:"token"`
	ExpireTime time.Time `json:"expire_time"`
}

type RequestLoginDto struct {
	Login    string `json:"login" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type ResponseLoginDto struct {
	Status string           `json:"status"`
	Time   time.Time        `json:"time"`
	Answer LoginResponseDto `json:"answer"`
}

type LoginResponseDto struct {
	ID        string  `json:"id"`
	FirstName string  `json:"first_name"`
	LastName  string  `json:"last_name"`
	Username  string  `json:"username"`
	Email     string  `json:"email"`
	CompanyId string  `json:"company_id"`
	RoleId    string  `json:"role_id"`
	Auth      JwtAuth `json:"auth"`
}

func MapUserWithAuthToLoginDto(userAuth *domain.UserWithAuth) LoginResponseDto {
	return LoginResponseDto{
		ID:        userAuth.User.ID,
		FirstName: userAuth.User.FirstName,
		LastName:  userAuth.User.LastName,
		Username:  userAuth.User.Username,
		Email:     userAuth.User.Email,
		CompanyId: userAuth.User.CompanyId,
		RoleId:    userAuth.User.RoleId,
		Auth: JwtAuth{
			Token:      userAuth.Token,
			ExpireTime: userAuth.ExpiresAt,
		},
	}
}

type RequestRefreshTokenDto struct {
	Token string `json:"token" binding:"required"`
}

type ResponseRefreshTokenDto struct {
	Status string           `json:"status"`
	Time   time.Time        `json:"time"`
	Answer LoginResponseDto `json:"answer"`
}
