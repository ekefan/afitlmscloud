package user

import (
	"database/sql"
	"errors"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
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
