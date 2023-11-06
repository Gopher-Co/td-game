package models

import "slices"

type Wave struct {
	Swarms      map[EnemySwarm]struct{}
	CurrentTime Frames
}

func (w *Wave) CallEnemies() []*Enemy {
	es := make([]*Enemy, 0, len(w.Swarms))
	for k := range w.Swarms {
		if e := k.Call(); e != nil {
			es = append(es, e)
		}
	}

	return slices.Clip(es)
}

func (w *Wave) Ended() bool {
	for k := range w.Swarms {
		if !k.Ended() {
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

	// MaxCalls is a maximal amount of enemies that can be called.
	MaxCalls int

	// CurCalls is the current amount of enemies called.
	CurCalls int
}

func (s *EnemySwarm) Ended() bool {
	return s.CurCalls == s.MaxCalls
}

func (s *EnemySwarm) Call() *Enemy {
	return nil
}
