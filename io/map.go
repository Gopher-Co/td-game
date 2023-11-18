package io

import "github.com/gopher-co/td-game/models"

func LoadMapConfigs() ([]models.MapConfig, error) {
	mcfgs, err := ReadConfigs[models.MapConfig]("./Maps", ".json")
	if err != nil {
		return nil, err
	}

	for k := range mcfgs {
		if err := mcfgs[k].InitImage(); err != nil {
			return nil, err
		}
	}

	return mcfgs, nil
}