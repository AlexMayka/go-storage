package rpCompany

import (
	"context"
	"database/sql"
	"errors"
	"go-storage/internal/domain"
	pkgErrors "go-storage/pkg/errors"
	"time"
)

type RepositoryCompany struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *RepositoryCompany {
	return &RepositoryCompany{db: db}
}

func (r *RepositoryCompany) Create(ctx context.Context, c *domain.Company) (*domain.Company, error) {
	createDate := time.Now()

	_, err := r.db.ExecContext(ctx, QueryCreateCompany, c.ID, c.Name, c.Path, c.Description, createDate)
	if err != nil {
		return nil, pkgErrors.Database("unable to insert company")
	}

	return &domain.Company{
		ID:          c.ID,
		Name:        c.Name,
		Path:        c.Path,
		Description: c.Description,
		CreatedAt:   createDate,
	}, nil
}

func (r *RepositoryCompany) GetCompanyById(ctx context.Context, id string) (*domain.Company, error) {
	var company domain.Company
	row := r.db.QueryRowContext(ctx, QueryGetCompanyById, id)

	if err := row.Scan(&company.ID, &company.Name, &company.Path, &company.Description, &company.CreatedAt); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, pkgErrors.NotFound("company not found")
		}
		return nil, pkgErrors.NotFound("company not found")
	}

	return &company, nil
}

func (r *RepositoryCompany) GetAllCompanies(ctx context.Context) ([]*domain.Company, error) {
	var companies []*domain.Company
	rows, err := r.db.QueryContext(ctx, QueryGetCompanies)

	if err != nil {
		return nil, pkgErrors.Database("unable to query all companies")
	}
	defer func(rows *sql.Rows) {
		_ = rows.Close()
	}(rows)

	for rows.Next() {
		var company domain.Company
		if err := rows.Scan(&company.ID, &company.Name, &company.Path, &company.Description, &company.CreatedAt); err != nil {
			return nil, pkgErrors.Database("unable to query all companies")
		}
		companies = append(companies, &company)
	}

	if err := rows.Err(); err != nil {
		return nil, pkgErrors.Database("unable to query all companies")
	}

	return companies, nil
}

func (r *RepositoryCompany) DeleteCompany(ctx context.Context, id string) error {
	_, err := r.db.ExecContext(ctx, QueryDeleteCompanies, id)
	if err != nil {
		return pkgErrors.Database("unable to delete company")
	}
	return nil
}
