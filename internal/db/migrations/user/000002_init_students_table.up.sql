CREATE TABLE students (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL UNIQUE,
    courses TEXT[], -- list of course IDs (can be UUID or strings)
    biometric_template TEXT NOT NULL,
    updated_at TIMESTAMP DEFAULT now()
);

CREATE TABLE eligibilities (
    id BIGSERIAL PRIMARY KEY,
    course_id BIGINT NOT NULL,
    student_id BIGINT NOT NULL REFERENCES students(id) ON DELETE CASCADE,
    eligibility FLOAT NOT NULL,
    min_eligibility FLOAT NOT NULL,
    updated_at TIMESTAMP DEFAULT now(),

    UNIQUE(course_id, student_id)
);
