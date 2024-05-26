// Package main provides the entry point of the game.
package main

import (
	"errors"
	"fmt"
	"io/fs"
	"log"
	_ "net/http/pprof"
	"os"
	"time"

	"github.com/gopher-co/td-game/models/coopstate"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"

	"github.com/gopher-co/td-game/io"
	"github.com/gopher-co/td-game/models/gamestate"
	"github.com/gopher-co/td-game/models/general"
	"github.com/gopher-co/td-game/models/menustate"
	"github.com/gopher-co/td-game/models/replaystate"
	"github.com/gopher-co/td-game/ui"
)

var pprof = func() {}

// Game implements ebiten.Game interface.
type Game struct {
	s       general.State
	fscreen bool
}

// Update updates the game state by one tick.
func (g *Game) Update() error {
	if inpututil.IsKeyJustPressed(ebiten.KeyF11) {
		g.fscreen = !g.fscreen
		ebiten.SetFullscreen(g.fscreen)
	}
	if g.s.End() {
		switch g.s.(type) {
		case *gamestate.GameState:
			gs := g.s.(*gamestate.GameState)
			Replays = append(Replays, gs.Watcher)
			if gs.Win {
				go func() {
					PlayerState.LevelsComplete[gs.LevelName] = struct{}{}
					if err := io.SaveStats(PlayerState); err != nil {
						log.Println("save unsuccessful")
					}
				}()
			}
			g.s = menustate.New(PlayerState, Levels, Replays, general.Widgets(UI))
		case *menustate.MenuState:
			ms := g.s.(*menustate.MenuState)
			if ms.Next != "" && ms.Host == nil {
				g.s = gamestate.New(Levels[ms.Next], Maps, Enemies, Towers, PlayerState, general.Widgets(UI))
			} else if ms.Next != "" && ms.Host != nil {
				g.s = coopstate.ClientState{}
			} else if ms.NextReplay != -1 {
				r := Replays[ms.NextReplay]
				g.s = replaystate.New(r, Levels[r.Name], Maps, Towers, Enemies, general.Widgets(UI))
			}
		case *replaystate.ReplayState:
			g.s = menustate.New(PlayerState, Levels, Replays, general.Widgets(UI))
		default:
			panic(fmt.Sprintf("type %T must be handled", g.s))
		}
	}
	return g.s.Update()
}

var t = time.NewTicker(time.Second / 90)

// Draw draws the game screen by one frame.
func (g *Game) Draw(screen *ebiten.Image) {
	g.s.Draw(screen)
	<-t.C
}

// Layout returns the game screen size.
func (g *Game) Layout(_, _ int) (screenWidth, screenHeight int) {
	return 1920, 1080
}

// main is the entry point of the game.
func main() {
	ebiten.SetWindowSize(1280, 720)
	ebiten.SetWindowTitle("Go Build, Go Defend!")

	// load maps
	mcfgs, err := io.LoadMapConfigs()
	if err != nil {
		log.Fatalln(err)
	}

	for k := range mcfgs {
		Maps[mcfgs[k].Name] = &mcfgs[k]
	}

	// load levels
	lcfgs, err := io.LoadLevelConfigs()
	if err != nil {
		log.Fatalln(err)
	}

	for k := range lcfgs {
		lcfgs[k].Order = k + 1
		Levels[lcfgs[k].LevelName] = &lcfgs[k]
	}

	// load enemies
	ecfgs, err := io.LoadEnemyConfigs()
	if err != nil {
		log.Fatalln(err)
	}

	for k := range ecfgs {
		Enemies[ecfgs[k].Name] = &ecfgs[k]
	}

	// load towers
	tcfgs, err := io.LoadTowerConfigs()
	if err != nil {
		log.Fatalln(err)
	}

	for k := range tcfgs {
		Towers[tcfgs[k].Name] = &tcfgs[k]
	}

	// load ui
	uicfg, err := io.LoadUIConfig()
	if err != nil {
		log.Fatalln(err)
	}
	uicfg.Colors[ui.MenuMainLogoImage] = "menu_logo"
	uicfg.Colors[ui.MenuMainImage] = "menu_main"

	for k, v := range uicfg.Colors {
		UI[k], err = ui.InitImage(v)
		if err != nil {
			log.Fatalf("image not loaded %v:%v, error: %v", k, v, err)
		}
	}

	if err := os.Mkdir("Replays", 0o777); err != nil && !errors.Is(err, fs.ErrExist) {
		log.Fatalln("replays system is broken:", err)
	}

	replays, err := io.LoadReplays()
	if err != nil {
		log.Fatalln("Replays not loaded:", err)
	}
	Replays = replays

	PlayerState, err = io.LoadStats()
	if err != nil {
		log.Fatalln("Stats not loaded:", err)
	}

	if err := PlayerState.Valid(Levels); err != nil {
		log.Fatalln("Invalid player stats:", err)
	}
	// LEVEL LOADING
	menu := menustate.New(PlayerState, Levels, Replays, general.Widgets(UI))
	game := &Game{s: menu}

	// pprof
	pprof()

	log.Println("Starting game...")
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
