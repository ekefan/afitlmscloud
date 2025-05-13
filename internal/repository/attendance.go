package repository

import (
	"context"
	"database/sql"
	"time"

	db "github.com/ekefan/afitlmscloud/internal/db/sqlc"
)

type AttendanceRepository interface {
	CreateAttendanceSession(ctx context.Context, arg AttendanceSessionParams) error
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

type AttendanceSessionParams struct {
	AttendanceData []LectureAttendanceParams
	db.CreateLectureSessionParams
}

type LectureAttendanceParams struct {
	SessionID      int64     `json:"session_id,omitempty"`
	StudentID      int64     `json:"student_id" binding:"required"`
	AttendanceTime time.Time `json:"attendance_time" binding:"required"`
	Attended       bool      `json:"attended" binding:"required"`
}

func (as *attendanceStore) CreateAttendanceSession(ctx context.Context, arg AttendanceSessionParams) error {
	err := execTx(ctx, as.dbConn, func(q *db.Queries) error {
		sessionID, err := q.CreateLectureSession(ctx, db.CreateLectureSessionParams{
			CourseCode:  arg.CourseCode,
			LecturerID:  arg.LecturerID,
			SessionDate: arg.SessionDate,
		})
		if err != nil {
			return err
		}

		for _, attendanceData := range arg.AttendanceData {
			attendanceData.SessionID = sessionID
			err := q.CreateLectureAttendance(ctx, db.CreateLectureAttendanceParams{
				SessionID:      attendanceData.SessionID,
				StudentID:      attendanceData.StudentID,
				AttendanceTime: attendanceData.AttendanceTime,
				Attended:       attendanceData.Attended,
			})
			if err != nil {
				return err
			}
		}
		return nil
	})
	return err
}
