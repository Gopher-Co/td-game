package ingame

import (
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/icza/gox/imagex/colorx"

	"github.com/gopher-co/td-game/models/config"
	"github.com/gopher-co/td-game/models/general"
)

// Enemy is an entity moving on the Path and trying to
// achieve its end to deal damage to the Player.
//
// In the beginning MaxHealth == State.Health.
type Enemy struct {
	// Name is a name of the enemy.
	Name string

	// State is a state of the enemy.
	State EnemyState

	// Path is a path of the enemy.
	Path Path

	// MaxHealth is a maximal health of the enemy.
	MaxHealth int

	// Vrms is a root mean square speed of the enemy.
	Vrms general.Coord

	// Damage is a damage of the enemy.
	Damage int

	// MoneyAward is a money award for killing the enemy.
	MoneyAward int

	// Weaknesses is a list of weaknesses of the enemy.
	Weaknesses map[general.TypeAttack]Weakness

	// Strengths is a list of strengths of the enemy.
	Strengths map[general.TypeAttack]Strength

	// Image is an image of the enemy.
	Image *ebiten.Image
}

// NewEnemy creates a new entity of Enemy.
//
// Panics if cfg.Color is not a correct hex-string of "#xxxxxx".
func NewEnemy(cfg *config.Enemy, path Path) *Enemy {
	en := &Enemy{
		Name: cfg.Name,
		State: EnemyState{
			CurrPoint: -1,
			Pos:       path[0],
			Health:    cfg.MaxHealth,
			Dead:      false,
		},
		Path:       path,
		MaxHealth:  cfg.MaxHealth,
		Vrms:       cfg.Vrms,
		Damage:     cfg.Damage,
		MoneyAward: cfg.MoneyAward,
		Weaknesses: map[general.TypeAttack]Weakness{},
		Strengths:  map[general.TypeAttack]Strength{},
	}

	for _, v := range cfg.Strengths {
		en.Strengths[v.T] = Strength(v)
	}

	for _, v := range cfg.Weaknesses {
		en.Weaknesses[v.T] = Weakness(v)
	}
	if cfg.Image() != nil {
		en.Image = cfg.Image()
	} else {
		clr, err := colorx.ParseHexColor(cfg.Name)
		if err != nil {
			panic(err)
		}

		img := ebiten.NewImage(32, 32)
		vector.DrawFilledCircle(img, 16, 16, 16, clr, true)
	}

	en.changeDirection()

	return en
}

// DealDamageToPlayer returns the final damage to the player.
func (e *Enemy) DealDamageToPlayer() int {
	if e.State.Dead && e.State.FinalDamage > 0 {
		fd := e.State.FinalDamage
		e.State.FinalDamage = 0
		return fd
	}
	return 0
}

// DealDamage decreases the health of the enemy on dmg points.
// If health is less than dmg, health will become zero.
func (e *Enemy) DealDamage(dmg int) {
	e.State.Health = max(0, e.State.Health-dmg)
}

// Draw draws the enemy on the screen.
func (e *Enemy) Draw(screen *ebiten.Image) {
	geom := ebiten.GeoM{}
	geom.Translate(float64(e.State.Pos.X-float32(e.Image.Bounds().Dx()/2)), float64(e.State.Pos.Y-float32(e.Image.Bounds().Dy()/2)))
	screen.DrawImage(e.Image, &ebiten.DrawImageOptions{GeoM: geom})
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
		e.State.FinalDamage = e.Damage
		e.State.PassPath = true
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

	e.State.Vx = dX / general.Coord(frameTime)
	e.State.Vy = dY / general.Coord(frameTime)
	e.State.TimeNextPointLeft = frameTime
}

// Die marks the enemy dead.
func (e *Enemy) Die() {
	e.State.Dead = true
}

// FinalDamage returns the final damage depending on weaknesses
// and strengths of the enemy.
func (e *Enemy) FinalDamage(t general.TypeAttack, dmg int) int {
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
	if e.State.Dead {
		return
	}

	if e.State.Health == 0 {
		e.Die()
		return
	}

	// calculate velocities on the first iteration
	if e.State.CurrPoint == -1 {
		e.changeDirection()
	}
	e.move()
	if e.State.TimeNextPointLeft == 0 {
		e.changeDirection()
	}
}

// move moves the enemy to the next point.
func (e *Enemy) move() {
	e.State.Pos.X += e.State.Vx
	e.State.Pos.Y += e.State.Vy
	e.State.TimeNextPointLeft--
}

// EnemyState is a struct
type EnemyState struct {
	// CurrPoint is a current point in Path.
	CurrPoint int

	// Pos is a current position of the enemy.
	Pos general.Point

	// Vx is a velocity on X-axis.
	Vx general.Coord

	// Vy is a velocity on Y-axis.
	Vy general.Coord

	// Health is a current health of the enemy.
	Health int

	// TimeNextPointLeft is a time left to the next point.
	TimeNextPointLeft general.Frames

	// Dead is a flag that shows if the enemy is dead.
	Dead bool

	// PassPath is a flag that shows if the enemy has passed the path.
	PassPath bool

	// FinalDamage is a final damage to the player.
	FinalDamage int
}

// Weakness stores effects that are detrimental to the enemy
type Weakness struct {
	T      general.TypeAttack
	IncDmg int
}

// IncDamage returns increased damage.
func (w Weakness) IncDamage(damage int) int {
	return damage + w.IncDmg
}

// Strength stores effects that are useful to the enemy
type Strength struct {
	T      general.TypeAttack
	DecDmg int
}

// DecDamage returns decreased damage.
func (w Strength) DecDamage(damage int) int {
	return max(damage-w.DecDmg, 0)
}
