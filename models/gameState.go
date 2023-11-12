package models

import "github.com/ebitenui/ebitenui"

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
	Map         Map
	TowersToBuy map[string]TowerConfig
	EnemyToCall map[string]EnemyConfig
	Ended       bool
	State       CurrentState
	UI          *ebitenui.UI
	LastWave    int
	CurrentWave int
	Waves       []Wave
	Time        Frames
}
