package graphics

import (
	"image/color"
	"math/rand"
)

// Vibrant color that is not too bright or too dark.
func GenerateHiveColor() color.Color {
	r := rand.Float64()*0.5 + 0.25
	g := rand.Float64()*0.5 + 0.25
	b := rand.Float64()*0.5 + 0.25

	return color.RGBA{
		R: uint8(r * 255),
		G: uint8(g * 255),
		B: uint8(b * 255),
		A: 255,
	}
}
