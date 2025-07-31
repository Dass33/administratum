-- +goose up
DROP TABLE user_tables;

CREATE TABLE user_tables (
    user_id UUID NOT NULL,
    table_id UUID NOT NULL,
    permission TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    PRIMARY KEY (user_id, table_id),
    CONSTRAINT fk_user_tables_user_id
        FOREIGN KEY (user_id)
        REFERENCES users(id)
        ON DELETE CASCADE,
    CONSTRAINT fk_user_tables_table_id
        FOREIGN KEY (table_id)
        REFERENCES tables(id)
        ON DELETE CASCADE
);

-- +goose Down

DROP TABLE user_tables;

CREATE TABLE user_tables (
    user_id UUID NOT NULL,
    table_id UUID NOT NULL,
    sheet_id UUID,
    permission TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    PRIMARY KEY (user_id, table_id),
    CONSTRAINT fk_user_tables_user_id
        FOREIGN KEY (user_id)
        REFERENCES users(id)
        ON DELETE CASCADE,
    CONSTRAINT fk_user_tables_table_id
        FOREIGN KEY (table_id)
        REFERENCES tables(id)
        ON DELETE CASCADE,
    CONSTRAINT fk_user_tables_sheet_id
        FOREIGN KEY (sheet_id)
        REFERENCES sheets(id)
        ON DELETE SET NULL
);
