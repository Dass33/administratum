-- name: GetOldestBranchFromTable :one
select * from branches
where table_id = ?
order by created_at asc
limit 1;