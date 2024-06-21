package menustate

import (
	"github.com/ebitenui/ebitenui"
	"github.com/hajimehoshi/ebiten/v2"

	"github.com/gopher-co/td-game/models/config"
	"github.com/gopher-co/td-game/models/coopstate"
	"github.com/gopher-co/td-game/models/general"
	"github.com/gopher-co/td-game/models/ingame"
	"github.com/gopher-co/td-game/replay"
)

// MenuState is a struct that represents the state of the menu.
type MenuState struct {
	// Replays is a list of replays.
	Replays []*replay.Watcher

	// Levels is a list of levels.
	Levels map[string]*config.Level

	// Ended is true if the menu is ended.
	Ended bool

	// UI is a UI of the menu.
	UI *ebitenui.UI

	// Next is a name of the next level.
	Next string

	// NextReplay is an index of the next replay.
	NextReplay int

	// State is a state of the player.
	State *ingame.PlayerState

	Host coopstate.GameHostClient

	Stream coopstate.GameHost_JoinLobbyClient
}

// New creates a new entity of MenuState.
func New(state *ingame.PlayerState, configs map[string]*config.Level, replays []*replay.Watcher, widgets general.Widgets) *MenuState {
	ms := &MenuState{
		Levels:     configs,
		Ended:      false,
		UI:         nil,
		Next:       "",
		Replays:    replays,
		NextReplay: -1,
		State:      state,
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
