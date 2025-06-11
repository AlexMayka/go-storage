package hdUser

import (
	"go-storage/internal/domain"
	"time"
)

func ToDomainCreate(dto RegistrationUserDto) *domain.User {
	return &domain.User{
		FirstName:  dto.FirstName,
		SecondName: dto.SecondName,
		LastName:   dto.LastName,
		Username:   dto.Username,
		Email:      dto.Email,
		Phone:      dto.Phone,
		Password:   dto.Password,
		CompanyId:  dto.CompanyId,
		RoleId:     dto.RoleId,
	}
}

func ToResponseCreate(user *domain.User) *ResponseRegisterUserDto {
	return &ResponseRegisterUserDto{
		Status: "success",
		Time:   time.Now(),
		Answer: RegisterUserDto{
			ID:         user.ID,
			FirstName:  user.FirstName,
			SecondName: user.SecondName,
			LastName:   user.LastName,
			Username:   user.Username,
			Email:      user.Email,
			Phone:      user.Phone,
			CompanyId:  user.CompanyId,
			RoleId:     user.RoleId,
		},
	}
}
