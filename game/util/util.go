package util

import (
	"math"
	"math/rand"
	"runtime"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/jaypipes/ghw"
	"github.com/nfnt/resize"
)

type Number interface {
	int | uint | int8 | uint8 | int16 | uint16 | int32 | uint32 | int64 | uint64 | float32 | float64
}

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

func DistanceSquared[T Number](aX T, aY T, bX T, bY T) T {
	return (aX-bX)*(aX-bX) + (aY-bY)*(aY-bY)
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

func ForEachSquareLatticeArea(x int, y int, radius int, f func(int, int)) {
	for i := int(x - radius); i <= int(x+radius); i++ {
		for j := int(y - radius); j <= int(y+radius); j++ {
			f(i, j)
		}
	}
}

func ForEachLatticeLine(x1 int, y1 int, x2 int, y2 int, f func(int, int)) {
	dx := x2 - x1
	dy := y2 - y1

	if dx == 0 {
		for j := y1; j <= y2; j++ {
			f(x1, j)
		}
		return
	}

	if dy == 0 {
		for i := x1; i <= x2; i++ {
			f(i, y1)
		}
		return
	}

	if math.Abs(float64(dx)) > math.Abs(float64(dy)) {
		for i := x1; i <= x2; i++ {
			j := int(float64(dy)/float64(dx)*float64(i-x1) + float64(y1))
			f(i, j)
		}
	} else {
		for j := y1; j <= y2; j++ {
			i := int(float64(dx)/float64(dy)*float64(j-y1) + float64(x1))
			f(i, j)
		}
	}
}

func GetNormalizedLine(x1 float64, y1 float64, x2 float64, y2 float64) (float64, float64) {
	dx := x2 - x1
	dy := y2 - y1
	mag := math.Sqrt(dx*dx + dy*dy)
	return dx / mag, dy / mag
}

func GetNormalizedVector(x float64, y float64) (float64, float64) {
	mag := math.Sqrt(x*x + y*y)
	return x / mag, y / mag
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

func Pick[T any](x []T) T {
	return x[rand.Intn(len(x))]
}
