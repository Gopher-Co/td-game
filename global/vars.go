package global

import "github.com/gopher-co/td-game/models"

var (
	GlobalUI      = make(models.UI)
	GlobalMaps    = make(map[string]*models.MapConfig)
	GlobalLevels  = make(map[string]*models.LevelConfig)
	GlobalTowers  = make(map[string]*models.TowerConfig)
	GlobalEnemies = make(map[string]*models.EnemyConfig)
)
