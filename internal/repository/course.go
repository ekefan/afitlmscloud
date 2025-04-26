package repository

import (
	"context"
	"database/sql"

	db "github.com/ekefan/afitlmscloud/internal/db/sqlc"
)

type CourseRepository interface {
	CreateCourse(ctx context.Context, arg db.CreateCourseParams) (db.Course, error)
	DropCourse(ctx context.Context, arg db.DropCourseParams) (sql.Result, error)
	RegisterCourse(ctx context.Context, arg db.RegisterCourseParams) error
	GetStudentEligibilityForAllCourses(ctx context.Context, studentID int64) ([]db.GetStudentEligibilityForAllCoursesRow, error)
	UnassignLecturerFromCourse(ctx context.Context, arg db.UnassignLecturerFromCourseParams) (sql.Result, error)
	AssignLecturerToCourse(ctx context.Context, arg db.AssignLecturerToCourseParams) error
	GetLecturerAvailabilityForAllCourses(ctx context.Context, lecturerID int64) ([]db.GetLecturerAvailabilityForAllCoursesRow, error)
}

var _ CourseRepository = (*courseStrore)(nil)

type courseStrore struct {
	dbConn *sql.DB
	*db.Queries
}

func NewCourseStore(dbConn *sql.DB) CourseRepository {
	return &courseStrore{
		dbConn:  dbConn,
		Queries: db.New(dbConn),
	}
}
