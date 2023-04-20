package gameplay

import "github.com/yohamta/donburi"

type PositionData struct {
	x, y float64
}

type VelocityData struct {
	x, y float64
}

var Position = donburi.NewComponentType[PositionData]()
var Velocity = donburi.NewComponentType[VelocityData]()
