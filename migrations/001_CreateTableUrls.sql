-- +goose Up 
create table urls ( 
    url_id serial primary key,
    short varchar unique,
    original varchar unique
); 

-- +goose Down 
drop table urls; 