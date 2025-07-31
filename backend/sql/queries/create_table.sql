-- name: CreateTable :one
INSERT INTO tables (id, game_url, created_at, updated_at)
VALUES (
    gen_random_uuid(),
    ?,
    datetime('now'),
    datetime('now')
)
RETURNING *; 