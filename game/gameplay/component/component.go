package comp

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/yohamta/donburi"
)

type PositionData struct {
	X, Y float64
}

type VelocityData struct {
	X, Y float64
}

type Activity uint8

const (
	Wandering Activity = iota
	Delivering
	Attacking
	Fleeing
	Mining
)

type CellActivityData struct {
	Activity Activity
}

type SpeedData struct {
	Speed float64
}

type SpriteData struct {
	Sprite *ebiten.Image
}

type VerticesData struct {
	Vertices []ebiten.Vertex
}

type IndicesData struct {
	Indices []uint16
}

type GridData struct {
	Grid      [][]float32 // doubles down as grid health
	DirtyMask [][]bool
	TypeMask  [][]bool // true if wall, false if food
}

type ImageData struct {
	Img *ebiten.Image
}

type ColorData struct {
	R uint8
	G uint8
	B uint8
}

type HiveData struct {
	SpawnCooldown  float64
	SpawnCount     int
	SpawnCountdown float64
}

type ParentData struct {
	Id donburi.Entity
}

// Tags
var Cell = donburi.NewTag()

// Components
var Position = donburi.NewComponentType[PositionData]()
var Velocity = donburi.NewComponentType[VelocityData]()
var CellActivity = donburi.NewComponentType[CellActivityData]()
var Speed = donburi.NewComponentType[SpeedData]()
var Sprite = donburi.NewComponentType[SpriteData]()
var Color = donburi.NewComponentType[ColorData]()
var Vertices = donburi.NewComponentType[VerticesData]()
var Indices = donburi.NewComponentType[IndicesData]()
var Grid = donburi.NewComponentType[GridData]()
var Image = donburi.NewComponentType[ImageData]()
var Hive = donburi.NewComponentType[HiveData]()
var Parent = donburi.NewComponentType[ParentData]()
