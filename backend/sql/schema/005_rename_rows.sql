-- +goose up
ALTER TABLE rows RENAME TO column_data;

-- +goose down
ALTER TABLE column_data RENAME TO rows;
