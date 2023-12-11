package io

import (
	"fmt"

	"github.com/gopher-co/td-game/replay"
)

// LoadReplays loads replays from the Replays directory.
func LoadReplays() ([]*replay.Watcher, error) {
	rs, err := ReadConfigs[replay.Watcher]("./Replays", ".json")
	if err != nil {
		return nil, fmt.Errorf("couldn't read replays: %w", err)
	}

	prs := make([]*replay.Watcher, len(rs))
	for i := 0; i < len(rs); i++ {
		prs[i] = &rs[i]
	}

	return prs, nil
}
