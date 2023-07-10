package system

import (
	"autocell/game/config"
	comp "autocell/game/gameplay/component"
	ent "autocell/game/gameplay/entity"
	"autocell/game/util"
	"math"
	"math/rand"

	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/ecs"
)

func wanderingAI(ecs *ecs.ECS, grid *comp.GridData, cellEntry *donburi.Entry) {
	pos := comp.Position.Get(cellEntry)
	vel := comp.Velocity.Get(cellEntry)
	cell := comp.Cell.Get(cellEntry)
	class := cell.Class

	newVelX := vel.X
	newVelY := vel.Y
	comp.Sprite.Get(cellEntry).Opacity = 0.5

	var (
		randomTurningRate        float32 = 0.003
		randomPheromoneRate      float32 = 0.075
		randomPheromoneCheckRate float32 = 0.1 // was 0.6
	)

	var randVal float32 = rand.Float32()
	var isRandomTurning bool = randVal < randomTurningRate
	var isRandomPheromoneDrop bool = randVal < randomPheromoneRate
	var isRandomPheromoneCheck bool = randVal < randomPheromoneCheckRate
	var lastPheromoneX float64 = cell.LastPheromoneX
	var lastPheromoneY float64 = cell.LastPheromoneY

	if isRandomPheromoneCheck {
		_, newVelX, newVelY = calcPheromone(comp.Wandering, comp.Speed.Get(cellEntry).Speed, vel, pos, ecs, cellEntry)
	}

	if isRandomTurning {
		// Action: Randomly change direction, random between -15 deg to 15 deg
		deltaAngle := (rand.Float64() - 0.5) * math.Pi / 2

		// nudge newVel a bit by deltaAngle
		newVelX = math.Cos(deltaAngle)*newVelX - math.Sin(deltaAngle)*newVelY
		newVelY = math.Sin(deltaAngle)*newVelX + math.Cos(deltaAngle)*newVelY
	}

	// get pheromne intensity at current position

	if isRandomPheromoneDrop && class == comp.Wanderer && util.DistanceSquared(pos.X, pos.Y, lastPheromoneX, lastPheromoneY) > 400 {
		ent.CreatePheromoneEntity(ecs.World, &ent.CreatePheromoneOptions{
			X:         pos.X,
			Y:         pos.Y,
			Activity:  comp.Wandering,
			HiveID:    comp.Parent.Get(cellEntry).Id,
			Intensity: 1,
		})
		cell.LastPheromoneX = pos.X
		cell.LastPheromoneY = pos.Y
	}

	// CheckAction: Check if cell is going towards a wall
	nextX := int((pos.X + newVelX/60 + config.Game.TileSize/2) / config.Game.TileSize)
	nextY := int((pos.Y + newVelY/60 + config.Game.TileSize/2) / config.Game.TileSize)

	if nextX >= int(config.Game.Width/config.Game.TileSize)-1 ||
		nextX <= 0 ||
		nextY >= int(config.Game.Height/config.Game.TileSize)-1 ||
		nextY <= 0 ||
		grid.Grid[nextX][nextY] > 0 {
		vel.X = -vel.X
		vel.Y = -vel.Y
	} else {
		vel.X = newVelX
		vel.Y = newVelY
	}

	// PostAction: Move cell, update internal state
	pos.X += vel.X * 1 / 60
	pos.Y += vel.Y * 1 / 60
}

// TODO: THIS IS FUCKING SUS !!!!!!!!!!!
func calcPheromone(activity comp.Activity, speed float64, vel *comp.VelocityData, pos *comp.PositionData, ecs *ecs.ECS, cellEntry *donburi.Entry) (bool, float64, float64) {
	isActionDone := false
	newVelX := vel.X
	newVelY := vel.Y

	_ = math.Sqrt(vel.X*vel.X + vel.Y*vel.Y)
	normX, normY := util.GetNormalizedVector(newVelX, newVelY)

	sensorRadius := 20
	sensorSize := 8.0
	sensorAngle := 3.14159 * 1 / 6 // 45 deg

	// create 3 vector from normX, normY. 1 is rotated -30 deg, 1 is not rotated, 1 is rotated 30 deg
	vx1 := math.Cos(-sensorAngle)*normX - math.Sin(-sensorAngle)*normY
	vy1 := math.Sin(-sensorAngle)*normX + math.Cos(-sensorAngle)*normY
	vx2 := normX
	vy2 := normY
	vx3 := math.Cos(sensorAngle)*normX - math.Sin(sensorAngle)*normY
	vy3 := math.Sin(sensorAngle)*normX + math.Cos(sensorAngle)*normY

	// create 3 points from 3 vectors, with sensorRadius distance from pos, offset by sensorSize/2
	x1 := pos.X + vx1*float64(sensorRadius) - sensorSize/2
	y1 := pos.Y + vy1*float64(sensorRadius) - sensorSize/2
	x2 := pos.X + vx2*float64(sensorRadius) - sensorSize/2
	y2 := pos.Y + vy2*float64(sensorRadius) - sensorSize/2
	x3 := pos.X + vx3*float64(sensorRadius) - sensorSize/2
	y3 := pos.Y + vy3*float64(sensorRadius) - sensorSize/2

	intensity1 := 0.0
	intensity2 := 0.0
	intensity3 := 0.0

	pheromoneQuadTree.ForEach(x1, y1, sensorSize, sensorSize, func(x, y, w, h float64, otherEntity donburi.Entity) {
		otherEntry := ecs.World.Entry(otherEntity)
		phero := comp.Pheromone.Get(otherEntry)

		if phero.Activity == activity {
			intensity1 += phero.Intensity
		}
	})

	pheromoneQuadTree.ForEach(x2, y2, sensorSize, sensorSize, func(x, y, w, h float64, otherEntity donburi.Entity) {
		otherEntry := ecs.World.Entry(otherEntity)
		phero := comp.Pheromone.Get(otherEntry)

		if phero.Activity == activity {
			intensity2 += phero.Intensity
		}
	})

	pheromoneQuadTree.ForEach(x3, y3, sensorSize, sensorSize, func(x, y, w, h float64, otherEntity donburi.Entity) {
		otherEntry := ecs.World.Entry(otherEntity)
		phero := comp.Pheromone.Get(otherEntry)

		if phero.Activity == activity {
			intensity3 += phero.Intensity
		}
	})

	nextX := x2 + sensorSize/2
	nextY := y2 + sensorSize/2

	// if all 3 are 0, then just do nothing
	if intensity1 == 0 && intensity2 == 0 && intensity3 == 0 {
		return isActionDone, newVelX, newVelY
	}

	// find most intense pheromone
	if intensity1 > intensity2 && intensity1 > intensity3 {
		nextX = x1 + sensorSize/2
		nextY = y1 + sensorSize/2
	} else if intensity3 > intensity2 && intensity3 > intensity1 {
		nextX = x3 + sensorSize/2
		nextY = y3 + sensorSize/2
	}

	if math.IsNaN(nextX) || math.IsNaN(nextY) {
		return false, vel.X, vel.Y
	}

	nextNormX, nextNormY := util.GetNormalizedLine(pos.X, pos.Y, nextX, nextY)

	newVelX = nextNormX * speed
	newVelY = nextNormY * speed

	_ = newVelX
	_ = newVelY

	isActionDone = true

	return isActionDone, newVelX, newVelY
}
