package ingame

import (
	"github.com/hajimehoshi/ebiten/v2"

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
		if v.Sold || !v.State.IsTurnedOn {
			continue
		}
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
	geom := ebiten.GeoM{}
	geom.Scale(float64(1500)/float64(m.Image.Bounds().Dx()), float64(1080)/float64(m.Image.Bounds().Dy()))
	screen.DrawImage(m.Image, &ebiten.DrawImageOptions{GeoM: geom})
	m.Path.Draw(screen)
	for _, p := range m.Projectiles {
		if p.dead {
			continue
		}
		p.Draw(screen)
	}

	for _, t := range m.Towers {
		if t.Sold {
			continue
		}
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
func (p Path) Draw(_ *ebiten.Image) {
	//for i := 0; i < len(p)-1; i++ {
	//	drawLine(screen, p[i], p[i+1])
	//}
}
