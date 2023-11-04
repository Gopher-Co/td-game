package models

import "image"

type Upgrade struct {
	UpgradeConfig
	Image     image.Image
	OpenLevel int
}

type UpgradeConfig struct {
	DeltaDamage      int
	DeltaSpeedAttack Frames
	DeltaRadius      Coord
}
