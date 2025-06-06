package hdCompany

import (
	"go-storage/internal/domain"
	"time"
)

func ToDomainCreate(dto RequestRegisterCompanyDto) *domain.Company {
	return &domain.Company{
		Name:        dto.Name,
		Description: dto.Desc,
	}
}

func ToDomainUpdate(dto RequestUpdateCompanyDto) *domain.Company {
	return &domain.Company{
		Name:        dto.Name,
		Description: dto.Description,
	}
}

func ToResponseCompany(company *domain.Company) *ResponseCompanyDto {
	return &ResponseCompanyDto{
		Status: "success",
		Time:   company.CreatedAt,
		Answer: &CompanyDto{Id: company.ID, Name: company.Name, Path: company.Path, Description: company.Description},
	}
}

func ToResponseCompanies(companies []*domain.Company) *ResponseCompaniesDto {
	var answerCompanies = make([]*CompanyDto, len(companies))

	for index, company := range companies {
		answerCompanies[index] = &CompanyDto{
			Id:          company.ID,
			Name:        company.Name,
			Path:        company.Path,
			Description: company.Description,
		}
	}

	return &ResponseCompaniesDto{
		Status: "success",
		Time:   time.Now(),
		Answer: answerCompanies,
	}
}

func ToResponseDelete() *ResponseDeleteCompanyDto {
	return &ResponseDeleteCompanyDto{
		Status: "success",
		Time:   time.Now(),
	}
}
