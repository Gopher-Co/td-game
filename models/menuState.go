package models

import "github.com/ebitenui/ebitenui"

// MenuState is a struct that represents the state of the menu.
type MenuState struct {
	Ended bool
	UI    *ebitenui.UI
}
