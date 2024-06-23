package coopstate

import (
	"context"

	image2 "github.com/ebitenui/ebitenui/image"
	"github.com/ebitenui/ebitenui/widget"
	"golang.org/x/image/colornames"

	"github.com/gopher-co/td-game/models/config"
	"github.com/gopher-co/td-game/replay"
)

// handleSpeed handles the speed button click.
func (s *GameState) handleSpeed(args *widget.ButtonClickedEventArgs) {
	if s.speedUp {
		_, _ = s.cli.SlowGameDown(s.ctx, &SlowGameDownRequest{})
		args.Button.Image = &widget.ButtonImage{
			Idle: image2.NewNineSliceColor(colornames.Cornflowerblue),
		}
		return
	}

	_, _ = s.cli.SpeedGameUp(s.ctx, &SpeedGameUpRequest{})
	args.Button.Image = &widget.ButtonImage{
		Idle: image2.NewNineSliceColor(colornames.Greenyellow),
	}
}

// handleStart handles the start button click.
func (s *GameState) handleStart(args *widget.ButtonClickedEventArgs) {
	b := args.Button
	if !b.GetWidget().Disabled {
		_, _ = s.cli.StartNewWave(context.Background(), &StartNewWaveRequest{})
		b.GetWidget().Disabled = true
	}
}

// handleMenu handles the menu button click.
func (s *GameState) handleMenu(_ *widget.ButtonClickedEventArgs) {
	_ = s.stream.CloseSend()
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
	_, _ = s.cli.UpgradeTower(s.ctx, &UpgradeTowerRequest{
		Tower: &TowerId{Id: int64(s.chosenTower.Index)},
	})

	s.Watcher.Append(s.Time, replay.UpgradeTower, replay.InfoUpgradeTower{
		Index: s.chosenTower.Index,
	})
}

// handleTurning handles the turning button click.
func (s *GameState) handleTurning(args *widget.ButtonClickedEventArgs) {
	btn := args.Button

	if s.chosenTower.State.IsTurnedOn {
		_, _ = s.cli.TurnTowerOff(s.ctx, &TurnTowerOffRequest{Tower: &TowerId{Id: int64(s.chosenTower.Index)}})
		btn.Text().Label = "OFF"
		btn.Image = &widget.ButtonImage{
			Idle: image2.NewNineSliceColor(colornames.Indianred),
		}

		s.Watcher.Append(s.Time, replay.TurnOff, replay.InfoTurnOffTower{
			Index: s.chosenTower.Index,
		})

		return
	}

	_, _ = s.cli.TurnTowerOn(s.ctx, &TurnTowerOnRequest{Tower: &TowerId{Id: int64(s.chosenTower.Index)}})
	btn.Text().Label = "ON"
	btn.Image = &widget.ButtonImage{
		Idle: image2.NewNineSliceColor(colornames.Lawngreen),
	}

	s.Watcher.Append(s.Time, replay.TurnOn, replay.InfoTurnOnTower{
		Index: s.chosenTower.Index,
	})
}

// handleTuneFirst handles the tune first button click.
func (s *GameState) handleTuneFirst(_ *widget.ButtonClickedEventArgs) {
	_, _ = s.cli.ChangeTowerAimType(s.ctx, &ChangeTowerAimTypeRequest{
		Tower:      &TowerId{Id: int64(s.chosenTower.Index)},
		NewAimType: int32(TuneTowerRequest_AIM_TOWER_AT_FIRST),
	})
	s.Watcher.Append(s.Time, replay.TuneFirst, replay.InfoTuneFirst{
		Index: s.chosenTower.Index,
	})
}

// handleTuneStrong handles the tune strong button click.
func (s *GameState) handleTuneStrong(_ *widget.ButtonClickedEventArgs) {
	_, _ = s.cli.ChangeTowerAimType(s.ctx, &ChangeTowerAimTypeRequest{
		Tower:      &TowerId{Id: int64(s.chosenTower.Index)},
		NewAimType: int32(TuneTowerRequest_AIM_TOWER_AT_STRONG),
	})

	s.Watcher.Append(s.Time, replay.TuneStrong, replay.InfoTuneStrong{
		Index: s.chosenTower.Index,
	})
}

// handleTuneWeak handles the tune weak button click.
func (s *GameState) handleTuneWeak(_ *widget.ButtonClickedEventArgs) {
	_, _ = s.cli.ChangeTowerAimType(s.ctx, &ChangeTowerAimTypeRequest{
		Tower:      &TowerId{Id: int64(s.chosenTower.Index)},
		NewAimType: int32(TuneTowerRequest_AIM_TOWER_AT_LAST),
	})

	s.Watcher.Append(s.Time, replay.TuneWeak, replay.InfoTuneWeak{
		Index: s.chosenTower.Index,
	})
}

// handleSell handles the sell button click.
func (s *GameState) handleSell(_ *widget.ButtonClickedEventArgs) {
	_, _ = s.cli.SellTower(s.ctx, &SellTowerRequest{
		Tower: &TowerId{Id: int64(s.chosenTower.Index)},
	})

	s.Watcher.Append(s.Time, replay.SellTower, replay.InfoSellTower{
		Index: s.chosenTower.Index,
	})

	s.showTowerMenu()
}
