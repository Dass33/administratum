-- name: SwapColumns :exec
UPDATE columns 
SET order_index = CASE 
    WHEN columns.id = ?1 THEN ?3
    WHEN columns.id = ?2 THEN ?4
    ELSE columns.order_index
END,
updated_at = datetime('now')
WHERE columns.id IN (?1, ?2);

-- name: SwapColumnsWithPermissionCheck :execrows
UPDATE columns 
SET order_index = CASE 
    WHEN columns.id = ?1 THEN ?3
    WHEN columns.id = ?2 THEN ?4
    ELSE columns.order_index
END,
updated_at = datetime('now')
WHERE columns.id IN (?1, ?2)
  AND columns.sheet_id IN (
    SELECT sheets.id FROM sheets
    JOIN branches ON sheets.branch_id = branches.id
    JOIN user_tables ON branches.table_id = user_tables.table_id
    WHERE user_tables.user_id = ?5 
      AND user_tables.permission IN ('owner', 'contributor')
  );
