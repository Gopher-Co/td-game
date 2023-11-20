package models

import (
	"github.com/hajimehoshi/ebiten/v2"
)

// State is an interface that represents a state.
type State interface {
	Drawable[*ebiten.Image]
	Update() error
	loadUI(widgets Widgets)
	End() bool
}

// Widgets represents a collection of widgets.
type Widgets map[string]*ebiten.Image
