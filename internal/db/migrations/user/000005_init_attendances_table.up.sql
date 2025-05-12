CREATE TABLE lecture_sessions (
    id BIGSERIAL PRIMARY KEY,
    course_code VARCHAR(10) NOT NULL,
    lecturer_id BIGINT NOT NULL,
    session_date DATE NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT  now()
);

CREATE TABLE lecture_attendance (
    session_id BIGINT NOT NULL,
    student_id BIGINT NOT NULL,
    attendance_time TIME NOT NULL,
    attended BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMP NOT NULL DEFAULT now(),

    FOREIGN KEY (session_id) REFERENCES lecture_sessions(id) ON DELETE CASCADE,
    PRIMARY KEY (session_id, student_id)
);