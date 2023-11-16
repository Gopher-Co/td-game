package models

import "slices"

type GameRule []*Wave

func NewGameRule(config []WaveConfig) GameRule {
	grs := make(GameRule, len(config))
	for i := range config {
		grs[i] = NewWave(&config[i])
	}

	return grs
}

// Wave is a set of the enemies.
type Wave struct {
	Swarms []*EnemySwarm
}

func NewWave(config *WaveConfig) *Wave {
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
		if e := v.Update(); e != "" {
			es = append(es, e)
		}
	}

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
	Timeout Frames

	// Interval is time between calls.
	Interval Frames

	// CurrTime is current time relatively the swarm's start.
	CurrTime Frames

	// MaxCalls is a maximal amount of enemies that can be called.
	MaxCalls int

	// CurCalls is the current amount of enemies called.
	CurCalls int
}

func NewEnemySwarm(config *EnemySwarmConfig) *EnemySwarm {
	return &EnemySwarm{
		EnemyName: config.EnemyName,
		Timeout:   config.Timeout,
		Interval:  config.Interval,
		CurrTime:  config.CurrTime,
		MaxCalls:  config.MaxCalls,
		CurCalls:  config.CurCalls,
	}
}

// Ended returns true if maximum calls amount exceeded.
func (s *EnemySwarm) Ended() bool {
	return s.CurCalls == s.MaxCalls
}

// Update increases time of EnemySwarm and returns a new enemy id
// if it's time for it.
func (s *EnemySwarm) Update() string {
	defer func() { s.CurrTime++ }()

	if s.CurrTime-s.Timeout >= 0 && (s.CurrTime-s.Timeout)%s.Interval == 0 {
		return s.EnemyName
	}
	return ""
}
