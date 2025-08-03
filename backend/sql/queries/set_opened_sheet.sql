-- name: SetOpenedSheet :exec
UPDATE users
SET opened_sheet = ?
WHERE id = ?;
