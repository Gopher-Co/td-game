package coopstate

import (
	"fmt"
	"image/color"

	image2 "github.com/ebitenui/ebitenui/image"
	"github.com/ebitenui/ebitenui/widget"
	"golang.org/x/image/colornames"

	"github.com/gopher-co/td-game/models/general"
	"github.com/gopher-co/td-game/models/ingame"
	"github.com/gopher-co/td-game/ui/font"
)

// newTowerMenuUI creates a new tower menu UI.
func (s *GameState) newTowerMenuUI(widgets general.Widgets) *widget.Container {
	root := widget.NewContainer(
		widget.ContainerOpts.Layout(widget.NewGridLayout(
			widget.GridLayoutOpts.Columns(1),
			widget.GridLayoutOpts.Stretch([]bool{true}, []bool{false, false, false, false}),
		)),
	)

	cInfo = root

	info := s.textContainer(widgets)
	upgrades := s.upgradesContainer(widgets)
	tuning := s.tuningContainer(widgets)
	sell := s.sellContainer(widgets)

	root.AddChild(info)
	root.AddChild(upgrades)
	root.AddChild(tuning)
	root.AddChild(sell)

	return root
}

// textContainer creates a container that contains the name of the tower.
func (s *GameState) textContainer(_ general.Widgets) *widget.Container {
	root := widget.NewContainer(
		widget.ContainerOpts.Layout(widget.NewGridLayout(
			widget.GridLayoutOpts.Columns(1),
			widget.GridLayoutOpts.Stretch([]bool{true}, []bool{false}),
			widget.GridLayoutOpts.Padding(widget.Insets{
				Top: 50,
			}),
		)),
	)

	name := widget.NewText(
		widget.TextOpts.Text("NAME", font.TTF64, color.White),
		widget.TextOpts.MaxWidth(400),
		widget.TextOpts.Position(widget.TextPositionCenter, widget.TextPositionStart),
	)

	s.uiUpdater.Append(func() {
		if s.chosenTower == nil {
			return
		}
		name.Label = s.chosenTower.Name
	})

	root.AddChild(name)

	return root
}

// upgradesContainer creates a container that contains the upgrades of the tower.
func (s *GameState) upgradesContainer(_ general.Widgets) *widget.Container {
	root := widget.NewContainer(
		widget.ContainerOpts.Layout(widget.NewGridLayout(
			widget.GridLayoutOpts.Columns(1),
			widget.GridLayoutOpts.Stretch([]bool{true}, []bool{false, false, false}),
		)),
	)

	btn := widget.NewButton(
		widget.ButtonOpts.Image(&widget.ButtonImage{
			Idle:     image2.NewNineSliceColor(color.RGBA{R: 0x99, G: 0xe7, B: 0xa9, A: 0xff}),
			Hover:    image2.NewNineSliceColor(color.RGBA{R: 0xa9, G: 0xee, B: 0xae, A: 0xff}),
			Pressed:  image2.NewNineSliceColor(color.RGBA{R: 0x89, G: 0xd7, B: 0x99, A: 0xff}),
			Disabled: image2.NewNineSliceColor(color.RGBA{R: 0x66, G: 0x05, B: 0x28, A: 0xff}),
		}),
		widget.ButtonOpts.Text("UPGRADE", font.TTF32, &widget.ButtonTextColor{
			Idle:     color.White,
			Disabled: color.Black,
		}),
		widget.ButtonOpts.WidgetOpts(widget.WidgetOpts.MinSize(0, 100)),
		widget.ButtonOpts.ClickedHandler(s.handleUpgrade),
	)

	level := widget.NewText(
		widget.TextOpts.Text("Level", font.TTF48, color.White),
	)

	info := s.textUpgradeInfo()

	s.uiUpdater.Append(func() {
		if s.chosenTower == nil {
			return
		}

		level.Label = fmt.Sprintf("Level %d", s.chosenTower.UpgradesBought+1)

		// all the upgrades are bought
		if s.chosenTower.UpgradesBought >= len(s.chosenTower.Upgrades) {
			c := info.Children()
			insertValues(c[0].(*widget.Text), s.chosenTower.Damage, 0, "Damage")
			insertValues(c[1].(*widget.Text), int(s.chosenTower.Radius), 0, "Radius")
			insertValues(c[2].(*widget.Text), s.chosenTower.SpeedAttack, 0, "Speed")
			insertValues(c[3].(*widget.Text), int(s.chosenTower.ProjectileVrms), 0, "ProjSpeed")

			btn.Text().Label = "SOLD OUT"
			btn.GetWidget().Disabled = true

			return
		}

		c := info.Children()
		u := s.chosenTower.Upgrades[s.chosenTower.UpgradesBought]
		insertValues(c[0].(*widget.Text), s.chosenTower.Damage, u.DeltaDamage, "Damage")
		insertValues(c[1].(*widget.Text), int(s.chosenTower.Radius), int(u.DeltaRadius), "Radius")
		insertValues(c[2].(*widget.Text), s.chosenTower.SpeedAttack, u.DeltaSpeedAttack, "Speed")
		insertValues(c[3].(*widget.Text), int(s.chosenTower.ProjectileVrms), 0, "ProjSpeed")

		openLevel := s.chosenTower.Upgrades[s.chosenTower.UpgradesBought].OpenLevel
		_, ok := s.PlayerState.LevelsComplete[openLevel]

		if !ok && openLevel != "" {
			btn.Text().Label = "Complete level to unlock:\n" + s.chosenTower.Upgrades[s.chosenTower.UpgradesBought].OpenLevel
		} else {
			btn.Text().Label = fmt.Sprintf("UPGRADE ($%d)", s.chosenTower.Upgrades[s.chosenTower.UpgradesBought].Price)
		}

		if s.PlayerMapState.Money < s.chosenTower.Upgrades[s.chosenTower.UpgradesBought].Price ||
			!ok && openLevel != "" {
			btn.GetWidget().Disabled = true
		} else {
			btn.GetWidget().Disabled = false
		}
	})

	root.AddChild(level)
	root.AddChild(btn)
	root.AddChild(info)

	return root
}

// textUpgradeInfo creates a container that contains the info about the tower.
func (s *GameState) textUpgradeInfo() *widget.Container {
	root := widget.NewContainer(
		widget.ContainerOpts.Layout(widget.NewRowLayout(
			widget.RowLayoutOpts.Direction(widget.DirectionVertical),
			widget.RowLayoutOpts.Spacing(5),
		)),
	)

	textDamage := widget.NewText(
		widget.TextOpts.Text(fmt.Sprintf(`Damage: %d[color=00FF00]+0[/color]`, 0), font.TTF40, color.White),
		widget.TextOpts.ProcessBBCode(true),
	)
	textRadius := widget.NewText(
		widget.TextOpts.Text(fmt.Sprintf(`Radius: %.0f[color=FF0000]+0[/color]`, 0.), font.TTF40, color.White),
		widget.TextOpts.ProcessBBCode(true),
	)
	textSpeed := widget.NewText(
		widget.TextOpts.Text(fmt.Sprintf(`Speed: %.0f[color=FF0000]+0[/color]`, 0.), font.TTF40, color.White),
		widget.TextOpts.ProcessBBCode(true),
	)
	textProjSpeed := widget.NewText(
		widget.TextOpts.Text(fmt.Sprintf(`ProjSpeed: %.0f[color=AAAAAA]+0[/color]`, 0.), font.TTF40, color.White),
		widget.TextOpts.ProcessBBCode(true),
	)

	root.AddChild(textDamage)
	root.AddChild(textRadius)
	root.AddChild(textSpeed)
	root.AddChild(textProjSpeed)

	return root
}

// insertValues inserts values into the text.
func insertValues(c *widget.Text, v, deltav int, s string) {
	if deltav > 0 {
		c.Label = fmt.Sprintf("%s: %d[color=00FF00]+%d[/color]", s, v, deltav)
	} else if deltav < 0 {
		c.Label = fmt.Sprintf("%s: %d[color=FF0000]%d[/color]", s, v, deltav)
	} else {
		c.Label = fmt.Sprintf("%s: %d", s, v)
	}
}

// tuningContainer creates a container that contains the tuning of the tower.
func (s *GameState) tuningContainer(_ general.Widgets) *widget.Container {
	root := widget.NewContainer(
		widget.ContainerOpts.Layout(widget.NewGridLayout(
			widget.GridLayoutOpts.Columns(1),
			widget.GridLayoutOpts.Stretch([]bool{true}, []bool{false, true}),
		)),
	)

	btnTurn := widget.NewButton(
		widget.ButtonOpts.Image(&widget.ButtonImage{
			Idle: image2.NewNineSliceColor(colornames.Lawngreen),
		}),
		widget.ButtonOpts.Text("ON", font.TTF54, &widget.ButtonTextColor{
			Idle: color.White,
		}),
		widget.ButtonOpts.ClickedHandler(s.handleTurning),
	)

	s.uiUpdater.Append(func() {
		if s.chosenTower == nil {
			return
		}

		if s.chosenTower.State.IsTurnedOn {
			btnTurn.Text().Label = "ON"
			btnTurn.Image = &widget.ButtonImage{
				Idle: image2.NewNineSliceColor(colornames.Lawngreen),
			}
		} else {
			btnTurn.Text().Label = "OFF"
			btnTurn.Image = &widget.ButtonImage{
				Idle: image2.NewNineSliceColor(colornames.Indianred),
			}
		}
	})

	root.AddChild(btnTurn)
	root.AddChild(s.radio())

	return root
}

// radio creates a radio group.
func (s *GameState) radio() *widget.Container {

	root := widget.NewContainer(
		widget.ContainerOpts.Layout(widget.NewRowLayout(
			widget.RowLayoutOpts.Direction(widget.DirectionVertical),
			widget.RowLayoutOpts.Spacing(20),
			widget.RowLayoutOpts.Padding(widget.Insets{Top: 40}),
		)),
	)

	first := widget.NewButton(
		widget.ButtonOpts.TextPadding(widget.Insets{
			Top:    10,
			Left:   0,
			Right:  0,
			Bottom: 10,
		}),
		widget.ButtonOpts.WidgetOpts(widget.WidgetOpts.LayoutData(widget.RowLayoutData{
			Position: widget.RowLayoutPositionCenter,
			Stretch:  true,
		})),
		widget.ButtonOpts.Text("First", font.TTF40, &widget.ButtonTextColor{
			Idle: color.White,
		}),
		widget.ButtonOpts.Image(&widget.ButtonImage{
			Idle:    image2.NewNineSliceColor(color.RGBA{R: 0x7f, G: 0x27, B: 0xd7, A: 0xff}),
			Hover:   image2.NewNineSliceColor(color.RGBA{R: 0x9a, G: 0x3b, B: 0xea, A: 0xff}),
			Pressed: image2.NewNineSliceColor(color.RGBA{R: 0x6a, G: 0x16, B: 0xc2, A: 0xff}),
		}),
		widget.ButtonOpts.ClickedHandler(s.handleTuneFirst),
	)

	strong := widget.NewButton(
		widget.ButtonOpts.TextPadding(widget.Insets{
			Top:    10,
			Left:   0,
			Right:  0,
			Bottom: 10,
		}),
		widget.ButtonOpts.WidgetOpts(widget.WidgetOpts.LayoutData(widget.RowLayoutData{
			Position: widget.RowLayoutPositionCenter,
			Stretch:  true,
		})),
		widget.ButtonOpts.Text("Strong", font.TTF40, &widget.ButtonTextColor{
			Idle: color.White,
		}),
		widget.ButtonOpts.Image(&widget.ButtonImage{
			Idle:    image2.NewNineSliceColor(color.RGBA{R: 0x7f, G: 0x27, B: 0xd7, A: 0xff}),
			Hover:   image2.NewNineSliceColor(color.RGBA{R: 0x9a, G: 0x3b, B: 0xea, A: 0xff}),
			Pressed: image2.NewNineSliceColor(color.RGBA{R: 0x6a, G: 0x16, B: 0xc2, A: 0xff}),
		}),
		widget.ButtonOpts.ClickedHandler(s.handleTuneStrong),
	)

	weak := widget.NewButton(
		widget.ButtonOpts.TextPadding(widget.Insets{
			Top:    10,
			Left:   0,
			Right:  0,
			Bottom: 10,
		}),
		widget.ButtonOpts.WidgetOpts(widget.WidgetOpts.LayoutData(widget.RowLayoutData{
			Position: widget.RowLayoutPositionCenter,
			Stretch:  true,
		})),
		widget.ButtonOpts.Text("Weak", font.TTF40, &widget.ButtonTextColor{
			Idle: color.White,
		}),
		widget.ButtonOpts.Image(&widget.ButtonImage{
			Idle:    image2.NewNineSliceColor(color.RGBA{R: 0x7f, G: 0x27, B: 0xd7, A: 0xff}),
			Hover:   image2.NewNineSliceColor(color.RGBA{R: 0x9a, G: 0x3b, B: 0xea, A: 0xff}),
			Pressed: image2.NewNineSliceColor(color.RGBA{R: 0x6a, G: 0x16, B: 0xc2, A: 0xff}),
		}),
		widget.ButtonOpts.ClickedHandler(s.handleTuneWeak),
	)

	root.AddChild(first)
	root.AddChild(strong)
	root.AddChild(weak)

	r := widget.NewRadioGroup(
		widget.RadioGroupOpts.Elements(first, strong, weak),
	)

	s.uiUpdater.Append(func() {
		if s.chosenTower == nil {
			return
		}

		switch s.chosenTower.State.AimType {
		case ingame.First:
			r.SetActive(first)
		case ingame.Strongest:
			r.SetActive(strong)
		case ingame.Weakest:
			r.SetActive(weak)
		}
	})

	return root
}

// sellContainer creates a container that contains the sell button.
func (s *GameState) sellContainer(_ general.Widgets) *widget.Container {
	root := widget.NewContainer(
		widget.ContainerOpts.Layout(widget.NewGridLayout(
			widget.GridLayoutOpts.Columns(1),
			widget.GridLayoutOpts.Stretch([]bool{true}, []bool{true}),
			widget.GridLayoutOpts.Padding(widget.Insets{Top: 40}),
		)),
	)

	btnSell := widget.NewButton(
		widget.ButtonOpts.Image(&widget.ButtonImage{
			Idle: image2.NewNineSliceColor(color.RGBA{R: 0xff, G: 0x66, B: 0x66, A: 0xff}),
		}),
		widget.ButtonOpts.Text("SELL", font.TTF64, &widget.ButtonTextColor{
			Idle: color.White,
		}),
		widget.ButtonOpts.ClickedHandler(s.handleSell),
		widget.ButtonOpts.TextPadding(widget.Insets{Top: 10, Bottom: 10}),
	)

	root.AddChild(btnSell)

	return root
}
