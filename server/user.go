package server

func (s *Server) registerUserRoutes() {
	ur := s.router.Group("/users")

	ur.POST("/", s.userService.CreateUser)
	ur.GET("/:id", s.userService.GetUser)
	ur.PUT("/:id", s.userService.UpdateUser)
	ur.DELETE("/:id", s.userService.DeleteUser)
}
