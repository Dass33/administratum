-- name: AddColumn :one
INSERT INTO columns (id, name, type, required, sheet_id, created_at, updated_at, source_column_id)
VALUES (
    gen_random_uuid(),
    ?,
    ?,
    ?,
    ?,
    datetime('now'),
    datetime('now'),
    ?
)
RETURNING *;
