create type commented_on as enum (
	'accomodation',
	'transport',
	'activities'
);

create table item_comments (
	uuid uuid primary key default gen_random_uuid(),
	updated_at timestamp not null default now(),
	created_at timestamp not null default now(),

	created_by text not null,
	content text not null,

	commented_on_ commented_on not null,
	commented_on_uuid uuid not null
);
