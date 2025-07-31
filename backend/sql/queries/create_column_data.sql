-- name: CreateColumnData :one
INSERT INTO column_data (id, idx, value, column_id, created_at, updated_at)
VALUES (
    gen_random_uuid(),
    ?,
    ?,
    ?,
    datetime('now'),
    datetime('now')
)
RETURNING *; 