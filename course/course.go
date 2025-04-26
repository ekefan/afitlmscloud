package course

import (
	"context"
	"errors"
	"log/slog"

	db "github.com/ekefan/afitlmscloud/internal/db/sqlc"
	"github.com/ekefan/afitlmscloud/internal/repository"
)

var (
	ErrFailedToRegisterCourse                    = errors.New("failed to register course")
	ErrFailedToDropCourse                        = errors.New("failed to drop course")
	ErrFailedToGetStudentEligbilityForAllCourses = errors.New("failed to get student eligibility for all courses")
	ErrFailedToAssignCourse                      = errors.New("failed to assign course")
	ErrFailedToUnAssignCourse                    = errors.New("failed to unassign course")
	ErrFailedToGetAvailabilityForAllCourses      = errors.New("failed to get availability for all courses")
	// ErrActiveLecturerAlreadySetForThisCourse     = errors.New("failed to set active lecturer")
	ErrFailedToSetActiveLecturer = errors.New("failed to set active lecturer")
)

type CourseService struct {
	repo repository.CourseRepository
}

type Eligibility struct {
	CourseCode  string
	CourseName  string
	Eligibility float64
}

type Availability struct {
	CourseCode   string
	CourseName   string
	Availability float64
}

func NewCourseService(courseRepository repository.CourseRepository) *CourseService {
	return &CourseService{
		repo: courseRepository,
	}
}

type UserCourseData struct {
	UserID     int64
	CourseCode string
}

func (csvc *CourseService) RegisterCouse(ctx context.Context, data UserCourseData) error {
	if err := csvc.repo.RegisterCourse(ctx, db.RegisterCourseParams{
		StudentID:  data.UserID,
		CourseCode: data.CourseCode,
	}); err != nil {
		slog.Error("failed to perform course operation", "error", err)
		return ErrFailedToRegisterCourse
	}
	return nil
}

func (csvc *CourseService) DropCourses(ctx context.Context, data UserCourseData) error {
	res, err := csvc.repo.DropCourse(ctx, db.DropCourseParams{
		StudentID:  data.UserID,
		CourseCode: data.CourseCode,
	})
	numOfRowsAffected, rerr := res.RowsAffected()
	if rerr != nil {
		slog.Error("failed to get the number of rows affected for dropping course", "error", err)
		return rerr
	}
	if numOfRowsAffected == 0 {
		slog.Error("no course found to drop", "studentId", data.UserID, "courseCode", data.CourseCode)
		return ErrFailedToDropCourse
	}
	if err != nil {
		slog.Error("failed to perform course operation", "error", err)
		return ErrFailedToDropCourse
	}
	return nil
}

func (csvc *CourseService) GetStudentEligibilityForAllCourses(ctx context.Context, studentID int64) ([]Eligibility, error) {
	res, err := csvc.repo.GetStudentEligibilityForAllCourses(ctx, studentID)
	if err != nil {
		slog.Error("failed to get student eligibility for all courses", "error", err)
		return []Eligibility{}, ErrFailedToGetStudentEligbilityForAllCourses
	}

	eligibility := []Eligibility{}
	for _, e := range res {
		eligibility = append(eligibility, Eligibility{
			CourseCode:  e.CourseCode,
			CourseName:  e.CourseName,
			Eligibility: e.Eligibility,
		})
	}
	return eligibility, nil
}

func (csvc *CourseService) AssignLecturerToCourse(ctx context.Context, data UserCourseData) error {
	if err := csvc.repo.AssignLecturerToCourse(ctx, db.AssignLecturerToCourseParams{
		CourseCode: data.CourseCode,
		LecturerID: data.UserID,
	}); err != nil {
		slog.Error("failed to assign lecturer to course", "error", err)
		return err
	}
	return nil
}

func (csvc *CourseService) UnassignLecturerFromCourse(ctx context.Context, data UserCourseData) error {
	res, err := csvc.repo.UnassignLecturerFromCourse(ctx, db.UnassignLecturerFromCourseParams{
		LecturerID: data.UserID,
		CourseCode: data.CourseCode,
	})
	numOfRowsAffected, rerr := res.RowsAffected()
	if rerr != nil {
		slog.Error("failed to get the number of rows affected for unassigning lecturer", "error", err)
		return rerr
	}
	if err != nil || numOfRowsAffected == 0 {
		slog.Error("failed to perform course operation", "error", err)
		return ErrFailedToDropCourse
	}
	return nil
}

func (csvc *CourseService) GetLecturerAvailabilityForAllCourses(ctx context.Context, lecturerID int64) ([]Availability, error) {
	res, err := csvc.repo.GetLecturerAvailabilityForAllCourses(ctx, lecturerID)
	if err != nil {
		slog.Error("failed to get student eligibility for all courses", "error", err)
		return []Availability{}, ErrFailedToGetStudentEligbilityForAllCourses
	}

	availability := []Availability{}
	for _, e := range res {
		availability = append(availability, Availability{
			CourseCode:   e.CourseCode,
			CourseName:   e.CourseName,
			Availability: e.Availability,
		})
	}
	return availability, nil
}

func (csvc *CourseService) SetActiveLecturer(ctx context.Context, lecturerID int64, courseCode string) error {
	err := csvc.repo.SetActiveLecturer(ctx, db.SetActiveLecturerParams{
		ActiveLecturerID: lecturerID,
		CourseCode:       courseCode})
	if err != nil {
		slog.Error("failed to set active lecturer id", "error", err)
		return ErrFailedToSetActiveLecturer
	}

	// TODO: Send ActiveLecturerEvent...
	return nil
}
