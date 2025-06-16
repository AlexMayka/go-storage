package domain

import "time"

type Folder struct {
	ID           string
	Name         string
	Path         Path
	ParentId     string
	CompanyId    string
	UserCreateId string
	CreatedAt    time.Time
	UpdatedAt    time.Time
	IsActive     bool
}
