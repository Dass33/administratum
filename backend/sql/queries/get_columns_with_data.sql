-- name: GetColumnsWithData :many
SELECT 
    c.id as column_id,
    c.name as column_name,
    c.type as column_type,
    c.required as column_required,
    cd.id as data_id,
    cd.idx as data_idx,
    cd.value as data_value
FROM columns c
LEFT JOIN column_data cd ON c.id = cd.column_id
WHERE EXISTS (
    SELECT 1
    FROM sheets s
    WHERE s.id = c.sheet_id AND s.id = ?
)
ORDER BY c.id, cd.idx; 
