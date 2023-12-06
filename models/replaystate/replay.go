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

type CurrentState int

const (
	Paused = CurrentState(iota)
	Running
)

type ReplayState struct {
	// Map is a map of the game.
	Map *ingame.Map

	// EnemyToCall is a map of enemies that can be called.
	EnemyToCall map[string]*config.Enemy

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

	speedUp bool

	cancel context.CancelFunc

	rw *replay.Watcher

	currAction int
}

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

func (r *ReplayState) Draw(screen *ebiten.Image) {
	if r.Ended == true {
		return
	}

	subScreen := screen.SubImage(image.Rect(0, 0, 1500, 1080))
	r.Map.Draw(subScreen.(*ebiten.Image))

	r.UI.Draw(screen)
}

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

func (r *ReplayState) setStateAfterWave() {
	r.Map.Enemies = []*ingame.Enemy{}
	r.Map.Projectiles = []*ingame.Projectile{}
	r.State = Running
	r.CurrentWave++
}

func (r *ReplayState) setStateAfterEnd() {
	ebiten.SetTPS(60)
	r.cancel()
	r.cancel = nil
}

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

func (r *ReplayState) End() bool {
	return r.Ended
}

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
			r.putTowerHandler(t, general.Point{general.Coord(info.X), general.Coord(info.Y)})
		case replay.Stop:
			r.Ended = true
			return
		default:
			log.Fatalln("Not handled type:", action.Type)
		}
		r.currAction++

	}
}

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
