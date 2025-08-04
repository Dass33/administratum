-- name: DeleteSheet :exec
DELETE FROM sheets 
WHERE id = ?; 