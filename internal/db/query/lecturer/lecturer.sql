-- name: CreateLecturer :one
INSERT INTO lecturers (
    user_id, biometric_template, courses, courses_actively_teaching
) VALUES (
    $1, $2, $3, $4
) RETURNING *;

-- name: GetLecturerByID :one
SELECT * FROM lecturers WHERE id = $1;

-- name: GetLecturerByUserID :one
SELECT * FROM lecturers WHERE user_id = $1;

-- name: UpdateLecturerCourses :one
UPDATE lecturers
SET courses = $2, courses_actively_teaching = $3, updated_at = now()
WHERE id = $1
RETURNING *;

-- name: DeleteLecturer :execresult
DELETE FROM lecturers WHERE id = $1;
