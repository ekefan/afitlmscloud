package server

func (s *Server) registerUserRoutes() {
	ur := s.router.Group("/users")
	student := s.router.Group("/users/students")
	lecturer := s.router.Group("/users/lecturers")
	ur.POST("/", s.userService.CreateUser)
	ur.GET("/:id", s.userService.GetUser)
	ur.PUT("/:id/password", s.userService.UpdateUserPassword)
	ur.PUT("/:id/email", s.userService.UpdateUserEmail)
	ur.DELETE("/:id", s.userService.DeleteUser)

	ur.POST("/auth", s.userService.LoginUser)
	ur.POST("/:id/enrollments", s.userService.EnrollUser)

	student.POST("/:id/course_registrations", s.userService.RegisterCourses)
	student.GET("/:id/eligibility", s.userService.CheckEligibilityForAllRegisteredCourses)
	student.DELETE("/:id/course_registrations/:course_code", s.userService.DropCoursesRegisteredByStudent)

	lecturer.POST("/:id/course_assignments", s.userService.AssignCourses)
	lecturer.GET("/:id/availability", s.userService.CheckAvailabilityForAllAssignedCourses)
	lecturer.DELETE("/:id/course_assignments/:course_code", s.userService.UnassignCourses)
	lecturer.PUT("/:id/course_assignments/:course_code", s.userService.SetActiveLecturer)
}
