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

func deliveringAI(ecs *ecs.ECS, grid *comp.GridData, cellEntry *donburi.Entry) {
	pos := comp.Position.Get(cellEntry)
	vel := comp.Velocity.Get(cellEntry)
	cell := comp.Cell.Get(cellEntry)
	class := cell.Class

	parentHive := comp.Parent.Get(cellEntry).Id
	hiveEntry := ecs.World.Entry(parentHive)
	hivePos := comp.Position.Get(hiveEntry)

	var (
		randomTurningRate        float32 = 1.0 / 60 * (1.0 / 100)
		randomPheromoneRate      float32 = 0.075
		randomPheromoneCheckRate float32 = 0.6 // was 0.6
	)

	var randVal float32 = rand.Float32()
	var isRandomTurning bool = randVal < randomTurningRate
	var isRandomPheromoneDrop bool = randVal < randomPheromoneRate
	var isRandomPheromoneCheck bool = randVal < randomPheromoneCheckRate
	var lastPheromoneX float64 = cell.LastPheromoneX
	var lastPheromoneY float64 = cell.LastPheromoneY

	newVelX := vel.X
	newVelY := vel.Y

	// set sprite brighter
	comp.Sprite.Get(cellEntry).Opacity = 1

	// Action: Check if cell is facing a 'go home' pheromone
	if isRandomPheromoneCheck {
		_, newVelX, newVelY = calcPheromone(comp.Wandering, comp.Speed.Get(cellEntry).Speed, vel, pos, ecs, cellEntry)
	}

	// Action: Randomly change direction
	if isRandomTurning {
		deltaAngle := (rand.Float64() - 0.5) * math.Pi / 2

		// nudge direction a bit by deltaAngle
		newVelX = math.Cos(deltaAngle)*newVelX - math.Sin(deltaAngle)*newVelY
		newVelY = math.Sin(deltaAngle)*newVelX + math.Cos(deltaAngle)*newVelY
	}

	// If hive is near enough
	if util.DistanceSquared(pos.X, pos.Y, hivePos.X, hivePos.Y) < 10000 {
		dy := hivePos.Y - pos.Y
		dx := hivePos.X - pos.X

		mag := math.Sqrt(dx*dx + dy*dy)

		// set velocity to point towards hive
		newVelX = dx / mag * comp.Speed.Get(cellEntry).Speed
		newVelY = dy / mag * comp.Speed.Get(cellEntry).Speed
	}

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
