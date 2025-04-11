-- name: CreateEligibility :one
INSERT INTO eligibility (
    course_id, student_id, value, min_value
) VALUES (
    $1, $2, $3, $4
) RETURNING *;

-- name: GetEligibility :one
SELECT * FROM eligibility
WHERE course_id = $1 AND student_id = $2;

-- name: ListEligibilityForStudent :many
SELECT * FROM eligibility
WHERE student_id = $1;

-- name: UpdateEligibility :one
UPDATE eligibility
SET value = $3, min_value = $4, updated_at = now()
WHERE course_id = $1 AND student_id = $2
RETURNING *;

-- name: DeleteEligibility :exec
DELETE FROM eligibility
WHERE course_id = $1 AND student_id = $2;
