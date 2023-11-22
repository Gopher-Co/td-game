package io

import "github.com/gopher-co/td-game/models"

// LoadEnemyConfigs loads enemy configs from the Enemies directory.
func LoadEnemyConfigs() ([]models.EnemyConfig, error) {
	ecfgs, err := ReadConfigs[models.EnemyConfig]("./Enemies", ".enm")
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
