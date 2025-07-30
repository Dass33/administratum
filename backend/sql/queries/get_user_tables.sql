-- name: GetUserTables :one
select * from user_tables
where user_id = ? and table_id = ?;
