-- +goose up
DROP TABLE column_data;

CREATE TABLE column_data (
    id UUID primary key,
    idx integer not null,
    value text,
    column_id UUID not null,
    created_at timestamp not null,
    updated_at timestamp not null,
    constraint fk_column_data_column_id
        foreign key (column_id)
        references columns(id)
        on delete cascade
);


-- +goose down
DROP TABLE column_data;

CREATE TABLE column_data (
    id UUID primary key,
    idx integer not null,
    value text,
    sheet_id UUID not null,
    created_at timestamp not null,
    updated_at timestamp not null,
    constraint fk_rows_sheet_id
        foreign key (sheet_id)
        references sheets(id)
        on delete cascade
);
