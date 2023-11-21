package models

// Upgrade is an entity stores useful effects for towers.
type Upgrade struct {

	// Price is a price of the upgrade.
	Price int

	// DeltaDamage is a delta of the damage.
	DeltaDamage int

	// DeltaSpeedAttack is a delta of the speed of the attack.
	DeltaSpeedAttack Frames

	// DeltaRadius is a delta of the radius.
	DeltaRadius Coord

	// OpenLevel is a level when the upgrade is opened.
	OpenLevel int
}

// NewUpgrade returns a new upgrade.
func NewUpgrade(config *UpgradeConfig) *Upgrade {
	return &Upgrade{
		Price:            config.Price,
		DeltaDamage:      config.DeltaDamage,
		DeltaSpeedAttack: config.DeltaSpeedAttack,
		DeltaRadius:      config.DeltaRadius,
		OpenLevel:        config.OpenLevel,
	}
}
