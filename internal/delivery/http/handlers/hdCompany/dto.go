package hdCompany

import "time"

type RequestRegisterCompanyDto struct {
	Name string `json:"name" binding:"required"`
	Desc string `json:"description" binding:"required"`
}

type RequestUpdateCompanyDto struct {
	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
}

type CompanyDto struct {
	Id          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Path        string `json:"path"`
}

type ResponseCompanyDto struct {
	Status string      `json:"status"`
	Time   time.Time   `json:"time"`
	Answer *CompanyDto `json:"answer"`
}

type ResponseCompaniesDto struct {
	Status string        `json:"status"`
	Time   time.Time     `json:"time"`
	Answer []*CompanyDto `json:"answer"`
}

type ResponseDeleteCompanyDto struct {
	Status string    `json:"status"`
	Time   time.Time `json:"time"`
}
