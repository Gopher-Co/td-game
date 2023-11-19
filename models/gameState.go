package models

import (
	"github.com/ebitenui/ebitenui"
	"github.com/hajimehoshi/ebiten/v2"
)

// CurrentState is an enum that represents the current state of the game.
type CurrentState int

const (
	// Running is the state when the game is running.
	Running CurrentState = iota
	// Paused is the state when the game is paused.
	Paused
)

// GameState is a struct that represents the state of the game.
type GameState struct {
	Map         *Map
	TowersToBuy map[string]*TowerConfig
	EnemyToCall map[string]*EnemyConfig
	Ended       bool
	State       CurrentState
	UI          *ebitenui.UI
	LastWave    int
	CurrentWave int
	GameRule    GameRule
	Time        Frames
}

func NewGameState(config *LevelConfig, maps map[string]*MapConfig, en map[string]*EnemyConfig, tw map[string]*TowerConfig, w Widgets) *GameState {
	gs := &GameState{
		Map:         NewMap(maps[config.MapName]),
		TowersToBuy: tw,
		EnemyToCall: en,
		Ended:       false,
		State:       Paused,
		UI:          nil, // loadUI loads it
		LastWave:    0,
		CurrentWave: -1,
		GameRule:    NewGameRule(config.GameRule),
		Time:        0,
	}

	gs.loadUI(w)
	gs.LastWave = len(gs.GameRule) - 1 // is it needed??

	return gs
}

func (s *GameState) Update() error {
	if s.Ended {
		return nil
	}

	if s.State == Paused {
		return nil
	}

	s.Map.Update()
	wave := s.GameRule[s.CurrentWave]
	if wave.Ended() && !s.Map.AreThereAliveEnemies() {
		s.State = Paused
		s.Map.Enemies = []*Enemy{}
		s.Map.Projectiles = []*Projectile{}
		if s.CurrentWave == len(s.GameRule)-1 {
			s.Ended = true
		}

		return nil
	}

	es := wave.CallEnemies()
	for _, str := range es {
		s.Map.Enemies = append(s.Map.Enemies, NewEnemy(s.EnemyToCall[str], s.Map.Path))
	}

	return nil
}

func (s *GameState) loadUI(widgets Widgets) {

}

func (s *GameState) End() bool {
	return s.Ended
}

func (s *GameState) NextState() State {
	return nil
}

func (s *GameState) Draw(screen *ebiten.Image) {
	s.Map.Draw(screen)
}
