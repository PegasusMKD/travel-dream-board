create type scrape_status as enum (
	'processing',
	'completed',
	'completed_by_ai',
	'failed'
);

create table scrape_audit (
	uuid uuid primary key default gen_random_uuid(),
	updated_at timestamp not null default now(),
	created_at timestamp not null default now(),

	url text not null,
	status scrape_status not null default 'processing',

	host text not null,

	title text,
	image_url text,
	description text
);
