CREATE DATABASE website;

CREATE TABLE IF NOT EXISTS public_suffix (
  suffix varchar(255) not null,
  primary key (suffix)
) engine=InnoDB default charset=utf8 collate=utf8_general_ci;