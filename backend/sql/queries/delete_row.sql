-- name: DeleteRow :exec
DELETE FROM column_data 
WHERE column_id IN (
    SELECT c.id 
    FROM columns c 
    WHERE c.sheet_id = ?1
) 
AND idx = ?2;

UPDATE column_data 
SET idx = idx - 1 
WHERE column_id IN (
    SELECT c.id 
    FROM columns c 
    WHERE c.sheet_id = ?1
) 
AND idx > ?2;
