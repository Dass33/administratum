-- +goose up
ALTER TABLE sheets
ADD COLUMN type TEXT NOT NULL DEFAULT 'list';

ALTER TABLE column_data
ADD COLUMN type TEXT;

-- +goose down
ALTER TABLE sheets
DROP COLUMN type;

ALTER TABLE column_data
DROP COLUMN type;
