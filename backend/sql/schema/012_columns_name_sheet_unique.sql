-- +goose up
CREATE UNIQUE INDEX columns_name_sheet_id_unique ON columns (name, sheet_id);

-- +goose down
DROP INDEX columns_name_sheet_id_unique; 