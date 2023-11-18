package io

import "github.com/gopher-co/td-game/models"

func LoadLevelConfigs() ([]models.LevelConfig, error) {
	lcfgs, err := ReadConfigs[models.LevelConfig]("./Levels", ".json")
	if err != nil {
		return nil, err
	}

	return lcfgs, nil
}
