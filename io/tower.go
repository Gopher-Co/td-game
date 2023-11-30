package io

import (
	"fmt"

	"github.com/gopher-co/td-game/models/config"
)

// LoadTowerConfigs loads tower configs from the Towers directory.
func LoadTowerConfigs() ([]config.Tower, error) {
	tcfgs, err := ReadConfigs[config.Tower]("./Towers", ".twr")
	if err != nil {
		return nil, fmt.Errorf("read tower config failed: %w", err)
	}

	for k := range tcfgs {
		if err := tcfgs[k].InitImage(); err != nil {
			return nil, fmt.Errorf("tower image init failed: %w", err)
		}
		if err := tcfgs[k].ProjectileConfig.InitImage(); err != nil {
			return nil, fmt.Errorf("projectile image init failed: %w", err)
		}
	}

	return tcfgs, nil
}
