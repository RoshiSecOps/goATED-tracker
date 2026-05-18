-- name: CreateUser :one
INSERT INTO users (id, created_at, updated_at, username, passwordhash)
VALUES (
    gen_random_uuid(),
    NOW(),
    NOW(),
    $1,
    $2
)
RETURNING *;

-- name: WipeUsers :exec
DELETE FROM users;

-- name: GetUserByName :one
SELECT * FROM users
WHERE username = $1;

-- name: GetUserIDByName :one
SELECT * FROM users
WHERE username = $1;