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

// Tags
var Hive = donburi.NewTag()
var Cell = donburi.NewTag()

// Components
var Position = donburi.NewComponentType[PositionData]()
var Velocity = donburi.NewComponentType[VelocityData]()
var Speed = donburi.NewComponentType[SpeedData]()
var Sprite = donburi.NewComponentType[SpriteData]()
var Vertices = donburi.NewComponentType[VerticesData]()
var Indices = donburi.NewComponentType[IndicesData]()
var Grid = donburi.NewComponentType[GridData]()
var Image = donburi.NewComponentType[ImageData]()
