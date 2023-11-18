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

func NewGameState(config *LevelConfig, en map[string]*EnemyConfig, tw map[string]*TowerConfig, w Widgets) *GameState {
	gs := &GameState{
		Map:         NewMap(GlobalMaps[config.MapName]),
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
	s.Map.Update()
	return nil
}

func (s *GameState) loadUI(widgets Widgets) {

}

func (s *GameState) End() bool {
	//TODO implement me
	panic("implement me")
}

func (s *GameState) NextState() State {
	return nil
}

func (s *GameState) Draw(screen *ebiten.Image) {
	s.Map.Draw(screen)
}
