package models

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/icza/gox/imagex/colorx"
)

const (
	EnemyImageWidth     = 24
	TowerImageWidth     = 32
	ProjectileImageWith = 16
	PathWidth           = 32
)

// Config structures are need to pass then to NewXXX functions.

// EnemyConfig is a config for enemy.
type EnemyConfig struct {
	Name       string     `json:"name"`
	MaxHealth  int        `json:"max_health"`
	Damage     int        `json:"damage"`
	Vrms       Coord      `json:"vrms"`
	MoneyAward int        `json:"money_award"`
	Strengths  []Strength `json:"strengths"`
	Weaknesses []Weakness `json:"weaknesses"`
	image      *ebiten.Image
}

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

func (c *EnemyConfig) Image() *ebiten.Image {
	return c.image
}

// TowerConfig is a config for tower.
type TowerConfig struct {
	Name               string           `json:"name"`
	Upgrades           []UpgradeConfig  `json:"upgrades"`
	Price              int              `json:"price"`
	Type               TypeAttack       `json:"type"`
	InitDamage         int              `json:"initial_damage"`
	InitRadius         Coord            `json:"initial_radius"`
	InitSpeedAttack    Frames           `json:"initial_speed_attack"`
	InitProjectileVrms Coord            `json:"init_projectile_speed"`
	ProjectileConfig   ProjectileConfig `json:"projectile_config"`
	OpenLevel          int              `json:"open_level"`
	image              *ebiten.Image
}

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

func (c *TowerConfig) Image() *ebiten.Image {
	return c.image
}

func (c *TowerConfig) InitUpgrades() []*Upgrade {
	ups := make([]*Upgrade, len(c.Upgrades))

	for i := 0; i < len(ups); i++ {
		ups[i] = NewUpgrade(&c.Upgrades[i])
	}

	return ups
}

// UpgradeConfig is a config for tower's upgrade.
type UpgradeConfig struct {
	Price            int    `json:"price"`
	DeltaDamage      int    `json:"delta_damage"`
	DeltaSpeedAttack Frames `json:"delta_speed_attack"`
	DeltaRadius      Coord  `json:"delta_radius"`
	OpenLevel        int    `json:"open_level"`
}

type ProjectileConfig struct {
	Name  string
	image *ebiten.Image
}

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

func (c *ProjectileConfig) Image() *ebiten.Image {
	return c.image
}

// LevelConfig is a config for level.
type LevelConfig struct {
	LevelName string         `json:"level_name"`
	Map       MapConfig      `json:"map"`
	GameRule  GameRuleConfig `json:"game_rule"`
}

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

	// CurrTime is current time relatively the swarm's start.
	CurrTime Frames `json:"curr_time"`

	// MaxCalls is a maximal amount of enemies that can be called.
	MaxCalls int `json:"max_calls"`

	// CurCalls is the current amount of enemies called.
	CurCalls int `json:"cur_calls"`
}

// UIConfig is a config for UI.
type UIConfig struct {
	// Colors contains hex-colors (e.g. "#AB0BA0") for each key in map
	Colors map[string]string `json:"colors"`
}

// MapConfig is a config for map.
type MapConfig struct {
	BackgroundColor string  `json:"background_color"`
	Path            []Point `json:"path"`
	image           *ebiten.Image
}

// InitImage initializes image from the temporary state of the entity.
func (c *MapConfig) InitImage() error {
	clr, err := colorx.ParseHexColor(c.BackgroundColor)
	if err != nil {
		return err
	}

	img := ebiten.NewImage(ebiten.WindowSize())
	img.Fill(clr)

	c.image = img

	return nil
}

// Image returns image.
func (c *MapConfig) Image() *ebiten.Image {
	return c.image
}
