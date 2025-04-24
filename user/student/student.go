package student

import (
	"context"
	"log/slog"

	"github.com/ekefan/afitlmscloud/course"
	"github.com/ekefan/afitlmscloud/internal/repository"
)

type StudentService struct {
	courseService *course.CourseService
	repo          repository.StudentRepository
}

func NewStudentService(courseService *course.CourseService, repo repository.StudentRepository) *StudentService {
	return &StudentService{
		courseService: courseService,
		repo:          repo,
	}

}

func (s *StudentService) RegisterCourses(ctx context.Context, studentID int64, courseCodes []string) error {
	for _, c := range courseCodes {
		err := s.courseService.RegisterCouse(ctx, course.StudentCourseData{
			CourseCode: c,
			StudentID:  studentID,
		})
		if err != nil {
			slog.Error("Handle error when registering courses")
			return err
		}
	}
	return nil
}

func (s *StudentService) DropCourses(ctx context.Context, studentID int64, courseCodes []string) error {
	for _, c := range courseCodes {
		err := s.courseService.DropCourses(ctx, course.StudentCourseData{
			CourseCode: c,
			StudentID:  studentID,
		})
		if err != nil {
			slog.Error("Handle error when dropping courses courses", "error", err)
			return err
		}
	}
	return nil
}

type StudentEligibility map[int64][]course.Eligibility

func (s *StudentService) CheckEligibilityStatus(ctx context.Context, studentID int64) (StudentEligibility, error) {
	courseEligibilities, err := s.courseService.GetStudentEligibilityForAllCourses(ctx, studentID)
	if err != nil {
		slog.Error("Handle error when getting eligibility", "error", err)
		return StudentEligibility{}, err
	}
	res := StudentEligibility{
		studentID: courseEligibilities,
	}
	return res, nil
}
