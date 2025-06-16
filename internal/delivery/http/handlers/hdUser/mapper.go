package hdUser

import (
	"go-storage/internal/domain"
	"time"
)

func ToDomainCreate(dto RequestRegistrationUserDto) *domain.User {
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

func ToResponseCreate(userAuth *domain.UserWithAuth) ResponseRegisterUserDto {
	userDto := MapUserWithAuthToDto(userAuth)
	return ResponseRegisterUserDto{
		Status: "success",
		Time:   time.Now(),
		Answer: userDto,
	}
}

func MapUserWithAuthToDto(userAuth *domain.UserWithAuth) RegisterUserDto {
	return RegisterUserDto{
		ID:         userAuth.User.ID,
		FirstName:  userAuth.User.FirstName,
		SecondName: userAuth.User.SecondName,
		LastName:   userAuth.User.LastName,
		Username:   userAuth.User.Username,
		Email:      userAuth.User.Email,
		Phone:      userAuth.User.Phone,
		CompanyId:  userAuth.User.CompanyId,
		RoleId:     userAuth.User.RoleId,
		Auth: JwtAuth{
			Token:      userAuth.Token,
			ExpireTime: userAuth.ExpiresAt,
		},
	}
}

func MapUserToDto(user *domain.User) UserDto {
	return UserDto{
		ID:         user.ID,
		FirstName:  user.FirstName,
		SecondName: user.SecondName,
		LastName:   user.LastName,
		Username:   user.Username,
		Email:      user.Email,
		Phone:      user.Phone,
		CompanyId:  user.CompanyId,
		RoleId:     user.RoleId,
		IsActive:   user.IsActive,
		CreatedAt:  user.CreatedAt,
		UpdatedAt:  user.UpdatedAt,
	}
}

func MapUsersToDto(users []*domain.User) []UserDto {
	result := make([]UserDto, 0, len(users))
	for _, user := range users {
		result = append(result, MapUserToDto(user))
	}
	return result
}

func ToDomainUpdate(dto RequestUpdateUserDto) *domain.User {
	return &domain.User{
		FirstName:  dto.FirstName,
		SecondName: dto.SecondName,
		LastName:   dto.LastName,
		Username:   dto.Username,
		Email:      dto.Email,
		Phone:      dto.Phone,
	}
}

func ToResponseUser(user *domain.User) *ResponseUserDto {
	return &ResponseUserDto{
		Status: "success",
		Time:   time.Now(),
		Answer: MapUserToDto(user),
	}
}

func ToResponseUsers(users []*domain.User) *ResponseUsersDto {
	return &ResponseUsersDto{
		Status: "success",
		Time:   time.Now(),
		Answer: MapUsersToDto(users),
	}
}

func ToResponseMessage(message string) *ResponseMessageDto {
	return &ResponseMessageDto{
		Status:  "success",
		Time:    time.Now(),
		Message: message,
	}
}
