-- name: UpdateColumnData :exec
UPDATE column_data
SET value = ?,
    updated_at = datetime('now')
WHERE id = ?;
