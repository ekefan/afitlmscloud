package user

import (
	"database/sql"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"github.com/lib/pq"
)

func init() {
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("rolesonly", rolesOnly)
	}
}

type EnrollmentReq struct {
	Roles             []string `json:"roles" binding:"required,min=1,rolesonly"`
	BioMetricTemplate string   `json:"biometric_template,omitempty"`
}

type AssignCourseReq struct {
	CourseCodes []string `json:"course_codes" binding:"required"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type RegisterCourseReq struct {
	CourseCodes []string `json:"course_codes" binding:"required"`
}

func (us *UserService) EnrollUser(ctx *gin.Context) {
	var req EnrollmentReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		slog.Error("invalid request")
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "invalid enrollment data",
			"error":   err.Error(),
		})
		return
	}

	userIDStr, ok := ctx.Params.Get("id")
	if !ok {
		slog.Error("invalid request")
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "invalid enrollment data",
		})
		return
	}
	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		slog.Error("invalid user ID format", "err", err)
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "invalid user ID format",
		})
		return
	}
	enrollmentData := EnrollmentData{
		Roles:             req.Roles,
		BioMetricTemplate: req.BioMetricTemplate,
		UserId:            userID,
	}
	if err := us.enrollUser(ctx, enrollmentData); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			slog.Error("bad request", "details", err)
			ctx.JSON(http.StatusBadRequest, gin.H{
				"message": "invalid enrollment data",
			})
			return
		}
		if errors.Is(err, ErrRolesViolatesRolesPolicy) {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": "invalid enrollment data, student can not be quality assurance admin",
			})
			return
		}
		if errors.Is(err, ErrNoBioMetricTemplate) {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": "invalid enrollment data, student or lecturer must have biometric template field",
			})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "internal server error",
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"message": "user enrolled successfully",
	})
}

func (us *UserService) LoginUser(ctx *gin.Context) {
	var data LoginRequest
	if err := ctx.ShouldBindJSON(&data); err != nil {
		slog.Error("invalid request")
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "invalid login data",
			"error":   err.Error(),
		})
		return
	}
	user, err := us.loginUser(ctx, data)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			slog.Error("bad request", "details", err)
			ctx.JSON(http.StatusNotFound, gin.H{
				"message": "no user with that email exist data",
			})
			return
		}
		if errors.Is(err, ErrIncorrectPassword) {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"message": "incorrect password",
			})
			return
		}
		slog.Error("failed to login user", "error", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "internal server error",
		})
		return
	}

	ctx.JSON(http.StatusOK, user)
}

func (us *UserService) UpdateUserEmail(ctx *gin.Context) {
	userId, err := strconv.Atoi(ctx.Param("id"))
	if userId < 1 || err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "valid integer user id is required",
		})
		return
	}

	var req ChangeUserEmailReq
	if err := ctx.BindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid request, but old and new emails must be provided",
		})
		return
	}

	user, err := us.changeUserEmail(ctx, ChangeUserEmailData{
		UserID:   int64(userId),
		NewEmail: req.NewEmail,
		OldEmail: req.OldEmail,
	})
	if err != nil {
		slog.Error("Unhandled error here", "detals", err)
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": "no user exists with such user id",
			})
			return
		}
		if pqErr, ok := err.(*pq.Error); ok {
			switch pqErr.Code.Name() {
			case "unique_violation":
				column := ""
				if pqErr.Constraint == "users_email_key" {
					column = "email"
				}
				if pqErr.Constraint == "users_sch_id_key" {
					column = "sch_id"
				}
				msg := fmt.Sprintf("%v already exists", column)
				ctx.JSON(http.StatusBadRequest, gin.H{
					"error": msg,
				})
			default:
				slog.Error("Unhandled pq error", "details", pqErr)
			}
			return
		}
		return
	}
	ctx.JSON(http.StatusOK, user)
}

func (us *UserService) UpdateUserPassword(ctx *gin.Context) {

	userId, err := strconv.Atoi(ctx.Param("id"))
	if userId < 1 || err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "valid integer user id is required",
		})
		return
	}
	var req ChangeUserPasswordReq
	if err := ctx.BindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid request",
		})
		return
	}

	// TODO: Hash passwords if provided,
	user, err := us.changeUserPassword(ctx, ChangeUserPasswordData{
		UserId:      int64(userId),
		OldPassword: req.OldPassword,
		NewPassword: req.NewPassword,
	})
	if err != nil {
		slog.Error("Unhandled error here", "detals", err)
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": "no user exists with such user id",
			})
			return
		}
		if pqErr, ok := err.(*pq.Error); ok {
			switch pqErr.Code.Name() {
			case "unique_violation":
				column := ""
				if pqErr.Constraint == "users_email_key" {
					column = "email"
				}
				if pqErr.Constraint == "users_sch_id_key" {
					column = "sch_id"
				}
				msg := fmt.Sprintf("%v already exists", column)
				ctx.JSON(http.StatusBadRequest, gin.H{
					"error": msg,
				})
			default:
				slog.Error("Unhandled pq error", "details", pqErr)
			}
			return
		}
		return
	}
	ctx.JSON(http.StatusOK, user)
}

func (us *UserService) RegisterCourses(ctx *gin.Context) {
	// TODO: write middleware to check student role...
	userId, err := strconv.Atoi(ctx.Param("id"))
	if userId < 1 || err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "valid integer user id is required",
		})
		return
	}

	var req RegisterCourseReq
	if err := ctx.BindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid request",
		})
		return
	}
	err = us.studentService.RegisterCourses(ctx, int64(userId), req.CourseCodes)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
}

func (us *UserService) CheckEligibilityForAllRegisteredCourses(ctx *gin.Context) {
	userId, err := strconv.Atoi(ctx.Param("id"))
	if userId < 1 || err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "valid integer user id is required",
		})
		return
	}

	studentEligibility, err := us.studentService.CheckEligibilityStatus(ctx, int64(userId))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	ctx.JSON(http.StatusOK, studentEligibility)
}

func (us *UserService) DropCoursesRegisteredByStudent(ctx *gin.Context) {
	// TODO: how do I validate that the course code is valid
	userId, idErr := strconv.Atoi(ctx.Param("id"))
	courseCode := ctx.Param("course_code")
	if userId < 1 || idErr != nil || courseCode == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "valid integer user id is required",
		})
		return
	}
	err := us.studentService.DropCourses(ctx, int64(userId), []string{courseCode})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error,
		})
	}
	ctx.Status(http.StatusAccepted)
}

func (us *UserService) AssignCourses(ctx *gin.Context) {
	userId, err := strconv.Atoi(ctx.Param("id"))
	if userId < 1 || err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "valid integer user id is required",
		})
		return
	}

	var req AssignCourseReq
	if err := ctx.BindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid request",
		})
		return
	}

	fmt.Println(len(req.CourseCodes), req.CourseCodes)
	err = us.lecturerService.AssignLecturerToCourse(ctx, int64(userId), req.CourseCodes)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	ctx.Status(http.StatusAccepted)
}

func (us *UserService) CheckAvailabilityForAllAssignedCourses(ctx *gin.Context) {
	userId, err := strconv.Atoi(ctx.Param("id"))
	if userId < 1 || err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "valid integer user id is required",
		})
		return
	}
	lecturerAvailability, err := us.lecturerService.CheckAvailabilityStatus(ctx, int64(userId))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	ctx.JSON(http.StatusOK, lecturerAvailability)
}

func (us *UserService) UnassignCourses(ctx *gin.Context) {
	userId, idErr := strconv.Atoi(ctx.Param("id"))
	courseCode := ctx.Param("course_code")
	if userId < 1 || idErr != nil || courseCode == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "valid integer user id is required",
		})
		return
	}

	err := us.lecturerService.UnassignLecturerFromCourse(ctx, int64(userId), []string{courseCode})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	ctx.Status(http.StatusAccepted)
}

func (us *UserService) SetActiveLecturer(ctx *gin.Context) {
	userId, idErr := strconv.Atoi(ctx.Param("id"))
	courseCode := ctx.Param("course_code")
	if userId < 1 || idErr != nil || courseCode == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "valid integer user id is required",
		})
		return
	}

	err := us.lecturerService.SetActiveLecturer(ctx, int64(userId), courseCode)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error,
		})
	}
	ctx.Status(http.StatusAccepted)
}
