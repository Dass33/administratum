-- name: CreateBranch :one
INSERT INTO branches (id, name, table_id, created_at, updated_at)
VALUES (
    gen_random_uuid(),
    ?,
    ?,
    datetime('now'),
    datetime('now')
)
RETURNING *; 