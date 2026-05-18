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
SELECT * FROM teams;

-- name: GetTeamByName :one
SELECT * FROM teams
Where teamname = $1;

-- name: WipeTeams :exec
DELETE FROM teams;

-- name: GetTeamIDByName :one
SELECT * FROM teams
WHERE teamname = $1;