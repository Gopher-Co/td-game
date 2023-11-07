package models

import "image"

type Aim int

const (
	Direct Aim = iota
	Splash
)

type Tower struct {
	Name           string
	Damage         int
	Type           TypeAttack
	Price          int
	Image          image.Image
	Radius         float64
	State          *TowerState
	SpeedAttack    Frames
	Upgrades       []Upgrade
	UpgradesBought int
	OpenLevel      int
}

type TowerState struct {
	AimType    Aim
	IsTurnedOn bool
	CoolDown   Frames
	Point      Point
	Aim        *Enemy
}

type Upgrade struct {
	Name        string
	Description string
	Config      UpgradeConfig
	Image       image.Image
	OpenLevel   int
}

type UpgradeConfig struct {
	DeltaDamage      int
	DeltaSpeedAttack Frames
	DeltaRadius      float64
}
