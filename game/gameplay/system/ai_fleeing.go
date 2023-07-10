package system

import (
	"autocell/game/config"
	comp "autocell/game/gameplay/component"
	"math"

	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/ecs"
)

func fleeingAI(ecs *ecs.ECS, grid *comp.GridData, cellEntry *donburi.Entry) {
	pos := comp.Position.Get(cellEntry)
	vel := comp.Velocity.Get(cellEntry)

	viewDist := 20.0

	// Action: Find nearest enemy cell position
	nearestEnemyX := -1000.0
	nearestEnemyY := -1000.0

	cellQuadTree.ForEach(pos.X-viewDist, pos.Y-viewDist, viewDist*2, viewDist*2, func(x, y, w, h float64, otherEntity donburi.Entity) {
		otherEntry := ecs.World.Entry(otherEntity)
		isFoundSoldier := comp.Cell.Get(otherEntry).Class == comp.Soldier
		isFoundEnemy := comp.Parent.Get(otherEntry).Id != comp.Parent.Get(cellEntry).Id

		if isFoundSoldier && isFoundEnemy {
			otherPos := comp.Position.Get(otherEntry)
			enemyX := otherPos.X
			enemyY := otherPos.Y

			currNearestDistSq := (nearestEnemyX-pos.X)*(nearestEnemyX-pos.X) + (nearestEnemyY-pos.Y)*(nearestEnemyY-pos.Y)
			newNearestDistSq := (enemyX-pos.X)*(enemyX-pos.X) + (enemyY-pos.Y)*(enemyY-pos.Y)

			if newNearestDistSq < currNearestDistSq {
				nearestEnemyX = enemyX
				nearestEnemyY = enemyY
			}
		}
	})

	// Determine new velocity
	// Calculate angle between cell and enemy, then move away from enemy
	angle := math.Atan2(pos.Y-nearestEnemyY, pos.X-nearestEnemyX)
	newVelX := math.Cos(angle) * comp.Speed.Get(cellEntry).Speed
	newVelY := math.Sin(angle) * comp.Speed.Get(cellEntry).Speed

	// CheckAction: Check if cell is going towards a wall
	nextX := int((pos.X + newVelX/60 + config.Game.TileSize/2) / config.Game.TileSize)
	nextY := int((pos.Y + newVelY/60 + config.Game.TileSize/2) / config.Game.TileSize)

	if nextX >= int(config.Game.Width/config.Game.TileSize) ||
		nextX < 0 ||
		nextY >= int(config.Game.Height/config.Game.TileSize) ||
		nextY < 0 ||
		grid.Grid[nextX][nextY] > 0 {
		vel.X = vel.Y
		vel.Y = -vel.X
	} else {
		vel.X = newVelX
		vel.Y = newVelY
	}

	// PostAction: Move cell, update internal state
	pos.X += vel.X * 1 / 60
	pos.Y += vel.Y * 1 / 60
}
