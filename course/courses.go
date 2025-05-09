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
	ErrFailedToSetActiveLecturer                 = errors.New("failed to set active lecturer")
	ErrFailedToGetStudentEligibilityList         = errors.New("failed to get student elibility list")
	// ErrActiveLecturerAlreadySetForThisCourse  = errors.New("failed to set active lecturer")
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

type CourseEligbilityListResp struct {
	CourseData  CourseData
	StudentData []db.GetAllStudentsEligibilityForCourseRow
}

type CourseData struct {
	Name       string `json:"name" binding:"required"`
	Faculty    string `json:"faculty" binding:"required"`
	Level      string `json:"level" binding:"required"`
	Department string `json:"department" binding:"required"`
}

func (csvc *CourseService) GetStudentEligibilityList(ctx context.Context, courseCode string) (CourseEligbilityListResp, error) {
	response := CourseEligbilityListResp{}
	dbres, err := csvc.repo.GetAllStudentsEligibilityForCourse(ctx, courseCode)
	if err != nil {
		slog.Error("failed to get student eligibility list", "error", err)
		return response, ErrFailedToGetStudentEligibilityList
	}
	courseDetails, err := csvc.repo.GetCourseMetaData(ctx, courseCode)
	if err != nil {
		slog.Error("failed to get student eligibility list", "error", err)
		return response, ErrFailedToGetStudentEligibilityList
	}

	response = CourseEligbilityListResp{
		CourseData: CourseData{
			Name:       courseDetails.Name,
			Faculty:    courseDetails.Faculty,
			Level:      courseDetails.Level,
			Department: courseDetails.Department,
		},
		StudentData: dbres,
	}
	return response, nil
}

type Course struct {
	CourseData
	CourseCode string `json:"course_code" binding:"required"`
}

func (csvc *CourseService) createCourses(ctx context.Context, course Course) (Course, error) {
	dbCourse, err := csvc.repo.CreateCourse(ctx, db.CreateCourseParams{
		Name:       course.Name,
		Department: course.Department,
		Faculty:    course.Faculty,
		Level:      course.Level,
		CourseCode: course.CourseCode,
	})
	if err != nil {
		slog.Error("failed to create a new course", "error", err)
		return Course{}, err
	}
	return Course{
		CourseData: CourseData{
			Name:       dbCourse.Name,
			Faculty:    dbCourse.Faculty,
			Level:      dbCourse.Level,
			Department: dbCourse.Department,
		},
		CourseCode: dbCourse.CourseCode,
	}, nil

}

func (csvc *CourseService) getCourse(ctx context.Context, courseCode string) (Course, error) {
	dbCourse, err := csvc.repo.GetCourse(ctx, courseCode)
	if err != nil {
		slog.Error("failed to get a course", "error", err)
		return Course{}, err
	}

	return Course{
		CourseData: CourseData{
			Name:       dbCourse.Name,
			Faculty:    dbCourse.Faculty,
			Level:      dbCourse.Level,
			Department: dbCourse.Department,
		},
		CourseCode: dbCourse.CourseCode,
	}, nil
}

func (csvc *CourseService) deleteCourse(ctx context.Context, courseCode string) error {
	result, err := csvc.repo.DeleteCourse(ctx, courseCode)
	if err != nil {
		slog.Error("failed to delete course", "error", err)
		return nil
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil || rowsAffected == 0 {
		slog.Error("failed to delete course", "error", err)
		return nil
	}

	return nil
}
