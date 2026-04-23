-- name: GetUserByEmail :one
SELECT * FROM users WHERE email = @email LIMIT 1;

-- name: GetUserById :one
SELECT * FROM users WHERE uuid = @uuid LIMIT 1;

-- name: UpsertUser :one
INSERT INTO users (email, name, avatar_url)
VALUES (@email, @name, @avatar_url)
ON CONFLICT (email)
DO UPDATE SET
    name = EXCLUDED.name,
    avatar_url = EXCLUDED.avatar_url,
    updated_at = now()
RETURNING *;

-- name: CreateGuestUser :one
INSERT INTO users (name)
VALUES (@name)
RETURNING *;
