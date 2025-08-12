-- name: DeleteBranch :exec
DELETE FROM branches 
WHERE id = ?;

-- name: DeleteBranchWithPermissionCheck :execrows
DELETE FROM branches 
WHERE id = ? 
  AND table_id IN (
    SELECT table_id FROM user_tables 
    WHERE user_id = ? 
      AND permission IN ('owner', 'contributor')
  ); 
