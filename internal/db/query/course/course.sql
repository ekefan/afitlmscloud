-- name: CreateCourse :one
INSERT INTO courses (
    name, 
    faculty, 
    department, 
    level, 
    course_code
) VALUES (
    $1, $2, $3, $4, $5
) RETURNING *;

-- name: GetCourseFiltered :many
SELECT * FROM courses
WHERE faculty = $1 AND department = $2 AND level = $3;


-- name: RegisterCourse :exec
INSERT INTO course_students (
    course_code,
    student_id
) VALUES (
    $1, $2
);

-- name: DropCourse :execresult
DELETE FROM course_students
WHERE course_code = $1 AND student_id = $2;