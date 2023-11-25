package gamestate

import (
	"fmt"
	"image"

	"github.com/ebitenui/ebitenui"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"

	"github.com/gopher-co/td-game/models/config"
	"github.com/gopher-co/td-game/models/general"
	"github.com/gopher-co/td-game/models/ingame"
)

// CurrentState is an enum that represents the current state of the game.
type CurrentState int

const (
	// Running is the state when the game is running.
	Running CurrentState = iota

	// Paused is the state when the game is paused.
	Paused

	// NextWaveReady is the state when the next wave is ready.
	NextWaveReady
)

// GameState is a struct that represents the state of the game.
type GameState struct {
	// Map is a map of the game.
	Map *ingame.Map

	// TowersToBuy is a map of towers that can be bought.
	TowersToBuy map[string]*config.Tower

	// EnemyToCall is a map of enemies that can be called.
	EnemyToCall map[string]*config.Enemy

	// Ended is a flag that represents if the game is ended.
	Ended bool

	// State is a current state of the game.
	State CurrentState

	// UI is a UI of the game.
	UI *ebitenui.UI

	// CurrentWave is a number of the current wave.
	CurrentWave int

	// GameRule is a game rule of the game.
	GameRule ingame.GameRule

	// Time is a time of the game.
	Time general.Frames

	// PlayerMapState is a state of the player on the map.
	PlayerMapState ingame.PlayerMapState

	tookTower *config.Tower
}

// New creates a new entity of GameState.
func New(
	config *config.Level,
	maps map[string]*config.Map,
	en map[string]*config.Enemy,
	tw map[string]*config.Tower,
	w general.Widgets,
) *GameState {
	gs := &GameState{
		Map:         ingame.NewMap(maps[config.MapName]),
		TowersToBuy: tw,
		EnemyToCall: en,
		State:       NextWaveReady,
		CurrentWave: -1,
		GameRule:    ingame.NewGameRule(config.GameRule),
		PlayerMapState: ingame.PlayerMapState{
			Health: 100,
			Money:  650,
		},
	}

	gs.loadUI(w)

	return gs
}

// Update updates the state of the game.
func (s *GameState) Update() error {
	if s.Ended {
		return nil
	}

	if s.State == Paused {
		return nil
	}

	s.UI.Update()

	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButton0) && s.tookTower != nil {
		x, y := ebiten.CursorPosition()
		pos := general.Point{general.Coord(x), general.Coord(y)}

		if x < 1500 && s.PlayerMapState.Money >= s.tookTower.Price {
			if t := ingame.NewTower(s.tookTower, pos, s.Map.Path); t != nil {
				s.PlayerMapState.Money -= s.tookTower.Price
				s.tookTower = nil
				s.Map.Towers = append(s.Map.Towers, t)
			}
		}
	}

	if s.State == NextWaveReady {
		return nil
	}

	s.Map.Update()
	wave := s.GameRule[s.CurrentWave]
	if wave.Ended() && !s.Map.AreThereAliveEnemies() {
		s.State = NextWaveReady
		s.Map.Enemies = []*ingame.Enemy{}
		s.Map.Projectiles = []*ingame.Projectile{}
		if s.CurrentWave == len(s.GameRule)-1 {
			s.Ended = true
		}

		return nil
	}

	es := wave.CallEnemies()
	for _, str := range es {
		s.Map.Enemies = append(s.Map.Enemies, ingame.NewEnemy(s.EnemyToCall[str], s.Map.Path))
	}

	for _, e := range s.Map.Enemies {
		if e.State.Dead {
			if e.State.PassPath {
				s.PlayerMapState.Health = max(s.PlayerMapState.Health-e.DealDamageToPlayer(), 0)
			} else {
				s.PlayerMapState.Money += e.MoneyAward
				e.MoneyAward = 0
			}
		}
	}

	return nil
}

// End returns true if the game is ended.
func (s *GameState) End() bool {
	return s.Ended
}

// Draw draws the game on the screen.
func (s *GameState) Draw(screen *ebiten.Image) {
	subScreen := screen.SubImage(image.Rect(0, 0, 1500, 1080))
	s.Map.Draw(subScreen.(*ebiten.Image))
	if s.CurrentWave >= 0 {
		ebitenutil.DebugPrintAt(screen, fmt.Sprintf("Wave %d", s.CurrentWave+1), 0, 1900)
	}

	s.UI.Draw(screen)
}

// loadUI loads UI.
func (s *GameState) loadUI(widgets general.Widgets) {
	s.UI = s.loadGameUI(widgets)
}
