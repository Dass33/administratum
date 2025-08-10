-- name: UpdateBranch :exec
UPDATE branches
SET name = ?,
    is_protected = ?,
    updated_at = datetime('now')
WHERE id = ?;
