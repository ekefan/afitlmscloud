-- name: GetLecturerAvailabilityForAllCourses :many
SELECT
    c.name As course_name,
    cl.availability,
    c.course_code
FROM course_lecturers cl
JOIN courses c ON c.course_code = cl.course_code
WHERE cl.lecturer_id = $1;

-- -- name: GetAllStudentsEligibilityForCourse :many
-- SELECT
--     cs.student_id
--     cs.eligibility
--     cs.updated_at
-- FROM course_students cs
-- WHERE cs.course_code = $1;