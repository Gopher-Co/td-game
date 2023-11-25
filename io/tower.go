package io

import (
	"github.com/gopher-co/td-game/models/config"
)

// LoadTowerConfigs loads tower configs from the Towers directory.
func LoadTowerConfigs() ([]config.Tower, error) {
	tcfgs, err := ReadConfigs[config.Tower]("./Towers", ".twr")
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
