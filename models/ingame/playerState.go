package ingame

import (
	"fmt"

	"github.com/gopher-co/td-game/models/config"
)

// PlayerState is a struct that represents a state of the player.
type PlayerState struct {
	// LevelsComplete is a set of levels that player has completed.
	LevelsComplete map[string]struct{} `json:"levels_complete"`
}

// Valid returns an error if the player's state is not valid.
func (ps *PlayerState) Valid(levels map[string]*config.Level) error {
	for k := range ps.LevelsComplete {
		if _, ok := levels[k]; !ok {
			return fmt.Errorf("level %v doesn't exist", k)
		}
	}

	return nil
}
