alter table books drop column author;
alter table books add column slug varchar(255) not null;
alter table books add column pages int not null;

alter table users add column created_at timestamp default Now();

alter table books_status add column created_at timestamp default Now();
alter table books_status add column returned_at timestamp default null;

