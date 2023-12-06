package replay_test

import (
	"bytes"
	"reflect"
	"testing"

	"github.com/gopher-co/td-game/replay"
)

func TestActions(t *testing.T) {
	rep := replay.Watcher{Actions: []replay.Action{
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

	buf := new(bytes.Buffer)
	if err := rep.Write(buf); err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	actualRep := replay.Watcher{}
	if err := actualRep.Read(buf); err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if !reflect.DeepEqual(rep, actualRep) {
		t.Errorf("got %#+v, expected %#+v", actualRep, rep)
	}
}
