package domain

import "time"

type User struct {
	ID         string
	FirstName  string
	SecondName string
	LastName   string
	Username   string
	Email      string
	Phone      string
	Password   string
	CompanyId  string
	RoleId     string
	LastLogin  time.Time
	CreatedAt  time.Time
	UpdatedAt  time.Time
	IsActive   bool
}
