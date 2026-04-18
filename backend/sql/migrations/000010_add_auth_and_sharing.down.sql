ALTER TABLE memories DROP CONSTRAINT fk_memories_uploaded_by;
ALTER TABLE memories ALTER COLUMN uploaded_by TYPE text USING NULL;

ALTER TABLE comments DROP CONSTRAINT fk_comments_created_by;
ALTER TABLE comments ALTER COLUMN created_by TYPE text USING NULL;

ALTER TABLE votes DROP CONSTRAINT fk_votes_voted_by;
ALTER TABLE votes ALTER COLUMN voted_by TYPE text USING NULL;

ALTER TABLE transport DROP COLUMN user_uuid;
ALTER TABLE activities DROP COLUMN user_uuid;
ALTER TABLE accomodations DROP COLUMN user_uuid;

ALTER TABLE boards DROP COLUMN user_uuid;

DROP TABLE share_tokens;
DROP TABLE users;
