package coopstate

import (
	"fmt"
	"github.com/gopher-co/td-game/models/config"
	"github.com/gopher-co/td-game/models/coopstate/models"
	"github.com/gopher-co/td-game/models/general"
	"github.com/gopher-co/td-game/models/ingame"
	"sync"
)

type Lib struct {
	Towers map[string]*config.Tower
}

type State struct {
	lib    *Lib
	mu     *sync.RWMutex
	Map    models.Map
	Player ingame.PlayerMapState
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

func (s *State) Update() {

}