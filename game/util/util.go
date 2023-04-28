package util

import (
	"math"
	"runtime"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/jaypipes/ghw"
	"github.com/nfnt/resize"
)

func Clamp(a float64, min float64, max float64) float64 {
	if a < min {
		return min
	}
	if a > max {
		return max
	}
	return a
}

func GpuInfo() (gpu string) {
	gpu = "Unknown GPU"
	info, err := ghw.GPU()
	if err != nil {
		return
	}

	for _, gc := range info.GraphicsCards {
		if gc.DeviceInfo != nil {
			return gc.DeviceInfo.Product.Name
		}
	}
	return
}

func GetGraphicsLibrary() ebiten.GraphicsLibrary {
	if runtime.GOOS == "windows" {
		return ebiten.GraphicsLibraryOpenGL
	}

	return ebiten.GraphicsLibraryAuto
}

func RangeInterpolate(a float64, aMin float64, aMax float64, bMin float64, bMax float64) float64 {
	return bMin + (a-aMin)*(bMax-bMin)/(aMax-aMin)
}

func Distance(aX float64, aY float64, bX float64, bY float64) float64 {
	return math.Sqrt((aX-bX)*(aX-bX) + (aY-bY)*(aY-bY))
}

func GetCircleLatticeArea(x float64, y float64, radius float64) [][]int {
	lattice := make([][]int, 0)

	for i := int(x - radius); i <= int(x+radius); i++ {
		for j := int(y - radius); j <= int(y+radius); j++ {
			if Distance(x, y, float64(i), float64(j)) <= radius {
				lattice = append(lattice, []int{i, j})
			}
		}
	}

	return lattice
}

func ResizeImage(img *ebiten.Image, w int, h int) *ebiten.Image {
	image := resize.Resize(uint(w), uint(h), img.SubImage(img.Bounds()), resize.NearestNeighbor)
	return ebiten.NewImageFromImage(image)
}

func FilePathToName(path string) string {
	name := path
	for i := len(path) - 1; i >= 0; i-- {
		if path[i] == '/' {
			name = path[i+1:]
			break
		}
	}

	// Remove extension
	for i := len(name) - 1; i >= 0; i-- {
		if name[i] == '.' {
			name = name[:i]
			break
		}
	}

	return name
}
