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
	const isActive = true
	createDate := time.Now()

	_, err := r.db.ExecContext(ctx, QueryCreateCompany, c.ID, c.Name, c.Path, c.Description, createDate, createDate, isActive)
	if err != nil {
		return nil, pkgErrors.Database("unable to insert company")
	}

	return &domain.Company{
		ID:          c.ID,
		Name:        c.Name,
		Path:        c.Path,
		Description: c.Description,
		CreatedAt:   createDate,
		UpdatedAt:   createDate,
		IsActive:    isActive,
	}, nil
}

func (r *RepositoryCompany) GetCompanyById(ctx context.Context, id string) (*domain.Company, error) {
	var company domain.Company
	row := r.db.QueryRowContext(ctx, QueryGetCompanyById, id)

	if err := row.Scan(&company.ID, &company.Name, &company.Path, &company.Description, &company.CreatedAt, &company.UpdatedAt, &company.IsActive); err != nil {
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
		if err := rows.Scan(&company.ID, &company.Name, &company.Path, &company.Description, &company.CreatedAt, &company.UpdatedAt, &company.IsActive); err != nil {
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

func (r *RepositoryCompany) UpdateIsActive(ctx context.Context, id string, on bool) error {
	query := QueryChangeIsActive
	_, err := r.db.ExecContext(ctx, query, on, id)
	if err != nil {
		return pkgErrors.Database("unable to delete company")
	}
	return nil
}

func (r *RepositoryCompany) Update(ctx context.Context, c *domain.Company) (*domain.Company, error) {
	updateDate := time.Now()
	_, err := r.db.ExecContext(ctx, QueryUpdateCompany, c.ID, c.Name, c.Path, c.Description, c.CreatedAt, updateDate, c.IsActive)
	if err != nil {
		return nil, pkgErrors.Database("unable to insert company")
	}

	return &domain.Company{
		ID:          c.ID,
		Name:        c.Name,
		Path:        c.Path,
		Description: c.Description,
		CreatedAt:   c.CreatedAt,
		UpdatedAt:   updateDate,
		IsActive:    c.IsActive,
	}, nil
}
