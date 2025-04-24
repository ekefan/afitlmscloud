package server

import (
	"database/sql"

	"github.com/ekefan/afitlmscloud/course"
	"github.com/ekefan/afitlmscloud/internal/repository"
	"github.com/ekefan/afitlmscloud/user"
	"github.com/ekefan/afitlmscloud/user/student"
	"github.com/gin-gonic/gin"
)

type Server struct {
	router      *gin.Engine
	userService *user.UserService
}

func NewServer(dbConn *sql.DB) *Server {
	courseService := course.NewCourseService(repository.NewCourseStore(dbConn))
	studentService := student.NewStudentService(courseService, repository.NewStudentStore(dbConn))

	server := &Server{
		router: gin.Default(),
		userService: user.NewUserService(
			repository.NewUserStore(dbConn),
			repository.NewStudentStore(dbConn),
			studentService,
			repository.NewLecturerStore(dbConn),
		),
	}
	server.registerUserRoutes()
	return server
}

func (s *Server) StartServer(addr ...string) {
	port := "8080"
	if len(addr) > 0 && addr[0] != "" {
		port = addr[0]
	}
	s.router.Run(":" + port)
}
