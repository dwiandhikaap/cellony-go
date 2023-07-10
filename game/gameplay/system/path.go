package system

/*
var lastNodeX = -1
var lastNodeY = -1

func PathNodeSpawningSystem(ecs *ecs.ECS) {
	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		cursorX, cursorY := camera.CursorWorldPosition()

		cx := int(cursorX)
		cy := int(cursorY)
		radius := 12.0

		isFirstNode := lastNodeX == -1 || lastNodeY == -1
		isFarEnough := util.DistanceSquared(lastNodeX, lastNodeY, cx, cy) > 100

		if !isFirstNode && !isFarEnough {
			return
		}

		mapQuery := donburi.NewQuery(
			filter.Contains(comp.Grid),
		)

		mapQuery.Each(ecs.World, func(entry *donburi.Entry) {
			grid := comp.Grid.Get(entry).Grid
			mask := comp.Grid.Get(entry).Mask

			// Outer circle, reduce by 0.1 each steps
			for i := 0; i < 15; i++ {
				r := radius * (1.5 + float64(i)*0.1)
				indices := util.GetCircleLatticeArea(cursorX/config.Game.TileSize, cursorY/config.Game.TileSize, r/config.Game.TileSize)
				for _, index := range indices {
					xIndex := int(index[0])
					yIndex := int(index[1])

					if xIndex < 0 || xIndex >= len(grid) || yIndex < 0 || yIndex >= len(grid[0]) {
						continue
					}

					if bitmask.HasBit(mask[xIndex][yIndex], comp.MarkedMask) {
						continue
					}

					if grid[xIndex][yIndex] <= 0 {
						continue
					}

					mask[xIndex][yIndex] = bitmask.SetBit(mask[xIndex][yIndex], comp.MarkedMask)
					mask[xIndex][yIndex] = bitmask.SetBit(mask[xIndex][yIndex], comp.DirtyMask)
				}
			}
		})

		lastNodeX, lastNodeY = cx, cy
		return
	} else if ebiten.IsMouseButtonPressed(ebiten.MouseButtonMiddle) {
		// Clear all marked map tile in radius
		cursorX, cursorY := camera.CursorWorldPosition()

		cx := int(cursorX)
		cy := int(cursorY)
		radius := 12.0

		mapQuery := donburi.NewQuery(
			filter.Contains(comp.Grid),
		)

		mapQuery.Each(ecs.World, func(entry *donburi.Entry) {
			grid := comp.Grid.Get(entry).Grid
			mask := comp.Grid.Get(entry).Mask

			// Outer circle, reduce by 0.1 each steps
			for i := 0; i < 15; i++ {
				r := radius * (1.5 + float64(i)*0.1)
				indices := util.GetCircleLatticeArea(cursorX/config.Game.TileSize, cursorY/config.Game.TileSize, r/config.Game.TileSize)
				for _, index := range indices {
					xIndex := int(index[0])
					yIndex := int(index[1])

					if xIndex < 0 || xIndex >= len(grid) || yIndex < 0 || yIndex >= len(grid[0]) {
						continue
					}

					if !bitmask.HasBit(mask[xIndex][yIndex], comp.MarkedMask) {
						continue
					}

					mask[xIndex][yIndex] = bitmask.ClearBit(mask[xIndex][yIndex], comp.MarkedMask)
					mask[xIndex][yIndex] = bitmask.SetBit(mask[xIndex][yIndex], comp.DirtyMask)
				}
			}
		})

		lastNodeX, lastNodeY = cx, cy
		return

	}

	lastNodeX, lastNodeY = -1, -1
}
*/
