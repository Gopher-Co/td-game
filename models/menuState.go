package models

import (
	"image/color"
	"log"
	"math"
	"os"
	"runtime"

	"github.com/ebitenui/ebitenui"
	"github.com/ebitenui/ebitenui/image"
	"github.com/ebitenui/ebitenui/widget"
	"github.com/golang/freetype/truetype"
	"github.com/hajimehoshi/ebiten/v2"
	"golang.org/x/image/colornames"
	"golang.org/x/image/font"
	"golang.org/x/image/font/gofont/goregular"

	"github.com/gopher-co/td-game/ui"
)

// MenuState is a struct that represents the state of the menu.
type MenuState struct {
	Levels []*LevelConfig
	Ended  bool
	UI     *ebitenui.UI
}

func NewMenuState(configs []*LevelConfig, widgets Widgets) *MenuState {
	ms := &MenuState{
		Levels: configs,
		Ended:  false,
		UI:     nil,
	}
	ms.loadUI(widgets)

	return ms
}

func (m *MenuState) Draw(image *ebiten.Image) {
	m.UI.Draw(image)
}

func (m *MenuState) Update() error {
	m.UI.Update()
	return nil
}

func (m *MenuState) loadUI(widgets Widgets) {
	m.UI = m.loadMainMenuUI(widgets)
}

func (m *MenuState) End() bool {
	return m.Ended
}

func (m *MenuState) NextState() State {
	return nil
}

func mustLoadFont(size float64) font.Face {
	ttfFont, err := truetype.Parse(goregular.TTF)
	if err != nil {
		panic(err)
	}

	return truetype.NewFace(ttfFont, &truetype.Options{
		Size:    size,
		DPI:     72,
		Hinting: font.HintingFull,
	})
}

func (m *MenuState) loadMainMenuUI(widgets Widgets) *ebitenui.UI {
	mainMenuImg := widgets[ui.MenuMainImage]
	bgImg := widgets[ui.MenuBackgroundImage]

	img := ebiten.NewImage(1280, 720)
	geom := ebiten.GeoM{}
	geom.Scale(1280./float64(mainMenuImg.Bounds().Dx()), 720./float64(mainMenuImg.Bounds().Dy()))
	img.DrawImage(widgets[ui.MenuMainImage], &ebiten.DrawImageOptions{GeoM: geom})

	menuBackground := image.NewNineSlice(bgImg, [3]int{0, bgImg.Bounds().Dx(), 0}, [3]int{0, bgImg.Bounds().Dy(), 0})

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
		widget.ContainerOpts.Layout(widget.NewGridLayout(
			widget.GridLayoutOpts.Columns(1),
			widget.GridLayoutOpts.Stretch([]bool{true}, []bool{false}),
			widget.GridLayoutOpts.Padding(widget.Insets{
				Top:    250,
				Left:   0,
				Right:  0,
				Bottom: 0,
			}),
		)),
		//widget.ContainerOpts.BackgroundImage(image.NewNineSliceColor(colornames.Red)),
	)

	logoImage := widget.NewGraphic(
		widget.GraphicOpts.Image(img),
		widget.GraphicOpts.WidgetOpts(widget.WidgetOpts.LayoutData(
			widget.GridLayoutData{
				HorizontalPosition: widget.GridLayoutPositionCenter,
			},
		)),
	)
	logo.AddChild(logoImage)
	root.AddChild(buttons)
	root.AddChild(logo)

	return &ebitenui.UI{Container: root}
}

func (m *MenuState) btn(widgets Widgets) *widget.Container {
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
			widget.GridLayoutOpts.Stretch([]bool{true}, []bool{false, false, false, false}),
		)),
		widget.ContainerOpts.BackgroundImage(image.NewNineSliceColor(colornames.Royalblue)),
	)

	fnt := mustLoadFont(72)

	text := widget.NewText(
		widget.TextOpts.Text("Go Build,\nGo Defend!", fnt, color.White),
	)

	btn1 := widget.NewButton(
		widget.ButtonOpts.WidgetOpts(widget.WidgetOpts.MinSize(600, 100)),
		widget.ButtonOpts.Image(&widget.ButtonImage{
			Idle: image.NewNineSlice(widgets[ui.MenuButtonPlayImage], [3]int{0, 1300, 0}, [3]int{0, 800, 0}),
		}),
		widget.ButtonOpts.Text("PLAY!", fnt, &widget.ButtonTextColor{Idle: color.White}),
		widget.ButtonOpts.ClickedHandler(func(_ *widget.ButtonClickedEventArgs) {
			m.UI = m.loadLevelMenuUI(widgets)
			runtime.GC()
		}),
	)
	btn2 := widget.NewButton(
		widget.ButtonOpts.WidgetOpts(widget.WidgetOpts.MinSize(600, 100)),
		widget.ButtonOpts.Image(&widget.ButtonImage{
			Idle: image.NewNineSlice(widgets[ui.MenuButtonPlayImage], [3]int{0, 1300, 0}, [3]int{0, 800, 0}),
		}),
		widget.ButtonOpts.Text("Replays", fnt, &widget.ButtonTextColor{Idle: color.White}),
	)
	btn3 := widget.NewButton(
		widget.ButtonOpts.WidgetOpts(widget.WidgetOpts.MinSize(600, 100)),
		widget.ButtonOpts.Image(&widget.ButtonImage{
			Idle: image.NewNineSlice(widgets[ui.MenuButtonPlayImage], [3]int{0, 1300, 0}, [3]int{0, 800, 0}),
		}),
		widget.ButtonOpts.Text("Exit", fnt, &widget.ButtonTextColor{Idle: color.White}),
		widget.ButtonOpts.ClickedHandler(func(_ *widget.ButtonClickedEventArgs) {
			os.Exit(0)
		}),
	)

	buttons.AddChild(text)
	buttons.AddChild(btn1)
	buttons.AddChild(btn2)
	buttons.AddChild(btn3)
	return buttons
}

func (m *MenuState) loadLevelMenuUI(widgets Widgets) *ebitenui.UI {
	bgImg := widgets[ui.MenuBackgroundImage]
	menuBackground := image.NewNineSlice(bgImg, [3]int{0, bgImg.Bounds().Dx(), 0}, [3]int{0, bgImg.Bounds().Dy(), 0})

	backBtn := widget.NewButton(
		widget.ButtonOpts.Image(&widget.ButtonImage{Idle: image.NewNineSliceColor(colornames.Aqua)}),
		widget.ButtonOpts.Text("<", mustLoadFont(128), &widget.ButtonTextColor{Idle: color.White}),
		widget.ButtonOpts.ClickedHandler(func(args *widget.ButtonClickedEventArgs) {
			m.UI = m.loadMainMenuUI(widgets)
			log.Println("ui")
			runtime.GC()
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
	root.AddChild(m.loadScrollingLevels(widgets))

	return &ebitenui.UI{Container: root}
}

func (m *MenuState) loadScrollingLevels(widgets Widgets) *widget.Container {
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
	for _, v := range m.Levels {
		cont := widget.NewContainer(
			widget.ContainerOpts.WidgetOpts(widget.WidgetOpts.LayoutData(widget.GridLayoutData{MaxWidth: 300})),
			widget.ContainerOpts.Layout(widget.NewGridLayout(
				widget.GridLayoutOpts.Columns(1),
				widget.GridLayoutOpts.Stretch([]bool{true}, []bool{false, true, false}),
			)),
		)
		text1 := widget.NewText(
			widget.TextOpts.Text("aboba", mustLoadFont(72), color.White),
		)
		text2 := widget.NewText(
			widget.TextOpts.Text(v.LevelName, mustLoadFont(72), color.White),
			widget.TextOpts.WidgetOpts(widget.WidgetOpts.MinSize(0, 600)),
		)
		btn := widget.NewButton(
			widget.ButtonOpts.Image(&widget.ButtonImage{Idle: image.NewNineSliceColor(colornames.Beige)}),
			widget.ButtonOpts.Text("Play", mustLoadFont(72), &widget.ButtonTextColor{Idle: color.Black}),
			widget.ButtonOpts.WidgetOpts(widget.WidgetOpts.LayoutData(widget.GridLayoutData{
				VerticalPosition: widget.GridLayoutPositionEnd,
			})),
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
			Idle: image.NewNineSliceColor(color.NRGBA{0x13, 0x1a, 0x22, 0xff}),
			Mask: image.NewNineSliceColor(color.NRGBA{0x13, 0x1a, 0x22, 0xff}),
		}),
	)
	root.AddChild(scrollContainer)

	pageSizeFunc := func() int {
		return int(math.Round(float64(scrollContainer.ContentRect().Dx()) / float64(content.GetWidget().Rect.Dx()) * 1000))
	}

	vSlider := widget.NewSlider(
		widget.SliderOpts.Direction(widget.DirectionHorizontal),
		widget.SliderOpts.MinMax(0, 1000),
		widget.SliderOpts.PageSizeFunc(pageSizeFunc),
		//On change update scroll location based on the Slider's value
		widget.SliderOpts.ChangedHandler(func(args *widget.SliderChangedEventArgs) {
			scrollContainer.ScrollLeft = float64(args.Slider.Current) / 1000
		}),
		widget.SliderOpts.Images(
			// Set the track images
			&widget.SliderTrackImage{
				Idle:  image.NewNineSliceColor(color.NRGBA{100, 100, 100, 255}),
				Hover: image.NewNineSliceColor(color.NRGBA{100, 100, 100, 255}),
			},
			// Set the handle images
			&widget.ButtonImage{
				Idle:    image.NewNineSliceColor(color.NRGBA{255, 100, 100, 255}),
				Hover:   image.NewNineSliceColor(color.NRGBA{255, 100, 100, 255}),
				Pressed: image.NewNineSliceColor(color.NRGBA{255, 100, 100, 255}),
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
