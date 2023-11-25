package gamestate

import (
	"fmt"
	"image"
	"image/color"

	"github.com/ebitenui/ebitenui"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/vector"

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

	speedUp bool
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
	} else if inpututil.IsMouseButtonJustPressed(ebiten.MouseButton2) {
		s.tookTower = nil
	}

	if s.State == NextWaveReady {
		return nil
	}

	s.Map.Update()
	wave := s.GameRule[s.CurrentWave]
	if wave.Ended() && !s.Map.AreThereAliveEnemies() {
		s.setStateAfterWave()
		return nil
	}

	s.updateRunning(wave)

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

	if s.tookTower != nil {
		s.drawTookImageBeforeCursor(screen)
	}

	s.UI.Draw(screen)
}

func (s *GameState) setStateAfterWave() {
	s.State = NextWaveReady
	s.Map.Enemies = []*ingame.Enemy{}
	s.Map.Projectiles = []*ingame.Projectile{}
	if s.CurrentWave == len(s.GameRule)-1 {
		s.Ended = true
	}
}

func (s *GameState) setStateAfterEnd() {
	ebiten.SetTPS(60)
}

func (s *GameState) updateRunning(wave *ingame.Wave) {
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
}

// loadUI loads UI.
func (s *GameState) loadUI(widgets general.Widgets) {
	s.UI = s.loadGameUI(widgets)
}

func (s *GameState) drawTookImageBeforeCursor(screen *ebiten.Image) {
	img := s.tookTower.Image()
	cx, cy := ebiten.CursorPosition()
	ix, iy := img.Bounds().Dx(), img.Bounds().Dy()

	if !ingame.CheckCollisionPath(general.Point{general.Coord(cx), general.Coord(cy)}, s.Map.Path) {
		vector.DrawFilledCircle(screen, float32(cx), float32(cy), s.tookTower.InitRadius, color.RGBA{0, 0, 0, 0x20}, true)
	} else {
		vector.DrawFilledCircle(screen, float32(cx), float32(cy), s.tookTower.InitRadius, color.RGBA{0xff, 0, 0, 0x20}, true)
	}
	geom := ebiten.GeoM{}
	geom.Translate(float64(cx-ix/2), float64(cy-iy/2))
	screen.DrawImage(img, &ebiten.DrawImageOptions{GeoM: geom})
}
