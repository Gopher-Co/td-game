package main

import (
	"fmt"
	"log"

	"github.com/ebitenui/ebitenui"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"

	"github.com/gopher-co/td-game/global"
	"github.com/gopher-co/td-game/io"
	"github.com/gopher-co/td-game/models"
)

var TempEnemy *models.Enemy

type Game struct {
	s       models.State
	UI      *ebitenui.UI
	fscreen bool
}

func (g *Game) Update() error {
	if inpututil.IsKeyJustPressed(ebiten.KeyF11) {
		g.fscreen = !g.fscreen
		ebiten.SetFullscreen(g.fscreen)
	}
	return g.s.Update()
}

func (g *Game) Draw(screen *ebiten.Image) {
	g.s.Draw(screen)
	ebitenutil.DebugPrint(screen, fmt.Sprintf("TPS: %f\n FPS %f\n", ebiten.ActualTPS(), ebiten.ActualFPS()))
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return 1920, 1080
}

func main() {
	ebiten.SetWindowSize(1280, 720)
	ebiten.SetWindowTitle("Hello, World!")
	var err error

	// load ui
	global.GlobalUI, err = io.LoadUIConfig()
	if err != nil {
		log.Fatalln(err)
	}
	ms := models.NewMenuState([]*models.LevelConfig{{
		LevelName: "ASDASD1",
		MapName:   "",
		GameRule:  nil,
	}, {
		LevelName: "ASDASD2",
		MapName:   "",
		GameRule:  nil,
	}, {
		LevelName: "ASDASD3",
		MapName:   "",
		GameRule:  nil,
	}, {
		LevelName: "ASDASD4",
		MapName:   "",
		GameRule:  nil,
	}, {
		LevelName: "ASDASD5",
		MapName:   "",
		GameRule:  nil,
	}, {
		LevelName: "ASDASD6",
		MapName:   "",
		GameRule:  nil,
	}, {
		LevelName: "ASDASD7",
		MapName:   "",
		GameRule:  nil,
	}, {
		LevelName: "ASDASD8",
		MapName:   "",
		GameRule:  nil,
	}}, models.Widgets(global.GlobalUI))

	log.Println("Starting game...")
	if err := ebiten.RunGame(&Game{s: ms}); err != nil {
		log.Fatal(err)
	}
}
