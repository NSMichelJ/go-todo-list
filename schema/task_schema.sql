create database if not exists task;

use task;

create table if not exists task(
    id bigint unsigned not null auto_increment,
    content varchar(255) not null,
    created datetime(3),
    primary key(id) 
)