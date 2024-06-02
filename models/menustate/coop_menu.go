package menustate

import (
	"fmt"
	"image/color"
	"regexp"

	"github.com/ebitenui/ebitenui"
	"github.com/ebitenui/ebitenui/image"
	"github.com/ebitenui/ebitenui/widget"

	"github.com/gopher-co/td-game/models/general"
	"github.com/gopher-co/td-game/ui"
	"github.com/gopher-co/td-game/ui/font"
)

var valid = regexp.MustCompile(`^[a-zA-Z0-9_. ]*$`).MatchString

func (m *MenuState) loadCoopMenuUI(widgets general.Widgets) *ebitenui.UI {
	bgImg := widgets[ui.MenuBackgroundImage]
	menuBackground := image.NewNineSliceSimple(bgImg, 0, 1)

	backBtn := widget.NewButton(
		widget.ButtonOpts.Image(&widget.ButtonImage{Idle: image.NewNineSliceSimple(widgets[ui.LevelMenuBackButtonImage], 0, 1)}),
		widget.ButtonOpts.Text("<", font.TTF128, &widget.ButtonTextColor{Idle: color.White}),
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

	createContainer := widget.NewContainer(
		widget.ContainerOpts.Layout(widget.NewGridLayout(
			widget.GridLayoutOpts.Columns(2),
			widget.GridLayoutOpts.Stretch([]bool{true, true, false}, []bool{true}),
			widget.GridLayoutOpts.Spacing(0, 0),
		)),
	)

	fields := widget.NewContainer(
		widget.ContainerOpts.Layout(widget.NewRowLayout(
			widget.RowLayoutOpts.Direction(widget.DirectionVertical),
			widget.RowLayoutOpts.Spacing(10),
		)),
	)
	var sName, sLevel string
	name := widget.NewTextInput(
		widget.TextInputOpts.Validation(func(newInputText string) (bool, *string) {
			if valid(newInputText) && len(newInputText) < 20 {
				return true, &newInputText
			}
			return false, nil
		}),
		widget.TextInputOpts.Face(font.TTF20),
		widget.TextInputOpts.Color(&widget.TextInputColor{
			Idle:          color.NRGBA{254, 255, 255, 255},
			Disabled:      color.NRGBA{R: 200, G: 200, B: 200, A: 255},
			Caret:         color.NRGBA{254, 255, 255, 255},
			DisabledCaret: color.NRGBA{R: 200, G: 200, B: 200, A: 255},
		}),
		widget.TextInputOpts.Placeholder("Nick"),
		widget.TextInputOpts.Image(&widget.TextInputImage{
			Idle:     image.NewNineSliceColor(color.NRGBA{R: 100, G: 100, B: 100, A: 255}),
			Disabled: image.NewNineSliceColor(color.NRGBA{R: 100, G: 100, B: 100, A: 255}),
		}),
		widget.TextInputOpts.CaretOpts(
			widget.CaretOpts.Size(font.TTF20, 2),
		),

		//This is called whenver there is a change to the text
		widget.TextInputOpts.ChangedHandler(func(args *widget.TextInputChangedEventArgs) {
			sName = args.InputText
		}),
		widget.TextInputOpts.WidgetOpts(
			//Set the layout information to center the textbox in the parent
			widget.WidgetOpts.LayoutData(widget.RowLayoutData{
				MaxWidth: 250,
				Stretch:  true,
			}),
		),
		widget.TextInputOpts.Padding(widget.Insets{
			Top:    20,
			Left:   10,
			Right:  10,
			Bottom: 20,
		}),
	)
	level := widget.NewTextInput(
		widget.TextInputOpts.Validation(func(newInputText string) (bool, *string) {
			if valid(newInputText) {
				return true, &newInputText
			}
			return false, nil
		}),
		widget.TextInputOpts.Face(font.TTF20),
		widget.TextInputOpts.Color(&widget.TextInputColor{
			Idle:          color.NRGBA{254, 255, 255, 255},
			Disabled:      color.NRGBA{R: 200, G: 200, B: 200, A: 255},
			Caret:         color.NRGBA{254, 255, 255, 255},
			DisabledCaret: color.NRGBA{R: 200, G: 200, B: 200, A: 255},
		}),
		widget.TextInputOpts.Placeholder("Level name"),
		widget.TextInputOpts.Image(&widget.TextInputImage{
			Idle:     image.NewNineSliceColor(color.NRGBA{R: 100, G: 100, B: 100, A: 255}),
			Disabled: image.NewNineSliceColor(color.NRGBA{R: 100, G: 100, B: 100, A: 255}),
		}),
		widget.TextInputOpts.CaretOpts(
			widget.CaretOpts.Size(font.TTF20, 2),
		),

		//This is called whenver there is a change to the text
		widget.TextInputOpts.ChangedHandler(func(args *widget.TextInputChangedEventArgs) {
			sLevel = args.InputText
		}),
		widget.TextInputOpts.Padding(widget.Insets{
			Top:    20,
			Left:   10,
			Right:  10,
			Bottom: 20,
		}),
		widget.TextInputOpts.WidgetOpts(
			//Set the layout information to center the textbox in the parent
			widget.WidgetOpts.LayoutData(widget.RowLayoutData{
				Stretch:   true,
				MaxWidth:  250,
				MaxHeight: 0,
			}),
		),
	)

	submit := widget.NewButton(
		widget.ButtonOpts.Text("Create", font.TTF32, &widget.ButtonTextColor{Idle: color.NRGBA{255, 255, 255, 255}}),
		widget.ButtonOpts.ClickedHandler(func(args *widget.ButtonClickedEventArgs) {
			fmt.Println(sName, sLevel)
		}),
		widget.ButtonOpts.Image(&widget.ButtonImage{Idle: image.NewNineSliceColor(color.Black)}),
	)

	fields.AddChild(name)
	fields.AddChild(level)
	fields.AddChild(submit)
	createContainer.AddChild(fields)
	root.AddChild(backBtn)
	root.AddChild(createContainer)
	//root.AddChild(m.loadScrollingReplays(widgets))

	return &ebitenui.UI{Container: root}
}
