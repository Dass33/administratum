-- +goose Up
CREATE TABLE users (
    id UUID primary key,
    created_at timestamp not null,
    updated_at timestamp not null,
    email text not null unique,
    hashed_password text
    not null default 'unset'
);

-- +goose Down
DROP TABLE users;
