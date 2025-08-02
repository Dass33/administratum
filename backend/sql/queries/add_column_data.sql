-- name: AddColumnData :one
INSERT INTO column_data (id, idx, value, column_id, created_at, updated_at)
VALUES (
    gen_random_uuid(),
    ?,
    ?,
    ?,
    datetime('now'),
    datetime('now')
)
RETURNING *;

-- name: UpdateSheetRowCountByColumn :exec
UPDATE sheets
SET row_count = CASE 
    WHEN row_count < ? THEN ?
    ELSE row_count
    END,
    updated_at = datetime('now')
WHERE id = (
    SELECT c.sheet_id 
    FROM columns c 
    WHERE c.id = ?
);
