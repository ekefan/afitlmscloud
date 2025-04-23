package repository

import (
	"context"
	"database/sql"

	db "github.com/ekefan/afitlmscloud/internal/db/sqlc"
)

type CourseRepository interface {
	GetCourseFiltered(ctx context.Context, arg db.GetCourseFilteredParams) ([]db.Course, error)
	CreateCourse(ctx context.Context, arg db.CreateCourseParams) (db.Course, error)
	DropCourse(ctx context.Context, arg db.DropCourseParams) (sql.Result, error)
	RegisterCourse(ctx context.Context, arg db.RegisterCourseParams) error
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
