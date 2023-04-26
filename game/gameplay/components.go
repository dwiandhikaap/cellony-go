package gameplay

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/yohamta/donburi"
)

type PositionData struct {
	x, y float64
}

type VelocityData struct {
	x, y float64
}

type SpeedData struct {
	speed float64
}

type SpriteData struct {
	sprite *ebiten.Image
}

type VerticesData struct {
	vertices []ebiten.Vertex
}

type IndicesData struct {
	indices []uint16
}

type GridData struct {
	grid      [][]float32 // doubles down as grid health
	dirtyMask [][]bool
}

type ImageData struct {
	img *ebiten.Image
}

type ColorData struct {
	r uint8
	g uint8
	b uint8
}

type HiveData struct {
	spawnCooldown  float64
	spawnCount     int
	spawnCountdown float64
}

// Tags
var Cell = donburi.NewTag()

// Components
var Position = donburi.NewComponentType[PositionData]()
var Velocity = donburi.NewComponentType[VelocityData]()
var Speed = donburi.NewComponentType[SpeedData]()
var Sprite = donburi.NewComponentType[SpriteData]()
var Color = donburi.NewComponentType[ColorData]()
var Vertices = donburi.NewComponentType[VerticesData]()
var Indices = donburi.NewComponentType[IndicesData]()
var Grid = donburi.NewComponentType[GridData]()
var Image = donburi.NewComponentType[ImageData]()
var Hive = donburi.NewComponentType[HiveData]()
