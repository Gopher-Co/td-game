package replaystate

import (
	"context"
	"image"
	"log"

	"github.com/ebitenui/ebitenui"
	"github.com/hajimehoshi/ebiten/v2"

	"github.com/gopher-co/td-game/models/config"
	"github.com/gopher-co/td-game/models/general"
	"github.com/gopher-co/td-game/models/ingame"
	"github.com/gopher-co/td-game/replay"
)

// CurrentState is a type that represents the current state of the game.
type CurrentState int

const (
	// Paused is a state that represents the game is paused.
	Paused = CurrentState(iota)

	// Running is a state that represents the game is running.
	Running
)

// ReplayState is a struct that represents the state of the game.
type ReplayState struct {
	// Map is a map of the game.
	Map *ingame.Map

	// EnemyToCall is a map of enemies that can be called.
	EnemyToCall map[string]*config.Enemy

	// TowerToBuy is a map of towers that can be bought.
	TowerToBuy map[string]*config.Tower

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

	// speedUp is a flag that represents if the game is speed up.
	speedUp bool

	// cancel is a function that cancels the context.
	cancel context.CancelFunc

	// rw is a watcher of the replay.
	rw *replay.Watcher

	// currAction is an index of the current action.
	currAction int
}

// New creates a new entity of ReplayState.
func New(
	w *replay.Watcher,
	cfg *config.Level,
	maps map[string]*config.Map,
	tw map[string]*config.Tower,
	en map[string]*config.Enemy,
	widgets general.Widgets,
) *ReplayState {
	ctx, cancel := context.WithCancel(context.Background())

	rs := &ReplayState{
		Map:            ingame.NewMap(maps[cfg.MapName]),
		EnemyToCall:    en,
		TowerToBuy:     tw,
		State:          Running,
		GameRule:       ingame.NewGameRule(cfg.GameRule),
		Time:           0,
		PlayerMapState: w.InitPlayerMapState,
		cancel:         cancel,
		rw:             w,
	}

	rs.UI = rs.loadUI(ctx, widgets)

	return rs
}

// Draw draws the game.
func (r *ReplayState) Draw(screen *ebiten.Image) {
	if r.Ended == true {
		return
	}

	subScreen := screen.SubImage(image.Rect(0, 0, 1500, 1080))
	r.Map.Draw(subScreen.(*ebiten.Image))

	r.UI.Draw(screen)
}

// Update updates the game.
func (r *ReplayState) Update() error {
	if r.Ended {
		return nil
	}

	if r.State == Paused {
		return nil
	}

	r.Action()

	r.UI.Update()

	r.Map.Update()

	wave := r.GameRule[r.CurrentWave]
	if wave.Ended() && !r.Map.AreThereAliveEnemies() {
		r.setStateAfterWave()
		if r.CurrentWave == len(r.GameRule) {
			r.Ended = true
			r.setStateAfterEnd()
		}
		return nil
	}

	r.updateRunning(wave)
	r.Time++
	return nil
}

// setStateAfterWave sets the state after the wave.
func (r *ReplayState) setStateAfterWave() {
	r.Map.Enemies = []*ingame.Enemy{}
	r.Map.Projectiles = []*ingame.Projectile{}
	r.State = Running
	r.CurrentWave++
}

// setStateAfterEnd sets the state after the end of the game.
func (r *ReplayState) setStateAfterEnd() {
	ebiten.SetTPS(60)
	r.cancel()
	r.cancel = nil
}

// updateRunning updates the game in the running state.
func (r *ReplayState) updateRunning(wave *ingame.Wave) {
	es := wave.CallEnemies()
	for _, str := range es {
		r.Map.Enemies = append(r.Map.Enemies, ingame.NewEnemy(r.EnemyToCall[str], r.Map.Path))
	}

	for _, e := range r.Map.Enemies {
		if e.State.Dead {
			if e.State.PassPath {
				r.PlayerMapState.Health = max(r.PlayerMapState.Health-e.DealDamageToPlayer(), 0)
			} else {
				r.PlayerMapState.Money += e.MoneyAward
				e.MoneyAward = 0
			}
		}
	}
}

// End returns true if the game is ended.
func (r *ReplayState) End() bool {
	return r.Ended
}

// Action performs the action.
func (r *ReplayState) Action() {
	for {
		action := r.rw.Actions[r.currAction]
		if r.Time != action.F {
			return
		}

		switch action.Type {
		case replay.PutTower:
			info := action.Info.(replay.InfoPutTower)
			t := r.TowerToBuy[info.Name]
			r.putTowerHandler(t, general.Point{X: general.Coord(info.X), Y: general.Coord(info.Y)})
		case replay.UpgradeTower:
			info := action.Info.(replay.InfoUpgradeTower)
			r.Map.Towers[info.Index].Upgrade(nil)
		case replay.TuneWeak:
			info := action.Info.(replay.InfoTuneWeak)
			r.Map.Towers[info.Index].State.AimType = ingame.Weakest
		case replay.TuneStrong:
			info := action.Info.(replay.InfoTuneStrong)
			r.Map.Towers[info.Index].State.AimType = ingame.Strongest
		case replay.TuneFirst:
			info := action.Info.(replay.InfoTuneFirst)
			r.Map.Towers[info.Index].State.AimType = ingame.First
		case replay.TurnOn:
			info := action.Info.(replay.InfoTurnOnTower)
			r.Map.Towers[info.Index].State.IsTurnedOn = true
		case replay.TurnOff:
			info := action.Info.(replay.InfoTurnOffTower)
			r.Map.Towers[info.Index].State.IsTurnedOn = false
		case replay.SellTower:
			info := action.Info.(replay.InfoSellTower)
			r.Map.Towers[info.Index].Sold = true
		case replay.Stop:
			r.Ended = true
			return
		default:
			log.Fatalln("Not handled type:", action.Type)
		}
		r.currAction++

	}
}

// putTowerHandler handles the put tower action.
func (r *ReplayState) putTowerHandler(tt *config.Tower, pos general.Point) *ingame.Tower {
	if pos.X < 1500 && r.PlayerMapState.Money >= tt.Price {
		if t := ingame.NewTower(tt, pos, r.Map.Path); t != nil {
			r.PlayerMapState.Money -= tt.Price
			r.Map.Towers = append(r.Map.Towers, t)

			return t
		}
	}
	return nil
}
