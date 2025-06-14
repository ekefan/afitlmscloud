package enrollment

import (
	"log/slog"

	"github.com/gin-gonic/gin"
)

func (es *EnrollmentService) Enroll(ctx *gin.Context) {
	var req FastAPIEnrollInitialRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(400, gin.H{
			"message": "invalid request",
			"error":   err.Error(),
		})
		return
	}
	resp, err := es.enroll(ctx, req)
	if err != nil {
		slog.Error("failed to enroll user", "error", err)
		ctx.JSON(500, gin.H{
			"message": "failed to enroll user",
			"error":   err.Error(),
		})
		return
	}

	ctx.JSON(200, gin.H{
		"data": resp,
	})
}
