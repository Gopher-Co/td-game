package models

import "github.com/gopher-co/td-game/models/ingame"

// Map represents a map.
type Map struct {
	Towers []*Tower
	Path   ingame.Path
}
