CREATE TABLE attendances (
    course_id BIGINT NOT NULL,
    student_id BIGINT NOT NULL,
    attendance_date DATE NOT NULL,
    status BOOLEAN NOT NULL DEFAULT false, -- true = present, false = absent
    created_at TIMESTAMP DEFAULT now(),

    PRIMARY KEY (course_id, student_id, attendance_date)
);

CREATE INDEX idx_attendance_course_date ON attendances (course_id, attendance_date);
