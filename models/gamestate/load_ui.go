package gamestate

import (
	"fmt"
	"image/color"

	"github.com/ebitenui/ebitenui"
	image2 "github.com/ebitenui/ebitenui/image"
	"github.com/ebitenui/ebitenui/widget"
	"golang.org/x/image/colornames"

	"github.com/gopher-co/td-game/models/general"
	"github.com/gopher-co/td-game/ui/font"
)

// loadGameUI loads UI of the game.
func (s *GameState) loadGameUI(widgets general.Widgets) *ebitenui.UI {
	root := widget.NewContainer(
		widget.ContainerOpts.Layout(widget.NewGridLayout(
			widget.GridLayoutOpts.Columns(2),
			widget.GridLayoutOpts.Stretch([]bool{false, true}, []bool{true}),
		)),
	)

	mapContainer := s.loadMapContainer(widgets)

	towerMenuContainer := s.loadTowerMenuContainer(widgets)

	root.AddChild(mapContainer)
	root.AddChild(towerMenuContainer)

	return &ebitenui.UI{Container: root}
}

// loadMapContainer loads a container that contains the map.
func (s *GameState) loadMapContainer(_ general.Widgets) *widget.Container {
	mapContainer := widget.NewContainer(
		widget.ContainerOpts.WidgetOpts(widget.WidgetOpts.MinSize(1500, 0)),
		widget.ContainerOpts.Layout(widget.NewStackedLayout()),
	)

	buttonContainer := widget.NewContainer(
		widget.ContainerOpts.Layout(widget.NewAnchorLayout()),
	)

	speedContainer := widget.NewContainer(
		widget.ContainerOpts.Layout(widget.NewAnchorLayout()),
	)

	waveContainer := widget.NewContainer(
		widget.ContainerOpts.Layout(widget.NewAnchorLayout()),
	)

	waveText := widget.NewText(
		widget.TextOpts.Text("", font.TTF64, color.White),
		widget.TextOpts.WidgetOpts(widget.WidgetOpts.LayoutData(widget.AnchorLayoutData{
			HorizontalPosition: 0,
			VerticalPosition:   widget.AnchorLayoutPositionEnd,
			StretchHorizontal:  false,
			StretchVertical:    false,
		})),
	)

	s.uiUpdater.Append(func() {
		if s.CurrentWave < 0 || s.CurrentWave >= len(s.GameRule) {
			waveText.Label = ""
			return
		}
		waveText.Label = fmt.Sprintf("Wave: %d/%d", s.CurrentWave+1, len(s.GameRule))
	})

	waveContainer.AddChild(waveText)

	backButton := widget.NewButton(
		widget.ButtonOpts.Image(&widget.ButtonImage{
			Idle: image2.NewNineSliceColor(colornames.Cornflowerblue),
		}),
		widget.ButtonOpts.TextPadding(widget.Insets{
			Top:    5,
			Left:   10,
			Right:  10,
			Bottom: 5,
		}),
		widget.ButtonOpts.Text("Menu", font.TTF64, &widget.ButtonTextColor{Idle: color.White}),
		widget.ButtonOpts.ClickedHandler(s.handleMenu),
		widget.ButtonOpts.WidgetOpts(widget.WidgetOpts.LayoutData(widget.AnchorLayoutData{
			HorizontalPosition: widget.AnchorLayoutPositionEnd,
			VerticalPosition:   0,
			StretchHorizontal:  false,
			StretchVertical:    false,
		})),
	)
	buttonContainer.AddChild(backButton)

	buttonGroup := widget.NewContainer(
		widget.ContainerOpts.Layout(widget.NewRowLayout(
			widget.RowLayoutOpts.Direction(widget.DirectionHorizontal),
		)),
		widget.ContainerOpts.WidgetOpts(widget.WidgetOpts.LayoutData(widget.AnchorLayoutData{
			HorizontalPosition: widget.AnchorLayoutPositionEnd,
			VerticalPosition:   widget.AnchorLayoutPositionEnd,
			StretchHorizontal:  false,
			StretchVertical:    false,
		})),
	)

	startButton := widget.NewButton(
		widget.ButtonOpts.Image(&widget.ButtonImage{
			Idle:     image2.NewNineSliceColor(colornames.Darkgreen),
			Disabled: image2.NewNineSliceColor(color.RGBA{R: 128, G: 10, B: 30, A: 180}),
		}),
		widget.ButtonOpts.Text("Start", font.TTF64, &widget.ButtonTextColor{
			Idle:     color.White,
			Disabled: color.RGBA{R: 0xff, G: 0xff, B: 0xff, A: 180},
		}),
		widget.ButtonOpts.ClickedHandler(s.handleStart),
	)

	s.uiUpdater.Append(func() {
		if s.State != Running && !s.End() {
			startButton.GetWidget().Disabled = false
		}
	})

	speedButton := widget.NewButton(
		widget.ButtonOpts.Image(&widget.ButtonImage{
			Idle: image2.NewNineSliceColor(colornames.Cornflowerblue),
		}),
		widget.ButtonOpts.Text(">>", font.TTF64, &widget.ButtonTextColor{Idle: color.White}),
		widget.ButtonOpts.ClickedHandler(s.handleSpeed),
	)

	buttonGroup.AddChild(startButton)
	buttonGroup.AddChild(speedButton)

	speedContainer.AddChild(buttonGroup)

	mapContainer.AddChild(waveContainer)
	mapContainer.AddChild(buttonContainer)
	mapContainer.AddChild(speedContainer)

	return mapContainer
}

var cMenu, cInfo *widget.Container

// showTowerMenu shows the tower menu.
func (s *GameState) showTowerMenu() {
	menu := s.UI.Container.Children()[1].(*widget.Container).Children()[2].(*widget.Container)
	menu.RemoveChildren()
	menu.AddChild(cMenu)
}

// showTowerInfoMenu shows the tower info menu.
func (s *GameState) showTowerInfoMenu() {
	menu := s.UI.Container.Children()[1].(*widget.Container).Children()[2].(*widget.Container)
	menu.RemoveChildren()
	menu.AddChild(cInfo)
}
