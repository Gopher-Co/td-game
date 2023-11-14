package models

import (
	"github.com/ebitenui/ebitenui"
	"github.com/hajimehoshi/ebiten/v2"
)

// MenuState is a struct that represents the state of the menu.
type MenuState struct {
	Ended bool
	UI    *ebitenui.UI
}

func NewMenuState() *MenuState {
	//TODO implement me
	panic("implement me")
}

func (m *MenuState) Draw(image *ebiten.Image) {
	//TODO implement me
	panic("implement me")
}

func (m *MenuState) Update() error {
	//TODO implement me
	panic("implement me")
}

func (m *MenuState) LoadUI(widgets Widgets) {
	//TODO implement me
	panic("implement me")
}

func (m *MenuState) End() bool {
	//TODO implement me
	panic("implement me")
}

func (m *MenuState) NextState() State {
	return NewGameState()
}
