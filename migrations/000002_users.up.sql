create table users (
    id uuid primary key,
    email varchar(255) unique not null,
);