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

-- name: GetCourse :one
SELECT * FROM courses
WHERE course_code = $1;

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


-- name: AssignLecturerToCourse :exec
INSERT INTO course_lecturers (
    course_code,
    lecturer_id
) VALUES (
    $1, $2
);

-- name: DeleteCourse :execresult
DELETE FROM courses
WHERE course_code = $1;

-- name: UnassignLecturerFromCourse :execresult
DELETE FROM course_lecturers
WHERE course_code = $1 AND lecturer_id = $2;

-- name: SetActiveLecturer :exec
UPDATE courses 
SET
    active_lecturer_id = $1
WHERE active_lecturer_id = 0 AND course_code = $2;

-- name: GetCourseMetaData :one
SELECT
    c.name,
    c.faculty,
    c.department,
    c.level
FROM courses c
WHERE c.course_code = $1;


-- name: UpdateLecturerAttendedCount :exec
UPDATE courses
SET lecturer_attended_count = lecturer_attended_count + 1
WHERE course_code = $1;


-- name: UpdateCourseNumberOfLecturesPerSemester :exec
UPDATE courses
SET num_of_lectures_per_semester = $2
WHERE course_code = $1;