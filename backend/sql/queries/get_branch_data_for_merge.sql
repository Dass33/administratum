-- name: GetBranchDataForMerge :many
SELECT 
    s.id as sheet_id,
    s.name as sheet_name,
    s.type as sheet_type,
    s.created_at as sheet_created_at,
    s.updated_at as sheet_updated_at,
    s.source_sheet_id,
    c.id as column_id,
    c.name as column_name,
    c.type as column_type,
    c.required as column_required,
    c.created_at as column_created_at,
    c.updated_at as column_updated_at,
    c.source_column_id,
    cd.id as column_data_id,
    cd.idx as column_data_idx,
    cd.value as column_data_value,
    cd.created_at as column_data_created_at,
    cd.updated_at as column_data_updated_at
FROM sheets s
LEFT JOIN columns c ON c.sheet_id = s.id
LEFT JOIN column_data cd ON cd.column_id = c.id
WHERE s.branch_id = ?
ORDER BY s.id, c.id, cd.idx;