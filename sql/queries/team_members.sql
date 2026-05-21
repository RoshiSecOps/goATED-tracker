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


-- name: GetTeamsForUser :many
SELECT teams.* 
FROM teams
JOIN team_members ON teams.id = team_members.team_id
WHERE team_members.user_id = $1;

-- name: CheckMembership :one
SELECT EXISTS (
    SELECT 1 FROM team_members
    WHERE user_id = $1 AND team_id = $2
);