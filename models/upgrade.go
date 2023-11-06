package models

// Upgrade is an entity stores useful effects for towers.
type Upgrade struct {
	UpgradeConfig

	// OpenLevel is a number of the level completing
	// of which opens the upgrade.
	OpenLevel int
}

// UpgradeConfig stores deltas for tower's parameters.
// Applied upgrade adds to corresponding tower parameters these deltas.
type UpgradeConfig struct {
	DeltaDamage      int
	DeltaSpeedAttack Frames
	DeltaRadius      Coord
}
