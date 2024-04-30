package coopstate

import "sync"

type Conns map[string]*State

type Server struct {
	once  sync.Once
	conns Conns
}

func NewServer() *Server {
	return &Server{once: sync.Once{}}
}

func (s *Server) CurrentState() any {
	return nil
}

func (s *Server) Init() {
	s.once.Do(func() {})
}

func (s *Server) TakeNewConnection(key string, conn any) {
	s.conns[key] = conn
}

func (s *Server) RemoveConnection(key string) {
	delete(s.conns, key)
}

func (s *Server) GetConnection(key string) any {
	return s.conns[key]
}
