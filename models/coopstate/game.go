package coopstate

import (
	"context"
	"image"
	"image/color"
	"log"
	maps2 "maps"
	"slices"
	"time"

	"github.com/ebitenui/ebitenui"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/vector"

	"github.com/gopher-co/td-game/models/config"
	"github.com/gopher-co/td-game/models/general"
	"github.com/gopher-co/td-game/models/ingame"
	"github.com/gopher-co/td-game/replay"
	"github.com/gopher-co/td-game/ui/updater"
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
	cli GameHostClient

	stream GameHost_JoinLobbyClient

	// LevelName is a name of the level.
	LevelName string

	// Map is a map of the game.
	Map *ingame.Map

	// TowersToBuy is a map of towers that can be bought.
	TowersToBuy map[string]*config.Tower

	// EnemyToCall is a map of enemies that can be called.
	EnemyToCall map[string]*config.Enemy

	// Ended is a flag that represents if the game is ended.
	Ended bool

	// Win is a flag that represents if the game is won.
	Win bool

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

	// tookTower is a tower that was taken from the right sidebar.
	tookTower *config.Tower

	// chosenTower is a tower that was chosen from the map.
	chosenTower *ingame.Tower

	// speedUp is a flag that represents if the game is speeded up.
	speedUp bool

	// Watcher is a watcher of the game.
	Watcher *replay.Watcher

	// PlayerState is a state of the player.
	PlayerState *ingame.PlayerState

	uiUpdater *updater.Updater

	ctx context.Context

	ch <-chan *JoinLobbyResponse
}

// New creates a new entity of GameState.
func New(
	level *config.Level,
	maps map[string]*config.Map,
	en map[string]*config.Enemy,
	tw map[string]*config.Tower,
	ps *ingame.PlayerState,
	w general.Widgets,
	cli GameHostClient,
	cli2 GameHost_JoinLobbyClient,
) *GameState {
	// remove all the unavailable towers
	tw2 := maps2.Clone(tw)
	filter(tw2, func(s string, c *config.Tower) bool {
		if c.OpenLevel == "" {
			return false
		}

		_, ok := ps.LevelsComplete[c.OpenLevel]
		return !ok
	})

	// creating gamestate from configs
	gs := &GameState{
		cli:         cli,
		stream:      cli2,
		LevelName:   level.LevelName,
		Map:         ingame.NewMap(maps[level.MapName]),
		TowersToBuy: tw2,
		EnemyToCall: en,
		State:       NextWaveReady,
		CurrentWave: -1,
		GameRule:    ingame.NewGameRule(level.GameRule),
		PlayerMapState: ingame.PlayerMapState{
			Health: 100,
			Money:  650,
		},
		Watcher: &replay.Watcher{
			Name: level.LevelName,
			InitPlayerMapState: ingame.PlayerMapState{
				Health: 100,
				Money:  650,
			},
			Actions: make([]replay.Action, 0, 2500),
		},
		PlayerState: ps,
		uiUpdater:   new(updater.Updater),
		ctx:         context.Background(),
	}

	ch := make(chan *JoinLobbyResponse)
	go func() {
		for {
			v, err := cli2.Recv()
			if err != nil {
				log.Println(err)
			}
			ch <- v
		}
	}()

	gs.ch = ch
	gs.UI = gs.loadGameUI(w)

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

	// if clicked on tower
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButton0) {
		s.rightSidebarHandle()
	}

	select {
	case v := <-s.ch:
		if msg := v.GetPutTower(); msg != nil {
			towerConfig := s.TowersToBuy[msg.TowerName]
			s.putTowerHandler(towerConfig, int(msg.Point.X), int(msg.Point.Y))
		} else if msg := v.GetStartNewWave(); msg != nil {
			s.State = Running
			s.CurrentWave++
		} else if msg := v.GetSpeedUp(); msg != nil {
			ebiten.SetTPS(180)
			s.speedUp = true
		} else if msg := v.GetSlowDown(); msg != nil {
			ebiten.SetTPS(60)
			s.speedUp = false
		} else if msg := v.GetUpgradeTower(); msg != nil {
			s.upgradeTowerHandler(s.findTowerByIndex(int(msg.Tower.Id)))
		} else if msg := v.GetTurnOn(); msg != nil {
			s.turnOnTowerHandler(s.findTowerByIndex(int(msg.Tower.Id)))
		} else if msg := v.GetTurnOff(); msg != nil {
			s.turnOffTowerHandler(s.findTowerByIndex(int(msg.Tower.Id)))
			s.findTowerByIndex(int(msg.Tower.Id)).State.IsTurnedOn = false
		} else if msg := v.GetSellTower(); msg != nil {
			s.sellTowerHandler(s.findTowerByIndex(int(msg.Tower.Id)))
		} else if msg := v.GetTuneTower(); msg != nil {
			switch msg.Aim {
			case TuneTowerRequest_AIM_TOWER_AT_FIRST:
				s.tuneFirstTowerHandler(s.findTowerByIndex(int(msg.Tower.Id)))
			case TuneTowerRequest_AIM_TOWER_AT_STRONG:
				s.tuneStrongTowerHandler(s.findTowerByIndex(int(msg.Tower.Id)))
			case TuneTowerRequest_AIM_TOWER_AT_LAST:
				s.tuneWeakTowerHandler(s.findTowerByIndex(int(msg.Tower.Id)))
			}
		}
	default:
	}

	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButton0) && s.tookTower != nil {
		x, y := ebiten.CursorPosition()
		s.cli.PutTower(s.ctx, &PutTowerRequest{
			TowerName: s.tookTower.Name,
			Point:     &Point{X: float32(x), Y: float32(y)},
		})
		s.tookTower = nil
	} else if inpututil.IsMouseButtonJustPressed(ebiten.MouseButton2) {
		s.tookTower = nil
	}

	s.UI.Update()
	s.uiUpdater.Update()

	if s.Ended {
		return nil
	}
	if s.State == NextWaveReady {
		return nil
	}

	s.Map.Update()

	wave := s.GameRule[s.CurrentWave]
	s.updateRunning(wave)

	if wave.Ended() && !s.Map.AreThereAliveEnemies() {
		s.setStateAfterWave()
		if s.CurrentWave == len(s.GameRule)-1 {
			s.Ended = true
			s.Win = true
			s.setStateAfterEnd()
		}
		return nil
	}

	s.Time++
	return nil
}

// End returns true if the game is ended.
func (s *GameState) End() bool {
	return s.Ended
}

// Draw draws the game on the screen.
func (s *GameState) Draw(screen *ebiten.Image) {
	if s.Ended {
		return
	}

	subScreen := screen.SubImage(image.Rect(0, 0, 1500, 1080))
	s.Map.Draw(subScreen.(*ebiten.Image))

	if s.tookTower != nil {
		s.drawTookImageBeforeCursor(screen)
	}

	s.UI.Draw(screen)
}

// setStateAfterWave sets the state after the wave.
func (s *GameState) setStateAfterWave() {
	s.State = NextWaveReady
	s.Map.Enemies = []*ingame.Enemy{}
	s.Map.Projectiles = []*ingame.Projectile{}
}

// clear clears the game state.
func (s *GameState) clear() {
	ebiten.SetTPS(60)
}

// setStateAfterEnd sets the state after the end of the game.
func (s *GameState) setStateAfterEnd() {
	s.clear()
	// replay save
	s.Watcher.Append(s.Time, replay.Stop, replay.InfoStop{Null: nil})

	timestamp := time.Now().Truncate(0).Format("2006-01-02T15_04_05")
	s.Watcher.Time = timestamp
	if err := replay.Save("./Replays/replay_"+timestamp+".json", s.Watcher); err != nil {
		log.Println("couldn't save replay:", err)
		return
	}
}

// updateRunning updates the state of the game when it is running.
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

// drawTookImageBeforeCursor draws the image of the tower that was taken from the right sidebar.
func (s *GameState) drawTookImageBeforeCursor(screen *ebiten.Image) {
	img := s.tookTower.Image()
	cx, cy := ebiten.CursorPosition()
	ix, iy := img.Bounds().Dx(), img.Bounds().Dy()

	if !ingame.CheckCollisionPath(general.Point{X: general.Coord(cx), Y: general.Coord(cy)}, s.Map.Path) {
		vector.DrawFilledCircle(screen, float32(cx), float32(cy), s.tookTower.InitRadius, color.RGBA{A: 0x20}, false)
	} else {
		vector.DrawFilledCircle(screen, float32(cx), float32(cy), s.tookTower.InitRadius, color.RGBA{R: 0xff, A: 0x20}, false)
	}

	geom := ebiten.GeoM{}
	geom.Translate(float64(cx-ix/2), float64(cy-iy/2))
	screen.DrawImage(img, &ebiten.DrawImageOptions{GeoM: geom})
}

// filter filters the map m by the function f.
func filter[K comparable, V any, M ~map[K]V](m M, f func(K, V) bool) {
	for k, v := range m {
		if f(k, v) {
			delete(m, k)
		}
	}
}

// rightSidebarHandle handles the right sidebar.
func (s *GameState) rightSidebarHandle() {
	x, _ := ebiten.CursorPosition()
	if x <= 1500 {
		ts := slices.Clone(s.Map.Towers)

		// choose the tower at the top of the screen
		slices.Reverse(ts)

		b := true
		for _, t := range ts {
			if b && t.IsClicked() {
				t.Chosen = true
				s.chosenTower = t
				b = !b

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

// putTowerHandler handles the putting of the tower.
func (s *GameState) putTowerHandler(tt *config.Tower, x, y int) *ingame.Tower {
	pos := general.Point{X: general.Coord(x), Y: general.Coord(y)}

	if x < 1500 && s.PlayerMapState.Money >= tt.Price {
		if t := ingame.NewTower(tt, pos, s.Map.Path); t != nil {
			s.tookTower = nil
			s.PlayerMapState.Money -= tt.Price
			s.Map.Towers = append(s.Map.Towers, t)

			return t
		}
	}
	return nil
}

// sellTowerHandler handles the selling of the tower.
func (s *GameState) sellTowerHandler(t *ingame.Tower) {
	p := t.Price
	for i := 0; i < t.UpgradesBought; i++ {
		p += t.Upgrades[i].Price
	}

	p = p * 7 / 10
	s.PlayerMapState.Money += p

	t.Sold = true

	s.Map.Towers = slices.DeleteFunc(s.Map.Towers, func(tower *ingame.Tower) bool {
		return tower == s.chosenTower
	})
	s.chosenTower = nil
}

// upgradeTowerHandler handles the upgrading of the tower.
func (s *GameState) upgradeTowerHandler(t *ingame.Tower) {
	if !t.Upgrade(s.PlayerState.LevelsComplete) {
		return
	}

	price := t.Upgrades[t.UpgradesBought-1].Price
	s.PlayerMapState.Money -= price
}

// turnOnTowerHandler handles the turning on of the tower.
func (s *GameState) turnOnTowerHandler(t *ingame.Tower) {
	t.State.IsTurnedOn = true
}

// turnOffTowerHandler handles the turning off of the tower.
func (s *GameState) turnOffTowerHandler(t *ingame.Tower) {
	t.State.IsTurnedOn = false
}

// tuneFirstTowerHandler handles the tuning of the tower.
// It sets the tower to aim at the first enemy.
func (s *GameState) tuneFirstTowerHandler(t *ingame.Tower) {
	t.State.AimType = ingame.First
}

// tuneStrongTowerHandler handles the tuning of the tower.
// It sets the tower to aim at the strongest enemy.
func (s *GameState) tuneStrongTowerHandler(t *ingame.Tower) {
	t.State.AimType = ingame.Strongest
}

// tuneWeakTowerHandler handles the tuning of the tower.
// It sets the tower to aim at the weakest enemy.
func (s *GameState) tuneWeakTowerHandler(t *ingame.Tower) {
	t.State.AimType = ingame.Weakest
}

// findTowerIndex finds the index of the tower t in the map of towers.
func (s *GameState) findTowerByIndex(i int) *ingame.Tower {
	for _, v := range s.Map.Towers {
		if v.Index == i {
			return v
		}
	}
	return nil
}
