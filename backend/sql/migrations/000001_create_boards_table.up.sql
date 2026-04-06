-- Enable UUID extension
create extension if not exists "uuid-ossp";

-- trigger to update updated_at timestamp
create or replace function update_updated_at_column()
returns trigger as $$
begin
    new.updated_at = now();
    return new;
end;
$$ language 'plpgsql';

-- Table: boards

create type boards_status as enum (
	'planning',
	'cancelled',
	'completed'
);

create table boards (
	uuid uuid primary key default gen_random_uuid(),

	updated_at timestamp not null default now(),
	created_at timestamp not null default now(),

	name text not null,
	details text,
	location_name text not null,
	starts_at date,
	lasts_until date,

	status boards_status not null default 'planning',

	thumbnail_url text
);
