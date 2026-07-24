-- name: CreateUser :one
INSERT INTO users (name,email,hashed_password)
VALUES (
    $1,
    $2,
    $3
)
RETURNING *;

-- name: UpdateUser :one
UPDATE users
SET 
    name = COALESCE(sqlc.narg('name'),name),
    email = COALESCE(sqlc.narg('email'),email),
    hashed_password = COALESCE(sqlc.narg('hashed_password'),hashed_password),
    updated_at = NOW()
WHERE id = $1
RETURNING *;


-- name: GetUserByID :one
SELECT * FROM users
WHERE id = $1;

-- name: GetUserByEmail :one
SELECT * FROM users
WHERE email = $1;