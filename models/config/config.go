package config

import (
	"fmt"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/icza/gox/imagex/colorx"

	"github.com/gopher-co/td-game/models/general"
	"github.com/gopher-co/td-game/ui"
)

const (
	// EnemyImageWidth is a width of the enemy image.
	EnemyImageWidth = 48

	// TowerImageWidth is a width of the tower image.
	TowerImageWidth = 64

	// ProjectileImageWith is a width of the projectile image.
	ProjectileImageWith = 32

	// PathWidth is a width of the path.
	PathWidth = 64
)

// Config structures are need to pass then to NewXXX functions.

// Enemy is a config for enemy.
type Enemy struct {
	// Name is a name of the enemy.
	Name string `json:"name"`

	// MaxHealth is a maximal health of the enemy.
	MaxHealth int `json:"max_health"`

	// Damage is a damage of the enemy.
	Damage int `json:"damage"`

	// Vrms is a root mean square speed of the enemy.
	Vrms general.Coord `json:"vrms"`

	// MoneyAward is a money award for killing the enemy.
	MoneyAward int `json:"money_award"`

	// Strengths is a list of strengths of the enemy.
	Strengths []Strength `json:"strengths"`

	// Weaknesses is a list of weaknesses of the enemy.
	Weaknesses []Weakness `json:"weaknesses"`

	image *ebiten.Image
}

// Strength is a config for strength.
type Strength struct {
	// T is a type of the strength.
	T general.TypeAttack `json:"type"`

	// DecDmg is a damage decrement of the strength.
	DecDmg int `json:"dec_dmg"`
}

// Weakness is a config for weakness.
type Weakness struct {
	// T is a type of the weakness.
	T general.TypeAttack `json:"type"`

	// IncDmg is a damage increment of the weakness.
	IncDmg int `json:"inc_dmg"`
}

// InitImage initializes image from the temporary state of the entity.
func (c *Enemy) InitImage() error {
	clr, err := colorx.ParseHexColor(c.Name)
	if err == nil {
		img := ebiten.NewImage(EnemyImageWidth, EnemyImageWidth)
		vector.DrawFilledRect(img, 0, 0, EnemyImageWidth, EnemyImageWidth, clr, true)
		c.image = img

		return nil
	}
	png, err := ui.InitPNG("./assets/" + c.Name)
	if err == nil {
		img := ebiten.NewImage(EnemyImageWidth, EnemyImageWidth)
		geom := ebiten.GeoM{}
		geom.Scale(float64(EnemyImageWidth)/float64(png.Bounds().Dx()), float64(EnemyImageWidth)/float64(png.Bounds().Dy()))
		img.DrawImage(png, &ebiten.DrawImageOptions{GeoM: geom})
		c.image = img

		return nil
	}

	return fmt.Errorf("image init failed: %w", err)
}

// Image returns image.
func (c *Enemy) Image() *ebiten.Image {
	return c.image
}

// Tower is a config for tower.
type Tower struct {
	// Name is a name of the tower.
	Name string `json:"name"`

	// Upgrades is a list of upgrades of the tower.
	Upgrades []Upgrade `json:"upgrades"`

	// Price is a price of the tower.
	Price int `json:"price"`

	// Type is a type of the tower attack.
	Type general.TypeAttack `json:"type"`

	// InitDamage is an initial damage of the tower.
	InitDamage int `json:"initial_damage"`

	// InitRadius is an initial radius of the tower.
	InitRadius general.Coord `json:"initial_radius"`

	// InitSpeedAttack is an initial speed attack of the tower.
	InitSpeedAttack general.Frames `json:"initial_speed_attack"`

	// InitProjectileVrms is an initial projectile vrms of the tower.
	InitProjectileVrms general.Coord `json:"init_projectile_speed"`

	// ProjectileConfig is a config for projectile.
	ProjectileConfig Projectile `json:"projectile_config"`

	// OpenLevel is a level when the tower can be opened.
	OpenLevel int `json:"open_level"`

	image *ebiten.Image
}

// InitImage initializes image from the temporary state of the entity.
func (c *Tower) InitImage() error {
	clr, err := colorx.ParseHexColor(c.Name)
	if err == nil {
		img := ebiten.NewImage(TowerImageWidth, TowerImageWidth)
		vector.DrawFilledRect(img, 0, 0, TowerImageWidth, TowerImageWidth, clr, true)
		c.image = img

		return nil
	}
	png, err := ui.InitPNG("./assets/" + c.Name)
	if err == nil {
		img := ebiten.NewImage(TowerImageWidth, TowerImageWidth)
		geom := ebiten.GeoM{}
		geom.Scale(float64(TowerImageWidth)/float64(png.Bounds().Dx()), float64(TowerImageWidth)/float64(png.Bounds().Dy()))
		img.DrawImage(png, &ebiten.DrawImageOptions{GeoM: geom})
		c.image = img

		return nil
	}

	return fmt.Errorf("image init failed: %w", err)
}

// Image returns image.
func (c *Tower) Image() *ebiten.Image {
	return c.image
}

// Upgrade is a config for tower's upgrade.
type Upgrade struct {
	// Price is a price of the upgrade.
	Price int `json:"price"`

	// DeltaDamage is a delta damage of the upgrade.
	DeltaDamage int `json:"delta_damage"`

	// DeltaSpeedAttack is a delta speed attack of the upgrade.
	DeltaSpeedAttack general.Frames `json:"delta_speed_attack"`

	// DeltaRadius is a delta radius of the upgrade.
	DeltaRadius general.Coord `json:"delta_radius"`

	// OpenLevel is a level when the upgrade can be opened.
	OpenLevel int `json:"open_level"`
}

// Projectile is a config for projectile.
type Projectile struct {
	Name  string
	image *ebiten.Image
}

// InitImage initializes image from the temporary state of the entity.
func (c *Projectile) InitImage() error {
	clr, err := colorx.ParseHexColor(c.Name)
	if err != nil {
		return err
	}

	img := ebiten.NewImage(ProjectileImageWith, ProjectileImageWith)
	vector.DrawFilledRect(img, 0, 0, ProjectileImageWith, ProjectileImageWith, clr, true)
	c.image = img

	return nil
}

// Image returns image.
func (c *Projectile) Image() *ebiten.Image {
	return c.image
}

// Level is a config for level.
type Level struct {
	// LevelName is a name of the level.
	LevelName string `json:"level_name"`

	// MapName is a name of the map.
	MapName string `json:"map_name"`

	// GameRule is a config for game rule.
	GameRule GameRule `json:"game_rule"`

	// needed for level numeration
	Order int `json:"-"`
}

// GameRule is a config for game rule.
type GameRule []Wave

// Wave is a config for wave.
type Wave struct {
	Swarms []EnemySwarm `json:"swarms"`
}

// EnemySwarm is a config for enemy swarm.
type EnemySwarm struct {
	// EnemyName is a name of the enemy.
	EnemyName string `json:"enemy_name"`

	// Timeout is the time when the first enemy can be called.
	Timeout general.Frames `json:"timeout"`

	// Interval is time between calls.
	Interval general.Frames `json:"interval"`

	// MaxCalls is a maximal amount of enemies that can be called.
	MaxCalls int `json:"max_calls"`
}

// UI is a config for GlobalUI.
type UI struct {
	// Colors contains hex-colors (e.g. "#AB0BA0") for each key in map
	Colors map[string]string `json:"colors"`
}

// Map is a config for map.
type Map struct {
	// Name is a name of the map.
	Name string `json:"name"`

	// BackgroundColor is a background color of the map.
	BackgroundColor string `json:"background_color"`

	// Path is a path of the map.
	Path []general.Point `json:"path"`

	image *ebiten.Image
}

// InitImage initializes image from the temporary state of the entity.
func (c *Map) InitImage() error {
	img, err := ui.InitColor(c.BackgroundColor)
	if err != nil {
		img, err = ui.InitPNG(c.BackgroundColor)
		if err != nil {
			return fmt.Errorf("couldn't load image for map: %w", err)
		}
	}

	geom := ebiten.GeoM{}
	geom.Scale(1500./float64(img.Bounds().Dx()), 1080./float64(img.Bounds().Dy()))

	c.image = ebiten.NewImage(1500, 1080)
	c.image.DrawImage(img, &ebiten.DrawImageOptions{GeoM: geom})

	return nil
}

// Image returns image.
func (c *Map) Image() *ebiten.Image {
	return c.image
}
