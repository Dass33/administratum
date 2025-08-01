-- name: UpdateColumn :exec
UPDATE columns
SET name = ?,
    type = ?,
    required = ?
WHERE id = ?;
