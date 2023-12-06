package io

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"os"

	"github.com/gopher-co/td-game/models/ingame"
)

func LoadStats() (*ingame.PlayerState, error) {
	f, err := os.Open("stats.json")
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			return &ingame.PlayerState{LevelsComplete: map[int]struct{}{}}, nil
		}

		return nil, fmt.Errorf("stats file can't be open: %w", err)
	}

	defer func() {
		_ = f.Close()
	}()

	stats := new(ingame.PlayerState)
	if err = json.NewDecoder(f).Decode(stats); err != nil {
		return nil, fmt.Errorf("stats json not parsed: %w", err)
	}

	return stats, nil
}

func SaveStats(stats *ingame.PlayerState) error {
	f, err := os.OpenFile("stats.json", os.O_WRONLY|os.O_SYNC|os.O_CREATE|os.O_TRUNC, 0o666)
	if err != nil && !errors.Is(err, fs.ErrExist) {
		return fmt.Errorf("stats file can't be open: %w", err)
	}

	defer func() { _ = f.Close() }()

	buf := bufio.NewWriter(f)
	if err := json.NewEncoder(buf).Encode(*stats); err != nil {
		return fmt.Errorf("unsuccessful save: %w", err)
	}

	return buf.Flush()
}
