package hdUser

import "time"

type RegistrationUserDto struct {
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

type RegisterUserDto struct {
	ID         string `json:"id" binding:"required"`
	FirstName  string `json:"first_name" binding:"required"`
	SecondName string `json:"second_name" binding:"required"`
	LastName   string `json:"last_name" binding:"required"`
	Username   string `json:"username" binding:"required"`
	Email      string `json:"email" binding:"required"`
	Phone      string `json:"phone,omitempty"`
	CompanyId  string `json:"company_id" binding:"required"`
	RoleId     string `json:"role_id" binding:"required"`
}

type ResponseRegisterUserDto struct {
	Status string          `json:"status"`
	Time   time.Time       `json:"time"`
	Answer RegisterUserDto `json:"answer"`
}
