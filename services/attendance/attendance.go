package attendance

import (
	"context"
	"errors"
	"log/slog"
	"time"

	db "github.com/ekefan/afitlmscloud/internal/db/sqlc"
	"github.com/ekefan/afitlmscloud/internal/repository"
	"github.com/ekefan/afitlmscloud/services/course"
)

var (
	ErrFailedToCreateAttendanceSession = errors.New("failed to create attendance session")
)

type AttendanceService struct {
	repo          repository.AttendanceRepository
	courseService *course.CourseService
}

func NewAttendanceService(courseService *course.CourseService, attendanceRepository repository.AttendanceRepository) *AttendanceService {
	return &AttendanceService{
		repo:          attendanceRepository,
		courseService: courseService,
	}
}

type AttendanceSession struct {
	CourseCode     string                               `json:"course_code"`
	LecturerID     int64                                `json:"lecturer_id"`
	SessionDate    time.Time                            `json:"session_date"` // parse to time.Time
	AttendanceData []repository.LectureAttendanceParams `json:"attendance_data"`
}

func (as *AttendanceService) createNewAttendanceSession(ctx context.Context, attendanceSession AttendanceSession) error {
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

	studentData := make([]repository.StudentAttendanceData, len(attendanceSession.AttendanceData))
	for i, data := range attendanceSession.AttendanceData {
		studentData[i] = repository.StudentAttendanceData{
			StudentID: data.StudentID,
			Attended:  data.Attended,
		}
	}
	lectureMetaDataUpdate := course.UpdateCourseLectureMetaData{
		CourseCode:               attendanceSession.CourseCode,
		LecturerID:               attendanceSession.LecturerID,
		StudentAttendanceRecords: studentData,
	}
	err = as.courseService.OnAttendanceSessionCreated(ctx, lectureMetaDataUpdate)
	if err != nil {
		slog.Error("failed to update availability and eligibility for users", "error", err)
		return err
	}
	return nil
}
