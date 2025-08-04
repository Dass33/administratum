-- name: RenameSheet :exec
UPDATE sheets 
SET name = ?, updated_at = CURRENT_TIMESTAMP
WHERE id = ?; 