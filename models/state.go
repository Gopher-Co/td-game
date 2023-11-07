package models

import "image"

type State interface {
	Update()
	LoadUI(widgets Widgets)
	End() bool
	NextState() *State
}

type Widgets struct {
	w map[string]image.Image
}
