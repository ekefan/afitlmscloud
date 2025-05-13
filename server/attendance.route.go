package server

func (s *Server) registerAttendanceRoutes() {
	cs := s.router.Group("/attendances")

	cs.POST("/", s.attendanceService.RecordAttendance)
}
