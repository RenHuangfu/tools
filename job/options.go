package job

type ServerOption func(o *Server)

func WithEndpoint(endpoint string) ServerOption {
	return func(s *Server) {
		s.endpoint = endpoint
	}
}

func WithDebug(debug bool) ServerOption {
	return func(s *Server) {
		s.debug = debug
	}
}

func WithName(name string) ServerOption {
	return func(s *Server) {
		s.name = name
	}
}
