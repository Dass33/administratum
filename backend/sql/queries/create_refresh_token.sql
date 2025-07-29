-- name: CreateRefreshToken :one
insert into refresh_tokens (token, created_at, updated_at, user_id, expires_at)
VALUES (
    ?,
    NOW(),
    NOW(),
    ?,
    ?
)
returning *;
