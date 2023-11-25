package io

import (
	"github.com/gopher-co/td-game/models/config"
)

// LoadMapConfigs loads map configs from the Maps directory.
func LoadMapConfigs() ([]config.Map, error) {
	mcfgs, err := ReadConfigs[config.Map]("./Maps", ".map")
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
