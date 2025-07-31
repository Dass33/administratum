-- +goose up
ALTER TABLE column_data
RENAME COLUMN index_num TO idx;

-- +goose down
ALTER TABLE column_data
RENAME COLUMN idx TO index_num;
