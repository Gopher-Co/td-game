package gamestate

import (
	"fmt"
	"image/color"

	"github.com/ebitenui/ebitenui"
	image2 "github.com/ebitenui/ebitenui/image"
	"github.com/ebitenui/ebitenui/widget"
	"github.com/hajimehoshi/ebiten/v2"
	"golang.org/x/image/colornames"

	"github.com/gopher-co/td-game/models/general"
	"github.com/gopher-co/td-game/models/ingame"
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
		widget.ButtonOpts.ClickedHandler(func(_ *widget.ButtonClickedEventArgs) {
			s.setStateAfterEnd()
			s.Ended = true
		}),
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
		widget.ButtonOpts.ClickedHandler(func(args *widget.ButtonClickedEventArgs) {
			b := args.Button
			if !b.GetWidget().Disabled {
				s.State = Running
				s.CurrentWave++
				b.GetWidget().Disabled = true
			}
		}),
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
		widget.ButtonOpts.ClickedHandler(func(args *widget.ButtonClickedEventArgs) {
			if s.speedUp {
				ebiten.SetTPS(60)
				s.speedUp = false
				args.Button.Image = &widget.ButtonImage{
					Idle: image2.NewNineSliceColor(colornames.Cornflowerblue),
				}
				return
			}

			ebiten.SetTPS(180)
			s.speedUp = true
			args.Button.Image = &widget.ButtonImage{
				Idle: image2.NewNineSliceColor(colornames.Greenyellow),
			}
		}),
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

// updateTowerUI updates the tower UI.
func (s *GameState) updateTowerUI(t *ingame.Tower) {
	menuCont := cInfo.Children()

	info := menuCont[0].(*widget.Container)
	text0 := info.Children()[0].(*widget.Text)
	text0.Label = t.Name

	upgrades := menuCont[1].(*widget.Container)
	text1, btn := upgrades.Children()[0].(*widget.Text), upgrades.Children()[1].(*widget.Button)
	text1.Label = fmt.Sprintf("Level %d", t.UpgradesBought+1)

	if s.chosenTower.UpgradesBought >= len(s.chosenTower.Upgrades) {
		btn.Text().Label = "SOLD OUT"
		btn.GetWidget().Disabled = true
	} else {
		openLevel := s.chosenTower.Upgrades[s.chosenTower.UpgradesBought].OpenLevel
		_, ok := s.PlayerState.LevelsComplete[openLevel]
		if !ok && openLevel != "" {
			btn.Text().Label = "Complete level to unlock:\n" + s.chosenTower.Upgrades[s.chosenTower.UpgradesBought].OpenLevel
			btn.GetWidget().Disabled = true
		} else {
			btn.Text().Label = fmt.Sprintf("UPGRADE ($%d)", s.chosenTower.Upgrades[s.chosenTower.UpgradesBought].Price)
			if s.PlayerMapState.Money < t.Upgrades[t.UpgradesBought].Price {
				btn.GetWidget().Disabled = true
			} else {
				btn.GetWidget().Disabled = false
			}
		}
	}

	tuning := menuCont[2].(*widget.Container)
	turnButton := tuning.Children()[0].(*widget.Button)
	if s.chosenTower.State.IsTurnedOn {
		turnButton.Text().Label = "ON"
		turnButton.Image = &widget.ButtonImage{
			Idle: image2.NewNineSliceColor(colornames.Lawngreen),
		}
	} else {
		turnButton.Text().Label = "OFF"
		turnButton.Image = &widget.ButtonImage{
			Idle: image2.NewNineSliceColor(colornames.Indianred),
		}
	}
}
