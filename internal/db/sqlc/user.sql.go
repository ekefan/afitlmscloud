// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.28.0
// source: user.sql

package db

import (
	"context"
	"database/sql"
)

const createUser = `-- name: CreateUser :one
INSERT INTO users (
    full_name,
    email,
    hashed_password,
    sch_id
) VALUES (
    $1, $2, $3, $4
) RETURNING id, full_name, email, sch_id, hashed_password, password_changed, updated_at, created_at
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

const getUserByID = `-- name: GetUserByID :one
SELECT id, full_name, email, sch_id, hashed_password, password_changed, updated_at, created_at FROM users
WHERE id = $1
`

func (q *Queries) GetUserByID(ctx context.Context, id int64) (User, error) {
	row := q.queryRow(ctx, q.getUserByIDStmt, getUserByID, id)
	var i User
	err := row.Scan(
		&i.ID,
		&i.FullName,
		&i.Email,
		&i.SchID,
		&i.HashedPassword,
		&i.PasswordChanged,
		&i.UpdatedAt,
		&i.CreatedAt,
	)
	return i, err
}

const updateUser = `-- name: UpdateUser :one
UPDATE users
SET
    full_name = $2,
    email = $3,
    hashed_password = $4,
    password_changed = $5,
    updated_at = now()
WHERE id = $1
RETURNING id, full_name, email, sch_id, hashed_password, password_changed, updated_at, created_at
`

type UpdateUserParams struct {
	ID              int64  `json:"id"`
	FullName        string `json:"full_name"`
	Email           string `json:"email"`
	HashedPassword  string `json:"hashed_password"`
	PasswordChanged bool   `json:"password_changed"`
}

func (q *Queries) UpdateUser(ctx context.Context, arg UpdateUserParams) (User, error) {
	row := q.queryRow(ctx, q.updateUserStmt, updateUser,
		arg.ID,
		arg.FullName,
		arg.Email,
		arg.HashedPassword,
		arg.PasswordChanged,
	)
	var i User
	err := row.Scan(
		&i.ID,
		&i.FullName,
		&i.Email,
		&i.SchID,
		&i.HashedPassword,
		&i.PasswordChanged,
		&i.UpdatedAt,
		&i.CreatedAt,
	)
	return i, err
}
