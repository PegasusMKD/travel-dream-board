create table memories (
	uuid uuid primary key default gen_random_uuid(),
	updated_at timestamp not null default now(),
	created_at timestamp not null default now(),

	uploaded_by text not null,
	board_uuid uuid not null references boards (uuid),

	image_url text not null
);
