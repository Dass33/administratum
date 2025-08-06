-- name: CreateSheet :one
INSERT INTO sheets (id, name, branch_id, created_at, updated_at)
VALUES (
    gen_random_uuid(),
    ?,
    ?,
    datetime('now'),
    datetime('now')
)
RETURNING *; 
