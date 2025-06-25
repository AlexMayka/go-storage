package rpUser

const QueryCreateUser = `
	insert into users (
	   id, 
	   first_name, second_name, last_name, username, 
	   email, phone, "password", 
	   company_id, role_id, 
	   last_login, created_at, updated_at, is_active
	)
	values ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14);
`

const QueryGetUserByID = `
	SELECT id, first_name, second_name, last_name, username, email, phone, "password", company_id, role_id, last_login, created_at, updated_at, is_active
	FROM users 
	WHERE id = $1 AND is_active = true
`

const QueryGetUserByEmail = `
	SELECT id, first_name, second_name, last_name, username, email, phone, "password", company_id, role_id, last_login, created_at, updated_at, is_active
	FROM users 
	WHERE email = $1 AND is_active = true
`

const QueryGetUserByUsername = `
	SELECT id, first_name, second_name, last_name, username, email, phone, "password", company_id, role_id, last_login, created_at, updated_at, is_active
	FROM users 
	WHERE username = $1 AND is_active = true
`

const QueryGetUsersByCompanyID = `
	SELECT id, first_name, second_name, last_name, username, email, phone, "password", company_id, role_id, last_login, created_at, updated_at, is_active
	FROM users 
	WHERE company_id = $1 AND is_active = true
`

const QueryGetUserByIDWithCompany = `
	SELECT id, first_name, second_name, last_name, username, email, phone, "password", company_id, role_id, last_login, created_at, updated_at, is_active
	FROM users 
	WHERE id = $1 AND company_id = $2 AND is_active = true
`

const QueryUpdateUser = `
	UPDATE users 
	SET first_name = $2, second_name = $3, last_name = $4, username = $5, email = $6, phone = $7, updated_at = $8
	WHERE id = $1 AND is_active = true
`

const QueryUpdateUserWithCompany = `
	UPDATE users 
	SET first_name = $2, second_name = $3, last_name = $4, username = $5, email = $6, phone = $7, updated_at = $8
	WHERE id = $1 AND company_id = $9 AND is_active = true
`

const QueryUpdatePassword = `
	UPDATE users 
	SET "password" = $2, updated_at = $3
	WHERE id = $1 AND is_active = true
`

const QueryUpdatePasswordWithCompany = `
	UPDATE users 
	SET "password" = $2, updated_at = $3
	WHERE id = $1 AND company_id = $4 AND is_active = true
`

const QueryUpdateIsActive = `
	UPDATE users 
	SET is_active = $2, updated_at = $3
	WHERE id = $1
`

const QueryUpdateIsActiveWithCompany = `
	UPDATE users 
	SET is_active = $2, updated_at = $3
	WHERE id = $1 AND company_id = $4
`

const QueryUpdateLastLogin = `
	UPDATE users 
	SET last_login = $2, updated_at = $3
	WHERE id = $1 AND is_active = true
`

const QueryUpdateUserRole = `
	UPDATE users 
	SET role_id = $2, updated_at = $3
	WHERE id = $1 AND is_active = true
`

const QueryGetAllUsers = `
	SELECT id, first_name, second_name, last_name, username, email, phone, "password", company_id, role_id, last_login, created_at, updated_at, is_active
	FROM users 
	WHERE is_active = true
	ORDER BY created_at DESC
`

const QueryUpdateUserCompany = `
	UPDATE users 
	SET company_id = $2, updated_at = $3
	WHERE id = $1 AND is_active = true
`
