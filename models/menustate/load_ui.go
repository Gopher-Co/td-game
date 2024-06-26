package menustate

import (
	"fmt"
	"image/color"
	"math"
	"os"
	"sort"

	"github.com/ebitenui/ebitenui"
	"github.com/ebitenui/ebitenui/image"
	"github.com/ebitenui/ebitenui/widget"
	"github.com/hajimehoshi/ebiten/v2"
	"golang.org/x/image/colornames"

	"github.com/gopher-co/td-game/models/general"
	"github.com/gopher-co/td-game/ui"
	"github.com/gopher-co/td-game/ui/font"
)

// loadMainMenuUI loads the main menu UI.
func (m *MenuState) loadMainMenuUI(widgets general.Widgets) *ebitenui.UI {
	bgImg := widgets[ui.MenuBackgroundImage]

	menuBackground := image.NewNineSliceSimple(bgImg, 0, 1)

	root := widget.NewContainer(
		widget.ContainerOpts.Layout(widget.NewGridLayout(
			widget.GridLayoutOpts.Columns(2),
			widget.GridLayoutOpts.Stretch([]bool{false, true}, []bool{true}),
			widget.GridLayoutOpts.Spacing(0, 0),
		)),
		widget.ContainerOpts.BackgroundImage(menuBackground),
	)

	buttons := m.btn(widgets)

	logo := widget.NewContainer(
		widget.ContainerOpts.Layout(widget.NewRowLayout(
			widget.RowLayoutOpts.Direction(widget.DirectionVertical),
			widget.RowLayoutOpts.Spacing(50),
			widget.RowLayoutOpts.Padding(widget.Insets{
				Top:    250,
				Left:   0,
				Right:  0,
				Bottom: 0,
			}),
		)),
	)

	// main menu image
	mainImgW := widgets[ui.MenuMainImage]
	mainImg := ebiten.NewImage(980, 920)
	geom := ebiten.GeoM{}
	geom.Scale(980/float64(mainImgW.Bounds().Dx()), 920/float64(mainImgW.Bounds().Dx()))
	mainImg.DrawImage(mainImgW, &ebiten.DrawImageOptions{GeoM: geom})
	mainImage := widget.NewGraphic(
		widget.GraphicOpts.Image(mainImg),
		widget.GraphicOpts.WidgetOpts(widget.WidgetOpts.LayoutData(
			widget.RowLayoutData{
				Stretch: true,
			},
		)),
	)

	logo.AddChild(mainImage)
	root.AddChild(buttons)
	root.AddChild(logo)

	return &ebitenui.UI{Container: root}
}

// btn returns the buttons.
func (m *MenuState) btn(widgets general.Widgets) *widget.Container {
	buttons := widget.NewContainer(
		widget.ContainerOpts.Layout(widget.NewGridLayout(
			widget.GridLayoutOpts.Columns(1),
			widget.GridLayoutOpts.Padding(widget.Insets{
				Top:    50,
				Left:   10,
				Right:  10,
				Bottom: 0,
			}),
			widget.GridLayoutOpts.Spacing(0, 50),
			widget.GridLayoutOpts.Padding(widget.Insets{Top: 50}),
			widget.GridLayoutOpts.Stretch([]bool{true}, []bool{false, false, false, false, false}),
		)),
		widget.ContainerOpts.BackgroundImage(image.NewNineSliceSimple(widgets[ui.MenuLeftSidebarImage], 0, 1)),
	)

	// logo image loading
	logoImg := widgets[ui.MenuMainLogoImage]

	logoImage := widget.NewGraphic(
		widget.GraphicOpts.Image(logoImg),
	)

	btn1 := widget.NewButton(
		widget.ButtonOpts.WidgetOpts(widget.WidgetOpts.MinSize(600, 100)),
		widget.ButtonOpts.Image(&widget.ButtonImage{
			Idle: image.NewNineSliceSimple(widgets[ui.MenuButtonPlayImage], 0, 1),
		}),
		widget.ButtonOpts.Text("PLAY!", font.TTF72, &widget.ButtonTextColor{Idle: color.White}),
		widget.ButtonOpts.ClickedHandler(func(_ *widget.ButtonClickedEventArgs) {
			m.UI = m.loadLevelMenuUI(widgets)
		}),
	)
	btn2 := widget.NewButton(
		widget.ButtonOpts.WidgetOpts(widget.WidgetOpts.MinSize(600, 100)),
		widget.ButtonOpts.Image(&widget.ButtonImage{
			Idle: image.NewNineSliceSimple(widgets[ui.MenuButtonReplaysImage], 0, 1),
		}),
		widget.ButtonOpts.Text("Replays", font.TTF72, &widget.ButtonTextColor{Idle: color.White}),
		widget.ButtonOpts.ClickedHandler(func(_ *widget.ButtonClickedEventArgs) {
			m.UI = m.loadReplaysMenuUI(widgets)
		}),
	)

	btn4 := widget.NewButton(
		widget.ButtonOpts.WidgetOpts(widget.WidgetOpts.MinSize(600, 100)),
		widget.ButtonOpts.Image(&widget.ButtonImage{
			Idle: image.NewNineSliceSimple(widgets[ui.MenuButtonExitImage], 0, 1),
		}),
		widget.ButtonOpts.Text("Exit", font.TTF72, &widget.ButtonTextColor{Idle: color.White}),
		widget.ButtonOpts.ClickedHandler(func(_ *widget.ButtonClickedEventArgs) {
			os.Exit(0)
		}),
	)

	btn3 := widget.NewButton(
		widget.ButtonOpts.WidgetOpts(widget.WidgetOpts.MinSize(600, 100)),
		widget.ButtonOpts.Image(&widget.ButtonImage{
			Idle: image.NewNineSliceSimple(widgets[ui.MenuButtonExitImage], 0, 1),
		}),
		widget.ButtonOpts.Text("Co-op", font.TTF72, &widget.ButtonTextColor{Idle: color.White}),
		widget.ButtonOpts.ClickedHandler(func(_ *widget.ButtonClickedEventArgs) {
			m.UI = m.loadCoopMenuUI(widgets)
		}),
	)

	buttons.AddChild(logoImage)
	buttons.AddChild(btn1)
	buttons.AddChild(btn2)
	buttons.AddChild(btn3)
	buttons.AddChild(btn4)
	return buttons
}

// loadLevelMenuUI loads the level menu UI.
func (m *MenuState) loadLevelMenuUI(widgets general.Widgets) *ebitenui.UI {
	bgImg := widgets[ui.MenuBackgroundImage]
	menuBackground := image.NewNineSliceSimple(bgImg, 0, 1)

	root := widget.NewContainer(
		widget.ContainerOpts.Layout(widget.NewGridLayout(
			widget.GridLayoutOpts.Columns(1),
			widget.GridLayoutOpts.Stretch([]bool{true}, []bool{false, true}),
			widget.GridLayoutOpts.Spacing(0, 0),
		)),
		widget.ContainerOpts.BackgroundImage(menuBackground),
	)

	infoContainer := widget.NewContainer(
		widget.ContainerOpts.Layout(widget.NewGridLayout(
			widget.GridLayoutOpts.Columns(2),
			widget.GridLayoutOpts.Stretch([]bool{false, true}, []bool{true}),
		)),
	)

	//defer font.TTF128.Close()

	backBtn := widget.NewButton(
		widget.ButtonOpts.Image(&widget.ButtonImage{Idle: image.NewNineSliceSimple(widgets[ui.LevelMenuBackButtonImage], 0, 1)}),
		widget.ButtonOpts.Text("<", font.TTF128, &widget.ButtonTextColor{Idle: color.White}),
		widget.ButtonOpts.ClickedHandler(func(args *widget.ButtonClickedEventArgs) {
			m.UI = m.loadMainMenuUI(widgets)
		}),
		widget.ButtonOpts.TextPadding(widget.Insets{
			Left:  35,
			Right: 35,
		}),
	)

	text := ""
	if len(m.State.LevelsComplete) == len(m.Levels) {
		text = "YAYY! YOU'VE COMPLETED ALL THE LEVELS"
	} else {
		text = fmt.Sprintf("Completed %d/%d levels", len(m.State.LevelsComplete), len(m.Levels))
	}

	//defer font.TTF64.Close()
	textCompleted := widget.NewText(
		widget.TextOpts.Text(text, font.TTF64, color.White),
		widget.TextOpts.Position(widget.TextPositionCenter, widget.TextPositionCenter),
		widget.TextOpts.ProcessBBCode(true),
	)

	infoContainer.AddChild(backBtn)
	infoContainer.AddChild(textCompleted)

	root.AddChild(infoContainer)
	root.AddChild(m.loadScrollingLevels(widgets))

	return &ebitenui.UI{Container: root}
}

// loadScrollingLevels loads the scrolling levels.
func (m *MenuState) loadScrollingLevels(_ general.Widgets) *widget.Container {
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

	// blackImg := image.NewNineSliceColor(color.Black)
	for i, k := range levels {
		k := k

		cont := widget.NewContainer(
			widget.ContainerOpts.WidgetOpts(widget.WidgetOpts.MinSize(400, 900)),
			widget.ContainerOpts.Layout(widget.NewGridLayout(
				widget.GridLayoutOpts.Columns(1),
				widget.GridLayoutOpts.Stretch([]bool{true}, []bool{false, true, false}),
			)),
		)
		text1 := widget.NewText(
			widget.TextOpts.Text(fmt.Sprintf("Level %d", i+1), font.TTF72, color.White),
			widget.TextOpts.Position(widget.TextPositionCenter, widget.TextPositionStart),
		)
		text2 := widget.NewText(
			widget.TextOpts.MaxWidth(400),
			widget.TextOpts.Text(m.Levels[k].LevelName, font.TTF36, color.White),
		)
		btn := widget.NewButton(
			widget.ButtonOpts.Image(&widget.ButtonImage{Idle: image.NewNineSliceColor(colornames.Beige)}),
			widget.ButtonOpts.Text("Play", font.TTF72, &widget.ButtonTextColor{Idle: color.Black}),
			widget.ButtonOpts.WidgetOpts(widget.WidgetOpts.LayoutData(widget.GridLayoutData{
				VerticalPosition: widget.GridLayoutPositionEnd,
			})),
			widget.ButtonOpts.ClickedHandler(func(args *widget.ButtonClickedEventArgs) {
				m.Ended = true
				m.Next = k
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
		widget.SliderOpts.MinMax(0, 100),
		widget.SliderOpts.PageSizeFunc(pageSizeFunc),
		// On change update scroll location based on the Slider's value
		widget.SliderOpts.ChangedHandler(func(args *widget.SliderChangedEventArgs) {
			scrollContainer.ScrollLeft = float64(args.Slider.Current) / 100
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
