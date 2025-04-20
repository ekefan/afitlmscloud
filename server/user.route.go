package server

func (s *Server) registerUserRoutes() {
	ur := s.router.Group("/users")

	ur.POST("/", s.userService.CreateUser)
	ur.GET("/:id", s.userService.GetUser)
	ur.PATCH("/:id/password", s.userService.UpdateUserPassword)
	ur.PATCH("/:id/email", s.userService.UpdateUserEmail)
	ur.DELETE("/:id", s.userService.DeleteUser)

	ur.POST("/:id/enrollment", s.userService.EnrollUser)
	ur.POST("/auth", s.userService.LoginUser)
}
