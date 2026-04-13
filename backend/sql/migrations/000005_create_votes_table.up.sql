create type voted_on as enum (
	'accomodation',
	'transport',
	'activities'
);

create table votes (
	uuid uuid primary key default gen_random_uuid(),
	updated_at timestamp not null default now(),
	created_at timestamp not null default now(),

	voted_by text not null,
	rank int not null,

	voted_on_ voted_on not null,
	voted_on_uuid uuid not null
);
