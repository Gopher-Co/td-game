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

func NewMenuState(widgets Widgets) *MenuState {
	ms := &MenuState{
		Ended: false,
		UI:    nil,
	}
	ms.loadUI(widgets)

	return ms
}

func (m *MenuState) Draw(image *ebiten.Image) {
	//TODO implement me
	panic("implement me")
}

func (m *MenuState) Update() error {
	//TODO implement me
	panic("implement me")
}

func (m *MenuState) loadUI(widgets Widgets) {

}

func (m *MenuState) End() bool {
	//TODO implement me
	panic("implement me")
}

func (m *MenuState) NextState() State {
	return nil
}
