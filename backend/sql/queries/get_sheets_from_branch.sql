-- name: GetSheetsFromBranch :many
select
    s.*
from branches b
join sheets s on b.id = s.branch_id and b.id = ?;
