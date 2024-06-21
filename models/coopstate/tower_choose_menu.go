package coopstate

import (
	"fmt"
	"image/color"
	"math"

	image2 "github.com/ebitenui/ebitenui/image"
	"github.com/ebitenui/ebitenui/widget"
	"golang.org/x/image/colornames"

	"github.com/gopher-co/td-game/models/config"
	"github.com/gopher-co/td-game/models/general"
	"github.com/gopher-co/td-game/ui/font"
)

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
				widget.GridLayoutOpts.Stretch([]bool{true}, []bool{false, false}),
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
				Idle: image2.NewNineSliceSimple(v.Image(), config.TowerImageWidth, config.TowerImageWidth),
			}),
			widget.ButtonOpts.WidgetOpts(widget.WidgetOpts.LayoutData(widget.GridLayoutData{
				MaxWidth:           100,
				MaxHeight:          100,
				HorizontalPosition: 0,
				VerticalPosition:   0,
			})),
			widget.ButtonOpts.WidgetOpts(widget.WidgetOpts.MinSize(200, 100)),
			widget.ButtonOpts.ClickedHandler(s.handleTowerTake(v)),
			widget.ButtonOpts.WidgetOpts(widget.WidgetOpts.LayoutData(widget.GridLayoutData{
				HorizontalPosition: widget.GridLayoutPositionCenter,
				MaxWidth:           config.TowerImageWidth,
				MaxHeight:          config.TowerImageWidth,
			})),
		)

		text := widget.NewText(
			widget.TextOpts.Text(fmt.Sprintf("%s $%d", v.Name, v.Price), font.TTF20, color.White),
			widget.TextOpts.Position(widget.TextPositionCenter, widget.TextPositionStart),
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
		// On change update scroll location based on the Slider's value
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

// loadTowerMenuContainer creates a tower menu container.
func (s *GameState) loadTowerMenuContainer(widgets general.Widgets) *widget.Container {
	root := widget.NewContainer(
		widget.ContainerOpts.Layout(widget.NewGridLayout(
			widget.GridLayoutOpts.Columns(1),
			widget.GridLayoutOpts.Stretch([]bool{true}, []bool{false, false, true}),
		)),
		widget.ContainerOpts.WidgetOpts(widget.WidgetOpts.LayoutData(widget.GridLayoutData{})),
		widget.ContainerOpts.BackgroundImage(image2.NewNineSliceColor(colornames.Blueviolet)),
	)

	health := widget.NewText(
		widget.TextOpts.Text(fmt.Sprintf("Health: %d", s.PlayerMapState.Health), font.TTF40, color.White),
		widget.TextOpts.Insets(widget.Insets{
			Top:    0,
			Left:   10,
			Right:  0,
			Bottom: 0,
		}),
	)

	money := widget.NewText(
		widget.TextOpts.Text(fmt.Sprintf("Money: %d", s.PlayerMapState.Money), font.TTF40, color.White),
		widget.TextOpts.Insets(widget.Insets{
			Top:    0,
			Left:   10,
			Right:  0,
			Bottom: 0,
		}),
	)

	s.uiUpdater.Append(func() {
		health.Label = fmt.Sprintf("Health: %d", s.PlayerMapState.Health)
		money.Label = fmt.Sprintf("Money: %d", s.PlayerMapState.Money)
	})

	// menu on the right side
	menu := widget.NewContainer(
		widget.ContainerOpts.Layout(widget.NewStackedLayout()),
	)

	scrollContainer := s.scrollCont(widgets)
	menuTower := s.newTowerMenuUI(widgets)
	cMenu = scrollContainer
	cInfo = menuTower

	menu.AddChild(scrollContainer)

	root.AddChild(health)
	root.AddChild(money)
	root.AddChild(menu)

	return root
}
