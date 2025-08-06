-- +goose up
ALTER TABLE sheets
DROP COLUMN row_count;

-- +goose down
ALTER TABLE sheets
ADD COLUMN row_count integer not null;
