-- name: RevokeRefreshToken :exec
update refresh_tokens
set revoked_at = datetime('now'), updated_at = datetime('now')

where token = ?;
