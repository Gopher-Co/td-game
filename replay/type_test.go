package replay_test

import (
	"encoding/json"
	"reflect"
	"testing"

	"github.com/gopher-co/td-game/replay"
)

func TestActions(t *testing.T) {
	rep := replay.Replay{Actions: []replay.Action{
		{
			F:    1,
			Type: replay.PutTower,
			Info: replay.InfoPutTower{
				Name: "123",
				X:    4,
				Y:    6,
			},
		},
		{
			F:    12,
			Type: replay.SellTower,
			Info: replay.InfoSellTower{
				Index: 0,
			},
		},
		{
			F:    27,
			Type: replay.TuneWeak,
			Info: replay.InfoTuneWeak{Index: 0},
		},
	}}

	var err error
	var b []byte
	if b, err = json.Marshal(rep); err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	actualRep := replay.Replay{[]replay.Action{{Info: replay.InfoPutTower{}}, {Info: replay.InfoSellTower{}}}}
	if err := json.Unmarshal(b, &actualRep); err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if !reflect.DeepEqual(rep, actualRep) {
		t.Errorf("got %#+v, expected %#+v", actualRep, rep)
	}
}
