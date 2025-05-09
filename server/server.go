package server

import (
	"database/sql"

	"github.com/ekefan/afitlmscloud/course"
	"github.com/ekefan/afitlmscloud/internal/repository"
	"github.com/ekefan/afitlmscloud/user"
	"github.com/ekefan/afitlmscloud/user/lecturer"
	"github.com/ekefan/afitlmscloud/user/student"
	"github.com/gin-gonic/gin"
)

type Server struct {
	router        *gin.Engine
	userService   *user.UserService
	courseService *course.CourseService
}

func NewServer(dbConn *sql.DB) *Server {
	courseService := course.NewCourseService(repository.NewCourseStore(dbConn))
	studentService := student.NewStudentService(courseService, repository.NewStudentStore(dbConn))
	lecturerService := lecturer.NewLecturerService(courseService, repository.NewLecturerStore(dbConn))

	server := &Server{
		router: gin.Default(),
		userService: user.NewUserService(
			repository.NewUserStore(dbConn),
			repository.NewStudentStore(dbConn),
			studentService,
			lecturerService,
		),
		courseService: courseService,
	}
	server.registerUserRoutes()
	server.registerCourseRoutes()
	return server
}

func (s *Server) StartServer(addr ...string) {
	port := "8080"
	if len(addr) > 0 && addr[0] != "" {
		port = addr[0]
	}
	s.router.Run(":" + port)
}

// TODO: handle authentications and authorization
// TODO: handle validation of user_id or access_id for student/lecturer/user/course_codes
