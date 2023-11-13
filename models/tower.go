package models

import "image"

// Aim is a type that represents the aim of a tower.
type Aim int

const (
	// DirectAim is a value of Aim that represents a direct aim.
	DirectAim = Aim(iota)

	// Splash is a value of Aim that represents a splash aim.
	Splash
)

// Tower is a struct that represents a tower.
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

// TowerState is a struct that represents the state of a tower.
type TowerState struct {
	AimType    Aim
	IsTurnedOn bool
	CoolDown   Frames
	Point      Point
	Aim        *Enemy
}
