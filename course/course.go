package course

import (
	"context"
	"errors"
	"fmt"

	db "github.com/ekefan/afitlmscloud/internal/db/sqlc"
	"github.com/ekefan/afitlmscloud/internal/repository"
)

var (
	ErrFailedToRegisterCourse                    = errors.New("failed to register course")
	ErrFailedToDropCourse                        = errors.New("failed to drop course")
	ErrFailedToGetStudentEligbilityForAllCourses = errors.New("faile to get student eligibility for all courses")
)

type Course struct {
	ID         int64
	Name       string
	CourseCode string
	Faculty    string
	Department string
	Level      string
}

type CourseService struct {
	repo repository.CourseRepository
}

func NewCourseService(courseRepository repository.CourseRepository) *CourseService {
	return &CourseService{
		repo: courseRepository,
	}
}

type StudentCourseData struct {
	StudentID  int64
	CourseCode string
}

func (csvc *CourseService) RegisterCouse(ctx context.Context, data StudentCourseData) error {
	if err := csvc.repo.RegisterCourse(ctx, db.RegisterCourseParams{
		StudentID:  data.StudentID,
		CourseCode: data.CourseCode,
	}); err != nil {
		fmt.Println("failed to perform course operation", err)
		return ErrFailedToRegisterCourse
	}
	return nil
}

func (csvc *CourseService) DropCourses(ctx context.Context, data StudentCourseData) error {
	res, err := csvc.repo.DropCourse(ctx, db.DropCourseParams{
		StudentID:  data.StudentID,
		CourseCode: data.CourseCode,
	})
	numOfRowsAffected, rerr := res.RowsAffected()
	if rerr != nil {
		fmt.Println("failed to get the number of rows affected")
		return rerr
	}
	if err != nil || numOfRowsAffected == 0 {
		fmt.Println("failed to perform course operation", err)
		return ErrFailedToDropCourse
	}
	return nil
}

type Eligibility struct {
	CourseCode  string
	CourseName  string
	Eligibility float64
}

func (csvc *CourseService) GetStudentEligibilityForAllCourses(ctx context.Context, studentID int64) ([]Eligibility, error) {
	res, err := csvc.repo.GetStudentEligibilityForAllCourses(ctx, studentID)
	if err != nil {
		fmt.Println("failed to get student eligibility for all courses")
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
