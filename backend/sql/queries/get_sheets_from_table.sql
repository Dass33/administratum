-- name: GetSheetsFromTable :many
select s.*
from tables t
join branches b on t.id = b.table_id
join sheets s on b.id = s.branch_id
where t.id = ?
order by s.updated_at desc; 
