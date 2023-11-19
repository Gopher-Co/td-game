package ui

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/icza/gox/imagex/colorx"
)

func InitImage(s string) (*ebiten.Image, error) {
	clr, err := colorx.ParseHexColor(s)
	if err != nil {
		return nil, err
	}

	img := ebiten.NewImage(ebiten.WindowSize())
	img.Fill(clr)

	return img, nil
}
