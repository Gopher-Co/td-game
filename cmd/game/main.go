package main

import (
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/ebitenui/ebitenui"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"

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
	cfg := &models.EnemyConfig{
		Name:       "#DEAD00",
		MaxHealth:  18,
		Damage:     0,
		MoneyAward: 0,
		Strengths:  nil,
		Weaknesses: nil,
	}

	m := models.NewMap(&models.MapConfig{
		BackgroundColor: "#AB0BA0",
		Path:            path,
	})

	_ = cfg.InitImage()
	cfg.Vrms = 1
	TempEnemy = models.NewEnemy(cfg, path)
	m.Enemies = append(m.Enemies, TempEnemy)
	cfg.Name = fmt.Sprintf("#%06x", rand.Intn(0x1000000))
	go func() {
		time.Sleep(3 * time.Second)
		_ = cfg.InitImage()
		cfg.Vrms = 2.7
		m.Enemies = append(m.Enemies, models.NewEnemy(cfg, path))
		cfg.Name = fmt.Sprintf("#%06x", rand.Intn(0x1000000))
	}()
	tcfg := &models.TowerConfig{
		Name:            "#000",
		Upgrades:        nil,
		Price:           0,
		Type:            0,
		InitDamage:      1,
		InitRadius:      200,
		InitSpeedAttack: 10,
		OpenLevel:       0,
	}
	if err := tcfg.InitImage(); err != nil {
		log.Fatalln(err)
	}

	t := models.NewTower(tcfg, models.Point{X: 290, Y: 120}, path)
	t.State.Aim = TempEnemy
	m.Towers = append(m.Towers, t)

	lcfg := models.NewGameState(&models.LevelConfig{}, nil, nil, nil)
	lcfg.Map = m

	if err := ebiten.RunGame(&Game{s: lcfg}); err != nil {
		log.Fatal(err)
	}
}
