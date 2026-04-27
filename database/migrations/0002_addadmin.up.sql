alter table users add column is_suspended BOOLEAN default false;
create type user_role as enum ('admin', 'employee');
alter table users 
ADD column role user_role default 'employee';