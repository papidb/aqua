-- create base tables
begin;
create table customers (
  id uuid primary key not null default gen_random_uuid(),
  name text not null unique,
  email text not null unique,
  created_at timestamptz not null default current_timestamp,
  updated_at timestamptz,
  deleted_at timestamptz
);
create table resources (
  id uuid primary key not null default gen_random_uuid(),
  name text not null,
  region text not null,
  created_at timestamptz not null default current_timestamp,
  updated_at timestamptz,
  deleted_at timestamptz
);
-- map one customer to many resources
alter table resources
add column customer_name text references customers(name) on delete cascade;
commit;