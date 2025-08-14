-- +goose Up
-- +goose StatementBegin
ALTER TABLE sheets ADD COLUMN source_sheet_id TEXT;
ALTER TABLE columns ADD COLUMN source_column_id TEXT;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE sheets DROP COLUMN source_sheet_id;
ALTER TABLE columns DROP COLUMN source_column_id;
-- +goose StatementEnd