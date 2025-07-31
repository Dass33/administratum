-- name: GetTablesFromUser :many
select
    t.id id,
    game_url,
    t.created_at created_at,
    t.updated_at updated_at,
    t.name name
from user_tables ut
join tables t on t.id = ut.table_id and ut.user_id = ?;
