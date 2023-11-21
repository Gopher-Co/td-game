// Package main provides the entry point of the game.
package main

import (
	"fmt"
	"log"
	"time"

	"github.com/ebitenui/ebitenui"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"

	"github.com/gopher-co/td-game/global"
	"github.com/gopher-co/td-game/io"
	"github.com/gopher-co/td-game/models"
)

// Game implements ebiten.Game interface.
type Game struct {
	s       models.State
	UI      *ebitenui.UI
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
		case *models.GameState:
			g.s = models.NewMenuState(global.Levels, models.Widgets(global.UI))
		case *models.MenuState:
			ms := g.s.(*models.MenuState)
			log.Println(ms.Next)
			g.s = models.NewGameState(global.Levels[ms.Next], global.Maps, global.Enemies, global.Towers, models.Widgets(global.UI))
		default:
			panic(fmt.Sprintf("type %T must be handled", g.s))
		}
	}
	return g.s.Update()
}

// Draw draws the game screen by one frame.
func (g *Game) Draw(screen *ebiten.Image) {
	g.s.Draw(screen)
	ebitenutil.DebugPrint(screen, fmt.Sprintf("TPS: %f\n FPS %f\n", ebiten.ActualTPS(), ebiten.ActualFPS()))
}

// Layout returns the game screen size.
func (g *Game) Layout(_, _ int) (screenWidth, screenHeight int) {
	return 1920, 1080
}

// main is the entry point of the game.
func main() {
	ebiten.SetWindowSize(1280, 720)
	ebiten.SetWindowTitle("Hello, World!")

	// load maps
	mcfgs, err := io.LoadMapConfigs()
	if err != nil {
		log.Fatalln(err)
	}

	for k := range mcfgs {
		global.Maps[mcfgs[k].Name] = &mcfgs[k]
	}

	// load levels
	lcfgs, err := io.LoadLevelConfigs()

	if err != nil {
		log.Fatalln(err)
	}

	for k := range lcfgs {
		global.Levels[lcfgs[k].LevelName] = &lcfgs[k]
	}

	// load enemies
	ecfgs, err := io.LoadEnemyConfigs()
	if err != nil {
		log.Fatalln(err)
	}

	for k := range ecfgs {
		global.Enemies[ecfgs[k].Name] = &ecfgs[k]
	}

	// load towers
	tcfgs, err := io.LoadTowerConfigs()
	if err != nil {
		log.Fatalln(err)
	}

	for k := range tcfgs {
		global.Towers[tcfgs[k].Name] = &tcfgs[k]
	}

	// load ui
	global.UI, err = io.LoadUIConfig()
	if err != nil {
		log.Fatalln(err)
	}

	// LEVEL LOADING
	game := &Game{s: models.NewMenuState(global.Levels, models.Widgets(global.UI))}

	//gs := models.NewGameState(global.GlobalLevels["Level 1"], global.GlobalMaps, global.GlobalEnemies, global.GlobalTowers, nil)
	go func() {
		for {
			time.Sleep(time.Second)
			if gs, ok := game.s.(*models.GameState); ok {
				switch gs.State {
				case models.NextWaveReady:
					gs.CurrentWave++
					fallthrough
				case models.Paused:
					gs.State = models.Running
				}
			}
		}
	}()
	// SIMULATE SOME STATE
	//gs.Map.Towers = append(gs.Map.Towers, models.NewTower(global.GlobalTowers["#e0983a"], models.Point{300, 350}, gs.Map.Path))

	log.Println("Starting game...")
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
