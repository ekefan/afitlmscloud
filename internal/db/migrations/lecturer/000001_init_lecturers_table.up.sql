CREATE TABLE lecturers (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL UNIQUE,
    biometric_template TEXT,
    courses TEXT[],
    courses_actively_teaching TEXT[],
    updated_at TIMESTAMP DEFAULT now(),

    UNIQUE(user_id, biometric_template)
);

CREATE TABLE availabilities (
    id BIGSERIAL PRIMARY KEY,
    course_id BIGINT NOT NULL,
    lecturer_id BIGINT NOT NULL REFERENCES lecturers(id) ON DELETE CASCADE,
    availability FLOAT NOT NULL,
    updated_at TIMESTAMP DEFAULT now(),

    UNIQUE(course_id, lecturer_id)
);
