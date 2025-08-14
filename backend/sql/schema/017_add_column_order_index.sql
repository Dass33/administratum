-- +goose Up
ALTER TABLE columns ADD COLUMN order_index INTEGER NOT NULL DEFAULT 0;

-- Populate order_index based on created_at to maintain current ordering
UPDATE columns 
SET order_index = (
    SELECT ROW_NUMBER() OVER (PARTITION BY sheet_id ORDER BY created_at) - 1
    FROM columns c2 
    WHERE c2.id = columns.id
);

-- +goose Down
ALTER TABLE columns DROP COLUMN order_index;
