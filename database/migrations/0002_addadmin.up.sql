alter table users add column is_suspended BOOLEAN default false;
create type new_user_role as enum ('admin', 'user');
alter table users 
ADD column role new_user_role default 'user';