-- name: DeleteTable :exec
DELETE FROM tables 
WHERE id = ?; 