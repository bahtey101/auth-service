-- +migrate Up
create table users (
    id uuid primary key,
    email varchar(255) unique not null
);
-- +migrate Up
insert into users (id, email) values ('cfe51c65-a4d6-4a31-88ac-6dbdb5616f33', 'bob@example.com');
-- +migrate Up
insert into users (id, email) values ('2e780b4b-15fe-4288-8221-03e5bee470df', 'alex@example.com');
-- +migrate Up
insert into users (id, email) values ('a17e90db-14c7-4992-b357-e43b4ea1253f', 'fred@example.com');
-- +migrate Down
drop table if exists users;