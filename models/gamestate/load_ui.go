package gamestate

import (
	"context"
	"fmt"
	"image/color"
	"time"

	"github.com/ebitenui/ebitenui"
	image2 "github.com/ebitenui/ebitenui/image"
	"github.com/ebitenui/ebitenui/widget"
	"github.com/hajimehoshi/ebiten/v2"
	"golang.org/x/image/colornames"

	"github.com/gopher-co/td-game/models/general"
	"github.com/gopher-co/td-game/models/ingame"
	"github.com/gopher-co/td-game/ui/loaders"
)

// loadGameUI loads UI of the game.
func (s *GameState) loadGameUI(ctx context.Context, widgets general.Widgets) *ebitenui.UI {
	root := widget.NewContainer(
		widget.ContainerOpts.Layout(widget.NewGridLayout(
			widget.GridLayoutOpts.Columns(2),
			widget.GridLayoutOpts.Stretch([]bool{false, true}, []bool{true}),
		)),
	)

	mapContainer := s.loadMapContainer(ctx, widgets)

	towerMenuContainer := s.loadTowerMenuContainer(ctx, widgets)

	root.AddChild(mapContainer)
	root.AddChild(towerMenuContainer)

	return &ebitenui.UI{Container: root}
}

func (s *GameState) loadMapContainer(ctx context.Context, widgets general.Widgets) *widget.Container {
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
		widget.TextOpts.Text("", loaders.FontTrueType(64), color.White),
		widget.TextOpts.WidgetOpts(widget.WidgetOpts.LayoutData(widget.AnchorLayoutData{
			HorizontalPosition: 0,
			VerticalPosition:   widget.AnchorLayoutPositionEnd,
			StretchHorizontal:  false,
			StretchVertical:    false,
		})),
	)

	go func() {
		t := time.NewTicker(time.Second / time.Duration(ebiten.TPS()))
		for {
			if ctx.Err() != nil {
				return
			}
			<-t.C
			if s.CurrentWave < 0 || s.CurrentWave >= len(s.GameRule) {
				waveText.Label = ""
				continue
			}
			waveText.Label = fmt.Sprintf("Wave: %d/%d", s.CurrentWave+1, len(s.GameRule))
		}
	}()

	waveContainer.AddChild(waveText)

	backButton := widget.NewButton(
		widget.ButtonOpts.Image(&widget.ButtonImage{
			Idle: image2.NewNineSliceColor(colornames.Cornflowerblue),
		}),
		widget.ButtonOpts.Text("<", loaders.FontTrueType(80), &widget.ButtonTextColor{Idle: color.White}),
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

	var speedButton *widget.Button
	speedButton = widget.NewButton(
		widget.ButtonOpts.Image(&widget.ButtonImage{
			Idle: image2.NewNineSliceColor(colornames.Cornflowerblue),
		}),
		widget.ButtonOpts.Text("<<", loaders.FontTrueType(60), &widget.ButtonTextColor{Idle: color.White}),
		widget.ButtonOpts.ClickedHandler(func(_ *widget.ButtonClickedEventArgs) {
			if s.speedUp {
				ebiten.SetTPS(60)
				s.speedUp = false
				speedButton.Image = &widget.ButtonImage{
					Idle: image2.NewNineSliceColor(colornames.Cornflowerblue),
				}
				return
			}

			ebiten.SetTPS(180)
			s.speedUp = true
			speedButton.Image = &widget.ButtonImage{
				Idle: image2.NewNineSliceColor(colornames.Greenyellow),
			}
		}),
		widget.ButtonOpts.WidgetOpts(widget.WidgetOpts.LayoutData(widget.AnchorLayoutData{
			HorizontalPosition: widget.AnchorLayoutPositionEnd,
			VerticalPosition:   widget.AnchorLayoutPositionEnd,
			StretchHorizontal:  false,
			StretchVertical:    false,
		})),
	)
	speedContainer.AddChild(speedButton)

	mapContainer.AddChild(waveContainer)
	mapContainer.AddChild(buttonContainer)
	mapContainer.AddChild(speedContainer)

	return mapContainer
}

var cMenu, cInfo *widget.Container

func (s *GameState) showTowerMenu() {
	menu := s.UI.Container.Children()[1].(*widget.Container).Children()[2].(*widget.Container)
	menu.RemoveChildren()
	menu.AddChild(cMenu)
}

func (s *GameState) showTowerInfoMenu() {
	menu := s.UI.Container.Children()[1].(*widget.Container).Children()[2].(*widget.Container)
	menu.RemoveChildren()
	menu.AddChild(cInfo)
}

func (s *GameState) updateTowerUI(t *ingame.Tower) {
	menuCont := cInfo.Children()

	info := menuCont[0].(*widget.Container)
	text0 := info.Children()[0].(*widget.Text)
	text0.Label = t.Name

	upgrades := menuCont[1].(*widget.Container)
	text1, btn := upgrades.Children()[0].(*widget.Text), upgrades.Children()[1].(*widget.Button)
	text1.Label = fmt.Sprintf("Level %d", t.UpgradesBought+1)

	if t.UpgradesBought >= len(t.Upgrades) {
		btn.Text().Label = "SOLD OUT"
	} else {
		btn.Text().Label = fmt.Sprintf("UPGRADE ($%d)", t.Upgrades[t.UpgradesBought].Price)
	}

	if t.UpgradesBought >= len(t.Upgrades) || s.PlayerMapState.Money < t.Upgrades[t.UpgradesBought].Price {
		btn.GetWidget().Disabled = true
	} else {
		btn.GetWidget().Disabled = false
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

	//sell := menuCont[3].(*widget.Container)
}
