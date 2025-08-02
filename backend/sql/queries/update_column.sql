-- name: UpdateColumn :exec
UPDATE columns
SET name = ?,
    type = ?,
    required = ?,
    updated_at = datetime('now')
WHERE id = ?;
