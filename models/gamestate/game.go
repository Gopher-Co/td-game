package gamestate

import (
	"context"
	"fmt"
	"image"
	"image/color"
	"log"
	"os"
	"runtime"
	"slices"
	"time"

	"github.com/ebitenui/ebitenui"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/vector"

	"github.com/gopher-co/td-game/models/config"
	"github.com/gopher-co/td-game/models/general"
	"github.com/gopher-co/td-game/models/ingame"
	"github.com/gopher-co/td-game/replay"
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

	chosenTower *ingame.Tower

	speedUp bool

	cancel context.CancelFunc

	rw *replay.Watcher
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
		rw: &replay.Watcher{Actions: make([]replay.Action, 0, 2500)},
	}

	ctx, cancel := context.WithCancel(context.Background())
	gs.UI = gs.loadGameUI(ctx, w)
	gs.cancel = cancel

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

	s.Time++

	// if clicked on tower
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButton0) {
		s.rightSidebarHandle()
	}

	// put tower on the map
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButton0) && s.tookTower != nil {
		s.putTowerHandler()
	} else if inpututil.IsMouseButtonJustPressed(ebiten.MouseButton2) {
		s.tookTower = nil
	}

	s.UI.Update()

	if s.State == NextWaveReady {
		return nil
	}

	s.Map.Update()

	wave := s.GameRule[s.CurrentWave]
	if wave.Ended() && !s.Map.AreThereAliveEnemies() {
		s.setStateAfterWave()
		if s.CurrentWave == len(s.GameRule)-1 {
			s.Ended = true
			s.setStateAfterEnd()
		}
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
	if s.Ended == true {
		return
	}
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
}

func (s *GameState) setStateAfterEnd() {
	defer runtime.GC()

	ebiten.SetTPS(60)
	c := s.UI.Container.Children()
	for k := range c {
		c[k] = nil
	}

	s.UI.Container = nil
	s.UI = nil
	s.cancel()
	s.cancel = nil

	f, err := os.OpenFile("./Replays/replay_"+time.Now().Truncate(0).Format("2006-01-02T15_04_05")+".json", os.O_CREATE|os.O_WRONLY|os.O_SYNC, 0o666)
	if err != nil {
		log.Println("replay file wasn't created:", err)
		return
	}

	if err := s.rw.Write(f); err != nil {
		log.Println("couldn't save replay:", err)
		return
	}
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

func (s *GameState) rightSidebarHandle() {
	x, _ := ebiten.CursorPosition()
	if x <= 1500 {
		ts := slices.Clone(s.Map.Towers)
		slices.Reverse(ts)
		b := true
		for _, t := range ts {
			if b && t.IsClicked() {
				t.Chosen = true
				s.chosenTower = t
				b = !b
				s.updateTowerUI(t)
				s.showTowerInfoMenu()
			} else {
				t.Chosen = false
			}
		}
		if b {
			s.showTowerMenu()
		}
	}
}

func (s *GameState) putTowerHandler() bool {
	tt := s.tookTower
	s.tookTower = nil

	x, y := ebiten.CursorPosition()
	pos := general.Point{general.Coord(x), general.Coord(y)}

	if x < 1500 && s.PlayerMapState.Money >= tt.Price {
		if t := ingame.NewTower(tt, pos, s.Map.Path); t != nil {
			s.PlayerMapState.Money -= tt.Price
			s.Map.Towers = append(s.Map.Towers, t)

			s.rw.Append(s.Time, replay.PutTower, replay.InfoPutTower{
				Name: tt.Name,
				X:    x,
				Y:    y,
			})

			return true
		}
	}
	return false
}

func (s *GameState) sellTowerHandler() {
	t := s.chosenTower
	p := t.Price
	for i := 0; i < t.UpgradesBought; i++ {
		p += t.Upgrades[i].Price
	}

	p = p * 7 / 10
	s.PlayerMapState.Money += p

	t.Sold = true
	s.chosenTower = nil

	s.rw.Append(s.Time, replay.SellTower, replay.InfoSellTower{
		Index: s.findTowerIndex(t),
	})
}

func (s *GameState) upgradeTowerHandler() {
	t := s.chosenTower
	if !t.Upgrade(map[int]struct{}{1: {}}) {
		return
	}

	price := t.Upgrades[t.UpgradesBought-1].Price
	s.PlayerMapState.Money -= price

	s.rw.Append(s.Time, replay.UpgradeTower, replay.InfoUpgradeTower{
		Index: s.findTowerIndex(t),
	})
}

func (s *GameState) turnOnTowerHandler() {
	t := s.chosenTower
	t.State.IsTurnedOn = true

	s.rw.Append(s.Time, replay.TurnOn, replay.InfoTurnOnTower{
		Index: s.findTowerIndex(t),
	})
}

func (s *GameState) turnOffTowerHandler() {
	t := s.chosenTower
	t.State.IsTurnedOn = false

	s.rw.Append(s.Time, replay.TurnOff, replay.InfoTurnOffTower{
		Index: s.findTowerIndex(t),
	})
}

func (s *GameState) tuneFirstTowerHandler() {
	t := s.chosenTower
	t.State.AimType = ingame.First

	s.rw.Append(s.Time, replay.TuneFirst, replay.InfoTuneFirst{
		Index: s.findTowerIndex(t),
	})
}

func (s *GameState) tuneStrongTowerHandler() {
	t := s.chosenTower
	t.State.AimType = ingame.Strongest

	s.rw.Append(s.Time, replay.TuneStrong, replay.InfoTuneStrong{
		Index: s.findTowerIndex(t),
	})
}

func (s *GameState) tuneWeakTowerHandler() {
	t := s.chosenTower
	t.State.AimType = ingame.Weakest

	s.rw.Append(s.Time, replay.TuneWeak, replay.InfoTuneWeak{
		Index: s.findTowerIndex(t),
	})
}

func (s *GameState) findTowerIndex(t *ingame.Tower) int {
	for k, v := range s.Map.Towers {
		if v == t {
			return k
		}
	}
	return -1
}
