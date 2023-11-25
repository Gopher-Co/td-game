package gamestate

import (
	"fmt"
	"image/color"
	"math"
	"time"

	"github.com/ebitenui/ebitenui"
	image2 "github.com/ebitenui/ebitenui/image"
	"github.com/ebitenui/ebitenui/widget"
	"golang.org/x/image/colornames"

	"github.com/gopher-co/td-game/models/general"
	"github.com/gopher-co/td-game/ui/loaders"
)

// loadGameUI loads UI of the game.
func (s *GameState) loadGameUI(widgets general.Widgets) *ebitenui.UI {
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

	mapContainer.AddChild(buttonContainer)

	towerMenuContainer := widget.NewContainer(
		widget.ContainerOpts.Layout(widget.NewGridLayout(
			widget.GridLayoutOpts.Columns(1),
			widget.GridLayoutOpts.Stretch([]bool{true}, []bool{false, false, true}),
		)),
		widget.ContainerOpts.WidgetOpts(widget.WidgetOpts.LayoutData(widget.GridLayoutData{})),
		widget.ContainerOpts.BackgroundImage(image2.NewNineSliceColor(colornames.Blueviolet)),
	)
	ttf := loaders.FontTrueType(40)
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
			<-time.After(time.Millisecond)
			health.Label = fmt.Sprintf("Health: %d", s.PlayerMapState.Health)
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
			<-time.After(time.Millisecond)
			money.Label = fmt.Sprintf("Money: %d", s.PlayerMapState.Money)
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

// scrollCont creates a scroll container.
func (s *GameState) scrollCont(_ general.Widgets) *widget.Container {
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
			widget.TextOpts.Text(v.Name, loaders.FontTrueType(20), color.White),
		)

		cont.AddChild(button)
		cont.AddChild(text)

		content.AddChild(cont)
	}

	scrollContainer := widget.NewScrollContainer(
		widget.ScrollContainerOpts.StretchContentWidth(),
		widget.ScrollContainerOpts.Content(content),
		widget.ScrollContainerOpts.Image(&widget.ScrollContainerImage{
			Idle: image2.NewNineSliceColor(color.NRGBA{R: 0x13, G: 0x1a, B: 0x22, A: 0xff}),
			Mask: image2.NewNineSliceColor(color.NRGBA{R: 0x13, G: 0x1a, B: 0x22, A: 0xff}),
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
				Idle:  image2.NewNineSliceColor(color.NRGBA{R: 100, G: 100, B: 100, A: 255}),
				Hover: image2.NewNineSliceColor(color.NRGBA{R: 100, G: 100, B: 100, A: 255}),
			},
			// Set the handle images
			&widget.ButtonImage{
				Idle:    image2.NewNineSliceColor(color.NRGBA{R: 255, G: 100, B: 100, A: 255}),
				Hover:   image2.NewNineSliceColor(color.NRGBA{R: 255, G: 100, B: 100, A: 255}),
				Pressed: image2.NewNineSliceColor(color.NRGBA{R: 255, G: 100, B: 100, A: 255}),
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
