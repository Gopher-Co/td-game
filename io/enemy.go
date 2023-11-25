package io

import (
	"github.com/gopher-co/td-game/models/config"
)

// LoadEnemyConfigs loads enemy configs from the Enemies directory.
func LoadEnemyConfigs() ([]config.Enemy, error) {
	ecfgs, err := ReadConfigs[config.Enemy]("./Enemies", ".enm")
	if err != nil {
		return nil, err
	}

	for k := range ecfgs {
		if err := ecfgs[k].InitImage(); err != nil {
			return nil, err
		}
	}

	return ecfgs, nil
}
