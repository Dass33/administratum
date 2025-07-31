-- +goose up
ALTER TABLE users
RENAME COLUMN opened_table TO opened_sheet;

-- +goose down
ALTER TABLE users
RENAME COLUMN opened_sheet TO opened_table;
