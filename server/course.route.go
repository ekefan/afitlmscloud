package server

func (s *Server) registerCourseRoutes() {
	cs := s.router.Group("/courses")

	cs.POST("/", s.courseService.CreateCourses)
	cs.GET("/", s.courseService.GetCourse)
	cs.DELETE("/", s.courseService.DeleteCourse)
}
