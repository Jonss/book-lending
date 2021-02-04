create table books_status(
    id bigserial primary key,
    book_id bigint,
    bearer_user_id bigint,
    status varchar(30) not null,
    constraint fk_book_status foreign key(book_id) references books(id),
    constraint fk_book_lender_user foreign key(bearer_user_id) references users(id)
);