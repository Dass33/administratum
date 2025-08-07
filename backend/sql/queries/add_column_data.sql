-- name: AddColumnData :one
INSERT INTO column_data (id, idx, value, type, column_id, created_at, updated_at)
VALUES (
    gen_random_uuid(),
    ?,
    ?,
    ?,
    (SELECT id FROM columns WHERE name = ? AND sheet_id = ?),
    datetime('now'),
    datetime('now')
)
RETURNING *;
