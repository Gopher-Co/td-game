package coopstate

import (
	"context"
	"errors"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/google/uuid"
	"golang.org/x/exp/maps"
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

	speedUp bool
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
				return &AwaitGameResponse{Level: s.levelName, Nicknames: maps.Keys(s.conns)}, nil
			}
		}
	}
}

// NewServer creates a new server.
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

// TakeNewConnection takes a new connection.
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

// PutTower puts a tower.
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

// StartNewWave starts a new wave.
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

func (s *Server) SpeedGameUp(_ context.Context, r *SpeedGameUpRequest) (*SpeedGameUpResponse, error) {
	if s.speedUp {
		return &SpeedGameUpResponse{Status: Status_OK}, nil
	}
	s.speedUp = true
	var errs error
	for _, state := range s.states {
		if err := state.Send(&JoinLobbyResponse{Response: &JoinLobbyResponse_SpeedUp{}}); err != nil {
			errs = errors.Join(errs, err)
		}
	}

	if errs != nil {
		return nil, errs
	}
	return &SpeedGameUpResponse{Status: Status_OK}, nil
}

func (s *Server) SlowGameDown(_ context.Context, r *SlowGameDownRequest) (*SlowGameDownResponse, error) {
	if !s.speedUp {
		return &SlowGameDownResponse{Status: Status_OK}, nil
	}
	s.speedUp = false
	var errs error
	for _, state := range s.states {
		if err := state.Send(&JoinLobbyResponse{Response: &JoinLobbyResponse_SlowDown{}}); err != nil {
			errs = errors.Join(errs, err)
		}
	}

	if errs != nil {
		return nil, errs
	}
	return &SlowGameDownResponse{Status: Status_OK}, nil
}

func (s *Server) UpgradeTower(_ context.Context, r *UpgradeTowerRequest) (*UpgradeTowerResponse, error) {
	var errs error
	for _, state := range s.states {
		if err := state.Send(&JoinLobbyResponse{Response: &JoinLobbyResponse_UpgradeTower{r}}); err != nil {
			errs = errors.Join(errs, err)
		}
	}

	if errs != nil {
		return nil, errs
	}

	return &UpgradeTowerResponse{Status: Status_OK}, nil
}

func (s *Server) ChangeTowerAimType(_ context.Context, r *ChangeTowerAimTypeRequest) (*ChangeTowerAimTypeResponse, error) {
	var errs error
	for _, state := range s.states {
		if err := state.Send(&JoinLobbyResponse{Response: &JoinLobbyResponse_TuneTower{TuneTower: &TuneTowerRequest{
			Tower: r.Tower,
			Aim:   TuneTowerRequest_Aim(r.NewAimType),
		}}}); err != nil {
			errs = errors.Join(errs, err)
		}
	}

	if errs != nil {
		return nil, errs
	}

	return &ChangeTowerAimTypeResponse{Status: Status_OK}, nil
}

func (s *Server) SellTower(_ context.Context, r *SellTowerRequest) (*SellTowerResponse, error) {
	var errs error
	for _, state := range s.states {
		if err := state.Send(&JoinLobbyResponse{Response: &JoinLobbyResponse_SellTower{r}}); err != nil {
			errs = errors.Join(errs, err)
		}
	}

	if errs != nil {
		return nil, errs
	}

	return &SellTowerResponse{Status: Status_OK}, nil
}

func (s *Server) TurnTowerOn(_ context.Context, r *TurnTowerOnRequest) (*TurnTowerOnResponse, error) {
	var errs error
	for _, state := range s.states {
		if err := state.Send(&JoinLobbyResponse{Response: &JoinLobbyResponse_TurnOn{r}}); err != nil {
			errs = errors.Join(errs, err)
		}
	}

	if errs != nil {
		return nil, errs
	}

	return &TurnTowerOnResponse{Status: Status_OK}, nil
}

func (s *Server) TurnTowerOff(_ context.Context, r *TurnTowerOffRequest) (*TurnTowerOffResponse, error) {
	var errs error
	for _, state := range s.states {
		if err := state.Send(&JoinLobbyResponse{Response: &JoinLobbyResponse_TurnOff{r}}); err != nil {
			errs = errors.Join(errs, err)
		}
	}

	if errs != nil {
		return nil, errs
	}

	return &TurnTowerOffResponse{Status: Status_OK}, nil
}
