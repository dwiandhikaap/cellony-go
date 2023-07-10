package system

import (
	"autocell/game/config"
	comp "autocell/game/gameplay/component"
	ent "autocell/game/gameplay/entity"
	"autocell/game/util"
	bitmask "autocell/lib/bit"
	"math/rand"

	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/ecs"
)

func miningAI(ecs *ecs.ECS, grid *comp.GridData, cellEntry *donburi.Entry) {
	pos := comp.Position.Get(cellEntry)
	vel := comp.Velocity.Get(cellEntry)
	cell := comp.Cell.Get(cellEntry)
	class := cell.Class

	newVelX := vel.X
	newVelY := vel.Y

	var (
		randomPheromoneRate float32 = 0.075
	)

	var randVal float32 = rand.Float32()
	var isRandomPheromoneDrop bool = randVal < randomPheromoneRate
	var isActionDone bool = false
	var lastPheromoneX float64 = cell.LastPheromoneX
	var lastPheromoneY float64 = cell.LastPheromoneY

	if isRandomPheromoneDrop && class == comp.Wanderer && util.DistanceSquared(pos.X, pos.Y, lastPheromoneX, lastPheromoneY) > 200 {
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
	nearestMarkedWallX := -1000
	nearestMarkedWallY := -1000
	nearestMarkedWallDistSq := 100000000

	if !isActionDone && class == comp.Wanderer {
		// Action: Check if cell near marked wall
		currX := int((pos.X + config.Game.TileSize/2) / config.Game.TileSize)
		currY := int((pos.Y + config.Game.TileSize/2) / config.Game.TileSize)

		util.ForEachSquareLatticeArea(currX, currY, 1, func(x, y int) {
			if grid.Grid[x][y] > 0 && !bitmask.HasBit(grid.Mask[x][y], comp.WallMask) {
				if x < 0 || y < 0 || x >= len(grid.Grid) || y >= len(grid.Grid[0]) {
					return
				}

				currNearestDistSq := (nearestMarkedWallX-currX)*(nearestMarkedWallX-currX) + (nearestMarkedWallY-currY)*(nearestMarkedWallY-currY)
				newNearestDistSq := (x-currX)*(x-currX) + (y-currY)*(y-currY)

				if newNearestDistSq < currNearestDistSq {
					nearestMarkedWallX = x
					nearestMarkedWallY = y
					nearestMarkedWallDistSq = newNearestDistSq
				}
			}
		})

		// Determine new velocity
		if nearestMarkedWallX != -1000 {
			//newVelX = float64(nearestMarkedWallX-currX) * comp.Speed.Get(cellEntry).Speed * 1.5
			//newVelY = float64(nearestMarkedWallY-currY) * comp.Speed.Get(cellEntry).Speed * 1.5

			//comp.CellActivity.Get(cellEntry).Activity = comp.Mining

			isActionDone = true
		}
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

	// Update wall health
	isWallNearEnough := nearestMarkedWallDistSq < 3
	if isActionDone && isWallNearEnough {
		grid.Grid[nearestMarkedWallX][nearestMarkedWallY] -= 0.003
		if grid.Grid[nearestMarkedWallX][nearestMarkedWallY] <= 0 {
			grid.Grid[nearestMarkedWallX][nearestMarkedWallY] = 0
			grid.Mask[nearestMarkedWallX][nearestMarkedWallY] = bitmask.ClearBit(grid.Mask[nearestMarkedWallX][nearestMarkedWallY], comp.MarkedMask)
		}
		grid.Mask[nearestMarkedWallX][nearestMarkedWallY] = bitmask.SetBit(grid.Mask[nearestMarkedWallX][nearestMarkedWallY], comp.DirtyMask)

		// If grid mask is food, add food to cell
		if !bitmask.HasBit(grid.Mask[nearestMarkedWallX][nearestMarkedWallY], comp.WallMask) {
			comp.Cell.Get(cellEntry).ResourceCarried += 0.02
		}
	}

	pos.X += vel.X * 1 / 60
	pos.Y += vel.Y * 1 / 60
}
