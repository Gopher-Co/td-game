package models

import "image"

// State is an interface that represents a state.
type State interface {
	Update()
	LoadUI(widgets Widgets)
	End() bool
	NextState() *State
}

// Widgets is a struct that represents a collection of widgets.
type Widgets struct {
	w map[string]image.Image
}
