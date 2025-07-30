-- +goose up
CREATE TABLE tables (
    id UUID primary key,
    game_url text,
    created_at timestamp not null,
    updated_at timestamp not null
);

CREATE TABLE branches (
    id UUID primary key,
    name text not null,
    table_id UUID not null,
    created_at timestamp not null,
    updated_at timestamp not null,
    constraint fk_branches_table_id
        foreign key (table_id)
        references tables(id)
        on delete cascade
);

CREATE TABLE sheets (
    id UUID primary key,
    name text not null,
    row_count integer not null,
    branch_id UUID not null,
    created_at timestamp not null,
    updated_at timestamp not null,
    constraint fk_sheets_branch_id
        foreign key (branch_id)
        references branches(id)
        on delete cascade
);

CREATE TABLE columns (
    id UUID primary key,
    name text not null,
    type text not null,
    required boolean not null default false,
    sheet_id UUID not null,
    created_at timestamp not null,
    updated_at timestamp not null,
    constraint fk_columns_sheet_id
        foreign key (sheet_id)
        references sheets(id)
        on delete cascade
);

CREATE TABLE rows (
    id UUID primary key,
    index_num integer not null,
    value text,
    sheet_id UUID not null,
    created_at timestamp not null,
    updated_at timestamp not null,
    constraint fk_rows_sheet_id
        foreign key (sheet_id)
        references sheets(id)
        on delete cascade
);

CREATE TABLE user_tables (
    user_id UUID,
    table_id UUID,
    branch_id UUID,
    sheet_id UUID,
    permission text not null,
    created_at timestamp not null,
    updated_at timestamp not null,
    primary key (user_id, table_id),
    constraint fk_user_tables_user_id
        foreign key (user_id)
        references users(id)
        on delete cascade,
    constraint fk_user_tables_table_id
        foreign key (table_id)
        references tables(id)
        on delete cascade,
    constraint fk_user_tables_branch_id
        foreign key (branch_id)
        references branches(id)
        on delete set null,
    constraint fk_user_tables_sheet_id
        foreign key (sheet_id)
        references sheets(id)
        on delete set null
);

-- +goose down
DROP TABLE user_tables;
DROP TABLE rows;
DROP TABLE columns;
DROP TABLE sheets;
DROP TABLE branches;
DROP TABLE tables;
