// Package global contains global variables used in the game.
package global

import "github.com/gopher-co/td-game/models"

var (
	// UI is a map of images used in the game.
	UI = make(models.UI)

	// Maps is a map of maps used in the game.
	Maps = make(map[string]*models.MapConfig)

	// Levels is a map of levels used in the game.
	Levels = make(map[string]*models.LevelConfig)

	// Towers is a map of towers used in the game.
	Towers = make(map[string]*models.TowerConfig)

	// Enemies is a map of enemies used in the game.
	Enemies = make(map[string]*models.EnemyConfig)
)
