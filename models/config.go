package models

// Config structures are need to pass then to NewXXX functions.

type EnemyConfig struct {
	Name       string     `json:"name"`
	MaxHealth  int        `json:"max_health"`
	Damage     int        `json:"damage"`
	MoneyAward int        `json:"money_award"`
	Strengths  []Strength `json:"strengths"`
	Weaknesses []Weakness `json:"weaknesses"`
}

type TowerConfig struct {
	Upgrades        []UpgradeConfig `json:"upgrades"`
	Color           string          `json:"color"`
	Price           int             `json:"price"`
	InitDamage      int             `json:"initial_damage"`
	InitRadius      Coord           `json:"initial_radius"`
	InitSpeedAttack Frames          `json:"initial_speed_attack"`
	OpenLevel       int             `json:"open_level"`
}

type UpgradeConfig struct {
	Price            int    `json:"price"`
	DeltaDamage      int    `json:"delta_damage"`
	DeltaSpeedAttack Frames `json:"delta_speed_attack"`
	DeltaRadius      Coord  `json:"delta_radius"`
	OpenLevel        int    `json:"open_level"`
}

type LevelConfig struct {
	LevelName string       `json:"level_name"`
	Map       MapConfig    `json:"map"`
	Waves     []WaveConfig `json:"waves"`
}

type WaveConfig struct {
	Swarms []EnemySwarmConfig `json:"swarms"`
}

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

type UIConfig struct {
	// Colors contains hex-colors (e.g. "#AB0BA0") for each key in map
	Colors map[string]string `json:"colors"`
}

type MapConfig struct {
	BackgroundColor string  `json:"background_color"`
	Path            []Point `json:"path"`
}
