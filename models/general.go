package models

import (
	"image"
)

type Coord = float32

type Point struct {
	X Coord
	Y Coord
}

type TypeAttack int

const (
	Common = TypeAttack(iota)
)

type Drawable[T image.Image] interface {
	Draw(image T)
}
