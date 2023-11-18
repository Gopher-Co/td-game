package main

import (
	"fmt"
	"log"

	"github.com/ebitenui/ebitenui"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"

	"github.com/gopher-co/td-game/global"
	"github.com/gopher-co/td-game/io"
	"github.com/gopher-co/td-game/models"
)

var TempEnemy *models.Enemy

type Game struct {
	s  models.State
	UI *ebitenui.UI
}

func (g *Game) Update() error {
	return g.s.Update()
}

func (g *Game) Draw(screen *ebiten.Image) {
	g.s.Draw(screen)
	ebitenutil.DebugPrint(screen, fmt.Sprintf("TPS: %f\n FPS %f\n", ebiten.ActualTPS(), ebiten.ActualFPS()))
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return outsideWidth, outsideHeight
}

func main() {
	ebiten.SetWindowSize(640, 480)
	ebiten.SetWindowTitle("Hello, World!")

	// load maps
	mcfgs, err := io.LoadMapConfigs()
	if err != nil {
		log.Fatalln(err)
	}

	for k := range mcfgs {
		global.GlobalMaps[mcfgs[k].Name] = &mcfgs[k]
	}

	// load levels
	lcfgs, err := io.LoadLevelConfigs()
	if err != nil {
		log.Fatalln(err)
	}

	for k := range lcfgs {
		global.GlobalLevels[lcfgs[k].LevelName] = &lcfgs[k]
	}

	// load enemies
	ecfgs, err := io.LoadEnemyConfigs()
	if err != nil {
		log.Fatalln(err)
	}

	for k := range ecfgs {
		global.GlobalEnemies[ecfgs[k].Name] = &ecfgs[k]
	}

	// load towers
	tcfgs, err := io.LoadTowerConfigs()
	if err != nil {
		log.Fatalln(err)
	}

	for k := range tcfgs {
		global.GlobalTowers[tcfgs[k].Name] = &tcfgs[k]
	}

	// load ui
	global.GlobalUI, err = io.LoadUIConfig()
	if err != nil {
		log.Fatalln(err)
	}
	log.Println(global.GlobalUI)

	// LEVEL LOADING
	gs := models.NewGameState(global.GlobalLevels["Level 1"], global.GlobalMaps, nil, nil, nil)

	// SIMULATE SOME STATE
	gs.Map.Enemies = append(gs.Map.Enemies, models.NewEnemy(global.GlobalEnemies["#ab0ba0"], gs.Map.Path))
	gs.Map.Towers = append(gs.Map.Towers, models.NewTower(global.GlobalTowers["#e0983a"], models.Point{300, 350}, gs.Map.Path))

	log.Println("Starting game...")
	if err := ebiten.RunGame(&Game{s: gs}); err != nil {
		log.Fatal(err)
	}
}
