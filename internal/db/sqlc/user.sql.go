// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.28.0
// source: user.sql

package db

import (
	"context"
	"database/sql"

	"github.com/lib/pq"
)

const createUser = `-- name: CreateUser :one
INSERT INTO users (
    full_name,
    email,
    hashed_password,
    sch_id
) VALUES (
    $1, $2, $3, $4
) RETURNING id, full_name, roles, enrolled, email, sch_id, hashed_password, password_changed, updated_at, created_at
`

type CreateUserParams struct {
	FullName       string `json:"full_name"`
	Email          string `json:"email"`
	HashedPassword string `json:"hashed_password"`
	SchID          string `json:"sch_id"`
}

func (q *Queries) CreateUser(ctx context.Context, arg CreateUserParams) (User, error) {
	row := q.queryRow(ctx, q.createUserStmt, createUser,
		arg.FullName,
		arg.Email,
		arg.HashedPassword,
		arg.SchID,
	)
	var i User
	err := row.Scan(
		&i.ID,
		&i.FullName,
		pq.Array(&i.Roles),
		&i.Enrolled,
		&i.Email,
		&i.SchID,
		&i.HashedPassword,
		&i.PasswordChanged,
		&i.UpdatedAt,
		&i.CreatedAt,
	)
	return i, err
}

const deleteUser = `-- name: DeleteUser :execresult
DELETE FROM users
WHERE id = $1
`

func (q *Queries) DeleteUser(ctx context.Context, id int64) (sql.Result, error) {
	return q.exec(ctx, q.deleteUserStmt, deleteUser, id)
}

const enrollUser = `-- name: EnrollUser :one
UPDATE users
SET
    roles = $2,
    enrolled = $3,
    updated_at = now()
WHERE id = $1
    AND users.enrolled IS DISTINCT FROM TRUE
RETURNING id, full_name, roles, enrolled, email, sch_id, hashed_password, password_changed, updated_at, created_at
`

type EnrollUserParams struct {
	ID       int64    `json:"id"`
	Roles    []string `json:"roles"`
	Enrolled bool     `json:"enrolled"`
}

func (q *Queries) EnrollUser(ctx context.Context, arg EnrollUserParams) (User, error) {
	row := q.queryRow(ctx, q.enrollUserStmt, enrollUser, arg.ID, pq.Array(arg.Roles), arg.Enrolled)
	var i User
	err := row.Scan(
		&i.ID,
		&i.FullName,
		pq.Array(&i.Roles),
		&i.Enrolled,
		&i.Email,
		&i.SchID,
		&i.HashedPassword,
		&i.PasswordChanged,
		&i.UpdatedAt,
		&i.CreatedAt,
	)
	return i, err
}

const getUserByEmail = `-- name: GetUserByEmail :one
SELECT id, full_name, roles, enrolled, email, sch_id, hashed_password, password_changed, updated_at, created_at FROM users
WHERE email = $1
`

func (q *Queries) GetUserByEmail(ctx context.Context, email string) (User, error) {
	row := q.queryRow(ctx, q.getUserByEmailStmt, getUserByEmail, email)
	var i User
	err := row.Scan(
		&i.ID,
		&i.FullName,
		pq.Array(&i.Roles),
		&i.Enrolled,
		&i.Email,
		&i.SchID,
		&i.HashedPassword,
		&i.PasswordChanged,
		&i.UpdatedAt,
		&i.CreatedAt,
	)
	return i, err
}

const getUserByID = `-- name: GetUserByID :one
SELECT id, full_name, roles, enrolled, email, sch_id, hashed_password, password_changed, updated_at, created_at FROM users
WHERE id = $1
`

func (q *Queries) GetUserByID(ctx context.Context, id int64) (User, error) {
	row := q.queryRow(ctx, q.getUserByIDStmt, getUserByID, id)
	var i User
	err := row.Scan(
		&i.ID,
		&i.FullName,
		pq.Array(&i.Roles),
		&i.Enrolled,
		&i.Email,
		&i.SchID,
		&i.HashedPassword,
		&i.PasswordChanged,
		&i.UpdatedAt,
		&i.CreatedAt,
	)
	return i, err
}

const updateUserEmail = `-- name: UpdateUserEmail :one
UPDATE users
SET 
    email = $1, -- new_email
    updated_at = now()
WHERE id = $2 AND email = $3 -- old_email
RETURNING id, full_name, roles, enrolled, email, sch_id, hashed_password, password_changed, updated_at, created_at
`

type UpdateUserEmailParams struct {
	NewEmail string `json:"new_email"`
	ID       int64  `json:"id"`
	OldEmail string `json:"old_email"`
}

func (q *Queries) UpdateUserEmail(ctx context.Context, arg UpdateUserEmailParams) (User, error) {
	row := q.queryRow(ctx, q.updateUserEmailStmt, updateUserEmail, arg.NewEmail, arg.ID, arg.OldEmail)
	var i User
	err := row.Scan(
		&i.ID,
		&i.FullName,
		pq.Array(&i.Roles),
		&i.Enrolled,
		&i.Email,
		&i.SchID,
		&i.HashedPassword,
		&i.PasswordChanged,
		&i.UpdatedAt,
		&i.CreatedAt,
	)
	return i, err
}

const updateUserPassword = `-- name: UpdateUserPassword :one
UPDATE users
SET 
    hashed_password = $2,
    password_changed = TRUE,
    updated_at = now()
WHERE id = $1
RETURNING id, full_name, roles, enrolled, email, sch_id, hashed_password, password_changed, updated_at, created_at
`

type UpdateUserPasswordParams struct {
	ID             int64  `json:"id"`
	HashedPassword string `json:"hashed_password"`
}

func (q *Queries) UpdateUserPassword(ctx context.Context, arg UpdateUserPasswordParams) (User, error) {
	row := q.queryRow(ctx, q.updateUserPasswordStmt, updateUserPassword, arg.ID, arg.HashedPassword)
	var i User
	err := row.Scan(
		&i.ID,
		&i.FullName,
		pq.Array(&i.Roles),
		&i.Enrolled,
		&i.Email,
		&i.SchID,
		&i.HashedPassword,
		&i.PasswordChanged,
		&i.UpdatedAt,
		&i.CreatedAt,
	)
	return i, err
}
