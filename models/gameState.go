package models

import "github.com/ebitenui/ebitenui"

type CurrentState int

const (
	Running CurrentState = iota
	Paused
)

type GameState struct {
	Map         Map
	TowersToBuy map[string]Tower
	Ended       bool
	State       CurrentState
	UI          *ebitenui.UI
	LastWave    int
	CurrentWave int
	Waves       []Wave
	Time        Frames
}
