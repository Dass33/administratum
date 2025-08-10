-- name: DeleteBranch :exec
DELETE FROM branches 
WHERE id = ?; 
