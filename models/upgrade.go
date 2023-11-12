package models

// Upgrade is an entity stores useful effects for towers.
type Upgrade struct {
	Price            int
	DeltaDamage      int
	DeltaSpeedAttack Frames
	DeltaRadius      Coord

	// OpenLevel is a number of the level completing
	// of which opens the upgrade.
	OpenLevel int
}
