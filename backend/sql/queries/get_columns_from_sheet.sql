-- name: GetColumnsFromSheet :many
SELECT c.*
FROM columns c
WHERE EXISTS (
    SELECT 1
    FROM sheets s
    WHERE s.id = c.sheet_id AND s.id = ?
);
