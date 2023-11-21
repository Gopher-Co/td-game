package models

import (
	"github.com/hajimehoshi/ebiten/v2"
)

// State is an interface that represents a state.
type State interface {
	// Drawable is an interface that represents a drawable object.
	Drawable[*ebiten.Image]

	// Update updates the state.
	Update() error

	// loadUI loads the UI.
	loadUI(widgets Widgets)

	// End returns true if the state is ended.
	End() bool
}

// Widgets represents a collection of widgets.
type Widgets map[string]*ebiten.Image
