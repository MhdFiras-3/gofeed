-- name: CreateRefreshToken :one
INSERT INTO refresh_tokens (token,expires_at,user_id)
VALUES (
    $1,
    $2,
    $3
)
RETURNING *;


-- name: GetUserFromRefreshToken :one
SELECT users.* FROM users
JOIN refresh_tokens on users.id = refresh_tokens.user_id
WHERE refresh_tokens.token = $1
AND revoked_at IS NULL
AND expires_at > NOW();

-- name: RevokeRefreshToken :exec
UPDATE refresh_tokens
SET revoked_at = NOW(), updated_at = NOW()
WHERE token = $1;