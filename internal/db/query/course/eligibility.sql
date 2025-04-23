-- name: GetStudentEligibilityForAllCourses :many
SELECT
    c.name As course_name,
    cs.eligibility,
    c.course_code
FROM course_students cs
JOIN courses c ON c.course_code = cs.course_code
WHERE cs.student_id = $1;

-- -- name: GetAllStudentsEligibilityForCourse :many
-- SELECT
--     cs.student_id
--     cs.eligibility
--     cs.updated_at
-- FROM course_students cs
-- WHERE cs.course_code = $1;