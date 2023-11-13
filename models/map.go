package models

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

// Map is a struct that represents a map.
type Map struct {
	Towers      map[*Tower]struct{}
	Enemies     map[*Enemy]struct{}
	Projectiles map[*Projectile]struct{}
	Path        Path
	Image       *ebiten.Image
}

func (m *Map) Draw(screen *ebiten.Image) {
	screen.DrawImage(m.Image, nil)
	m.Path.Draw(screen)
	for p := range m.Projectiles {
		p = p
		// todo: draw
	}

	for t := range m.Towers {
		t = t
		// todo: draw
	}

	for e := range m.Enemies {
		e.Draw(screen)
	}
}

// Path is a struct that represents a path.
type Path []Point

func (p Path) Draw(screen *ebiten.Image) {
	for i := 0; i < len(p)-1; i++ {
		drawLine(screen, p[i], p[i+1])
	}
}

func drawLine(screen *ebiten.Image, p1, p2 Point) {
	x1, y1, x2, y2 := p1.X, p1.Y, p2.X, p2.Y
	vector.DrawFilledCircle(screen, x2, y2, PathWidth/2, color.RGBA{12, 23, 34, 255}, false)
	vector.StrokeLine(screen, x1, y1, x2, y2, PathWidth, color.RGBA{12, 23, 34, 255}, false)
}
