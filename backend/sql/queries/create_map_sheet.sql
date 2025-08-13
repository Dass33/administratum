-- name: CreateMapSheet :one
INSERT INTO sheets (id, name, branch_id, type, created_at, updated_at)
VALUES (
    gen_random_uuid(),
    ?,
    ?,
    'map',
    datetime('now'),
    datetime('now')
)
RETURNING id;

-- name: CreateMapSheetColumns :exec
INSERT INTO columns (id, name, type, required, sheet_id, created_at, updated_at, source_column_id, order_index)
VALUES 
    (gen_random_uuid(), 'name', 'text', true, ?1, datetime('now'), datetime('now'), NULL, 0),
    (gen_random_uuid(), 'value', 'any', true, ?1, datetime('now'), datetime('now'), NULL, 1),
    (gen_random_uuid(), 'comment', 'text', false, ?1, datetime('now'), datetime('now'), NULL, 2);
