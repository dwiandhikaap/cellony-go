package ent

import (
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/filter"

	"autocell/game/config"
	comp "autocell/game/gameplay/component"
	"autocell/game/graphics"
	"autocell/game/util"
)

func CreateHiveEntity(world donburi.World, x, y float64, isPlayer bool) donburi.Entity {
	hive := world.Create(comp.Position, comp.Vertices, comp.Indices, comp.Color, comp.Hive)
	hiveEntry := world.Entry(hive)

	radius := 64.0

	comp.Hive.Get(hiveEntry).SpawnCooldown = 5
	comp.Hive.Get(hiveEntry).SpawnCountdown = 0
	comp.Hive.Get(hiveEntry).SpawnCount = 10

	comp.Hive.Get(hiveEntry).IsPlayer = isPlayer
	comp.Hive.Get(hiveEntry).Resource = 0

	comp.Hive.Get(hiveEntry).WandererOdd = 1
	comp.Hive.Get(hiveEntry).SoldierOdd = 0
	comp.Hive.Get(hiveEntry).WorkerOdd = 0

	//x := rand.Float64() * float64(config.Game.Width)
	//y := rand.Float64() * float64(config.Game.Height)

	// padding
	xPadding := 0.2 * config.Game.Width
	yPadding := 0.2 * config.Game.Height

	x = util.RangeInterpolate(x, 0.0, config.Game.Width, xPadding, float64(config.Game.Width)-xPadding)
	y = util.RangeInterpolate(y, 0.0, config.Game.Height, yPadding, float64(config.Game.Height)-yPadding)

	comp.Position.Get(hiveEntry).X = x
	comp.Position.Get(hiveEntry).Y = y

	color := graphics.GenerateHiveColor()
	vs, is := graphics.GeneratePolygonVertices(float32(x), float32(y), color, radius, 16, 0.0)

	comp.Vertices.Get(hiveEntry).Vertices = vs
	comp.Indices.Get(hiveEntry).Indices = is

	r, g, b, _ := color.RGBA()
	comp.Color.Get(hiveEntry).R = uint8(r >> 8)
	comp.Color.Get(hiveEntry).G = uint8(g >> 8)
	comp.Color.Get(hiveEntry).B = uint8(b >> 8)

	// adjust map near hive
	mapQuery := donburi.NewQuery(
		filter.Contains(comp.Grid),
	)

	mapQuery.Each(world, func(entry *donburi.Entry) {
		grid := comp.Grid.Get(entry).Grid
		//mask := comp.Grid.Get(entry).Mask

		// Outer circle, reduce by 0.1 each steps
		for i := 0; i < 15; i++ {
			r := radius * (1.5 + float64(i)*0.1)
			indices := util.GetCircleLatticeArea(x/config.Game.TileSize, y/config.Game.TileSize, r/config.Game.TileSize)
			delta := 1 / 15.0
			for _, index := range indices {
				xIndex := int(index[0])
				yIndex := int(index[1])

				if xIndex < 0 || xIndex >= len(grid) || yIndex < 0 || yIndex >= len(grid[0]) {
					continue
				}

				grid[xIndex][yIndex] = float32(util.Clamp(float64(grid[xIndex][yIndex])-delta, 0.0, 1.0))
				//mask[xIndex][yIndex] = bitmask.SetBit(mask[xIndex][yIndex], comp.DirtyMask)
			}
		}
	})

	return hive
}
