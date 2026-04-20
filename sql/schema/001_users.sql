-- +goose Up
CREATE TABLE users(
    id uuid primary key default gen_random_uuid(),
    created_at timestamp not null,
    updated_at timestamp not null,
    name text unique not null
);

-- +goose Down
DROP TABLE users;
