package models

import (
	"image"
)

// Coord is a type for coordinates.
type Coord = float32

// Frames is a type for frames.
type Frames = int

// Point is a struct that represents a point.
type Point struct {
	X Coord
	Y Coord
}

// TypeAttack is an enum that represents the type of attack.
type TypeAttack int

const (
	// Common is the type of attack that is common.
	Common = TypeAttack(iota)
)

// Drawable is an interface that represents a drawable object.
type Drawable[T image.Image] interface {
	Draw(image T)
}

// ImageConfigure is an interface that can initialize image
// from the temporary state of the entity (e.g. color.RGBA from color hex-string)
// and save the image to itself.
type ImageConfigure[T image.Image] interface {
	InitImage() error
	Image() T
}
