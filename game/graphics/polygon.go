package graphics

import (
	"image/color"
	"math"

	"github.com/PerformLine/go-stockutil/colorutil"
	"github.com/hajimehoshi/ebiten/v2"
)

func GeneratePolygonVertices(cx float32, cy float32, color color.Color, radius float64, sides int, rotation float64) ([]ebiten.Vertex, []uint16) {
	vs := []ebiten.Vertex{}
	is := []uint16{}

	r, g, b, _ := color.RGBA()
	var (
		r1 = (r + 1) / 256
		g1 = (g + 1) / 256
		b1 = (b + 1) / 256
	)

	h, s, l := colorutil.RgbToHsl(float64(r1), float64(g1), float64(b1))
	r2, g2, b2 := colorutil.HslToRgb(h, s, l*0.75)

	vs = append(vs, ebiten.Vertex{
		DstX:   cx,
		DstY:   cy,
		SrcX:   0,
		SrcY:   0,
		ColorR: (float32(r1) + 1) / 256,
		ColorG: (float32(g1) + 1) / 256,
		ColorB: (float32(b1) + 1) / 256,
		ColorA: 1,
	})

	for i := 0; i < sides; i++ {
		theta := float64(i) * 2 * math.Pi / float64(sides)
		x := cx + float32(radius*math.Cos(theta+rotation))
		y := cy + float32(radius*math.Sin(theta+rotation))
		vs = append(vs, ebiten.Vertex{
			DstX:   float32(x),
			DstY:   float32(y),
			SrcX:   0,
			SrcY:   0,
			ColorR: float32(r2) / 255,
			ColorG: float32(g2) / 255,
			ColorB: float32(b2) / 255,
			ColorA: 1,
		})
	}

	for i := 0; i < sides; i++ {
		is = append(is, 0)
		is = append(is, uint16(i+1))
		is = append(is, uint16((i+1)%sides+1))
	}

	return vs, is
}
