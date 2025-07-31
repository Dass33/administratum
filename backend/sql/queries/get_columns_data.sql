-- name: GetColumnsData :many
select * from column_data cd
where column_id = ?;
