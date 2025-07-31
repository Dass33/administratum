-- name: GetBranchesFromTable :many
select * from branches
where table_id = ?;
