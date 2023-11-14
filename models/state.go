package models

import (
	"image"

	"github.com/hajimehoshi/ebiten/v2"
)

// State is an interface that represents a state.
type State interface {
	Drawable[*ebiten.Image]
	Update() error
	LoadUI(widgets Widgets)
	End() bool
	NextState() *State
}

// Widgets is a struct that represents a collection of widgets.
type Widgets struct {
	w map[string]image.Image
}
