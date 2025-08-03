-- name: DeleteColumn :exec
DELETE FROM columns
WHERE name = ? AND sheet_id = ?;
