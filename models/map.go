package models

import "image"

type Map struct {
	Towers      map[*Tower]struct{}
	Enemies     map[*Enemy]struct{}
	Projectiles map[*Projectile]struct{}
	Path        Path
	Image       image.Image
}

type Path []Point
