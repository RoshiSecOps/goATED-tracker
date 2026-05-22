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

-- name: CheckFindingAccess :one
SELECT EXISTS (
    SELECT 1 FROM findings
    WHERE id = $1 AND pentest_id = $2
);

-- name: CloseFinding :exec
UPDATE findings
SET status = 'closed', updated_at = NOW()
WHERE id = $1 AND status != 'closed';