package replay

import (
	"bytes"
	"encoding/json"
	"io"
	"strconv"

	"github.com/gopher-co/td-game/models/general"
	"github.com/gopher-co/td-game/models/ingame"
)

// ActionType is a type of action.
type ActionType int

const (
	// PutTower is a type of action that represents putting a tower.
	PutTower = ActionType(iota)

	// SellTower is a type of action that represents selling a tower.
	SellTower

	// UpgradeTower is a type of action that represents upgrading a tower.
	UpgradeTower

	// TurnOff is a type of action that represents turning off a tower.
	TurnOff

	// TurnOn is a type of action that represents turning on a tower.
	TurnOn

	// TuneFirst is a type of action that represents tuning first.
	TuneFirst

	// TuneStrong is a type of action that represents tuning strong.
	TuneStrong

	// TuneWeak is a type of action that represents tuning weak.
	TuneWeak

	// Stop is a type of action that represents stopping the game.
	Stop
)

// Action is an entity that represents an action.
type Action struct {
	// F is a frame when the action is performed.
	F general.Frames `json:"f"`

	// Type is a type of the action.
	Type ActionType `json:"type"`

	// Info is an info of the action.
	Info any `json:"info"`
}

// UnmarshalJSON unmarshals the action.
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

// InfoPutTower is an info of the action that represents putting a tower.
type InfoPutTower struct {
	// Name is a name of the tower.
	Name string `json:"name"`

	// X is a x coordinate of the tower.
	X int `json:"x"`

	// Y is a y coordinate of the tower.
	Y int `json:"y"`
}

// InfoSellTower is an info of the action that represents selling a tower.
type InfoSellTower struct {
	// Index is an index of the tower.
	Index int `json:"index"`
}

// InfoUpgradeTower is an info of the action that represents upgrading a tower.
type InfoUpgradeTower struct {
	// Index is an index of the tower.
	Index int `json:"index"`
}

// InfoTurnOffTower is an info of the action that represents turning off a tower.
type InfoTurnOffTower struct {
	// Index is an index of the tower.
	Index int `json:"index"`
}

// InfoTurnOnTower is an info of the action that represents turning on a tower.
type InfoTurnOnTower struct {
	// Index is an index of the tower.
	Index int `json:"index"`
}

// InfoTuneFirst is an info of the action that represents tuning first.
type InfoTuneFirst struct {
	// Index is an index of the tower.
	Index int `json:"index"`
}

// InfoTuneStrong is an info of the action that represents tuning strong.
type InfoTuneStrong struct {
	// Index is an index of the tower.
	Index int `json:"index"`
}

// InfoTuneWeak is an info of the action that represents tuning weak.
type InfoTuneWeak struct {
	// Index is an index of the tower.
	Index int `json:"index"`
}

// InfoStop is an info of the action that represents stopping the game.
type InfoStop struct {
	// Null is a null.
	Null any `json:"null"`
}

// Watcher is an entity that represents a watcher.
type Watcher struct {
	// Name is a name of the watcher.
	Name string `json:"name"`

	// Time is a time of the watcher.
	Time string `json:"time"`

	// InitPlayerMapState is an initial player map state.
	InitPlayerMapState ingame.PlayerMapState `json:"init_player_map_state"`

	// Actions is a list of actions.
	Actions []Action `json:"actions"`
}

// Append appends an action to the watcher.
func (wt *Watcher) Append(f general.Frames, at ActionType, info any) {
	wt.Actions = append(wt.Actions, Action{
		F:    f,
		Type: at,
		Info: info,
	})
}

// Write writes the watcher.
func (wt *Watcher) Write(w io.Writer) error {
	enc := json.NewEncoder(w)
	return enc.Encode(*wt)
}

// Read reads the watcher.
func (wt *Watcher) Read(r io.Reader) error {
	return json.NewDecoder(r).Decode(wt)
}
