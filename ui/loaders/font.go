package loaders

import (
	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font"
	"golang.org/x/image/font/gofont/goregular"
)

// FontTrueType loads a font from the gofont package.
func FontTrueType(size float64) font.Face {
	ttfFont, err := truetype.Parse(goregular.TTF)
	if err != nil {
		panic(err)
	}

	return truetype.NewFace(ttfFont, &truetype.Options{
		Size:    size,
		DPI:     72,
		Hinting: font.HintingFull,
	})
}
