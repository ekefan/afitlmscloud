package student

import (
	"context"
	"log/slog"

	"github.com/ekefan/afitlmscloud/course"
	db "github.com/ekefan/afitlmscloud/internal/db/sqlc"
	"github.com/ekefan/afitlmscloud/internal/repository"
)

type StudentServiceRepo interface {
	repository.EligibilityRepository
	repository.StudentRepository
}
type StudentService struct {
	courseService course.CourseService
	repo          StudentServiceRepo
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

type EligibilityResult struct {
	StudentID   int64   `json:"student_id"`
	CourseName  string  `json:"course_name"`
	CourseCode  string  `json:"course_code"`
	ELigibility float64 `json:"eligibility"`
}

func (s *StudentService) checkEligibilityStatus(ctx context.Context, studentID int64, courses []course.Course) ([]EligibilityResult, error) {
	res := []EligibilityResult{}
	for _, c := range courses {
		eligibility, err := s.repo.GetEligibility(ctx, db.GetEligibilityParams{
			StudentID: studentID,
			CourseID:  c.ID,
		})
		if err != nil {
			slog.Error("Hanlde error when getting eligibility")
			return []EligibilityResult{}, err
		}

		result := EligibilityResult{
			StudentID:   studentID,
			CourseName:  c.Name,
			CourseCode:  c.CourseCode,
			ELigibility: eligibility.Eligibility,
		}
		res = append(res, result)
	}

	return res, nil
}
