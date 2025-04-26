CREATE TABLE lecturers (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL UNIQUE,
    biometric_template TEXT NOT NULL,
    courses TEXT[],
    courses_actively_teaching TEXT[],
    updated_at TIMESTAMP DEFAULT now(),
    weighted_availability FLOAT NOT NULL DEFAULT 0.00,

    UNIQUE(user_id, biometric_template)
);