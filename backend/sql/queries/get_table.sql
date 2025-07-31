-- name: GetTable :one
select * from tables
where id = ?;
