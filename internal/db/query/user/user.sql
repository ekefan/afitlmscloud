-- name: CreateUser :one
INSERT INTO users (
    full_name,
    email,
    hashed_password,
    sch_id
) VALUES (
    $1, $2, $3, $4
) RETURNING *;

-- name: UpdateUser :one
UPDATE users
SET
    full_name = $2,
    email = $3,
    hashed_password = $4,
    password_changed = $5,
    updated_at = now()
WHERE id = $1
RETURNING *;

-- name: EnrollUser :one
UPDATE users
SET
    roles = $2,
    enrolled = $3,
    updated_at = now()
WHERE id = $1
    AND users.enrolled IS DISTINCT FROM TRUE
RETURNING *;

-- name: GetUserByID :one
SELECT * FROM users
WHERE id = $1;


-- name: DeleteUser :execresult
DELETE FROM users
WHERE id = $1;