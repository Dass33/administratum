-- +goose up
CREATE TABLE refresh_tokens (
    token text primary key,
    created_at timestamp not null,
    updated_at timestamp not null,
    user_id UUID not null,
    expires_at timestamp not null,
    revoked_at timestamp,

    constraint fk_user_id
        foreign key (user_id)
        references users(id)
    on delete cascade
);

-- +goose down
DROP TABLE refresh_tokens;
