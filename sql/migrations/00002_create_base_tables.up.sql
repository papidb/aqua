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
  name text not null unique,
  region text not null,
  created_at timestamptz not null default current_timestamp,
  updated_at timestamptz,
  deleted_at timestamptz
);
-- add index to resources table to improve performance
CREATE INDEX idx_resources_id ON resources(id);

-- Join table to establish many-to-many relationship
create table customer_resources (
  customer_id uuid not null references customers(id) on delete cascade,
  resource_id uuid not null references resources(id) on delete cascade,
  created_at timestamptz not null DEFAULT current_timestamp,
  primary key (customer_id, resource_id)
);

-- add index to customer_resources table
CREATE INDEX idx_customer_resources_customer_id ON customer_resources(customer_id);
commit;