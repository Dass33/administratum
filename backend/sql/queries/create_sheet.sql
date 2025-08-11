-- name: CreateSheet :one
INSERT INTO sheets (id, name, type, branch_id, created_at, updated_at, source_sheet_id)
VALUES (
    gen_random_uuid(),
    ?,
    ?,
    ?,
    datetime('now'),
    datetime('now'),
    ?
)
RETURNING *; 
