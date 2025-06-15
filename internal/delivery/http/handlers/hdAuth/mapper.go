package hdAuth

import (
	"go-storage/internal/domain"
	"time"
)

func ToResponseLogin(userAuth *domain.UserWithAuth) *ResponseLoginDto {
	return &ResponseLoginDto{
		Status: "success",
		Time:   time.Now(),
		Answer: MapUserWithAuthToLoginDto(userAuth),
	}
}

func ToResponseRefreshToken(userAuth *domain.UserWithAuth) *ResponseRefreshTokenDto {
	return &ResponseRefreshTokenDto{
		Status: "success",
		Time:   time.Now(),
		Answer: MapUserWithAuthToLoginDto(userAuth),
	}
}
