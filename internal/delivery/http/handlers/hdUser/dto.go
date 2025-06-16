package hdUser

import "time"

type ResponseRegisterUserDto struct {
	Status string          `json:"status"`
	Time   time.Time       `json:"time"`
	Answer RegisterUserDto `json:"answer"`
}

type RequestRegistrationUserDto struct {
	FirstName  string `json:"first_name" binding:"required"`
	SecondName string `json:"second_name,omitempty"`
	LastName   string `json:"last_name" binding:"required"`
	Username   string `json:"username" binding:"required"`
	Email      string `json:"email" binding:"required"`
	Phone      string `json:"phone,omitempty"`
	Password   string `json:"password" binding:"required"`
	CompanyId  string `json:"company_id" binding:"required"`
	RoleId     string `json:"-"`
	RoleName   string `json:"role_name" binding:"required"`
}

type JwtAuth struct {
	Token      string    `json:"token"`
	ExpireTime time.Time `json:"expire_time"`
}

type RegisterUserDto struct {
	ID         string  `json:"id"`
	FirstName  string  `json:"first_name"`
	SecondName string  `json:"second_name,omitempty"`
	LastName   string  `json:"last_name"`
	Username   string  `json:"username"`
	Email      string  `json:"email"`
	Phone      string  `json:"phone,omitempty"`
	CompanyId  string  `json:"company_id"`
	RoleId     string  `json:"role_id"`
	Auth       JwtAuth `json:"auth"`
}

type UserDto struct {
	ID         string    `json:"id"`
	FirstName  string    `json:"first_name"`
	SecondName string    `json:"second_name,omitempty"`
	LastName   string    `json:"last_name"`
	Username   string    `json:"username"`
	Email      string    `json:"email"`
	Phone      string    `json:"phone,omitempty"`
	CompanyId  string    `json:"company_id"`
	RoleId     string    `json:"role_id"`
	IsActive   bool      `json:"is_active"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

type ResponseUserDto struct {
	Status string    `json:"status"`
	Time   time.Time `json:"time"`
	Answer UserDto   `json:"answer"`
}

type ResponseUsersDto struct {
	Status string    `json:"status"`
	Time   time.Time `json:"time"`
	Answer []UserDto `json:"answer"`
}

type RequestUpdateUserDto struct {
	FirstName  string `json:"first_name,omitempty"`
	SecondName string `json:"second_name,omitempty"`
	LastName   string `json:"last_name,omitempty"`
	Username   string `json:"username,omitempty"`
	Email      string `json:"email,omitempty"`
	Phone      string `json:"phone,omitempty"`
}

type RequestChangePasswordDto struct {
	OldPassword string `json:"old_password" binding:"required"`
	NewPassword string `json:"new_password" binding:"required"`
}

type RequestRefreshTokenDto struct {
	Token string `json:"token" binding:"required"`
}

type ResponseMessageDto struct {
	Status  string    `json:"status"`
	Time    time.Time `json:"time"`
	Message string    `json:"message"`
}

type RequestUpdateRoleDto struct {
	RoleId string `json:"role_id" binding:"required"`
}

type RequestTransferCompanyDto struct {
	CompanyId string `json:"company_id" binding:"required"`
}
