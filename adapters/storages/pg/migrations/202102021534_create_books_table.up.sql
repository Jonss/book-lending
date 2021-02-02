create table books (
    id bigserial primary key,
    title varchar(200) not null,
    author varchar(120) not null,
    owner_id bigint,
    created_at timestamp default Now(),
    constraint fk_user foreign key(owner_id) references users(id)
)