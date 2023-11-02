package models

import (
	"image"
	"math"
)

type Enemy struct {
	Name       string
	State      EnemyState
	Path       []Point
	MaxHealth  int
	Vrms       Coord
	Damage     int
	MoneyAward int
	Weaknesses map[TypeAttack]Weakness
	Strengths  map[TypeAttack]Strength
	Image      image.Image
}

// ChangeDirection направляет врага на новую точку, если возможно.
// Меняются скорость и время, через которое враг прибудет в неё.
// Так же выставляет врага в реальное положение на карте, поскольку
// инкрементно изменяющаяся координата склонна накапливать погрешности.
//
// Если враг пришёл в финальную точку - выставляет флаг Dead = true.
func (e *Enemy) changeDirection() {
	e.State.CurrPoint++
	e.State.Pos = e.Path[e.State.CurrPoint]

	if e.State.CurrPoint == len(e.Path)-1 {
		e.Die()
		return
	}

	curr := e.Path[e.State.CurrPoint]
	next := e.Path[e.State.CurrPoint+1]

	dX := next.X - curr.X
	dY := next.Y - curr.Y

	t := math.Hypot(float64(dX), float64(dY)) / float64(e.Vrms)
	frameTime := int(math.Round(t))

	e.State.Vx = dX / Coord(frameTime)
	e.State.Vy = dY / Coord(frameTime)
	e.State.TimeNextPointLeft = frameTime
}

func (e *Enemy) Die() {
	e.State.Dead = true
}

func (e *Enemy) FinalDamage(t TypeAttack, dmg int) int {
	for k, v := range e.Weaknesses {
		if k == t {
			return v.IncDamage(dmg)
		}
	}
	for k, v := range e.Strengths {
		if k == t {
			return v.DecDamage(dmg)
		}
	}

	return dmg
}

func (e *Enemy) move() {
	e.State.Pos.X += e.State.Vx
	e.State.Pos.Y += e.State.Vy
	e.State.TimeNextPointLeft--
}

func (e *Enemy) Update() {
	if e.State.CurrPoint == -1 {
		e.changeDirection()
	}
	e.move()
	if e.State.TimeNextPointLeft == 0 {
		e.changeDirection()
	}
}

type EnemyState struct {
	CurrPoint         int
	Pos               Point
	Vx                Coord
	Vy                Coord
	Health            int
	TimeNextPointLeft Frames
	Dead              bool
}

type Weakness struct {
	T      TypeAttack
	IncDmg int
}

func (w Weakness) IncDamage(damage int) int {
	return damage + w.IncDmg
}

type Strength struct {
	T      TypeAttack
	DecDmg int
}

func (w Strength) DecDamage(damage int) int {
	return max(damage-w.DecDmg, 0)
}
