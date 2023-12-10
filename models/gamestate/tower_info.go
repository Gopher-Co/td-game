package gamestate

import (
	"context"
	"fmt"
	"image/color"
	"strconv"
	"time"

	image2 "github.com/ebitenui/ebitenui/image"
	"github.com/ebitenui/ebitenui/widget"
	"github.com/hajimehoshi/ebiten/v2"
	"golang.org/x/image/colornames"

	"github.com/gopher-co/td-game/models/general"
	"github.com/gopher-co/td-game/models/ingame"
	"github.com/gopher-co/td-game/replay"
	"github.com/gopher-co/td-game/ui/loaders"
)

func (s *GameState) newTowerMenuUI(ctx context.Context, widgets general.Widgets) *widget.Container {
	root := widget.NewContainer(
		widget.ContainerOpts.Layout(widget.NewGridLayout(
			widget.GridLayoutOpts.Columns(1),
			widget.GridLayoutOpts.Stretch([]bool{true}, []bool{false, false, false, false}),
		)),
	)

	cInfo = root

	info := s.textContainer(widgets)
	upgrades := s.upgradesContainer(ctx, widgets)
	tuning := s.tuningContainer(ctx, widgets)
	sell := s.sellContainer(widgets)

	root.AddChild(info)
	root.AddChild(upgrades)
	root.AddChild(tuning)
	root.AddChild(sell)

	return root
}

func (s *GameState) textContainer(widgets general.Widgets) *widget.Container {
	root := widget.NewContainer(
		widget.ContainerOpts.Layout(widget.NewGridLayout(
			widget.GridLayoutOpts.Columns(1),
			widget.GridLayoutOpts.Stretch([]bool{true}, []bool{false}),
			widget.GridLayoutOpts.Padding(widget.Insets{
				Top: 50,
			}),
		)),
	)

	ttf64 := loaders.FontTrueType(64)
	defer ttf64.Close()
	name := widget.NewText(
		widget.TextOpts.Text("NAME", ttf64, color.White),
		widget.TextOpts.MaxWidth(400),
		widget.TextOpts.Position(widget.TextPositionCenter, widget.TextPositionStart),
	)

	root.AddChild(name)

	return root
}

func (s *GameState) upgradesContainer(ctx context.Context, widgets general.Widgets) *widget.Container {
	root := widget.NewContainer(
		widget.ContainerOpts.Layout(widget.NewGridLayout(
			widget.GridLayoutOpts.Columns(1),
			widget.GridLayoutOpts.Stretch([]bool{true}, []bool{false, false, false}),
		)),
	)

	var checkBlock func()

	ttf32 := loaders.FontTrueType(32)
	ttf48 := loaders.FontTrueType(48)

	btn := widget.NewButton(
		widget.ButtonOpts.Image(&widget.ButtonImage{
			Idle:     image2.NewNineSliceColor(color.RGBA{0x99, 0xe7, 0xa9, 0xff}),
			Hover:    image2.NewNineSliceColor(color.RGBA{0xa9, 0xee, 0xae, 0xff}),
			Pressed:  image2.NewNineSliceColor(color.RGBA{0x89, 0xd7, 0x99, 0xff}),
			Disabled: image2.NewNineSliceColor(color.RGBA{0x66, 0x05, 0x28, 0xff}),
		}),
		widget.ButtonOpts.Text("UPGRADE", ttf32, &widget.ButtonTextColor{
			Idle:     color.White,
			Disabled: color.Black,
		}),
		widget.ButtonOpts.WidgetOpts(widget.WidgetOpts.MinSize(0, 100)),
		widget.ButtonOpts.ClickedHandler(func(args *widget.ButtonClickedEventArgs) {
			s.upgradeTowerHandler(s.chosenTower)
			checkBlock()

			s.Watcher.Append(s.Time, replay.UpgradeTower, replay.InfoUpgradeTower{
				Index: s.findTowerIndex(s.chosenTower),
			})
		}),
	)

	level := widget.NewText(
		widget.TextOpts.Text("Level", ttf48, color.White),
	)

	checkBlock = func() {
		openLevel := s.chosenTower.Upgrades[s.chosenTower.UpgradesBought-1].OpenLevel
		_, ok := s.PlayerState.LevelsComplete[openLevel]

		level.Label = fmt.Sprintf("Level %d", s.chosenTower.UpgradesBought+1)
		if s.chosenTower.UpgradesBought >= len(s.chosenTower.Upgrades) {
			btn.Text().Label = "SOLD OUT"
		} else if !ok && openLevel > 0 {
			btn.Text().Label = `Complete Level "` + strconv.Itoa(s.chosenTower.Upgrades[s.chosenTower.UpgradesBought].OpenLevel) + `" to unlock`
		} else {
			btn.Text().Label = fmt.Sprintf("UPGRADE ($%d)", s.chosenTower.Upgrades[s.chosenTower.UpgradesBought].Price)
		}

		if s.chosenTower.UpgradesBought >= len(s.chosenTower.Upgrades) ||
			s.PlayerMapState.Money < s.chosenTower.Upgrades[s.chosenTower.UpgradesBought].Price ||
			!ok && openLevel > 0 {
			btn.GetWidget().Disabled = true
		}
	}

	info := s.textUpgradeInfo()

	go func() {
		t := time.NewTicker(time.Second / time.Duration(ebiten.TPS()))
		for {
			select {
			case <-ctx.Done():
				return
			case <-t.C:
			}

			if s.chosenTower == nil {
				continue
			}

			if s.chosenTower.UpgradesBought >= len(s.chosenTower.Upgrades) {
				c := info.Children()
				insertValues(c[0].(*widget.Text), s.chosenTower.Damage, 0, "Damage")
				insertValues(c[1].(*widget.Text), int(s.chosenTower.Radius), 0, "Radius")
				insertValues(c[2].(*widget.Text), s.chosenTower.SpeedAttack, 0, "Speed")
				insertValues(c[3].(*widget.Text), int(s.chosenTower.ProjectileVrms), 0, "ProjSpeed")
				continue
			}

			c := info.Children()
			u := s.chosenTower.Upgrades[s.chosenTower.UpgradesBought]
			insertValues(c[0].(*widget.Text), s.chosenTower.Damage, u.DeltaDamage, "Damage")
			insertValues(c[1].(*widget.Text), int(s.chosenTower.Radius), int(u.DeltaRadius), "Radius")
			insertValues(c[2].(*widget.Text), s.chosenTower.SpeedAttack, u.DeltaSpeedAttack, "Speed")
			insertValues(c[3].(*widget.Text), int(s.chosenTower.ProjectileVrms), 0, "ProjSpeed")

			if s.PlayerMapState.Money < s.chosenTower.Upgrades[s.chosenTower.UpgradesBought].Price {
				btn.GetWidget().Disabled = true
			} else {
				btn.GetWidget().Disabled = false
			}
		}
	}()

	root.AddChild(level)
	root.AddChild(btn)
	root.AddChild(info)

	return root
}

func (s *GameState) textUpgradeInfo() *widget.Container {
	root := widget.NewContainer(
		widget.ContainerOpts.Layout(widget.NewRowLayout(
			widget.RowLayoutOpts.Direction(widget.DirectionVertical),
			widget.RowLayoutOpts.Spacing(5),
		)),
	)

	ttf40 := loaders.FontTrueType(40)
	defer ttf40.Close()
	textDamage := widget.NewText(
		widget.TextOpts.Text(fmt.Sprintf(`Damage: %d[color=00FF00]+0[/color]`, 0), ttf40, color.White),
		widget.TextOpts.ProcessBBCode(true),
	)
	textRadius := widget.NewText(
		widget.TextOpts.Text(fmt.Sprintf(`Radius: %.0f[color=FF0000]+0[/color]`, 0.), ttf40, color.White),
		widget.TextOpts.ProcessBBCode(true),
	)
	textSpeed := widget.NewText(
		widget.TextOpts.Text(fmt.Sprintf(`Speed: %.0f[color=FF0000]+0[/color]`, 0.), ttf40, color.White),
		widget.TextOpts.ProcessBBCode(true),
	)
	textProjSpeed := widget.NewText(
		widget.TextOpts.Text(fmt.Sprintf(`ProjSpeed: %.0f[color=AAAAAA]+0[/color]`, 0.), ttf40, color.White),
		widget.TextOpts.ProcessBBCode(true),
	)

	root.AddChild(textDamage)
	root.AddChild(textRadius)
	root.AddChild(textSpeed)
	root.AddChild(textProjSpeed)

	return root
}

func insertValues(c *widget.Text, v, deltav int, s string) {
	if deltav > 0 {
		c.Label = fmt.Sprintf("%s: %d[color=00FF00]+%d[/color]", s, v, deltav)
	} else if deltav < 0 {
		c.Label = fmt.Sprintf("%s: %d[color=FF0000]%d[/color]", s, v, deltav)
	} else {
		c.Label = fmt.Sprintf("%s: %d", s, v)
	}
}

func (s *GameState) tuningContainer(ctx context.Context, widgets general.Widgets) *widget.Container {
	root := widget.NewContainer(
		widget.ContainerOpts.Layout(widget.NewGridLayout(
			widget.GridLayoutOpts.Columns(1),
			widget.GridLayoutOpts.Stretch([]bool{true}, []bool{false, true}),
		)),
	)

	ttf54 := loaders.FontTrueType(54)
	defer ttf54.Close()
	btnTurn := widget.NewButton(
		widget.ButtonOpts.Image(&widget.ButtonImage{
			Idle: image2.NewNineSliceColor(colornames.Lawngreen),
		}),
		widget.ButtonOpts.Text("ON", ttf54, &widget.ButtonTextColor{
			Idle: color.White,
		}),
		widget.ButtonOpts.ClickedHandler(func(args *widget.ButtonClickedEventArgs) {
			btn := args.Button

			if s.chosenTower.State.IsTurnedOn {
				s.turnOffTowerHandler(s.chosenTower)
				btn.Text().Label = "OFF"
				btn.Image = &widget.ButtonImage{
					Idle: image2.NewNineSliceColor(colornames.Indianred),
				}

				s.Watcher.Append(s.Time, replay.TurnOff, replay.InfoTurnOffTower{
					Index: s.findTowerIndex(s.chosenTower),
				})

				return
			}

			s.turnOnTowerHandler(s.chosenTower)
			btn.Text().Label = "ON"
			btn.Image = &widget.ButtonImage{
				Idle: image2.NewNineSliceColor(colornames.Lawngreen),
			}

			s.Watcher.Append(s.Time, replay.TurnOn, replay.InfoTurnOnTower{
				Index: s.findTowerIndex(s.chosenTower),
			})
		}),
	)
	root.AddChild(btnTurn)
	root.AddChild(s.radio(ctx))

	return root
}

func (s *GameState) radio(ctx context.Context) *widget.Container {
	ttf40 := loaders.FontTrueType(40)

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
		widget.ButtonOpts.Text("First", ttf40, &widget.ButtonTextColor{
			Idle: color.White,
		}),
		widget.ButtonOpts.Image(&widget.ButtonImage{
			Idle:    image2.NewNineSliceColor(color.RGBA{0x7f, 0x27, 0xd7, 0xff}),
			Hover:   image2.NewNineSliceColor(color.RGBA{0x9a, 0x3b, 0xea, 0xff}),
			Pressed: image2.NewNineSliceColor(color.RGBA{0x6a, 0x16, 0xc2, 0xff}),
		}),
		widget.ButtonOpts.ClickedHandler(func(args *widget.ButtonClickedEventArgs) {
			s.tuneFirstTowerHandler(s.chosenTower)

			s.Watcher.Append(s.Time, replay.TuneFirst, replay.InfoTuneFirst{
				Index: s.findTowerIndex(s.chosenTower),
			})
		}),
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
		widget.ButtonOpts.Text("Strong", ttf40, &widget.ButtonTextColor{
			Idle: color.White,
		}),
		widget.ButtonOpts.Image(&widget.ButtonImage{
			Idle:    image2.NewNineSliceColor(color.RGBA{0x7f, 0x27, 0xd7, 0xff}),
			Hover:   image2.NewNineSliceColor(color.RGBA{0x9a, 0x3b, 0xea, 0xff}),
			Pressed: image2.NewNineSliceColor(color.RGBA{0x6a, 0x16, 0xc2, 0xff}),
		}),
		widget.ButtonOpts.ClickedHandler(func(args *widget.ButtonClickedEventArgs) {
			s.tuneStrongTowerHandler(s.chosenTower)

			s.Watcher.Append(s.Time, replay.TuneStrong, replay.InfoTuneStrong{
				Index: s.findTowerIndex(s.chosenTower),
			})
		}),
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
		widget.ButtonOpts.Text("Weak", ttf40, &widget.ButtonTextColor{
			Idle: color.White,
		}),
		widget.ButtonOpts.Image(&widget.ButtonImage{
			Idle:    image2.NewNineSliceColor(color.RGBA{0x7f, 0x27, 0xd7, 0xff}),
			Hover:   image2.NewNineSliceColor(color.RGBA{0x9a, 0x3b, 0xea, 0xff}),
			Pressed: image2.NewNineSliceColor(color.RGBA{0x6a, 0x16, 0xc2, 0xff}),
		}),
		widget.ButtonOpts.ClickedHandler(func(args *widget.ButtonClickedEventArgs) {
			s.tuneWeakTowerHandler(s.chosenTower)

			s.Watcher.Append(s.Time, replay.TuneWeak, replay.InfoTuneWeak{
				Index: s.findTowerIndex(s.chosenTower),
			})
		}),
	)

	root.AddChild(first)
	root.AddChild(strong)
	root.AddChild(weak)

	r := widget.NewRadioGroup(
		widget.RadioGroupOpts.Elements(first, strong, weak),
	)
	go func() {
		t := time.NewTicker(time.Second / time.Duration(ebiten.TPS()))
		for {
			select {
			case <-ctx.Done():
				return
			case <-t.C:
			}
			if s.chosenTower == nil {
				continue
			}

			switch s.chosenTower.State.AimType {
			case ingame.First:
				r.SetActive(first)
			case ingame.Strongest:
				r.SetActive(strong)
			case ingame.Weakest:
				r.SetActive(weak)
			}
		}
	}()

	return root
}

func (s *GameState) sellContainer(widgets general.Widgets) *widget.Container {
	root := widget.NewContainer(
		widget.ContainerOpts.Layout(widget.NewGridLayout(
			widget.GridLayoutOpts.Columns(1),
			widget.GridLayoutOpts.Stretch([]bool{true}, []bool{true}),
			widget.GridLayoutOpts.Padding(widget.Insets{Top: 40}),
		)),
	)

	ttf64 := loaders.FontTrueType(64)
	btnSell := widget.NewButton(
		widget.ButtonOpts.Image(&widget.ButtonImage{
			Idle: image2.NewNineSliceColor(color.RGBA{0xff, 0x66, 0x66, 0xff}),
		}),
		widget.ButtonOpts.Text("SELL", ttf64, &widget.ButtonTextColor{
			Idle: color.White,
		}),
		widget.ButtonOpts.ClickedHandler(func(args *widget.ButtonClickedEventArgs) {
			s.Watcher.Append(s.Time, replay.SellTower, replay.InfoSellTower{
				Index: s.findTowerIndex(s.chosenTower),
			})
			s.sellTowerHandler(s.chosenTower)

			s.showTowerMenu()
		}),
		widget.ButtonOpts.TextPadding(widget.Insets{Top: 10, Bottom: 10}),
	)

	root.AddChild(btnSell)

	return root
}
