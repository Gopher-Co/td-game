package main

import (
	"fmt"
	"log"

	"github.com/ebitenui/ebitenui"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"

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

	path := models.Path{{-16, -16}, {200, 200}, {300, 240}, {500, 50}, {500, 350}, {300, 270}, {300, 500}}
	m := models.NewMap(&models.MapConfig{
		BackgroundColor: "#AB0BA0",
		Path:            path,
	})

	ecfgs, err := io.LoadEnemyConfigs()
	if err != nil {
		log.Fatalln(err)
	}

	for k := range ecfgs {
		m.Enemies = append(m.Enemies, models.NewEnemy(&ecfgs[k], path))
	}

	tcfgs, err := io.LoadTowerConfigs()
	if err != nil {
		log.Fatalln(err)
	}

	for k := range tcfgs {
		if t := models.NewTower(&tcfgs[k], models.Point{190, 100}, path); t != nil {
			m.Towers = append(m.Towers, t)
		}

	}

	lcfg := models.NewGameState(&models.LevelConfig{}, nil, nil, nil)
	lcfg.Map = m

	if err := ebiten.RunGame(&Game{s: lcfg}); err != nil {
		log.Fatal(err)
	}
}
