package replay

import (
	"bytes"
	"encoding/json"
	"io"
	"strconv"

	"github.com/gopher-co/td-game/models/general"
	"github.com/gopher-co/td-game/models/ingame"
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
	Stop
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
	case SellTower:
		info := InfoSellTower{}
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
	case Stop:
		info := InfoStop{}
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
	Name string `json:"name"`
	X    int    `json:"x"`
	Y    int    `json:"y"`
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

type InfoStop struct {
	Null any `json:"null"`
}

type Watcher struct {
	Name               string                `json:"name"`
	Time               string                `json:"time"`
	InitPlayerMapState ingame.PlayerMapState `json:"init_player_map_state"`
	Actions            []Action              `json:"actions"`
}

func (wt *Watcher) Append(f general.Frames, at ActionType, info any) {
	wt.Actions = append(wt.Actions, Action{
		F:    f,
		Type: at,
		Info: info,
	})
}

func (wt *Watcher) Write(w io.Writer) error {
	enc := json.NewEncoder(w)
	return enc.Encode(*wt)
}

func (wt *Watcher) Read(r io.Reader) error {
	return json.NewDecoder(r).Decode(wt)
}
