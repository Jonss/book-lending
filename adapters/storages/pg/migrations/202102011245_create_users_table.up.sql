create table users (
    id bigserial primary key,
    full_name varchar(120) not null,
    external_id uuid not null,
    email varchar(60) unique
);