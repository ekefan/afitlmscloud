CREATE TABLE students (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL UNIQUE,
    courses TEXT[], -- list of course IDs (can be UUID or strings)
    biometric_template BYTEA,
    updated_at TIMESTAMP DEFAULT now()
);

CREATE TABLE eligibility (
    id BIGSERIAL PRIMARY KEY,
    course_id BIGINT NOT NULL,
    student_id BIGINT NOT NULL,
    value INTEGER NOT NULL,
    min_value INTEGER NOT NULL,
    updated_at TIMESTAMP DEFAULT now(),

    UNIQUE(course_id, student_id)
);
