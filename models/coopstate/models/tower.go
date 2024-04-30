package models

import (
	"github.com/gopher-co/td-game/models/config"
	"github.com/gopher-co/td-game/models/general"
	"github.com/gopher-co/td-game/models/ingame"
)

type Tower struct {
	*ingame.Tower
	Whose string
}

func NewTower(tower *config.Tower, x, y general.Coord, path ingame.Path, whose string) *Tower {
	t := ingame.NewTower(tower, general.Point{X: x, Y: y}, path)

	return &Tower{
		Tower: t,
		Whose: whose,
	}
}
