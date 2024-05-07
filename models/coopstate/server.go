package coopstate

import (
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"log"
	"sync"
)

import "google.golang.org/grpc"

type States map[string]*State
type Conns map[string]struct{}

type Server struct {
	id        string
	once      sync.Once
	conns     Conns
	states    States
	levelName string
	size      int
	UnimplementedServerServer
}

func (s *Server) JoinLobby(ctx context.Context, in *JoinLobbyRequest) (*JoinLobbyResponse, error) {
	if err := s.TakeNewConnection(in.Player.Id.Nickname, in.Lobby.Name); err != nil {
		return &JoinLobbyResponse{Status: Status_ERROR}, nil
	}

	return &JoinLobbyResponse{Status: Status_OK}, nil
}

func NewServer(levelName string, size int) *grpc.Server {
	grpcServer := grpc.NewServer()
	s := &Server{
		id:        uuid.NewString(),
		conns:     make(Conns, size),
		states:    make(States, size),
		levelName: levelName,
		size:      size,
	}
	log.Println(s.id)
	RegisterServerServer(grpcServer, s)

	return grpcServer
}

func (s *Server) TakeNewConnection(nick, id string) error {
	if s.id != id {
		return errors.New("id mismatch")
	}

	if _, ok := s.conns[nick]; ok {
		return fmt.Errorf("nickname %s already taken", nick)
	}

	s.conns[nick] = struct{}{}

	return nil
}
