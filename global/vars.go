package global

import "github.com/gopher-co/td-game/models"

var (
	UI      = make(models.UI)
	Maps    = make(map[string]*models.MapConfig)
	Levels  = make(map[string]*models.LevelConfig)
	Towers  = make(map[string]*models.TowerConfig)
	Enemies = make(map[string]*models.EnemyConfig)
)
