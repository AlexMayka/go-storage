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
