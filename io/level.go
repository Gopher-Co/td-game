package io

import "github.com/gopher-co/td-game/models"

// LoadLevelConfigs loads level configs from the Levels directory.
func LoadLevelConfigs() ([]models.LevelConfig, error) {
	lcfgs, err := ReadConfigs[models.LevelConfig]("./Levels", ".lvl")
	if err != nil {
		return nil, err
	}

	return lcfgs, nil
}
