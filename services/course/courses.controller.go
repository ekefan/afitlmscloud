package course

import (
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (csvc *CourseService) CreateCourses(ctx *gin.Context) {
	var req Course

	if err := ctx.ShouldBindJSON(&req); err != nil {
		slog.Error("invalid request")
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "invalid course details",
			"error":   err.Error(),
		})
		return
	}

	newCourse, err := csvc.createCourses(ctx, req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "failed to create a new course",
			"error":   err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, newCourse)
}

func (csvc *CourseService) GetCourse(ctx *gin.Context) {
	courseCode, ok := ctx.GetQuery("course_code")
	if !ok {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "invalid query param",
			"error":   "no query with key: 'course_code' found",
		})
		return
	}

	course, err := csvc.getCourse(ctx, courseCode)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "failed to create a new course",
			"error":   err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, course)
}

func (csvc *CourseService) DeleteCourse(ctx *gin.Context) {
	courseCode, ok := ctx.GetQuery("course_code")
	if !ok {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "invalid query param",
			"error":   "no query with key: 'course_code' found",
		})
		return
	}
	err := csvc.deleteCourse(ctx, courseCode)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "failed to delete the course",
			"error":   err.Error(),
		})
		return
	}
	ctx.Status(http.StatusAccepted)
}
