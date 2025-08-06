-- +goose up
ALTER TABLE branches
ADD COLUMN is_protected boolean not null default false;

-- +goose down
ALTER TABLE branches
DROP COLUMN is_protected;
