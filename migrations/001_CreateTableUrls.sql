-- +goose Up 
create table Urls ( 
    url_id serial primary key,
    short varchar unique,
    original varchar unique
); 

-- +goose down 
drop table Urls; 