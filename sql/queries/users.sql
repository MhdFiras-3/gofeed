-- name: CreateUser :one
INSERT INTO users (name,email,password)
VALUES (
    $1,
    $2,
    $3
)
RETURNING *;

-- name: UpdateUserPw :exec
UPDATE users
SET password = $2,updated_at = NOW()
WHERE id = $1;

-- name: GetUserByID :one
SELECT * FROM users
WHERE id = $1;

-- name: GetUserByEmail :one
SELECT * FROM users
WHERE email = $1;