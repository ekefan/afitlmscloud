-- name: CreateAvailability :one
INSERT INTO availabilities (
    course_id, lecturer_id, availability
) VALUES (
    $1, $2, $3
) RETURNING *;

-- name: GetAvailability :one
SELECT * FROM availabilities
WHERE course_id = $1 AND lecturer_id = $2;

-- name: GetAvailabilityByCourseId :one
SELECT * FROM availabilities
WHERE course_id = $1;

-- name: ListAvailabilityForLecturer :many
SELECT * FROM availabilities
WHERE lecturer_id = $1;

-- name: UpdateAvailability :one
UPDATE availabilities
SET availability = $3, updated_at = now()
WHERE course_id = $1 AND lecturer_id = $2
RETURNING *;

-- name: DeleteAvailability :execresult
DELETE FROM availabilities
WHERE course_id = $1 AND lecturer_id = $2;
