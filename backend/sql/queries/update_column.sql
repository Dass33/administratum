-- name: UpdateColumn :exec
UPDATE columns
SET name = ?,
    type = ?,
    required = ?,
    updated_at = datetime('now')
WHERE id = ?;

-- name: UpdateColumnWithPermissionCheck :execrows
UPDATE columns
SET name = ?,
    type = ?,
    required = ?,
    updated_at = datetime('now')
WHERE columns.id = ? 
  AND columns.sheet_id IN (
    SELECT sheets.id FROM sheets
    JOIN branches ON sheets.branch_id = branches.id
    JOIN user_tables ON branches.table_id = user_tables.table_id
    WHERE user_tables.user_id = ? 
      AND user_tables.permission IN ('owner', 'contributor')
  );
