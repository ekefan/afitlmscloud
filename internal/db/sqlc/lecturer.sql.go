// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.28.0
// source: lecturer.sql

package db

import (
	"context"
	"database/sql"

	"github.com/lib/pq"
)

const createLecturer = `-- name: CreateLecturer :one
INSERT INTO lecturers (
    user_id, biometric_template, courses, courses_actively_teaching
) VALUES (
    $1, $2, $3, $4
) RETURNING id, user_id, biometric_template, courses, courses_actively_teaching, updated_at
`

type CreateLecturerParams struct {
	UserID                  int64          `json:"user_id"`
	BiometricTemplate       sql.NullString `json:"biometric_template"`
	Courses                 []string       `json:"courses"`
	CoursesActivelyTeaching []string       `json:"courses_actively_teaching"`
}

func (q *Queries) CreateLecturer(ctx context.Context, arg CreateLecturerParams) (Lecturer, error) {
	row := q.queryRow(ctx, q.createLecturerStmt, createLecturer,
		arg.UserID,
		arg.BiometricTemplate,
		pq.Array(arg.Courses),
		pq.Array(arg.CoursesActivelyTeaching),
	)
	var i Lecturer
	err := row.Scan(
		&i.ID,
		&i.UserID,
		&i.BiometricTemplate,
		pq.Array(&i.Courses),
		pq.Array(&i.CoursesActivelyTeaching),
		&i.UpdatedAt,
	)
	return i, err
}

const deleteLecturer = `-- name: DeleteLecturer :execresult
DELETE FROM lecturers WHERE id = $1
`

func (q *Queries) DeleteLecturer(ctx context.Context, id int64) (sql.Result, error) {
	return q.exec(ctx, q.deleteLecturerStmt, deleteLecturer, id)
}

const getLecturerByID = `-- name: GetLecturerByID :one
SELECT id, user_id, biometric_template, courses, courses_actively_teaching, updated_at FROM lecturers WHERE id = $1
`

func (q *Queries) GetLecturerByID(ctx context.Context, id int64) (Lecturer, error) {
	row := q.queryRow(ctx, q.getLecturerByIDStmt, getLecturerByID, id)
	var i Lecturer
	err := row.Scan(
		&i.ID,
		&i.UserID,
		&i.BiometricTemplate,
		pq.Array(&i.Courses),
		pq.Array(&i.CoursesActivelyTeaching),
		&i.UpdatedAt,
	)
	return i, err
}

const getLecturerByUserID = `-- name: GetLecturerByUserID :one
SELECT id, user_id, biometric_template, courses, courses_actively_teaching, updated_at FROM lecturers WHERE user_id = $1
`

func (q *Queries) GetLecturerByUserID(ctx context.Context, userID int64) (Lecturer, error) {
	row := q.queryRow(ctx, q.getLecturerByUserIDStmt, getLecturerByUserID, userID)
	var i Lecturer
	err := row.Scan(
		&i.ID,
		&i.UserID,
		&i.BiometricTemplate,
		pq.Array(&i.Courses),
		pq.Array(&i.CoursesActivelyTeaching),
		&i.UpdatedAt,
	)
	return i, err
}

const updateLecturerCourses = `-- name: UpdateLecturerCourses :one
UPDATE lecturers
SET courses = $2, courses_actively_teaching = $3, updated_at = now()
WHERE id = $1
RETURNING id, user_id, biometric_template, courses, courses_actively_teaching, updated_at
`

type UpdateLecturerCoursesParams struct {
	ID                      int64    `json:"id"`
	Courses                 []string `json:"courses"`
	CoursesActivelyTeaching []string `json:"courses_actively_teaching"`
}

func (q *Queries) UpdateLecturerCourses(ctx context.Context, arg UpdateLecturerCoursesParams) (Lecturer, error) {
	row := q.queryRow(ctx, q.updateLecturerCoursesStmt, updateLecturerCourses, arg.ID, pq.Array(arg.Courses), pq.Array(arg.CoursesActivelyTeaching))
	var i Lecturer
	err := row.Scan(
		&i.ID,
		&i.UserID,
		&i.BiometricTemplate,
		pq.Array(&i.Courses),
		pq.Array(&i.CoursesActivelyTeaching),
		&i.UpdatedAt,
	)
	return i, err
}
