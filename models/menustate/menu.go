package menustate

import (
	"github.com/ebitenui/ebitenui"
	"github.com/hajimehoshi/ebiten/v2"

	"github.com/gopher-co/td-game/models/config"
	"github.com/gopher-co/td-game/models/general"
	"github.com/gopher-co/td-game/replay"
)

// MenuState is a struct that represents the state of the menu.
type MenuState struct {
	Replays    []*replay.Watcher
	Levels     map[string]*config.Level
	Ended      bool
	UI         *ebitenui.UI
	Next       string
	NextReplay int
}

// New creates a new entity of MenuState.
func New(configs map[string]*config.Level, replays []*replay.Watcher, widgets general.Widgets) *MenuState {
	ms := &MenuState{
		Levels:     configs,
		Ended:      false,
		UI:         nil,
		Next:       "",
		Replays:    replays,
		NextReplay: -1,
	}
	ms.loadUI(widgets)

	return ms
}

// Draw draws the menu.
func (m *MenuState) Draw(image *ebiten.Image) {
	m.UI.Draw(image)
}

// Update updates the menu.
func (m *MenuState) Update() error {
	m.UI.Update()
	return nil
}

// End returns true if the menu is ended.
func (m *MenuState) End() bool {
	return m.Ended
}

// loadUI loads the UI.
func (m *MenuState) loadUI(widgets general.Widgets) {
	m.UI = m.loadMainMenuUI(widgets)
}
