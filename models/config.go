package models

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/icza/gox/imagex/colorx"
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

// EnemyConfig is a config for enemy.
type EnemyConfig struct {
	// Name is a name of the enemy.
	Name string `json:"name"`

	// MaxHealth is a maximal health of the enemy.
	MaxHealth int `json:"max_health"`

	// Damage is a damage of the enemy.
	Damage int `json:"damage"`

	// Vrms is a root mean square speed of the enemy.
	Vrms Coord `json:"vrms"`

	// MoneyAward is a money award for killing the enemy.
	MoneyAward int `json:"money_award"`

	// Strengths is a list of strengths of the enemy.
	Strengths []Strength `json:"strengths"`

	// Weaknesses is a list of weaknesses of the enemy.
	Weaknesses []Weakness `json:"weaknesses"`

	image *ebiten.Image
}

// InitImage initializes image from the temporary state of the entity.
func (c *EnemyConfig) InitImage() error {
	clr, err := colorx.ParseHexColor(c.Name)
	if err != nil {
		return err
	}

	img := ebiten.NewImage(EnemyImageWidth, EnemyImageWidth)
	vector.DrawFilledCircle(img, EnemyImageWidth/2, EnemyImageWidth/2, EnemyImageWidth/2, clr, true)
	c.image = img

	return nil
}

// Image returns image.
func (c *EnemyConfig) Image() *ebiten.Image {
	return c.image
}

// TowerConfig is a config for tower.
type TowerConfig struct {
	// Name is a name of the tower.
	Name string `json:"name"`

	// Upgrades is a list of upgrades of the tower.
	Upgrades []UpgradeConfig `json:"upgrades"`

	// Price is a price of the tower.
	Price int `json:"price"`

	// Type is a type of the tower attack.
	Type TypeAttack `json:"type"`

	// InitDamage is an initial damage of the tower.
	InitDamage int `json:"initial_damage"`

	// InitRadius is an initial radius of the tower.
	InitRadius Coord `json:"initial_radius"`

	// InitSpeedAttack is an initial speed attack of the tower.
	InitSpeedAttack Frames `json:"initial_speed_attack"`

	// InitProjectileVrms is an initial projectile vrms of the tower.
	InitProjectileVrms Coord `json:"init_projectile_speed"`

	// ProjectileConfig is a config for projectile.
	ProjectileConfig ProjectileConfig `json:"projectile_config"`

	// OpenLevel is a level when the tower can be opened.
	OpenLevel int `json:"open_level"`

	image *ebiten.Image
}

// InitImage initializes image from the temporary state of the entity.
func (c *TowerConfig) InitImage() error {
	if err := c.ProjectileConfig.InitImage(); err != nil {
		return err
	}

	clr, err := colorx.ParseHexColor(c.Name)
	if err != nil {
		return err
	}

	img := ebiten.NewImage(TowerImageWidth, TowerImageWidth)
	vector.DrawFilledRect(img, 0, 0, TowerImageWidth, TowerImageWidth, clr, true)
	c.image = img

	return nil
}

// Image returns image.
func (c *TowerConfig) Image() *ebiten.Image {
	return c.image
}

// InitUpgrades initializes upgrades from the temporary state of the entity.
func (c *TowerConfig) InitUpgrades() []*Upgrade {
	ups := make([]*Upgrade, len(c.Upgrades))

	for i := 0; i < len(ups); i++ {
		ups[i] = NewUpgrade(&c.Upgrades[i])
	}

	return ups
}

// UpgradeConfig is a config for tower's upgrade.
type UpgradeConfig struct {
	// Price is a price of the upgrade.
	Price int `json:"price"`

	// DeltaDamage is a delta damage of the upgrade.
	DeltaDamage int `json:"delta_damage"`

	// DeltaSpeedAttack is a delta speed attack of the upgrade.
	DeltaSpeedAttack Frames `json:"delta_speed_attack"`

	// DeltaRadius is a delta radius of the upgrade.
	DeltaRadius Coord `json:"delta_radius"`

	// OpenLevel is a level when the upgrade can be opened.
	OpenLevel int `json:"open_level"`
}

// ProjectileConfig is a config for projectile.
type ProjectileConfig struct {
	Name  string
	image *ebiten.Image
}

// InitImage initializes image from the temporary state of the entity.
func (c *ProjectileConfig) InitImage() error {
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
func (c *ProjectileConfig) Image() *ebiten.Image {
	return c.image
}

// LevelConfig is a config for level.
type LevelConfig struct {
	// LevelName is a name of the level.
	LevelName string `json:"level_name"`

	// MapName is a name of the map.
	MapName string `json:"map_name"`

	// GameRule is a config for game rule.
	GameRule GameRuleConfig `json:"game_rule"`
}

// GameRuleConfig is a config for game rule.
type GameRuleConfig []WaveConfig

// WaveConfig is a config for wave.
type WaveConfig struct {
	Swarms []EnemySwarmConfig `json:"swarms"`
}

// EnemySwarmConfig is a config for enemy swarm.
type EnemySwarmConfig struct {
	// EnemyName is a name of the enemy.
	EnemyName string `json:"enemy_name"`

	// Timeout is the time when the first enemy can be called.
	Timeout Frames `json:"timeout"`

	// Interval is time between calls.
	Interval Frames `json:"interval"`

	// MaxCalls is a maximal amount of enemies that can be called.
	MaxCalls int `json:"max_calls"`
}

// UIConfig is a config for GlobalUI.
type UIConfig struct {
	// Colors contains hex-colors (e.g. "#AB0BA0") for each key in map
	Colors map[string]string `json:"colors"`
}

// MapConfig is a config for map.
type MapConfig struct {
	// Name is a name of the map.
	Name string `json:"name"`

	// BackgroundColor is a background color of the map.
	BackgroundColor string `json:"background_color"`

	// Path is a path of the map.
	Path []Point `json:"path"`

	image *ebiten.Image
}

// InitImage initializes image from the temporary state of the entity.
func (c *MapConfig) InitImage() error {
	clr, err := colorx.ParseHexColor(c.BackgroundColor)
	if err != nil {
		return err
	}

	img := ebiten.NewImage(1920, 1080)
	img.Fill(clr)

	c.image = img

	return nil
}

// Image returns image.
func (c *MapConfig) Image() *ebiten.Image {
	return c.image
}
