package ingame

// PlayerState is a struct that represents a state of the player.
type PlayerState struct {
	// LevelsComplete is a set of levels that player has completed.
	LevelsComplete map[int]struct{}
}
