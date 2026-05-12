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
