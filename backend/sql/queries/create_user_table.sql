-- name: CreateUserTable :one
INSERT INTO user_tables (user_id, table_id, permission, created_at, updated_at)
VALUES (
    ?,
    ?,
    ?,
    datetime('now'),
    datetime('now')
)
RETURNING *; 