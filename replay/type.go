package replay

import (
	"bytes"
	"encoding/json"
	"strconv"

	"github.com/gopher-co/td-game/models/general"
)

type ActionType int

const (
	PutTower = ActionType(iota)
	SellTower
	UpgradeTower
	TurnOff
	TurnOn
	TuneFirst
	TuneStrong
	TuneWeak
)

type Action struct {
	F    general.Frames `json:"f"`
	Type ActionType     `json:"type"`
	Info any            `json:"info"`
}

func (a *Action) UnmarshalJSON(b []byte) error {
	indexFStart := 5
	indexFEnd := bytes.Index(b, []byte{','})
	f, err := strconv.Atoi(string(b[indexFStart:indexFEnd]))
	if err != nil {
		return err
	}

	a.F = f

	indexTypeStart := bytes.Index(b, []byte(`"type":`)) + 7
	indexTypeEnd := bytes.Index(b[indexTypeStart:], []byte{','}) + indexTypeStart
	actionType, err := strconv.Atoi(string(b[indexTypeStart:indexTypeEnd]))
	if err != nil {
		return err
	}

	a.Type = ActionType(actionType)

	i := bytes.Index(b, []byte(`"info":`)) + 7
	infob := b[i : len(b)-1]

	switch ActionType(actionType) {
	case PutTower:
		info := InfoPutTower{}
		if err := json.Unmarshal(infob, &info); err != nil {
			return err
		}
		a.Info = info
	case UpgradeTower:
		info := InfoUpgradeTower{}
		if err := json.Unmarshal(infob, &info); err != nil {
			return err
		}
		a.Info = info
	case TurnOff:
		info := InfoTurnOffTower{}
		if err := json.Unmarshal(infob, &info); err != nil {
			return err
		}
		a.Info = info
	case TurnOn:
		info := InfoTurnOnTower{}
		if err := json.Unmarshal(infob, &info); err != nil {
			return err
		}
		a.Info = info
	case TuneFirst:
		info := InfoTuneFirst{}
		if err := json.Unmarshal(infob, &info); err != nil {
			return err
		}
		a.Info = info
	case TuneStrong:
		info := InfoTuneStrong{}
		if err := json.Unmarshal(infob, &info); err != nil {
			return err
		}
		a.Info = info
	case TuneWeak:
		info := InfoTuneWeak{}
		if err := json.Unmarshal(infob, &info); err != nil {
			return err
		}
		a.Info = info
	default:
		return err
	}

	return nil
}

type InfoPutTower struct {
	Name string        `json:"name"`
	X    general.Coord `json:"x"`
	Y    general.Coord `json:"y"`
}

type InfoSellTower struct {
	Index int `json:"index"`
}

type InfoUpgradeTower struct {
	Index int `json:"index"`
}

type InfoTurnOffTower struct {
	Index int `json:"index"`
}

type InfoTurnOnTower struct {
	Index int `json:"index"`
}

type InfoTuneFirst struct {
	Index int `json:"index"`
}

type InfoTuneStrong struct {
	Index int `json:"index"`
}

type InfoTuneWeak struct {
	Index int `json:"index"`
}

type Replay struct {
	Actions []Action `json:"actions"`
}
