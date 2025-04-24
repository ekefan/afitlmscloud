package user

import (
	"context"
	"errors"
	"log/slog"
	"slices"

	db "github.com/ekefan/afitlmscloud/internal/db/sqlc"
	"github.com/go-playground/validator/v10"
)

type Roles []string

const (
	// roles
	studentRole = iota
	lecturerRole
	qaAdminRole
	courseAdminRole
)

const (
	// string rep of roles
	StudentRole     = "student"
	LecturerRole    = "lecturer"
	QaAdminRole     = "qa_admin"
	CourseAdminRole = "course_admin"
)

var (
	ErrRolesViolatesRolesPolicy = errors.New("roles violates the role policy")
	ErrNoBioMetricTemplate      = errors.New("students or lecturers must enroll with a biometric template")
)

func rolesToString(role int) string {
	switch role {
	case studentRole:
		return "student"
	case lecturerRole:
		return "lecturer"
	case qaAdminRole:
		return "qa_admin"
	case courseAdminRole:
		return "course_admin"
	}
	return ""
}

// validateUserRolesPolicy checks if roles follows the roles UserRolesPolicy
//
// UserRolesPolicy:
//
//	The policy ensures that all Roles do not contain a student and qaAdmin role
func (us *UserService) validateUserRolesPolicy(roles Roles) error {
	if slices.Contains(roles, rolesToString(studentRole)) &&
		slices.Contains(roles, rolesToString(qaAdminRole)) {
		return ErrRolesViolatesRolesPolicy
	}
	return nil
}

type EnrollmentData struct {
	Roles             []string `json:"roles"`
	BioMetricTemplate string   `json:"biometric_template,omitempty"`
	UserId            int64    `json:"user_id"`
}

// TODO: define sub-domains,
func (us *UserService) enrollUser(ctx context.Context, data EnrollmentData) error {
	if err := us.validateUserRolesPolicy(data.Roles); err != nil {
		return err
	}
	if (slices.Contains(data.Roles, rolesToString(studentRole)) ||
		slices.Contains(data.Roles, rolesToString(lecturerRole))) &&
		data.BioMetricTemplate == "" {
		return ErrNoBioMetricTemplate
	}

	_, err := us.userRepo.EnrollUser(ctx, db.EnrollUserParams{
		ID:       data.UserId,
		Roles:    data.Roles,
		Enrolled: true,
	})
	if err != nil {
		// TODO: handle anticipated errors
		return err
	}

	if slices.Contains(data.Roles, rolesToString(studentRole)) {
		slog.Info("creating a new student")
		_, err := us.studentRepo.CreateStudent(ctx, db.CreateStudentParams{
			UserID:            data.UserId,
			BiometricTemplate: data.BioMetricTemplate,
		})
		if err != nil {
			slog.Error("couldn't create a student", "details", err)
			// TODO: define error for not being able to create a user
			return err
		}
	}
	if slices.Contains(data.Roles, rolesToString(lecturerRole)) {
		slog.Info("creating a new lecturer")
		_, err := us.lecturerRepo.CreateLecturer(ctx, db.CreateLecturerParams{
			UserID:            data.UserId,
			BiometricTemplate: data.BioMetricTemplate,
		})
		if err != nil {
			slog.Error("couldn't create a lecturer", "details", err)
			return err
		}
	}
	return nil
}

var allowedRoles = map[string]bool{
	"student":      true,
	"lecturer":     true,
	"qa_admin":     true,
	"course_admin": true,
}

func rolesOnly(fl validator.FieldLevel) bool {
	roles, ok := fl.Field().Interface().([]string)
	if !ok {
		return false
	}
	for _, role := range roles {
		if !allowedRoles[role] {
			return false
		}
	}
	return true
}
