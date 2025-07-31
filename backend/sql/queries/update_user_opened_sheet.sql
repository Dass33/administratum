-- name: UpdateUserOpenedSheet :one
UPDATE users
set opened_sheet = ?, updated_at = datetime('now')
where id = ?
RETURNING *; 