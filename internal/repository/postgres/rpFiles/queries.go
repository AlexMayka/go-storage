package rpFiles

const QueryCreateFile = `
INSERT INTO files (
    id, name, type, full_path, parent_id, company_id, user_created,
    mime_type, size, hash, storage_path,
    created_at, updated_at, is_active
) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14)
`

const QueryGetFile = `
SELECT id, name, type, full_path, parent_id, company_id, user_created,
       mime_type, size, hash, storage_path,
       created_at, updated_at, is_active
FROM files 
WHERE id = $1 AND company_id = $2 AND is_active = true
`

const QueryGetFileByPath = `
SELECT id, name, type, full_path, parent_id, company_id, user_created,
       mime_type, size, hash, storage_path,
       created_at, updated_at, is_active
FROM files 
WHERE full_path = $1 AND company_id = $2 AND is_active = true
`

const QueryGetFolderContents = `
SELECT id, name, type, full_path, parent_id, company_id, user_created,
       mime_type, size, hash, storage_path,
       created_at, updated_at, is_active
FROM files 
WHERE parent_id = $1 AND company_id = $2 AND is_active = true
ORDER BY type DESC, name ASC
`

const QueryGetFolderContentsByPath = `
SELECT id, name, type, full_path, parent_id, company_id, user_created,
       mime_type, size, hash, storage_path,
       created_at, updated_at, is_active
FROM files 
WHERE full_path LIKE $1 AND company_id = $2 AND is_active = true
  AND full_path != $3
ORDER BY type DESC, name ASC
`

const QueryGetFolderContentsByType = `
SELECT id, name, type, full_path, parent_id, company_id, user_created,
       mime_type, size, hash, storage_path,
       created_at, updated_at, is_active
FROM files 
WHERE parent_id = $1 AND company_id = $2 AND type = $3 AND is_active = true
ORDER BY name ASC
`

const QueryUpdateFile = `
UPDATE files 
SET name = $2, full_path = $3, parent_id = $4, mime_type = $5, 
    size = $6, hash = $7, storage_path = $8, updated_at = $9
WHERE id = $1 AND company_id = $10 AND is_active = true
`

const QueryUpdateFileName = `
UPDATE files 
SET name = $2, full_path = $3, updated_at = $4
WHERE id = $1 AND company_id = $5 AND is_active = true
`

const QueryUpdateFileParent = `
UPDATE files 
SET parent_id = $2, full_path = $3, updated_at = $4
WHERE id = $1 AND company_id = $5 AND is_active = true
`

const QueryDeleteFile = `
UPDATE files 
SET is_active = false, updated_at = $3
WHERE id = $1 AND company_id = $2 AND is_active = true
`

const QueryDeleteFolder = `
UPDATE files 
SET is_active = false, updated_at = $3
WHERE full_path = $1 AND company_id = $2 AND is_active = true
`

const QueryMoveFolderAndContents = `
UPDATE files 
SET full_path = REPLACE(full_path, $1, $2), updated_at = $5
WHERE full_path LIKE $3 AND company_id = $4 AND is_active = true
`

const QueryCreateFolder = `
INSERT INTO files (
    id, name, type, full_path, parent_id, company_id, user_created,
    created_at, updated_at, is_active
) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
`

const QueryGetFolder = `
SELECT id, name, type, full_path, parent_id, company_id, user_created,
       mime_type, size, hash, storage_path,
       created_at, updated_at, is_active
FROM files 
WHERE full_path = $1 AND company_id = $2 AND type = 'folder' AND is_active = true
`
