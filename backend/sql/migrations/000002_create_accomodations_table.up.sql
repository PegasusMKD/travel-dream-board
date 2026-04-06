create type accomodations_status as enum (
	'considering',
	'finalist',
	'rejected',
	'booked',
	'completed'
);

create table accomodations (
	uuid uuid primary key default gen_random_uuid(),
	updated_at timestamp not null default now(),
	created_at timestamp not null default now(),

	url text not null,

	title text not null,
	image_url text,
	notes text,
	status accomodations_status not null default 'considering',

	booking_reference text,
	selected boolean not null default false,

	board_uuid uuid not null references boards (uuid)
);
