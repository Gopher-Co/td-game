package io

import (
	"fmt"

	"github.com/gopher-co/td-game/models"
	"github.com/gopher-co/td-game/ui"
)

var ErrUINotOnce = fmt.Errorf("there must be only 1 file in ./UI")

func LoadUIConfig() (models.UI, error) {
	uicfg, err := ReadConfigs[map[string]string]("./UI", ".json")
	if err != nil {
		return nil, err
	}
	if len(uicfg) != 1 {
		return nil, ErrUINotOnce
	}

	uis := make(models.UI)
	for k, v := range uicfg[0] {
		img, err := ui.InitImage(v)
		if err != nil {
			return nil, err
		}
		uis[k] = img
	}

	return uis, nil
}
