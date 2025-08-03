-- name: UpdateSheetRowCount :exec
UPDATE sheets
SET row_count = CASE 
    WHEN row_count > ? THEN row_count 
    ELSE ? 
END,
    updated_at = datetime('now')
WHERE id = ?;
