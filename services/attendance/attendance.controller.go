package attendance

import (
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (as *AttendanceService) RecordAttendance(ctx *gin.Context) {
	var req AttendanceSession

	if err := ctx.ShouldBindJSON(&req); err != nil {
		slog.Error("invalid request")
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "invalid course details",
			"error":   err.Error(),
		})
		return
	}

	err := as.createNewAttendanceSession(ctx, req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "failed to record a record attendance for lecture session",
			"error":   err.Error(),
		})
		return
	}
}
