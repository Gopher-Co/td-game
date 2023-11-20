package models

import (
	"fmt"
	"image"
	"image/color"
	"math"
	"time"

	"github.com/ebitenui/ebitenui"
	image2 "github.com/ebitenui/ebitenui/image"
	"github.com/ebitenui/ebitenui/widget"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"golang.org/x/image/colornames"
)

// CurrentState is an enum that represents the current state of the game.
type CurrentState int

const (
	// Running is the state when the game is running.
	Running CurrentState = iota
	// Paused is the state when the game is paused.
	Paused
	NextWaveReady
)

// GameState is a struct that represents the state of the game.
type GameState struct {
	Map            *Map
	TowersToBuy    map[string]*TowerConfig
	EnemyToCall    map[string]*EnemyConfig
	Ended          bool
	State          CurrentState
	UI             *ebitenui.UI
	LastWave       int
	CurrentWave    int
	GameRule       GameRule
	Time           Frames
	PlayerMapState PlayerMapState
	tookTower      *TowerConfig
}

func NewGameState(config *LevelConfig, maps map[string]*MapConfig, en map[string]*EnemyConfig, tw map[string]*TowerConfig, w Widgets) *GameState {
	gs := &GameState{
		Map:         NewMap(maps[config.MapName]),
		TowersToBuy: tw,
		EnemyToCall: en,
		Ended:       false,
		State:       NextWaveReady,
		UI:          nil, // loadUI loads it
		LastWave:    0,
		CurrentWave: -1,
		GameRule:    NewGameRule(config.GameRule),
		Time:        0,
		PlayerMapState: PlayerMapState{
			Health: 100,
			Money:  650,
		},
	}

	gs.loadUI(w)
	gs.LastWave = len(gs.GameRule) - 1 // is it needed??

	return gs
}

func (s *GameState) Update() error {
	if s.Ended {
		return nil
	}

	if s.State == Paused {
		return nil
	}

	s.UI.Update()

	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButton0) && s.tookTower != nil {
		x, y := ebiten.CursorPosition()
		pos := Point{Coord(x), Coord(y)}

		if x < 1500 && s.PlayerMapState.Money >= s.tookTower.Price {
			if t := NewTower(s.tookTower, pos, s.Map.Path); t != nil {
				s.PlayerMapState.Money -= s.tookTower.Price
				s.tookTower = nil
				s.Map.Towers = append(s.Map.Towers, t)
			}
		}
	}

	if s.State == NextWaveReady {
		return nil
	}

	s.Map.Update()
	wave := s.GameRule[s.CurrentWave]
	if wave.Ended() && !s.Map.AreThereAliveEnemies() {
		s.State = NextWaveReady
		s.Map.Enemies = []*Enemy{}
		s.Map.Projectiles = []*Projectile{}
		if s.CurrentWave == len(s.GameRule)-1 {
			s.Ended = true
		}

		return nil
	}

	es := wave.CallEnemies()
	for _, str := range es {
		s.Map.Enemies = append(s.Map.Enemies, NewEnemy(s.EnemyToCall[str], s.Map.Path))
	}

	for _, e := range s.Map.Enemies {
		if e.State.Dead {
			if e.State.PassPath {
				s.PlayerMapState.Health = max(s.PlayerMapState.Health-e.DealDamageToPlayer(), 0)
			} else {
				s.PlayerMapState.Money += e.MoneyAward
				e.MoneyAward = 0
			}
		}
	}

	return nil
}

func (s *GameState) loadUI(widgets Widgets) {
	s.UI = s.loadGameUI(widgets)
}

func (s *GameState) End() bool {
	return s.Ended
}

func (s *GameState) Draw(screen *ebiten.Image) {
	subScreen := screen.SubImage(image.Rect(0, 0, 1500, 1080))
	s.Map.Draw(subScreen.(*ebiten.Image))
	if s.CurrentWave >= 0 {
		ebitenutil.DebugPrintAt(screen, fmt.Sprintf("Wave %d", s.CurrentWave+1), 0, 1900)
	}

	s.UI.Draw(screen)
}

func (s *GameState) loadGameUI(widgets Widgets) *ebitenui.UI {
	root := widget.NewContainer(
		widget.ContainerOpts.Layout(widget.NewGridLayout(
			widget.GridLayoutOpts.Columns(2),
			widget.GridLayoutOpts.Stretch([]bool{false, true}, []bool{true}),
		)),
	)

	mapContainer := widget.NewContainer(
		widget.ContainerOpts.WidgetOpts(widget.WidgetOpts.MinSize(1500, 0)),
		widget.ContainerOpts.Layout(widget.NewStackedLayout()),
	)

	buttonContainer := widget.NewContainer(
		widget.ContainerOpts.Layout(widget.NewAnchorLayout()),
	)

	backButton := widget.NewButton(
		widget.ButtonOpts.Image(&widget.ButtonImage{
			Idle: image2.NewNineSliceColor(colornames.Cornflowerblue),
		}),
		widget.ButtonOpts.Text("<", mustLoadFont(80), &widget.ButtonTextColor{Idle: color.White}),
		widget.ButtonOpts.ClickedHandler(func(_ *widget.ButtonClickedEventArgs) {
			s.State = Paused
		}),
		widget.ButtonOpts.WidgetOpts(widget.WidgetOpts.LayoutData(widget.AnchorLayoutData{
			HorizontalPosition: widget.AnchorLayoutPositionEnd,
			VerticalPosition:   0,
			StretchHorizontal:  false,
			StretchVertical:    false,
		})),
	)
	buttonContainer.AddChild(backButton)

	mapContainer.AddChild(buttonContainer)

	towerMenuContainer := widget.NewContainer(
		widget.ContainerOpts.Layout(widget.NewGridLayout(
			widget.GridLayoutOpts.Columns(1),
			widget.GridLayoutOpts.Stretch([]bool{true}, []bool{false, false, true}),
		)),
		widget.ContainerOpts.WidgetOpts(widget.WidgetOpts.LayoutData(widget.GridLayoutData{})),
		widget.ContainerOpts.BackgroundImage(image2.NewNineSliceColor(colornames.Blueviolet)),
	)
	ttf := mustLoadFont(40)
	health := widget.NewText(
		widget.TextOpts.Text(fmt.Sprintf("Health: %d", s.PlayerMapState.Health), ttf, color.White),
		widget.TextOpts.Insets(widget.Insets{
			Top:    0,
			Left:   10,
			Right:  0,
			Bottom: 0,
		}),
	)
	go func() {
		for {
			select {
			case <-time.After(time.Millisecond):
				health.Label = fmt.Sprintf("Health: %d", s.PlayerMapState.Health)
			}
		}
	}()

	money := widget.NewText(
		widget.TextOpts.Text(fmt.Sprintf("Money: %d", s.PlayerMapState.Money), ttf, color.White),
		widget.TextOpts.Insets(widget.Insets{
			Top:    0,
			Left:   10,
			Right:  0,
			Bottom: 0,
		}),
	)

	go func() {
		for {
			select {
			case <-time.After(time.Millisecond):
				money.Label = fmt.Sprintf("Money: %d", s.PlayerMapState.Money)
			}
		}
	}()

	scrollContainer := s.scrollCont(widgets)

	towerMenuContainer.AddChild(health)
	towerMenuContainer.AddChild(money)
	towerMenuContainer.AddChild(scrollContainer)

	root.AddChild(mapContainer)
	root.AddChild(towerMenuContainer)

	return &ebitenui.UI{Container: root}
}

func (s *GameState) scrollCont(widgets Widgets) *widget.Container {
	root := widget.NewContainer(
		widget.ContainerOpts.Layout(widget.NewGridLayout(
			widget.GridLayoutOpts.Columns(2),
			widget.GridLayoutOpts.Stretch([]bool{true, false}, []bool{true}),
		)),
	)

	content := widget.NewContainer(
		widget.ContainerOpts.Layout(widget.NewRowLayout(
			widget.RowLayoutOpts.Direction(widget.DirectionVertical),
			widget.RowLayoutOpts.Spacing(20),
		)),
	)

	for _, v := range s.TowersToBuy {
		v := v
		cont := widget.NewContainer(
			widget.ContainerOpts.Layout(widget.NewGridLayout(
				widget.GridLayoutOpts.Columns(1),
				widget.GridLayoutOpts.Stretch([]bool{false}, []bool{true, false}),
				widget.GridLayoutOpts.Padding(widget.Insets{
					Top:    10,
					Left:   10,
					Right:  10,
					Bottom: 10,
				}),
			)),
			widget.ContainerOpts.WidgetOpts(widget.WidgetOpts.LayoutData(widget.RowLayoutData{
				Position: widget.RowLayoutPositionCenter,
			})),
		)

		button := widget.NewButton(
			widget.ButtonOpts.Image(&widget.ButtonImage{
				Idle: image2.NewNineSliceSimple(v.Image(), 0, 1),
			}),
			widget.ButtonOpts.WidgetOpts(widget.WidgetOpts.LayoutData(widget.GridLayoutData{
				MaxWidth:           100,
				MaxHeight:          100,
				HorizontalPosition: 0,
				VerticalPosition:   0,
			})),
			widget.ButtonOpts.WidgetOpts(widget.WidgetOpts.MinSize(200, 100)),
			widget.ButtonOpts.ClickedHandler(func(_ *widget.ButtonClickedEventArgs) {
				s.tookTower = v
			}),
		)
		text := widget.NewText(
			widget.TextOpts.Text(v.Name, mustLoadFont(20), color.White),
		)

		cont.AddChild(button)
		cont.AddChild(text)

		content.AddChild(cont)
	}

	scrollContainer := widget.NewScrollContainer(
		widget.ScrollContainerOpts.StretchContentWidth(),
		widget.ScrollContainerOpts.Content(content),
		widget.ScrollContainerOpts.Image(&widget.ScrollContainerImage{
			Idle: image2.NewNineSliceColor(color.NRGBA{0x13, 0x1a, 0x22, 0xff}),
			Mask: image2.NewNineSliceColor(color.NRGBA{0x13, 0x1a, 0x22, 0xff}),
		}),
	)

	root.AddChild(scrollContainer)

	pageSizeFunc := func() int {
		return int(math.Round(float64(scrollContainer.ContentRect().Dy()) / float64(content.GetWidget().Rect.Dy()) * 1000))
	}

	vSlider := widget.NewSlider(
		widget.SliderOpts.Direction(widget.DirectionVertical),
		widget.SliderOpts.MinMax(0, 1000),
		widget.SliderOpts.PageSizeFunc(pageSizeFunc),
		//On change update scroll location based on the Slider's value
		widget.SliderOpts.ChangedHandler(func(args *widget.SliderChangedEventArgs) {
			scrollContainer.ScrollTop = float64(args.Slider.Current) / 1000
		}),
		widget.SliderOpts.Images(
			// Set the track images
			&widget.SliderTrackImage{
				Idle:  image2.NewNineSliceColor(color.NRGBA{100, 100, 100, 255}),
				Hover: image2.NewNineSliceColor(color.NRGBA{100, 100, 100, 255}),
			},
			// Set the handle images
			&widget.ButtonImage{
				Idle:    image2.NewNineSliceColor(color.NRGBA{255, 100, 100, 255}),
				Hover:   image2.NewNineSliceColor(color.NRGBA{255, 100, 100, 255}),
				Pressed: image2.NewNineSliceColor(color.NRGBA{255, 100, 100, 255}),
			},
		),
	)

	scrollContainer.GetWidget().ScrolledEvent.AddHandler(func(args interface{}) {
		a := args.(*widget.WidgetScrolledEventArgs)
		p := pageSizeFunc() / 3
		if p < 1 {
			p = 1
		}
		vSlider.Current -= int(math.Round(a.Y * float64(p)))
	})

	root.AddChild(vSlider)

	return root
}

//func (s *GameState) putNewTower(root *widget.Container, tower *Tower)
