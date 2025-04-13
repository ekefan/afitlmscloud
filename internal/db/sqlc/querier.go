// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.28.0

package db

import (
	"context"
	"database/sql"
)

type Querier interface {
	CreateAvailability(ctx context.Context, arg CreateAvailabilityParams) (Availability, error)
	CreateEligibility(ctx context.Context, arg CreateEligibilityParams) (Eligibility, error)
	CreateLecturer(ctx context.Context, arg CreateLecturerParams) (Lecturer, error)
	CreateStudent(ctx context.Context, arg CreateStudentParams) (Student, error)
	CreateUser(ctx context.Context, arg CreateUserParams) (User, error)
	DeleteAvailability(ctx context.Context, arg DeleteAvailabilityParams) (sql.Result, error)
	DeleteEligibility(ctx context.Context, arg DeleteEligibilityParams) (sql.Result, error)
	DeleteLecturer(ctx context.Context, id int64) (sql.Result, error)
	DeleteStudent(ctx context.Context, id int64) (sql.Result, error)
	DeleteUser(ctx context.Context, id int64) (sql.Result, error)
	GetAvailability(ctx context.Context, arg GetAvailabilityParams) (Availability, error)
	GetAvailabilityByCourseId(ctx context.Context, courseID int64) (Availability, error)
	GetEligibility(ctx context.Context, arg GetEligibilityParams) (Eligibility, error)
	GetEligibilityByCourseId(ctx context.Context, courseID int64) (Eligibility, error)
	GetLecturerByID(ctx context.Context, id int64) (Lecturer, error)
	GetLecturerByUserID(ctx context.Context, userID int64) (Lecturer, error)
	GetStudentByID(ctx context.Context, id int64) (Student, error)
	GetStudentByUserID(ctx context.Context, userID int64) (Student, error)
	GetUserByID(ctx context.Context, id int64) (User, error)
	ListAvailabilityForLecturer(ctx context.Context, lecturerID int64) ([]Availability, error)
	ListEligibilityForStudent(ctx context.Context, studentID int64) ([]Eligibility, error)
	SetMinEligibility(ctx context.Context, arg SetMinEligibilityParams) (Eligibility, error)
	UpdateAvailability(ctx context.Context, arg UpdateAvailabilityParams) (Availability, error)
	UpdateEligibility(ctx context.Context, arg UpdateEligibilityParams) (Eligibility, error)
	UpdateLecturerCourses(ctx context.Context, arg UpdateLecturerCoursesParams) (Lecturer, error)
	UpdateStudentCourses(ctx context.Context, arg UpdateStudentCoursesParams) (Student, error)
	UpdateUser(ctx context.Context, arg UpdateUserParams) (User, error)
}

var _ Querier = (*Queries)(nil)
