-- name: GetSheetsFromBranch :many
select
    s.id id,
    s.name name,
    row_count,
    s.updated_at updated_at,
    s.created_at created_at
from branches b
join sheets s on b.id = s.branch_id and b.id = ?;
