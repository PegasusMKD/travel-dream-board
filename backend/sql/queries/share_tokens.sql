-- name: CreateShareToken :one
INSERT INTO share_tokens (token, board_uuid)
VALUES (@token, @board_uuid)
RETURNING *;

-- name: GetShareToken :one
SELECT * FROM share_tokens WHERE token = @token LIMIT 1;

-- name: DeleteShareToken :exec
DELETE FROM share_tokens WHERE token = @token AND board_uuid = @board_uuid;

-- name: GetShareTokensForBoard :many
SELECT * FROM share_tokens WHERE board_uuid = @board_uuid ORDER BY created_at DESC;
