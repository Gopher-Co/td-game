package io

import "github.com/gopher-co/td-game/models"

func LoadTowerConfigs() ([]models.TowerConfig, error) {
	tcfgs, err := ReadConfigs[models.TowerConfig]("./Towers", ".json")
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