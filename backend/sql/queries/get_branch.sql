-- name: GetBranch :one
select * from branches
where id = ?;
