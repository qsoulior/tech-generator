package httpserver

import "time"

type OptionFunc func(*Server)

func ReadTimeout(timeout time.Duration) OptionFunc {
	return func(s *Server) {
		s.server.ReadTimeout = timeout
	}
}

func WriteTimeout(timeout time.Duration) OptionFunc {
	return func(s *Server) {
		s.server.WriteTimeout = timeout
	}
}
