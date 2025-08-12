-- name: UpdateColumnData :exec
UPDATE column_data
SET value = ?,
    updated_at = datetime('now')
WHERE id = ?;

-- name: UpdateColumnDataWithPermissionCheck :execrows
UPDATE column_data
SET value = ?,
    updated_at = datetime('now')
WHERE column_data.id = ? 
  AND column_data.column_id IN (
    SELECT columns.id FROM columns
    JOIN sheets ON columns.sheet_id = sheets.id
    JOIN branches ON sheets.branch_id = branches.id
    JOIN user_tables ON branches.table_id = user_tables.table_id
    WHERE user_tables.user_id = ? 
      AND user_tables.permission IN ('owner', 'contributor')
  );
