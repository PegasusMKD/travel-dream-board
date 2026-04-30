alter table transport
	add column outbound_departing_location text,
	add column outbound_arriving_location text,
	add column outbound_departing_at timestamptz,
	add column outbound_arriving_at timestamptz,
	add column inbound_departing_location text,
	add column inbound_arriving_location text,
	add column inbound_departing_at timestamptz,
	add column inbound_arriving_at timestamptz;
