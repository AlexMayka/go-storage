package rpFolder


const QueryGetFolderByName = `
	select id, name, path, parent_id, company_id, user_create, created_at, updated_at, is_active
	from folders
	where name = $1 and company_id = $2 and is_active = true;
`

const QueryGetFolderById = `
	select id, name, path, parent_id, company_id, user_create, created_at, updated_at, is_active
	from folders
	where id = $1 and company_id = $2 and is_active = true;
`

const QueryGetFolderByPath = `
	select id, name, path, parent_id, company_id, user_create, created_at, updated_at, is_active
	from folders
	where path = $1 and company_id = $2 and is_active = true;
`

const QueryGetFoldersByParentId = `
	select id, name, path, parent_id, company_id, user_create, created_at, updated_at, is_active
	from folders 
	where parent_id = $1 and company_id = $2 and is_active = true;
`

const QueryGetFoldersByCompanyId = `
	select id, name, path, parent_id, company_id, user_create, created_at, updated_at, is_active
	from folders
	where company_id = $1 and is_active = true;
`


const QueryUpdateFolder = `
	update folders 
	set 
		name = $1,
		path = $2,
		parent_id = $3,
		updated_at = $4,
		is_active = $5
	where
		id = $6 and is_active = true;
`

const QueryInsertFolder = `
	insert into folders (id, name, path, parent_id, company_id, user_create, created_at, updated_at, is_active)
	values ($1, $2, $3, $4, $5, $6, $7, $8, $9)
`
