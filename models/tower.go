package models

import (
	"image"
	"log"
	"math"
)

// Aim is a type that represents the enemy that tower attacks.
type Aim int

const (
	First = Aim(iota)
	Weakest
	Strongest
)

// Tower is a struct that represents a tower.
type Tower struct {
	Name           string
	Damage         int
	Type           TypeAttack
	Price          int
	Image          image.Image
	Radius         Coord
	State          TowerState
	SpeedAttack    Frames
	Upgrades       []Upgrade
	UpgradesBought int
}

func NewTower(config *TowerConfig, pos Point, path Path) *Tower {
	for i := 0; i < len(path)-1; i++ {
		if checkCollision(Point{pos.X, pos.Y}, path[i], path[i+1]) {
			return nil
		}
	}

	return &Tower{
		Name:   config.Name,
		Damage: config.InitDamage,
		Type:   config.Type,
		Price:  config.Price,
		Image:  config.Image(),
		Radius: config.InitRadius,
		State: TowerState{
			AimType:    First,
			IsTurnedOn: true,
			CoolDown:   config.InitSpeedAttack,
			Point:      pos,
			Aim:        nil,
		},
		SpeedAttack:    config.InitSpeedAttack,
		Upgrades:       config.Upgrades,
		UpgradesBought: 0,
	}
}

func checkCollision(p, p1, p2 Point) bool {
	tx1, ty1, tx2, ty2 := p1.X, p1.Y, p2.X, p2.Y
	x1, x2 := float64(min(tx1, tx2)), float64(max(tx1, tx2))
	y1, y2 := float64(min(ty1, ty2)), float64(max(ty1, ty2))

	// 		***
	// 		(x1,y1)**        z
	// 	 y	*b       *****
	// 		*	         a*****
	// 		()****************(x2,y2)
	//					x
	z := math.Hypot(x2-x1, y2-y1)

	sina := (y2 - y1) / z
	cosa := 1 - math.Pow(sina, 2)

	dx := PathWidth * cosa
	dy := PathWidth * sina

	x1 -= dx
	x2 += dx
	y1 -= dy
	y2 += dy

	A := Point{Coord(x1 - dy), Coord(y1 + dx)}
	B := Point{Coord(x1 + dy), Coord(y1 - dx)}
	C := Point{Coord(x2 - dy), Coord(y2 + dx)}
	D := Point{Coord(x2 + dy), Coord(y2 - dx)}

	sc := func(p1, p2 Point) Coord {
		return p1.X*p2.X + p1.Y*p2.Y
	}

	AM := Point{p.X - A.X, p.Y - A.Y}
	AB := Point{B.X - A.X, B.Y - A.Y}
	AD := Point{D.X - A.X, D.Y - A.Y}
	log.Println(A, B, C, D)
	return 0 <= sc(AM, AB) && sc(AM, AB) < sc(AB, AB) &&
		0 <= sc(AM, AD) && sc(AM, AD) < sc(AD, AD)
}

// TowerState is a struct that represents the state of a tower.
type TowerState struct {
	AimType    Aim
	IsTurnedOn bool
	CoolDown   Frames
	Point      Point
	Aim        *Enemy
}
