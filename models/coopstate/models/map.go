package models

import "github.com/gopher-co/td-game/models/ingame"

// Map represents a map.
type Map struct {
	// Towers is a list of towers that can be built on the map.
	Towers []*Tower
	Path   ingame.Path
}
