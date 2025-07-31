-- +goose up
DROP TABLE users;

CREATE TABLE users (
    id UUID PRIMARY KEY,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    email TEXT NOT NULL UNIQUE,
    hashed_password TEXT NOT NULL DEFAULT 'unset',
    opened_sheet UUID REFERENCES sheets(id) ON DELETE SET NULL
);

-- +goose down
DROP TABLE users;

CREATE TABLE users (
    id UUID primary key,
    created_at timestamp not null,
    updated_at timestamp not null,
    email text not null unique,
    hashed_password text not null default 'unset' ,
    opened_sheet UUID REFERENCES user_tables(table_id) ON DELETE SET NULL
)
