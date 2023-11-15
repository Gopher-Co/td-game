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

	m := models.NewMap(&models.MapConfig{
		BackgroundColor: "#AB0BA0",
		Path:            path,
	})
	for i := 0; i < 100; i++ {
		_ = cfg.InitImage()
		cfg.Vrms = 1 + rand.Float32()*5
		m.Enemies = append(m.Enemies, models.NewEnemy(cfg, path))
		cfg.Name = fmt.Sprintf("#%06x", rand.Intn(0x1000000))
	}
	t := models.NewTower(&models.TowerConfig{
		Name:            "#472398",
		Upgrades:        nil,
		Price:           0,
		Type:            0,
		InitDamage:      1,
		InitRadius:      10,
		InitSpeedAttack: 10,
		OpenLevel:       0,
	}, models.Point{X: 100, Y: 100}, path)
	m.Towers = append(m.Towers, t)

	ebiten.SetWindowSize(640, 480)
	ebiten.SetWindowTitle("Hello, World!")
	if err := ebiten.RunGame(&Game{s: models.NewGameState(m, nil, nil, nil, nil)}); err != nil {
		log.Fatal(err)
	}
}
