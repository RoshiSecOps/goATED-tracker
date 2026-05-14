-- name: CreateTeam :one
INSERT INTO teams (id, created_at, updated_at, teamname)
VALUES (
    gen_random_uuid(),
    NOW(),
    NOW(),
    $1
)
RETURNING *;

-- name: GetAllTeams :many
SELECT * from teams;

-- name: GetTeamByName :one
SELECT * from teams
Where teamname = $1;

-- name: WipeTeams :exec
DELETE FROM teams;