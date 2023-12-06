package ingame

import (
	"github.com/gopher-co/td-game/models/config"
	"github.com/gopher-co/td-game/models/general"
)

// Upgrade is an entity stores useful effects for towers.
type Upgrade struct {
	// Price is a price of the upgrade.
	Price int

	// DeltaDamage is a delta of the damage.
	DeltaDamage int

	// DeltaSpeedAttack is a delta of the speed of the attack.
	DeltaSpeedAttack general.Frames

	// DeltaRadius is a delta of the radius.
	DeltaRadius general.Coord

	// OpenLevel is a level when the upgrade is opened.
	OpenLevel int
}

// NewUpgrade returns a new upgrade.
func NewUpgrade(config *config.Upgrade) *Upgrade {
	return &Upgrade{
		Price:            config.Price,
		DeltaDamage:      config.DeltaDamage,
		DeltaSpeedAttack: config.DeltaSpeedAttack,
		DeltaRadius:      config.DeltaRadius,
		OpenLevel:        config.OpenLevel,
	}
}
