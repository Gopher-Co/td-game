package models

import (
	"github.com/hajimehoshi/ebiten/v2"
)

// Projectile is an entity generated by towers that flies to the enemy
// and deals the damage to it.
// Projectiles never misses the enemy and achieves the aim when TTL is equal to zero.
type Projectile struct {
	// Pos is a position of the projectile.
	Pos Point

	// Vrms is a root mean square speed of the projectile.
	Vrms Coord

	// Vx is a speed of the projectile on the X axis.
	Vx Coord

	// Vy is a speed of the projectile on the Y axis.
	Vy Coord

	// Type is a type of the projectile.
	Type TypeAttack

	// Damage is a damage of the projectile.
	Damage int

	// TTL is a time to live of the projectile.
	TTL Frames

	// TargetEnemy is an enemy that the projectile is flying to.
	TargetEnemy *Enemy

	// Image is an image of the projectile.
	Image *ebiten.Image

	// dead is a flag that shows if the projectile is dead.
	dead bool
}

// Update updates the projectile.
func (p *Projectile) Update() {
	p.move()
	if !p.dead && p.TTL == 0 {
		p.EnemyHit()
		p.dead = true
	}
}

// Draw draws the projectile.
func (p *Projectile) Draw(screen *ebiten.Image) {
	geom := ebiten.GeoM{}
	geom.Translate(float64(p.Pos.X-float32(p.Image.Bounds().Dx()/2)), float64(p.Pos.Y-float32(p.Image.Bounds().Dy()/2)))
	screen.DrawImage(p.Image, &ebiten.DrawImageOptions{GeoM: geom})
}

// move moves the projectile on the map.
func (p *Projectile) move() {
	p.Pos.X += p.Vx
	p.Pos.Y += p.Vy
	p.TTL = max(p.TTL-1, 0)
}

// EnemyHit checks if the projectile hit the enemy and returns true if it is.
func (p *Projectile) EnemyHit() {
	p.TargetEnemy.DealDamage(p.TargetEnemy.FinalDamage(p.Type, p.Damage))
}
