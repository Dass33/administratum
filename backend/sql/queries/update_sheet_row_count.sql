-- name: UpdateSheetRowCount :exec
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
