-- +migrate Up
create table tokens (
    user_id uuid constraint tokens_user_id_fkey references users(id),
    user_ip inet not null unique,
    token text not null unique
);
-- +migrate Down
drop table if exists tokens;