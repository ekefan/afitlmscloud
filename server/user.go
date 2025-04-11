package server

func (s *Server) registerUserRoutes() {
	ur := s.router.Group("/users")

	ur.POST("/", s.userService.CreateUser)
	ur.GET("/", s.userService.GetUser)
	ur.PUT("/", s.userService.UpdateUser)
	ur.DELETE("/", s.userService.DeleteUser)
}
