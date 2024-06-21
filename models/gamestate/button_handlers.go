package gamestate

import (
	image2 "github.com/ebitenui/ebitenui/image"
	"github.com/ebitenui/ebitenui/widget"
	"github.com/hajimehoshi/ebiten/v2"
	"golang.org/x/image/colornames"

	"github.com/gopher-co/td-game/models/config"
	"github.com/gopher-co/td-game/replay"
)

// handleSpeed handles the speed button click.
func (s *GameState) handleSpeed(args *widget.ButtonClickedEventArgs) {
	if s.speedUp {
		ebiten.SetTPS(60)
		s.speedUp = false
		args.Button.Image = &widget.ButtonImage{
			Idle: image2.NewNineSliceColor(colornames.Cornflowerblue),
		}
		return
	}

	ebiten.SetTPS(180)
	s.speedUp = true
	args.Button.Image = &widget.ButtonImage{
		Idle: image2.NewNineSliceColor(colornames.Greenyellow),
	}
}

// handleStart handles the start button click.
func (s *GameState) handleStart(args *widget.ButtonClickedEventArgs) {
	b := args.Button
	if !b.GetWidget().Disabled {
		s.State = Running
		s.CurrentWave++
		b.GetWidget().Disabled = true
	}
}

// handleMenu handles the menu button click.
func (s *GameState) handleMenu(_ *widget.ButtonClickedEventArgs) {
	s.setStateAfterEnd()
	s.Ended = true
}

// handleTowerTake handles the tower take button click.
func (s *GameState) handleTowerTake(v *config.Tower) func(eventArgs *widget.ButtonClickedEventArgs) {
	return func(args *widget.ButtonClickedEventArgs) {
		if s.PlayerMapState.Money >= v.Price {
			s.tookTower = v
		}
	}
}

// handleUpgrade handles the upgrade button click.
func (s *GameState) handleUpgrade(_ *widget.ButtonClickedEventArgs) {
	s.upgradeTowerHandler(s.chosenTower)

	s.Watcher.Append(s.Time, replay.UpgradeTower, replay.InfoUpgradeTower{
		Index: s.findTowerIndex(s.chosenTower),
	})
}

func (s *GameState) handleTurning(args *widget.ButtonClickedEventArgs) {
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
}

func (s *GameState) handleTuneFirst(_ *widget.ButtonClickedEventArgs) {
	s.tuneFirstTowerHandler(s.chosenTower)

	s.Watcher.Append(s.Time, replay.TuneFirst, replay.InfoTuneFirst{
		Index: s.findTowerIndex(s.chosenTower),
	})
}

func (s *GameState) handleTuneStrong(_ *widget.ButtonClickedEventArgs) {
	s.tuneStrongTowerHandler(s.chosenTower)

	s.Watcher.Append(s.Time, replay.TuneStrong, replay.InfoTuneStrong{
		Index: s.findTowerIndex(s.chosenTower),
	})
}

func (s *GameState) handleTuneWeak(_ *widget.ButtonClickedEventArgs) {
	s.tuneWeakTowerHandler(s.chosenTower)

	s.Watcher.Append(s.Time, replay.TuneWeak, replay.InfoTuneWeak{
		Index: s.findTowerIndex(s.chosenTower),
	})
}

func (s *GameState) handleSell(_ *widget.ButtonClickedEventArgs) {
	s.sellTowerHandler(s.chosenTower)

	s.Watcher.Append(s.Time, replay.SellTower, replay.InfoSellTower{
		Index: s.findTowerIndex(s.chosenTower),
	})

	s.showTowerMenu()
}
