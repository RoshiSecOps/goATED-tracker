-- name: CreateTeamMember :one
INSERT INTO team_members (id, created_at, updated_at, user_id, team_id)
VALUES (
    gen_random_uuid(),
    NOW(),
    NOW(),
    $1,
    $2
)
RETURNING *;

-- name: GetAllTeamMember :many
SELECT * FROM team_members;

-- name: WipeTeamMember :exec
DELETE FROM team_members;