-- name: RenameTable :exec
UPDATE tables 
SET name = ?, updated_at = CURRENT_TIMESTAMP
WHERE id = ?; 