package coopstate

import (
	"context"
	"errors"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/google/uuid"
	"google.golang.org/grpc"
)

// States represents a map of states.
type States map[string]GameHost_JoinLobbyServer

// Conns represents a map of connections.
type Conns map[string]struct{}

// Server represents a server.
type Server struct {
	// id is the server ID.
	id string
	// once is a sync.Once.
	once sync.Once
	// conns is a map of connections.
	conns Conns
	// states is a map of states.
	states States
	// levelName is the level name.
	levelName string
	// size is the size of the server.
	size int
	// gamestate is the game state.
	gamestate GameState
	// UnimplementedGameHostServer is an unimplemented game host server.
	UnimplementedGameHostServer
}

// JoinLobby joins the lobby.
func (s *Server) JoinLobby(in *JoinLobbyRequest, ss GameHost_JoinLobbyServer) error {
	s.states[in.Player.Id.Nickname] = ss
	if err := s.TakeNewConnection(in.Player.Id.Nickname, in.Lobby.Name); err != nil {
		return err
	}

	select {
	case <-ss.Context().Done():
		return nil
	}
}

// AwaitGame awaits the game.
func (s *Server) AwaitGame(ctx context.Context, _ *AwaitGameRequest) (*AwaitGameResponse, error) {
	for {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case <-time.After(time.Second):
			if len(s.conns) == s.size {
				return &AwaitGameResponse{Level: s.levelName}, nil
			}
		}
	}
}

func NewServer(levelName string, size int) (*grpc.Server, string) {
	grpcServer := grpc.NewServer()
	s := &Server{
		id:        uuid.NewString()[:8],
		conns:     make(Conns, size),
		states:    make(States, size),
		levelName: levelName,
		size:      size,
	}
	log.Println(s.id)
	RegisterGameHostServer(grpcServer, s)

	return grpcServer, s.id
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

func (s *Server) PutTower(_ context.Context, r *PutTowerRequest) (*PutTowerResponse, error) {
	var errs error
	for _, state := range s.states {
		if err := state.Send(&JoinLobbyResponse{Response: &JoinLobbyResponse_PutTower{r}}); err != nil {
			errs = errors.Join(errs, err)
		}
	}

	if errs != nil {
		return nil, errs
	}
	return &PutTowerResponse{Status: Status_OK}, nil
}

func (s *Server) StartNewWave(_ context.Context, r *StartNewWaveRequest) (*StartNewWaveResponse, error) {
	var errs error
	for _, state := range s.states {
		if err := state.Send(&JoinLobbyResponse{Response: &JoinLobbyResponse_StartNewWave{r}}); err != nil {
			errs = errors.Join(errs, err)
		}
	}

	if errs != nil {
		return nil, errs
	}
	return &StartNewWaveResponse{Status: Status_OK}, nil
}
