package io

import (
	"fmt"

	"github.com/gopher-co/td-game/models/config"
)

// ErrUINotOnce is returned when there are more than 1 files in ./UI.
var ErrUINotOnce = fmt.Errorf("there must be only 1 file in ./UI")

// LoadUIConfig loads UI configs from the UI directory.
func LoadUIConfig() (config.UI, error) {
	uicfgs, err := ReadConfigs[config.UI]("./UI", ".ui")
	if err != nil {
		return config.UI{}, fmt.Errorf("ui config read error: %w", err)
	}
	if len(uicfgs) != 1 {
		return config.UI{}, ErrUINotOnce
	}

	return uicfgs[0], err
}
