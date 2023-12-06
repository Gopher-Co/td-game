package ingame

import (
	"slices"

	"github.com/gopher-co/td-game/models/config"
	"github.com/gopher-co/td-game/models/general"
)

// GameRule is a set of waves.
type GameRule []*Wave

// NewGameRule returns a new GameRule.
func NewGameRule(config []config.Wave) GameRule {
	grs := make(GameRule, len(config))
	for i := range config {
		grs[i] = NewWave(&config[i])
	}

	return grs
}

// Wave is a set of the enemies.
type Wave struct {
	// Swarms is a set of the enemy swarms.
	Swarms []*EnemySwarm

	// Time is a current time of the wave.
	Time general.Frames
}

// NewWave returns a new Wave.
func NewWave(config *config.Wave) *Wave {
	swarms := make([]*EnemySwarm, len(config.Swarms))
	for i := 0; i < len(swarms); i++ {
		swarms[i] = NewEnemySwarm(&config.Swarms[i])
	}

	return &Wave{Swarms: swarms}
}

// CallEnemies returns a slice of ids of enemies that are
// supposed to appear on the map next frame.
func (w *Wave) CallEnemies() []string {
	es := make([]string, 0, len(w.Swarms))
	for _, v := range w.Swarms {
		if v.Ended() {
			continue
		}
		if e := v.Update(w.Time); e != "" {
			es = append(es, e)
		}
	}

	w.Time++
	return slices.Clip(es)
}

// Ended returns true if all the swarms are ended.
func (w *Wave) Ended() bool {
	for _, v := range w.Swarms {
		if !v.Ended() {
			return false
		}
	}
	return true
}

// EnemySwarm contains the rules for calling the next enemy.
// Enemies are called in the same interval limited times.
// EnemySwarm can call only one type of the enemy.
type EnemySwarm struct {
	// EnemyName is a name of the enemy.
	EnemyName string

	// Timeout is the time when the first enemy can be called.
	Timeout general.Frames

	// Interval is time between calls.
	Interval general.Frames

	// CurrTime is current time relatively the swarm's start.
	CurrTime general.Frames

	// MaxCalls is a maximal amount of enemies that can be called.
	MaxCalls int

	// CurCalls is the current amount of enemies called.
	CurCalls int
}

// NewEnemySwarm returns a new EnemySwarm.
func NewEnemySwarm(config *config.EnemySwarm) *EnemySwarm {
	return &EnemySwarm{
		EnemyName: config.EnemyName,
		Timeout:   config.Timeout,
		Interval:  config.Interval,
		MaxCalls:  config.MaxCalls,
		CurCalls:  0,
	}
}

// Ended returns true if maximum calls amount exceeded.
func (s *EnemySwarm) Ended() bool {
	return s.CurCalls == s.MaxCalls
}

// Update increases time of EnemySwarm and returns a new enemy id
// if it's time for it.
func (s *EnemySwarm) Update(t general.Frames) string {
	if t-s.Timeout >= 0 && (t-s.Timeout)%s.Interval == 0 {
		s.CurCalls++
		return s.EnemyName
	}
	return ""
}
