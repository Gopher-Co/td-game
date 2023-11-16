package models

// Upgrade is an entity stores useful effects for towers.
type Upgrade struct {
	Price            int
	DeltaDamage      int
	DeltaSpeedAttack Frames
	DeltaRadius      Coord
	OpenLevel        int
}

func NewUpgrade(config *UpgradeConfig) *Upgrade {
	return &Upgrade{
		Price:            config.Price,
		DeltaDamage:      config.DeltaDamage,
		DeltaSpeedAttack: config.DeltaSpeedAttack,
		DeltaRadius:      config.DeltaRadius,
		OpenLevel:        config.OpenLevel,
	}
}
