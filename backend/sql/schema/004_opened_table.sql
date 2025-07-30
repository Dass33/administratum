-- +goose up
ALTER TABLE users
ADD COLUMN opened_table UUID REFERENCES user_tables(table_id) ON DELETE SET NULL;

-- +goose down
ALTER TABLE users
DROP COLUMN opened_table;
