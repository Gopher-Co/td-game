package io

import "github.com/gopher-co/td-game/models"

// LoadTowerConfigs loads tower configs from the Towers directory.
func LoadTowerConfigs() ([]models.TowerConfig, error) {
	tcfgs, err := ReadConfigs[models.TowerConfig]("./Towers", ".twr")
	if err != nil {
		return nil, err
	}

	for k := range tcfgs {
		if err := tcfgs[k].InitImage(); err != nil {
			return nil, err
		}
		if err := tcfgs[k].ProjectileConfig.InitImage(); err != nil {
			return nil, err
		}
	}

	return tcfgs, nil
}
