package rpChunkedUpload

const QueryCreateChunkedUpload = `
INSERT INTO chunked_uploads (
    id, file_name, total_size, chunk_size, total_chunks, 
    uploaded_chunks, uploaded_size, status, company_id, user_created,
    parent_path, target_path, mime_type,
    created_at, updated_at, expires_at
) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16)
`

const QueryGetChunkedUpload = `
SELECT id, file_name, total_size, chunk_size, total_chunks,
       uploaded_chunks, uploaded_size, status, company_id, user_created,
       parent_path, target_path, mime_type,
       created_at, updated_at, expires_at
FROM chunked_uploads 
WHERE id = $1 AND company_id = $2
`

const QueryUpdateChunkedUpload = `
UPDATE chunked_uploads 
SET uploaded_chunks = $2, uploaded_size = $3, status = $4, updated_at = $5
WHERE id = $1 AND company_id = $6
`

const QueryDeleteChunkedUpload = `
DELETE FROM chunked_uploads 
WHERE id = $1 AND company_id = $2
`

const QueryAddChunk = `
INSERT INTO upload_chunks (
    upload_id, chunk_index, size, etag, uploaded, uploaded_at, retries
) VALUES ($1, $2, $3, $4, $5, $6, $7)
ON CONFLICT (upload_id, chunk_index) 
DO UPDATE SET 
    size = EXCLUDED.size,
    etag = EXCLUDED.etag,
    uploaded = EXCLUDED.uploaded,
    uploaded_at = EXCLUDED.uploaded_at,
    retries = EXCLUDED.retries
`

const QueryGetUploadProgress = `
SELECT cu.id, cu.file_name, cu.total_size, cu.chunk_size, cu.total_chunks,
       cu.uploaded_chunks, cu.uploaded_size, cu.status, cu.company_id, cu.user_created,
       cu.parent_path, cu.target_path, cu.mime_type,
       cu.created_at, cu.updated_at, cu.expires_at,
       COALESCE(
           json_agg(
               json_build_object(
                   'index', uc.chunk_index,
                   'size', uc.size,
                   'etag', uc.etag,
                   'uploaded', uc.uploaded,
                   'uploaded_at', uc.uploaded_at,
                   'retries', uc.retries
               ) ORDER BY uc.chunk_index
           ) FILTER (WHERE uc.chunk_index IS NOT NULL),
           '[]'::json
       ) as chunks
FROM chunked_uploads cu
LEFT JOIN upload_chunks uc ON cu.id = uc.upload_id
WHERE cu.id = $1
GROUP BY cu.id
`

const QueryCleanupExpiredUploads = `
DELETE FROM chunked_uploads 
WHERE expires_at < NOW() 
   OR (status = 'failed' AND updated_at < NOW() - INTERVAL '1 day')
   OR (status = 'completed' AND updated_at < NOW() - INTERVAL '7 days')
`

const QueryGetExpiredUploads = `
SELECT id, company_id, file_name 
FROM chunked_uploads 
WHERE expires_at < NOW() 
   OR (status = 'failed' AND updated_at < NOW() - INTERVAL '1 day')
LIMIT 100
`
