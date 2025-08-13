-- name: AddColumn :one
INSERT INTO columns (id, name, type, required, sheet_id, created_at, updated_at, source_column_id, order_index)
VALUES (
    gen_random_uuid(),
    ?1,
    ?2,
    ?3,
    ?4,
    datetime('now'),
    datetime('now'),
    ?5,
    (select COALESCE(MAX(order_index + 1), 0) from columns where sheet_id = ?4)
)
RETURNING *;
