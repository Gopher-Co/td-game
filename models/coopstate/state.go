package coopstate

import (
	"fmt"
	"sync"

	"github.com/gopher-co/td-game/models/config"
	"github.com/gopher-co/td-game/models/coopstate/models"
	"github.com/gopher-co/td-game/models/general"
	"github.com/gopher-co/td-game/models/ingame"
)

// Lib is a library of towers.
type Lib struct {
	// Towers is a map of towers.
	Towers map[string]*config.Tower
}

// State is a state of the game.
type State struct {
	// lib is a library of towers.
	lib *Lib
	// mu is a mutex.
	mu *sync.RWMutex
	// Map is a map.
	Map models.Map
	// Player is a player.
	Player ingame.PlayerMapState
	// Global is a global state.
	Global ingame.PlayerState
}

func NewState(lib *Lib) *State {
	return &State{lib: lib}
}

func (s *State) PutTower(x, y general.Coord, towerName, playerName string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.Player.Money < s.lib.Towers[towerName].Price {
		return fmt.Errorf("not enough money to buy tower %s", towerName)
	}

	s.Player.Money -= s.lib.Towers[towerName].Price

	s.Map.Towers = append(s.Map.Towers, models.NewTower(
		s.lib.Towers[towerName],
		x, y,
		s.Map.Path,
		playerName,
	))

	return nil
}

func (s *State) UpgradeTower(index int, playerName string) error {
	if s.Map.Towers[index].Whose != playerName {
		return fmt.Errorf("not your tower %s", playerName)
	}

	upgrade := s.Map.Towers[index].NextUpgrade()
	if upgrade == nil {
		return fmt.Errorf("no upgrade for tower %s", playerName)
	}

	if s.Player.Money < upgrade.Price {
		return fmt.Errorf("not enough money to buy upgrade %s", playerName)
	}

	s.Player.Money -= upgrade.Price
	s.Map.Towers[index].Upgrade(s.Global.LevelsComplete)

	return nil
}

func (s *State) TuneTower(index int, playerName string, aim ingame.Aim) error {
	if s.Map.Towers[index].Whose != playerName {
		return fmt.Errorf("not your tower %s", playerName)
	}

	s.Map.Towers[index].State.AimType = aim
	return nil
}

func (s *State) TurnOnTower(index int, playerName string) error {
	if s.Map.Towers[index].Whose != playerName {
		return fmt.Errorf("not your tower %s", playerName)
	}

	s.Map.Towers[index].State.IsTurnedOn = true
	return nil
}

func (s *State) TurnOffTower(index int, playerName string) error {
	if s.Map.Towers[index].Whose != playerName {
		return fmt.Errorf("not your tower %s", playerName)
	}

	s.Map.Towers[index].State.IsTurnedOn = false
	return nil
}
