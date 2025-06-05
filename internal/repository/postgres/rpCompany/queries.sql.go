package rpCompany

const QueryCreateCompany = `
	INSERT INTO companies (id, "name", storage_path, description, created_at)	
	VALUES ($1, $2, $3, $4, $5)
`

const QueryGetCompanyById = `
	SELECT id, "name", storage_path, description, created_at FROM companies WHERE id = $1
`

const QueryGetCompanies = `
	SELECT id, "name", storage_path, description, created_at FROM companies
`

const QueryDeleteCompanies = `
	DELETE FROM companies
	WHERE id=$1::uuid;
`

const QueryUpdateCompanyName = `

`
