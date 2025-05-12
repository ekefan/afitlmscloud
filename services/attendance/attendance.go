package attendance

import (
	"context"
	"errors"
	"log/slog"
	"time"

	db "github.com/ekefan/afitlmscloud/internal/db/sqlc"
	"github.com/ekefan/afitlmscloud/internal/repository"
)

var (
	ErrFailedToCreateAttendanceSession = errors.New("failed to create attendance session")
)

type AttendanceService struct {
	repo repository.AttendanceRepository
}

func NewAttendanceService(attendanceRepository repository.AttendanceRepository) *AttendanceService {
	return &AttendanceService{
		repo: attendanceRepository,
	}
}

type AttendanceSession struct {
	CourseCode     string                               `json:"course_code"`
	LecturerID     int64                                `json:"lecturer_id"`
	SessionDate    time.Time                            `json:"session_date"` // parse to time.Time
	AttendanceData []repository.LectureAttendanceParams `json:"attendance_data"`
}

func (as *AttendanceService) CreateNewAttendanceSession(ctx context.Context, attendanceSession AttendanceSession) error {
	err := as.repo.CreateAttendanceSession(ctx, repository.AttendanceSessionParams{
		AttendanceData: attendanceSession.AttendanceData,
		CreateLectureSessionParams: db.CreateLectureSessionParams{
			CourseCode:  attendanceSession.CourseCode,
			LecturerID:  attendanceSession.LecturerID,
			SessionDate: attendanceSession.SessionDate,
		},
	})

	if err != nil {
		slog.Error("failed to create attendance session", "error", err)
		return ErrFailedToCreateAttendanceSession
	}
	return nil
}
