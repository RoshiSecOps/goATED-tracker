-- name: CreateTeam :one
INSERT INTO teams (id, created_at, updated_at, teamname)
VALUES (
    gen_random_uuid(),
    NOW(),
    NOW(),
    $1
)
RETURNING *;