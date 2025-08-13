-- name: GetColumnsWithDataBySheet :many
SELECT 
    c.id as column_id,
    c.name as column_name,
    c.type as column_type,
    c.required as column_required,
    c.order_index as column_order_index,
    cd.id as data_id,
    cd.idx as data_idx,
    cd.value as data_value,
    cd.type as data_type
FROM columns c
LEFT JOIN column_data cd ON c.id = cd.column_id
WHERE c.sheet_id = ?
ORDER BY c.order_index, cd.idx;