-- name: UpdateUserTable :exec
UPDATE user_tables
set permission = ?, updated_at = datetime('now')
where user_id = ? and table_id = ?
RETURNING *; 
