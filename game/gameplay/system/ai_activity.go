package system

import (
	"autocell/game/config"
	comp "autocell/game/gameplay/component"
	"autocell/game/util"
	bitmask "autocell/lib/bit"

	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/ecs"
)

func activityAI(ecs *ecs.ECS, grid *comp.GridData, cellEntry *donburi.Entry) {
	cell := comp.Cell.Get(cellEntry)
	pos := comp.Position.Get(cellEntry)
	class := cell.Class

	parentHive := comp.Parent.Get(cellEntry).Id
	parentEntry := ecs.World.Entry(parentHive)
	isPlayer := comp.Hive.Get(parentEntry).IsPlayer

	viewDist := 24.0

	if class == comp.Wanderer {
		if comp.CellActivity.Get(cellEntry).Activity != comp.Delivering {
			if cell.ResourceCarried >= 1 {
				comp.CellActivity.Get(cellEntry).Activity = comp.Delivering
				vel := comp.Velocity.Get(cellEntry)
				vel.X = -vel.X
				vel.Y = -vel.Y
			} else if isPlayer {
				currX := int((pos.X + config.Game.TileSize/2) / config.Game.TileSize)
				currY := int((pos.Y + config.Game.TileSize/2) / config.Game.TileSize)

				nearWall := false

				util.ForEachSquareLatticeArea(currX, currY, 2, func(x, y int) {
					if x < 0 || y < 0 || x >= len(grid.Grid) || y >= len(grid.Grid[0]) {
						return
					}

					if grid.Grid[x][y] > 0 && !bitmask.HasBit(grid.Mask[x][y], comp.WallMask) {
						nearWall = true
					}
				})

				// Determine new velocity
				if nearWall {
					comp.CellActivity.Get(cellEntry).Activity = comp.Mining
				} else {
					comp.CellActivity.Get(cellEntry).Activity = comp.Wandering
				}
			} else {
				comp.CellActivity.Get(cellEntry).Activity = comp.Wandering
			}
		}

		if comp.CellActivity.Get(cellEntry).Activity == comp.Delivering {
			// check if hive near enough
			hiveX := comp.Position.Get(parentEntry).X
			hiveY := comp.Position.Get(parentEntry).Y

			distSq := (hiveX-pos.X)*(hiveX-pos.X) + (hiveY-pos.Y)*(hiveY-pos.Y)
			if distSq < 64*64 {
				comp.CellActivity.Get(cellEntry).Activity = comp.Wandering

				// add resource to hive
				comp.Hive.Get(parentEntry).Resource += cell.ResourceCarried
				cell.ResourceCarried = 0

				//println("Delivered resource to hive. Current resource:", comp.Hive.Get(parentEntry).Resource)
			}

		}

	} else if class == comp.Soldier {
		enemySoldierCount := 0
		cellQuadTree.KNearest(pos.X, pos.Y, viewDist, 1, func(x, y, w, h float64, otherEntity donburi.Entity) {
			otherEntry := ecs.World.Entry(otherEntity)
			isFoundEnemy := comp.Parent.Get(otherEntry).Id != comp.Parent.Get(cellEntry).Id

			if isFoundEnemy {
				enemySoldierCount++
			}
		})

		if enemySoldierCount > 0 {
			comp.CellActivity.Get(cellEntry).Activity = comp.Attacking
		} else {
			comp.CellActivity.Get(cellEntry).Activity = comp.Wandering
		}
	}
}
