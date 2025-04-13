package repository

import (
	"context"
	"database/sql"

	db "github.com/ekefan/afitlmscloud/internal/db/sqlc"
)

type StudentRepository interface {
	CreateStudent(ctx context.Context, arg db.CreateStudentParams) (db.Student, error)
	DeleteStudent(ctx context.Context, id int64) (sql.Result, error)
	GetStudentByID(ctx context.Context, id int64) (db.Student, error)
	GetStudentByUserID(ctx context.Context, userID int64) (db.Student, error)
	UpdateStudentCourses(ctx context.Context, arg db.UpdateStudentCoursesParams) (db.Student, error)
}

var _ StudentRepository = (*StudentStore)(nil)

type StudentStore struct {
	dbConn *sql.DB
	*db.Queries
}

func NewStudentStore(dbConn *sql.DB) StudentRepository {
	return &StudentStore{
		dbConn:  dbConn,
		Queries: db.New(dbConn),
	}
}
