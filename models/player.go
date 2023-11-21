package models

// PlayerMapState is a state of the player during the game.
type PlayerMapState struct {
	// Health is a player's health.
	Health int

	// Money is a player's money.
	Money int
}

// Dead returns true if player's health is equal to zero.
func (s *PlayerMapState) Dead() bool {
	return s.Health == 0
}
