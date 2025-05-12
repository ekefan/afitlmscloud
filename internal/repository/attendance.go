package repository

import (
	"context"
	"database/sql"

	db "github.com/ekefan/afitlmscloud/internal/db/sqlc"
)

type AttendanceRepository interface {
	CreateLectureAttendance(ctx context.Context, arg db.CreateLectureAttendanceParams) error
	CreateLectureSession(ctx context.Context, arg db.CreateLectureSessionParams) (int64, error)
	GetLectureAttendance(ctx context.Context, sessionID int64) ([]db.LectureAttendance, error)
	GetLectureSession(ctx context.Context, courseCode string) ([]db.LectureSession, error)
}

var _ AttendanceRepository = (*attendanceStore)(nil)

type attendanceStore struct {
	dbConn *sql.DB
	*db.Queries
}

func NewAttendanceStore(dbConn *sql.DB) AttendanceRepository {
	return &attendanceStore{
		dbConn:  dbConn,
		Queries: db.New(dbConn),
	}
}
