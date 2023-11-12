package models

import "image"

// Map is a struct that represents a map.
type Map struct {
	Towers      map[*Tower]struct{}
	Enemies     map[*Enemy]struct{}
	Projectiles map[*Projectile]struct{}
	Path        Path
	Image       image.Image
}

// Path is a struct that represents a path.
type Path []Point
