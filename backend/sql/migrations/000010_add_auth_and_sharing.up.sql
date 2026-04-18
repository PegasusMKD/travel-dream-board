CREATE TABLE users (
	uuid uuid primary key default gen_random_uuid(),
	updated_at timestamp not null default now(),
	created_at timestamp not null default now(),

	email text unique,
	name text not null,
	avatar_url text
);

CREATE TABLE share_tokens (
    token text primary key,
    board_uuid uuid not null references boards(uuid) on delete cascade,
    created_at timestamp not null default now()
);

-- Add user_uuid to boards
ALTER TABLE boards ADD COLUMN user_uuid uuid references users(uuid) on delete cascade;

-- Add user_uuid to items
ALTER TABLE accomodations ADD COLUMN user_uuid uuid not null references users(uuid) on delete cascade;
ALTER TABLE activities ADD COLUMN user_uuid uuid not null references users(uuid) on delete cascade;
ALTER TABLE transport ADD COLUMN user_uuid uuid not null references users(uuid) on delete cascade;

-- Update votes
-- Delete old data if any, since we can't safely cast text names to uuids
DELETE FROM votes;
ALTER TABLE votes ALTER COLUMN voted_by TYPE uuid USING NULL;
ALTER TABLE votes ADD CONSTRAINT fk_votes_voted_by FOREIGN KEY (voted_by) REFERENCES users(uuid) ON DELETE CASCADE;

-- Update comments
DELETE FROM comments;
ALTER TABLE comments ALTER COLUMN created_by TYPE uuid USING NULL;
ALTER TABLE comments ADD CONSTRAINT fk_comments_created_by FOREIGN KEY (created_by) REFERENCES users(uuid) ON DELETE CASCADE;

-- Update memories
DELETE FROM memories;
ALTER TABLE memories ALTER COLUMN uploaded_by TYPE uuid USING NULL;
ALTER TABLE memories ADD CONSTRAINT fk_memories_uploaded_by FOREIGN KEY (uploaded_by) REFERENCES users(uuid) ON DELETE CASCADE;
