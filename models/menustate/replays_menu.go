package menustate

import (
	"fmt"
	"image/color"
	"math"
	"sort"

	"github.com/ebitenui/ebitenui"
	"github.com/ebitenui/ebitenui/image"
	"github.com/ebitenui/ebitenui/widget"
	"golang.org/x/image/colornames"

	"github.com/gopher-co/td-game/models/general"
	"github.com/gopher-co/td-game/ui"
	"github.com/gopher-co/td-game/ui/loaders"
)

func (m *MenuState) loadReplaysMenuUI(widgets general.Widgets) *ebitenui.UI {
	bgImg := widgets[ui.MenuBackgroundImage]
	menuBackground := image.NewNineSliceSimple(bgImg, 0, 1)

	backBtn := widget.NewButton(
		widget.ButtonOpts.Image(&widget.ButtonImage{Idle: image.NewNineSliceColor(colornames.Aqua)}),
		widget.ButtonOpts.Text("<", loaders.FontTrueType(128), &widget.ButtonTextColor{Idle: color.White}),
		widget.ButtonOpts.ClickedHandler(func(args *widget.ButtonClickedEventArgs) {
			m.UI = m.loadMainMenuUI(widgets)
		}),
	)

	root := widget.NewContainer(
		widget.ContainerOpts.Layout(widget.NewGridLayout(
			widget.GridLayoutOpts.Columns(1),
			widget.GridLayoutOpts.Stretch([]bool{true}, []bool{false, true}),
			widget.GridLayoutOpts.Spacing(0, 0),
		)),
		widget.ContainerOpts.BackgroundImage(menuBackground),
	)

	root.AddChild(backBtn)
	root.AddChild(m.loadScrollingReplays(widgets))

	return &ebitenui.UI{Container: root}
}

func (m *MenuState) loadScrollingReplays(_ general.Widgets) *widget.Container {
	root := widget.NewContainer(
		widget.ContainerOpts.Layout(widget.NewGridLayout(
			widget.GridLayoutOpts.Columns(1),
			widget.GridLayoutOpts.Stretch([]bool{true}, []bool{true, false}),
		)),
	)

	content := widget.NewContainer(
		widget.ContainerOpts.Layout(widget.NewRowLayout(
			widget.RowLayoutOpts.Direction(widget.DirectionHorizontal),
			widget.RowLayoutOpts.Spacing(20),
		)),
	)

	levels := make([]string, 0, len(m.Levels))
	for k := range m.Levels {
		levels = append(levels, k)
	}
	sort.Strings(levels)
	ttf72 := loaders.FontTrueType(72)
	ttf36 := loaders.FontTrueType(36)
	// blackImg := image.NewNineSliceColor(color.Black)
	for k, v := range m.Replays {
		k := k
		v := v

		cont := widget.NewContainer(
			widget.ContainerOpts.WidgetOpts(widget.WidgetOpts.MinSize(400, 900)),
			widget.ContainerOpts.Layout(widget.NewGridLayout(
				widget.GridLayoutOpts.Columns(1),
				widget.GridLayoutOpts.Stretch([]bool{true}, []bool{false, true, false}),
			)),
		)
		text1 := widget.NewText(
			widget.TextOpts.Text("Replay", ttf72, color.White),
		)
		text2 := widget.NewText(
			widget.TextOpts.MaxWidth(400),
			widget.TextOpts.Text(fmt.Sprintf("Level: %s\nTimestamp: %s", v.Name, v.Time), ttf36, color.White),
		)
		btn := widget.NewButton(
			widget.ButtonOpts.Image(&widget.ButtonImage{Idle: image.NewNineSliceColor(colornames.Beige)}),
			widget.ButtonOpts.Text("Watch", ttf72, &widget.ButtonTextColor{Idle: color.Black}),
			widget.ButtonOpts.WidgetOpts(widget.WidgetOpts.LayoutData(widget.GridLayoutData{
				VerticalPosition: widget.GridLayoutPositionEnd,
			})),
			widget.ButtonOpts.ClickedHandler(func(args *widget.ButtonClickedEventArgs) {
				m.Ended = true
				m.NextReplay = k
			}),
		)

		cont.AddChild(text1)
		cont.AddChild(text2)
		cont.AddChild(btn)

		content.AddChild(cont)
	}

	scrollContainer := widget.NewScrollContainer(
		widget.ScrollContainerOpts.StretchContentWidth(),
		widget.ScrollContainerOpts.Content(content),
		widget.ScrollContainerOpts.Image(&widget.ScrollContainerImage{
			Idle: image.NewNineSliceColor(color.NRGBA{R: 0x13, G: 0x1a, B: 0x22, A: 0xff}),
			Mask: image.NewNineSliceColor(color.NRGBA{R: 0x13, G: 0x1a, B: 0x22, A: 0xff}),
		}),
	)
	root.AddChild(scrollContainer)

	pageSizeFunc := func() int {
		return int(math.Round(float64(scrollContainer.ContentRect().Dx()) / float64(content.GetWidget().Rect.Dx()) * 100))
	}

	vSlider := widget.NewSlider(
		widget.SliderOpts.Direction(widget.DirectionHorizontal),
		widget.SliderOpts.MinMax(0, 500),
		widget.SliderOpts.PageSizeFunc(pageSizeFunc),
		// On change update scroll location based on the Slider's value
		widget.SliderOpts.ChangedHandler(func(args *widget.SliderChangedEventArgs) {
			scrollContainer.ScrollLeft = float64(args.Slider.Current) / 500
		}),
		widget.SliderOpts.Images(
			// Set the track images
			&widget.SliderTrackImage{
				Idle:  image.NewNineSliceColor(color.NRGBA{R: 100, G: 100, B: 100, A: 255}),
				Hover: image.NewNineSliceColor(color.NRGBA{R: 100, G: 100, B: 100, A: 255}),
			},
			// Set the handle images
			&widget.ButtonImage{
				Idle:    image.NewNineSliceColor(color.NRGBA{R: 255, G: 100, B: 100, A: 255}),
				Hover:   image.NewNineSliceColor(color.NRGBA{R: 255, G: 100, B: 100, A: 255}),
				Pressed: image.NewNineSliceColor(color.NRGBA{R: 255, G: 100, B: 100, A: 255}),
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
