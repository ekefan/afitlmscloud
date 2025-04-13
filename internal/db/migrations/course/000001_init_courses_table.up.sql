CREATE TABLE courses (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    faculty TEXT NOT NULL,
    department TEXT NOT NULL,
    level TEXT NOT NULL,
    code VARCHAR(10) UNIQUE NOT NULL,
    num_of_lectures_per_term INT NOT NULL,
    active_lecturer_id BIGINT,
);

CREATE INDEX idx_courses_faculty ON courses(faculty);
CREATE INDEX idx_courses_department ON courses(department);
CREATE INDEX idx_courses_level ON courses(level);
CREATE INDEX idx_courses_faculty_department_level ON courses(faculty, department, level);

CREATE TABLE course_lecturers (
    course_id BIGINT NOT NULL,
    lecturer_id BIGINT NOT NULL,
    
    PRIMARY KEY (course_id),
    
    FOREIGN KEY (course_id) REFERENCES courses(id) ON DELETE CASCADE,
);

CREATE TABLE course_registered_students (
    course_id BIGINT NOT NULL,
    student_id BIGINT NOT  NULL,

    PRIMARY KEY (course_id),
    FOREIGN KEY (course_id) REFERENCES courses(id) ON DELETE CASCADE,
)