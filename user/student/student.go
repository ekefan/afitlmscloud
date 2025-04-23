package student

import (
	"context"
	"log/slog"

	"github.com/ekefan/afitlmscloud/course"
	"github.com/ekefan/afitlmscloud/internal/repository"
)

type StudentService struct {
	courseService course.CourseService
	repo          repository.StudentRepository
}

func NewStudentService(courseService course.CourseService) *StudentService {
	return &StudentService{
		courseService: courseService,
	}

}
func (s *StudentService) registerCourses(ctx context.Context, studentID int64, courses []course.Course) error {
	for _, c := range courses {
		err := s.courseService.RegisterCouse(ctx, course.StudentCourseData{
			CourseCode: c.CourseCode,
			StudentID:  studentID,
		})
		if err != nil {
			slog.Error("Handle error when registering courses")
			return err
		}
	}
	return nil
}

func (s *StudentService) dropCourses(ctx context.Context, studentID int64, courses []course.Course) error {
	for _, c := range courses {
		err := s.courseService.DropCourses(ctx, course.StudentCourseData{
			CourseCode: c.CourseCode,
			StudentID:  studentID,
		})
		if err != nil {
			slog.Error("Handle error when registering courses")
			return err
		}
	}
	return nil
}

type StudentEligibility map[int64][]course.Eligibility

func (s *StudentService) checkEligibilityStatus(ctx context.Context, studentID int64, courses []course.Course) (StudentEligibility, error) {
	courseEligibilities, err := s.courseService.GetStudentEligibilityForAllCourses(ctx, studentID)
	if err != nil {
		slog.Error("Hanlde error when getting eligibility")
		return StudentEligibility{}, err
	}
	res := StudentEligibility{
		studentID: courseEligibilities,
	}
	return res, nil
}
