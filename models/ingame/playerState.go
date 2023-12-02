package ingame

import (
	"fmt"

	"github.com/gopher-co/td-game/models/config"
)

// PlayerState is a struct that represents a state of the player.
type PlayerState struct {
	// LevelsComplete is a set of levels that player has completed.
	LevelsComplete map[int]struct{} `json:"levels_complete"`
}

func (ps *PlayerState) Valid(levels map[string]*config.Level) error {
	mn, mx := 0, len(levels)-1
	for k := range ps.LevelsComplete {
		mn = min(mn, k)
		mx = max(mx, k)
	}

	if mn < 0 {
		return fmt.Errorf("min level is less than zero: %v", mn)
	}

	if mx > len(levels)-1 {
		return fmt.Errorf("incorrect max level: %v", mx)
	}

	return nil
}
