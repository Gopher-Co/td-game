// Package global contains global variables used in the game.
package main

import (
	"github.com/gopher-co/td-game/models/config"
	"github.com/gopher-co/td-game/models/ingame"
	"github.com/gopher-co/td-game/replay"
)

var (
	// UI is a map of images used in the game.
	UI = ingame.UI{}

	// Maps is a map of maps used in the game.
	Maps = make(map[string]*config.Map)

	// Levels is a map of levels used in the game.
	Levels = make(map[string]*config.Level)

	// Towers is a map of towers used in the game.
	Towers = make(map[string]*config.Tower)

	// Enemies is a map of enemies used in the game.
	Enemies = make(map[string]*config.Enemy)

	Replays []*replay.Watcher
)
