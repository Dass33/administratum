-- name: ChangeGameUrl :exec
UPDATE tables 
SET game_url = ?, updated_at = CURRENT_TIMESTAMP
WHERE id = ?; 
