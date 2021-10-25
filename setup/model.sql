drop table if exists users;

create table users (
    id serial not null primary key,
    first_name varchar(34) not null,
    last_name varchar(48) not null,
    username varchar(50) not null,
    age int not null
);