package main

import (
	"fmt"
	"log"
	"math/rand"

	"github.com/ebitenui/ebitenui"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"

	"github.com/gopher-co/td-game/models"
)

type Game struct {
	path  models.Path
	es    []*models.Enemy
	State models.State
	UI    *ebitenui.UI
}

func (g *Game) Update() error {
	for _, e := range g.es {
		if !e.State.Dead {
			e.Update()
		}
	}

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	g.path.Draw(screen)

	for _, e := range g.es {
		if !e.State.Dead {
			e.Draw(screen)
		}
	}
	ebitenutil.DebugPrint(screen, fmt.Sprintf("TPS: %f\n FPS %f\n", ebiten.ActualTPS(), ebiten.ActualFPS()))

}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return 640, 480
}

func main() {
	path := models.Path{{-16, -16}, {200, 200}, {300, 240}, {500, 50}, {500, 350}, {300, 270}, {300, 490}}
	cfg := &models.EnemyConfig{
		Name:       "#DEAD00",
		MaxHealth:  1,
		Damage:     0,
		MoneyAward: 0,
		Strengths:  nil,
		Weaknesses: nil,
	}

	_ = cfg.InitImage()

	var es []*models.Enemy
	for i := 0; i < 10; i++ {
		cfg.Vrms = 1 + rand.Float32()*10
		es = append(es, models.NewEnemy(cfg, path))
	}
	ebiten.SetWindowSize(640, 480)
	ebiten.SetWindowTitle("Hello, World!")
	if err := ebiten.RunGame(&Game{es: es, path: path}); err != nil {
		log.Fatal(err)
	}
}
