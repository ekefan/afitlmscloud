-- name: CreateStudent :one
INSERT INTO students (
    user_id, courses, biometric_template
) VALUES (
    $1, $2, $3
) RETURNING *;

-- name: GetStudentByID :one
SELECT * FROM students WHERE id = $1;

-- name: GetStudentByUserID :one
SELECT * FROM students WHERE user_id = $1;

-- name: ListStudents :many
SELECT * FROM students ORDER BY id;

-- name: UpdateStudentCourses :one
UPDATE students
SET courses = $2, updated_at = now()
WHERE id = $1
RETURNING *;

-- name: UpdateStudentBiometric :one
UPDATE students
SET biometric_template = $2, updated_at = now()
WHERE id = $1
RETURNING *;

-- name: DeleteStudent :exec
DELETE FROM students WHERE id = $1;
