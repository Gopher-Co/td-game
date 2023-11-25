package ingame

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"

	"github.com/gopher-co/td-game/models/config"
	"github.com/gopher-co/td-game/models/general"
)

// Map is a struct that represents a map.
type Map struct {
	// Towers on the map now.
	Towers []*Tower

	// Enemies on the map now.
	Enemies []*Enemy

	// Projectiles on the map now.
	Projectiles []*Projectile

	// Path is a path of the map.
	Path Path

	// Image is an image of the map.
	Image *ebiten.Image
}

// NewMap creates a new entity of Map.
func NewMap(config *config.Map) *Map {
	m := &Map{
		Path:  config.Path,
		Image: config.Image(),
	}

	return m
}

// Update updates the map.
func (m *Map) Update() {
	for _, v := range m.Enemies {
		v.Update()
	}
	for _, v := range m.Towers {
		v.Update()
		v.TakeAim(m.Enemies)
		if p := v.Launch(); p != nil {
			m.Projectiles = append(m.Projectiles, p)
		}
	}
	for _, v := range m.Projectiles {
		v.Update()
	}
}

// Draw draws the map.
func (m *Map) Draw(screen *ebiten.Image) {
	screen.DrawImage(m.Image, nil)
	m.Path.Draw(screen)
	for _, p := range m.Projectiles {
		if p.dead {
			continue
		}
		p.Draw(screen)
	}

	for _, t := range m.Towers {
		t.Draw(screen)
	}

	for _, e := range m.Enemies {
		if !e.State.Dead {
			e.Draw(screen)
		}
	}
}

// AreThereAliveEnemies returns true if there are alive enemies on the map.
func (m *Map) AreThereAliveEnemies() bool {
	for _, e := range m.Enemies {
		if !e.State.Dead {
			return true
		}
	}

	return false
}

// Path is a struct that represents a path.
type Path []general.Point

// Draw draws the path.
func (p Path) Draw(screen *ebiten.Image) {
	for i := 0; i < len(p)-1; i++ {
		drawLine(screen, p[i], p[i+1])
	}
}

// drawLine draws a line between two points.
func drawLine(screen *ebiten.Image, p1, p2 general.Point) {
	x1, y1, x2, y2 := p1.X, p1.Y, p2.X, p2.Y
	vector.DrawFilledCircle(screen, x2, y2, config.PathWidth/2, color.RGBA{R: 12, G: 23, B: 34, A: 255}, true)
	vector.StrokeLine(screen, x1, y1, x2, y2, config.PathWidth, color.RGBA{R: 12, G: 23, B: 34, A: 255}, true)
}
