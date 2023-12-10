package replaystate

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
	"github.com/gopher-co/td-game/ui/loaders"
)

func (r *ReplayState) loadUI(ctx context.Context, widgets general.Widgets) *ebitenui.UI {
	root := widget.NewContainer(
		widget.ContainerOpts.Layout(widget.NewGridLayout(
			widget.GridLayoutOpts.Columns(2),
			widget.GridLayoutOpts.Stretch([]bool{false, true}, []bool{true}),
		)),
	)

	mapContainer := r.loadMapContainer(ctx, widgets)

	towerMenuContainer := r.loadTowerMenuContainer(ctx, widgets)

	root.AddChild(mapContainer)
	root.AddChild(towerMenuContainer)

	return &ebitenui.UI{Container: root}
}

func (r *ReplayState) loadMapContainer(ctx context.Context, widgets general.Widgets) *widget.Container {
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

	ttf64 := loaders.FontTrueType(64)
	defer ttf64.Close()

	waveText := widget.NewText(
		widget.TextOpts.Text("", ttf64, color.White),
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
			if r.CurrentWave < 0 || r.CurrentWave >= len(r.GameRule) {
				waveText.Label = ""
				continue
			}
			waveText.Label = fmt.Sprintf("Wave: %d/%d", r.CurrentWave+1, len(r.GameRule))
		}
	}()

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
		widget.ButtonOpts.Text("Menu", ttf64, &widget.ButtonTextColor{Idle: color.White}),
		widget.ButtonOpts.ClickedHandler(func(_ *widget.ButtonClickedEventArgs) {
			r.setStateAfterEnd()
			r.Ended = true
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
		widget.ButtonOpts.Text(">>", ttf64, &widget.ButtonTextColor{Idle: color.White}),
		widget.ButtonOpts.ClickedHandler(func(_ *widget.ButtonClickedEventArgs) {
			if r.speedUp {
				ebiten.SetTPS(60)
				r.speedUp = false
				speedButton.Image = &widget.ButtonImage{
					Idle: image2.NewNineSliceColor(colornames.Cornflowerblue),
				}
				return
			}

			ebiten.SetTPS(180)
			r.speedUp = true
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

func (r *ReplayState) loadTowerMenuContainer(ctx context.Context, widgets general.Widgets) *widget.Container {
	root := widget.NewContainer(
		widget.ContainerOpts.Layout(widget.NewGridLayout(
			widget.GridLayoutOpts.Columns(1),
			widget.GridLayoutOpts.Stretch([]bool{true}, []bool{false, false, true}),
		)),
		widget.ContainerOpts.WidgetOpts(widget.WidgetOpts.LayoutData(widget.GridLayoutData{})),
		widget.ContainerOpts.BackgroundImage(image2.NewNineSliceColor(colornames.Blueviolet)),
	)

	ttf40 := loaders.FontTrueType(40)
	defer ttf40.Close()

	health := widget.NewText(
		widget.TextOpts.Text(fmt.Sprintf("Health: %d", r.PlayerMapState.Health), ttf40, color.White),
		widget.TextOpts.Insets(widget.Insets{
			Top:    0,
			Left:   10,
			Right:  0,
			Bottom: 0,
		}),
	)
	go func() {
		ticker := time.NewTicker(time.Second / time.Duration(ebiten.ActualTPS()))
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
			}

			health.Label = fmt.Sprintf("Health: %d", r.PlayerMapState.Health)
		}
	}()

	money := widget.NewText(
		widget.TextOpts.Text(fmt.Sprintf("Money: %d", r.PlayerMapState.Money), ttf40, color.White),
		widget.TextOpts.Insets(widget.Insets{
			Top:    0,
			Left:   10,
			Right:  0,
			Bottom: 0,
		}),
	)

	go func() {
		ticker := time.NewTicker(time.Second / time.Duration(ebiten.ActualTPS()))
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
			}

			money.Label = fmt.Sprintf("Money: %d", r.PlayerMapState.Money)
		}
	}()

	// menu on the right side
	menu := widget.NewContainer(
		widget.ContainerOpts.Layout(widget.NewStackedLayout()),
	)

	root.AddChild(health)
	root.AddChild(money)
	root.AddChild(menu)

	return root
}
