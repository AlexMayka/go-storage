package rpCompany

const QueryCreateCompany = `
	INSERT INTO companies (id, "name", storage_path, description, created_at, update_at, is_active)	
	VALUES ($1, $2, $3, $4, $5, $6, $7)
`

const QueryGetCompanyById = `
	SELECT id, "name", storage_path, description, created_at, update_at, is_active FROM companies WHERE id = $1 and is_active = true
`

const QueryGetCompanies = `
	SELECT id, "name", storage_path, description, created_at, update_at, is_active FROM companies WHERE is_active = true
`

const QueryChangeIsActive = `
	UPDATE companies
	SET is_active = $1
	WHERE id = $2
`

const QueryUpdateCompany = `
	UPDATE companies
	SET id = $1, name = $2, storage_path = $3, description = $4, created_at = $5, update_at = $6, is_active = $7
	WHERE id = $1 AND is_active = true
`

const QueryDeleteCompanies = `
	DELETE FROM companies
	WHERE id=$1;
`
