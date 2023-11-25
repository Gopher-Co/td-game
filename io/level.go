package io

import (
	"github.com/gopher-co/td-game/models/config"
)

// LoadLevelConfigs loads level configs from the Levels directory.
func LoadLevelConfigs() ([]config.Level, error) {
	lcfgs, err := ReadConfigs[config.Level]("./Levels", ".lvl")
	if err != nil {
		return nil, err
	}

	return lcfgs, nil
}
