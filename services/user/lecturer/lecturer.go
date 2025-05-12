package lecturer

import (
	"context"
	"errors"
	"log/slog"

	db "github.com/ekefan/afitlmscloud/internal/db/sqlc"
	"github.com/ekefan/afitlmscloud/internal/repository"
	"github.com/ekefan/afitlmscloud/services/course"
)

type LecturerService struct {
	courseService *course.CourseService
	repo          repository.LecturerRepository
}

type Lecturer struct {
	ID int64
}

var (
	ErrFailedToCreateLecturer = errors.New("failed to create lecturer")
)

func NewLecturerService(courseService *course.CourseService, repo repository.LecturerRepository) *LecturerService {
	return &LecturerService{
		courseService: courseService,
		repo:          repo,
	}
}

func (ls *LecturerService) CreateLecturer(ctx context.Context, lecturerID int64, biometricTemplate string) (Lecturer, error) {
	lecturer, err := ls.repo.CreateLecturer(ctx, db.CreateLecturerParams{
		UserID:            lecturerID,
		BiometricTemplate: biometricTemplate,
	})
	if err != nil {
		slog.Error("failed to create a new lecturer", "error", err)
		return Lecturer{}, ErrFailedToCreateLecturer
	}
	return Lecturer{
		ID: lecturer.UserID,
	}, nil
}

// TODO: make a one batch process.
func (ls *LecturerService) AssignLecturerToCourse(ctx context.Context, lecturerID int64, courseCodes []string) error {
	for _, c := range courseCodes {
		err := ls.courseService.AssignLecturerToCourse(ctx, course.UserCourseData{
			CourseCode: c,
			UserID:     lecturerID,
		})
		if err != nil {
			slog.Error("Handle error when assigning courses")
			return err
		}
	}
	return nil
}

func (ls *LecturerService) UnassignLecturerFromCourse(ctx context.Context, lecturerID int64, courseCodes []string) error {
	for _, c := range courseCodes {
		err := ls.courseService.UnassignLecturerFromCourse(ctx, course.UserCourseData{
			CourseCode: c,
			UserID:     lecturerID,
		})
		if err != nil {
			slog.Error("Handle error when unassigning courses courses", "error", err)
			return err
		}
	}
	return nil
}

type LecturerAvailability map[int64][]course.Availability

func (ls *LecturerService) CheckAvailabilityStatus(ctx context.Context, lectuerID int64) (LecturerAvailability, error) {
	courseAvailabilities, err := ls.courseService.GetLecturerAvailabilityForAllCourses(ctx, lectuerID)
	if err != nil {
		slog.Error("Handle error when getting eligibility", "error", err)
		return LecturerAvailability{}, err
	}
	res := LecturerAvailability{
		lectuerID: courseAvailabilities,
	}
	return res, nil
}

func (ls *LecturerService) SetActiveLecturer(ctx context.Context, lecturerID int64, courseCode string) error {
	err := ls.courseService.SetActiveLecturer(ctx, lecturerID, courseCode)
	if err != nil {
		return err
	}
	return nil
}
