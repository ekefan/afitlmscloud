-- name: CreateEligibility :one
INSERT INTO eligibilities (
    course_id, student_id, eligibility, min_eligibility
) VALUES (
    $1, $2, $3, $4
) RETURNING *;

-- name: GetEligibility :one
SELECT * FROM eligibilities
WHERE course_id = $1 AND student_id = $2;

-- name: GetEligibilityByCourseId :one
SELECT * FROM eligibilities
WHERE course_id = $1;



-- name: ListEligibilityForStudent :many
SELECT * FROM eligibilities
WHERE student_id = $1;

-- name: UpdateEligibility :one
UPDATE eligibilities
SET eligibility = $3, updated_at = now()
WHERE course_id = $1 AND student_id = $2
RETURNING *;

-- name: SetMinEligibility :one
UPDATE eligibilities
SET eligibility = $3, updated_at = now()
WHERE course_id = $1 AND student_id = $2
RETURNING *;

-- name: DeleteEligibility :execresult
DELETE FROM eligibilities
WHERE course_id = $1 AND student_id = $2;
