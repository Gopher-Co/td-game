package ingame

import (
	"cmp"
	"image/color"
	"math"
	"slices"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"

	"github.com/gopher-co/td-game/models/config"
	"github.com/gopher-co/td-game/models/general"
)

// Aim is a type that represents the enemy that tower attacks.
type Aim int

const (
	// First is a type of aim that represents the first enemy.
	First = Aim(iota)

	// Weakest is a type of aim that represents the weakest enemy.
	Weakest

	// Strongest is a type of aim that represents the strongest enemy.
	Strongest
)

// Tower is a struct that represents a tower.
type Tower struct {
	Index int
	// Name is a name of the tower.
	Name string

	// Damage is a damage of the tower.
	Damage int

	// Type is a type of the tower.
	Type general.TypeAttack

	// Price is a price of the tower.
	Price int

	// Image is an image of the tower.
	Image *ebiten.Image

	// Radius is a radius of the tower.
	Radius general.Coord

	// State is a state of the tower.
	State TowerState

	// SpeedAttack is a speed of the tower's attack.
	SpeedAttack general.Frames

	// ProjectileVrms is a root mean square speed of the tower's projectile.
	ProjectileVrms general.Coord

	// ProjectileImage is an image of the tower's projectile.
	ProjectileImage *ebiten.Image

	// Upgrades is a list of upgrades of the tower.
	Upgrades []*Upgrade

	// UpgradesBought is a number of upgrades bought.
	UpgradesBought int

	// Chosen is a flag that shows if the tower is chosen.
	Chosen bool

	// Sold is a flag that shows if the tower is sold.
	Sold bool
}

var globalIndex int

// NewTower creates a new entity of Tower.
func NewTower(config *config.Tower, pos general.Point, path Path) *Tower {
	if CheckCollisionPath(pos, path) {
		return nil
	}

	initState := TowerState{
		AimType:    First,
		IsTurnedOn: true,
		CoolDown:   0,
		Pos:        pos,
		Aim:        nil,
	}

	t := &Tower{
		Index:           globalIndex,
		Name:            config.Name,
		Damage:          config.InitDamage,
		Type:            config.Type,
		Price:           config.Price,
		Image:           config.Image(),
		Radius:          config.InitRadius,
		State:           initState,
		SpeedAttack:     config.InitSpeedAttack,
		ProjectileVrms:  config.InitProjectileVrms,
		ProjectileImage: config.ProjectileConfig.Image(),
		UpgradesBought:  0,
	}
	globalIndex++

	t.initUpgrades(config.Upgrades)

	return t
}

// Launch launches a projectile from the tower.
func (t *Tower) Launch() *Projectile {
	if t.Sold || t.State.CoolDown != 0 || t.State.Aim == nil {
		return nil
	}
	t.State.CoolDown = (20 * 60) / t.SpeedAttack // N projectiles in 20 seconds (if TPS=60)

	p := &Projectile{
		Pos:         t.State.Pos,
		Vrms:        t.ProjectileVrms,
		Vx:          0,
		Vy:          0,
		Type:        t.Type,
		Damage:      t.Damage,
		TTL:         0,
		TargetEnemy: t.State.Aim,
		Image:       t.ProjectileImage,
	}
	target := p.TargetEnemy.State.Pos
	z := math.Hypot(float64(target.X-p.Pos.X), float64(target.Y-p.Pos.Y))

	ttl := math.Round(z / float64(p.Vrms))
	p.TTL = int(ttl)
	p.Vx = general.Coord(float64(target.X-p.Pos.X) / ttl)
	p.Vy = general.Coord(float64(target.Y-p.Pos.Y) / ttl)

	return p
}

// Upgrade upgrades the tower.
func (t *Tower) Upgrade(complete map[string]struct{}) bool {

	if t.UpgradesBought == len(t.Upgrades) {
		return false
	}

	upg := t.Upgrades[t.UpgradesBought]
	if complete != nil && upg.OpenLevel != "" {
		if _, ok := complete[upg.OpenLevel]; !ok {
			return false
		}
	}

	t.SpeedAttack += upg.DeltaSpeedAttack
	t.Damage += upg.DeltaDamage
	t.Radius += upg.DeltaRadius

	t.UpgradesBought++

	return true
}

// NextUpgrade returns the next upgrade.
func (t *Tower) NextUpgrade() *Upgrade {
	if len(t.Upgrades) == t.UpgradesBought {
		return nil
	}

	return t.Upgrades[t.UpgradesBought]
}

// Update updates the tower.
func (t *Tower) Update() {
	if t.Sold {
		return
	}
	t.State.CoolDown = max(t.State.CoolDown-1, 0)
}

// Draw draws the tower.
func (t *Tower) Draw(screen *ebiten.Image) {
	if t.Sold {
		return
	}
	if t.Chosen {
		vector.DrawFilledCircle(screen, t.State.Pos.X, t.State.Pos.Y, t.Radius, color.RGBA{A: 0x20}, true)
	}

	geom := ebiten.GeoM{}
	geom.Translate(float64(t.State.Pos.X-float32(t.Image.Bounds().Dx()/2)), float64(t.State.Pos.Y-float32(t.Image.Bounds().Dy()/2)))
	screen.DrawImage(t.Image, &ebiten.DrawImageOptions{GeoM: geom})
}

// TakeAim takes aim at the enemy.
func (t *Tower) TakeAim(e1 []*Enemy) {
	switch t.State.AimType {
	case First:
		t.takeAimFirst(e1)
	case Strongest:
		t.takeAimStrong(e1)
	case Weakest:
		t.takeAimWeak(e1)
	}
}

// takeAimFirst takes aim at the first enemy.
func (t *Tower) takeAimFirst(e1 []*Enemy) {
	enemies := slices.Clone(e1)
	enemies = slices.DeleteFunc(enemies, func(e *Enemy) bool {
		tx, ty, ex, ey := t.State.Pos.X, t.State.Pos.Y, e.State.Pos.X, e.State.Pos.Y
		return e.State.Dead || (tx-ex)*(tx-ex)+(ty-ey)*(ty-ey) > t.Radius*t.Radius
	})

	if len(enemies) == 0 {
		t.State.Aim = nil
		return
	}

	e := slices.MaxFunc(enemies, func(a, b *Enemy) int {
		if a.State.CurrPoint > b.State.CurrPoint {
			return 1
		} else if a.State.CurrPoint < b.State.CurrPoint {
			return -1
		}
		if general.Coord(a.State.TimeNextPointLeft)*a.Vrms < general.Coord(b.State.TimeNextPointLeft)*b.Vrms {
			return 1
		} else if general.Coord(a.State.TimeNextPointLeft)*a.Vrms > general.Coord(b.State.TimeNextPointLeft)*b.Vrms {
			return -1
		}
		return 0
	})

	t.State.Aim = e
}

// takeAimWeak takes aim at the weak enemy.
func (t *Tower) takeAimWeak(e1 []*Enemy) {
	enemies := slices.Clone(e1)
	enemies = slices.DeleteFunc(enemies, func(e *Enemy) bool {
		tx, ty, ex, ey := t.State.Pos.X, t.State.Pos.Y, e.State.Pos.X, e.State.Pos.Y
		return e.State.Dead || (tx-ex)*(tx-ex)+(ty-ey)*(ty-ey) > t.Radius*t.Radius
	})

	if len(enemies) == 0 {
		t.State.Aim = nil
		return
	}

	e := slices.MaxFunc(enemies, func(a, b *Enemy) int {
		return cmp.Compare(b.State.Health, a.State.Health)
	})

	t.State.Aim = e
}

// takeAimStrong takes aim at the strong enemy.
func (t *Tower) takeAimStrong(e1 []*Enemy) {
	enemies := slices.Clone(e1)
	enemies = slices.DeleteFunc(enemies, func(e *Enemy) bool {
		tx, ty, ex, ey := t.State.Pos.X, t.State.Pos.Y, e.State.Pos.X, e.State.Pos.Y
		return e.State.Dead || (tx-ex)*(tx-ex)+(ty-ey)*(ty-ey) > t.Radius*t.Radius
	})

	if len(enemies) == 0 {
		t.State.Aim = nil
		return
	}

	e := slices.MaxFunc(enemies, func(a, b *Enemy) int {
		return cmp.Compare(a.State.Health, b.State.Health)
	})

	t.State.Aim = e
}

// CheckCollisionPath checks if the path collides with the tower.
func CheckCollisionPath(pos general.Point, path Path) bool {
	for i := 0; i < len(path)-1; i++ {
		if checkCollision(general.Point{X: pos.X, Y: pos.Y}, path[i], path[i+1]) {
			return true
		}
	}

	return false
}

// checkCollision checks if the point collides with the line segment.
func checkCollision(p, p1, p2 general.Point) bool {
	x1, x2 := float64(p1.X), float64(p2.X)
	y1, y2 := float64(p1.Y), float64(p2.Y)

	sign := func(x float64) float64 {
		if math.Signbit(x) {
			return -1
		}
		return 1
	}

	// 		***
	// 		(x1,y1)**        z
	// 	 y	*b       *****
	// 		*	         a*****
	// 		()****************(x2,y2)
	//					x
	z := math.Hypot(x2-x1, y2-y1)

	sina := (y2 - y1) / z
	cosa := sign(x2-x1) * math.Sqrt(1-math.Pow(sina, 2))

	dx := config.PathWidth / 2 * cosa
	dy := config.PathWidth / 2 * sina

	x1 -= dx
	x2 += dx
	y1 -= dy
	y2 += dy

	A := general.Point{X: general.Coord(x1 - dy), Y: general.Coord(y1 + dx)}
	B := general.Point{X: general.Coord(x1 + dy), Y: general.Coord(y1 - dx)}
	D := general.Point{X: general.Coord(x2 + dy), Y: general.Coord(y2 - dx)}

	sc := func(p1, p2 general.Point) general.Coord {
		return p1.X*p2.X + p1.Y*p2.Y
	}

	AM := general.Point{X: p.X - A.X, Y: p.Y - A.Y}
	AB := general.Point{X: B.X - A.X, Y: B.Y - A.Y}
	AD := general.Point{X: D.X - A.X, Y: D.Y - A.Y}
	return 0 < sc(AM, AB) && sc(AM, AB) < sc(AB, AB) &&
		0 < sc(AM, AD) && sc(AM, AD) < sc(AD, AD)
}

// IsClicked checks if the tower is clicked.
func (t *Tower) IsClicked() bool {
	cx, cy := ebiten.CursorPosition()
	x1, x2 := int(t.State.Pos.X-config.TowerImageWidth/2), int(t.State.Pos.X+config.TowerImageWidth/2)
	y1, y2 := int(t.State.Pos.Y-config.TowerImageWidth/2), int(t.State.Pos.Y+config.TowerImageWidth/2)

	return cx >= x1 && cx <= x2 && cy >= y1 && cy <= y2 // inside the square
}

// initUpgrades initializes the upgrades of the tower.
func (t *Tower) initUpgrades(cfg []config.Upgrade) {
	ups := make([]*Upgrade, len(cfg))

	for i := 0; i < len(ups); i++ {
		ups[i] = NewUpgrade(&cfg[i])
	}

	t.Upgrades = ups
}

// TowerState is a struct that represents the state of a tower.
type TowerState struct {
	// AimType is a type of aim.
	AimType Aim

	// IsTurnedOn is a flag that shows if the tower is turned on.
	IsTurnedOn bool

	// CoolDown is a cool down of the tower.
	CoolDown general.Frames

	// Pos is a position of the tower.
	Pos general.Point

	// Aim is an enemy that the tower is aiming at.
	Aim *Enemy
}
