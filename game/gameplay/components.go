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

type SpriteData struct {
	sprite *ebiten.Image
}

var Position = donburi.NewComponentType[PositionData]()
var Velocity = donburi.NewComponentType[VelocityData]()
var Sprite = donburi.NewComponentType[SpriteData]()
