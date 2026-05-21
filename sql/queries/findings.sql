-- name: AddFindingToPentest :one
INSERT INTO findings (id, created_at, updated_at, title, status, severity, severity_score,
file, at_line, description, pentest_id)
VALUES (
    gen_random_uuid(),
    NOW(),
    NOW(),
    $1,
    $2,
    $3,
    $4,
    $5,
    $6,
    $7,
    $8
)
RETURNING *;

-- name: GetFindingsForPentest :many
SELECT * FROM findings
WHERE pentest_id = $1;

-- name: GetAllFindings :many
SELECT * FROM findings;

-- name: WipeFindings :exec
DELETE FROM findings;