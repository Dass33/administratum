-- name: DeleteTable :exec
DELETE FROM tables 
WHERE id = ?;

-- name: DeleteTableWithPermissionCheck :execrows
DELETE FROM tables 
WHERE id = ? 
  AND id IN (
    SELECT table_id FROM user_tables 
    WHERE user_id = ? 
      AND permission IN ('owner', 'contributor')
  ); 