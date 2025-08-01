-- name: UpdateColumnData :exec
UPDATE column_data
SET value = ?
WHERE id = ?;
