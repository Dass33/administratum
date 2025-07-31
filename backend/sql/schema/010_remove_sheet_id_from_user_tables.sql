-- +goose up
ALTER TABLE user_tables DROP COLUMN sheet_id;

-- +goose down
ALTER TABLE user_tables ADD COLUMN sheet_id UUID REFERENCES sheets(id) ON DELETE SET NULL; 