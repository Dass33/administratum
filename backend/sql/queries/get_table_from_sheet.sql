-- name: GetTableFromSheet :one
select t.* from sheets s
join branches b on s.branch_id = b.id
join tables t on t.id = b.table_id
where s.id = ?;
