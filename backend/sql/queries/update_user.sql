-- name: UpdateUser :one
UPDATE users
set email = ?, hashed_password = ?, updated_at = datetime('now')
where id = ?
RETURNING *;
