-- +goose up
ALTER TABLE tables
ADD COLUMN name text not null;

-- +goose down
ALTER TABLE tables
DROP COLUMN name;
