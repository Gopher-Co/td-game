package models

import (
	"image"
	"math"
)

// Enemy is an entity moving on the Path and trying to
// achieve its end to deal damage to the Player.
//
// In the beginning MaxHealth == State.Health.
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

// DealDamage decreases the health of the enemy on dmg points.
// If health is less than dmg, health will become zero.
func (e *Enemy) DealDamage(dmg int) {
	e.State.Health = max(0, e.State.Health-dmg)
}

// changeDirection directs the enemy to a new point, if possible.
// Changes the speed and time after which the enemy will arrive at it.
// Also puts the enemy in a real position on the map, since
// an incrementally changing coordinate tends to accumulate errors.
//
// If the enemy has reached the final point, sets the flag Dead = true.
func (e *Enemy) changeDirection() {
	// mark that the enemy achieved next point in Path
	e.State.CurrPoint++

	// set a real state
	// errors may have accumulated
	e.State.Pos = e.Path[e.State.CurrPoint]

	// event on achieving the end
	if e.State.CurrPoint == len(e.Path)-1 {
		// todo deal damage to player
		e.Die()
		return
	}

	// calculating new Vx, Vy and TimeNextPointLeft
	curr := e.Path[e.State.CurrPoint]
	next := e.Path[e.State.CurrPoint+1]

	dX := next.X - curr.X
	dY := next.Y - curr.Y

	t := math.Hypot(float64(dX), float64(dY)) / float64(e.Vrms) // t = S / Vrms
	frameTime := int(math.Round(t))

	e.State.Vx = dX / Coord(frameTime)
	e.State.Vy = dY / Coord(frameTime)
	e.State.TimeNextPointLeft = frameTime
}

// Die marks the enemy dead.
func (e *Enemy) Die() {
	e.State.Dead = true
}

// FinalDamage returns the final damage depending on weaknesses
// and strengths of the enemy.
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

// Update is a universal method for updating enemy's state by itself.
func (e *Enemy) Update() {
	// calculate velocities on the first iteration
	if e.State.CurrPoint == -1 {
		e.changeDirection()
	}
	e.move()
	if e.State.TimeNextPointLeft == 0 {
		e.changeDirection()
	}
}

func (e *Enemy) move() {
	e.State.Pos.X += e.State.Vx
	e.State.Pos.Y += e.State.Vy
	e.State.TimeNextPointLeft--
}

// EnemyState is a struct
type EnemyState struct {
	CurrPoint         int
	Pos               Point
	Vx                Coord
	Vy                Coord
	Health            int
	TimeNextPointLeft Frames
	Dead              bool
}

// Weakness stores effects that are detrimental to the enemy
type Weakness struct {
	T      TypeAttack
	IncDmg int
}

// IncDamage returns increased damage.
func (w Weakness) IncDamage(damage int) int {
	return damage + w.IncDmg
}

// Strength stores effects that are useful to the enemy
type Strength struct {
	T      TypeAttack
	DecDmg int
}

// DecDamage returns decreased damage.
func (w Strength) DecDamage(damage int) int {
	return max(damage-w.DecDmg, 0)
}
