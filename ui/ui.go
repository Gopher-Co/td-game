// Package ui provides tools for working with the user interface.
package ui

import (
	"image/color"
	"image/png"
	"os"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/icza/gox/imagex/colorx"
)

var (
	// TextColor is a color of the text.
	TextColor color.Color
)

// InitImage initializes an image.
func InitImage(s string) (*ebiten.Image, error) {
	return initColor(s)
}

// initColor initializes an image with a color.
func initColor(s string) (*ebiten.Image, error) {
	clr, err := colorx.ParseHexColor(s)
	if err != nil {
		return nil, err
	}

	img := ebiten.NewImage(ebiten.WindowSize())
	img.Fill(clr)

	return img, nil
}

// initPNG initializes an image with a PNG file.
func initPNG(s string) (*ebiten.Image, error) {
	f, err := os.Open(s)
	if err != nil {
		return nil, err
	}

	img, err := png.Decode(f)
	if err != nil {
		return nil, err
	}

	return ebiten.NewImageFromImage(img), nil
}
