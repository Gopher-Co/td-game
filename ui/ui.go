// Package ui provides tools for working with the user interface.
package ui

import (
	"fmt"
	"image/png"
	"os"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/icza/gox/imagex/colorx"
)

// InitImage initializes an image.
func InitImage(s string) (*ebiten.Image, error) {
	img, err := InitColor(s)
	if err == nil {
		return img, nil
	}

	return InitPNG("assets/" + s)
}

// InitColor initializes an image with a color.
func InitColor(s string) (*ebiten.Image, error) {
	clr, err := colorx.ParseHexColor(s)
	if err != nil {
		return nil, fmt.Errorf("invalid color %s: %w", s, err)
	}

	img := ebiten.NewImage(1, 1)
	img.Fill(clr)

	return img, nil
}

// InitPNG initializes an image with a PNG file.
func InitPNG(s string) (*ebiten.Image, error) {
	f, err := os.Open(s + ".png")
	if err != nil {
		return nil, fmt.Errorf("open png image failed: %w", err)
	}
	defer f.Close()

	img, err := png.Decode(f)
	if err != nil {
		return nil, fmt.Errorf("png decode failed: %w", err)
	}

	return ebiten.NewImageFromImage(img), nil
}
