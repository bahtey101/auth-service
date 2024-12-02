-- +migrate Up
create table tokens (
    user_id uuid constraint tokens_user_id_fkey references users(id),
    user_ip inet not null unique,
    token text not null unique,
    created_at timestamp not null default now()
);

-- +migrate Up
create index tokens_crid on tokens (created_at, user_id);

-- +migrate Down
drop table if exists tokens;

-- +migrate Down
drop index if exists tokens_crid;