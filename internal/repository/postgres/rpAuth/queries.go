package rpAuth

const QueryGetRoleById = `
	SELECT id, "name"
	FROM roles
	WHERE id = $1
`

const QueryGetRoleByName = `
	SELECT id, "name"
	FROM roles
	WHERE name = $1
`

const QueryGetPermissionById = `
	SELECT id, "name"
	FROM permissions
	WHERE id = $1
`

const GetPermissionByIds = `
	SELECT id, name
	FROM permissions
	WHERE id = ANY($1)
`

const QueryGetPermissionsIdByRoleId = `
	SELECT role_id, permission_id
	FROM role_permissions
	WHERE role_id = $1
`
