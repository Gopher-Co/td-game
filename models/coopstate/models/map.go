package models

import "github.com/gopher-co/td-game/models/ingame"

type Map struct {
	Towers []*Tower
	Path   ingame.Path
}
