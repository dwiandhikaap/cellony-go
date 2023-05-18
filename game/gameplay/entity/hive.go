package ent

import (
	"math/rand"

	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/filter"

	"cellony/game/config"
	comp "cellony/game/gameplay/component"
	"cellony/game/graphics"
	"cellony/game/util"
)

func CreateHiveEntity(world donburi.World) donburi.Entity {
	hive := world.Create(comp.Position, comp.Vertices, comp.Indices, comp.Color, comp.Hive)
	hiveEntry := world.Entry(hive)

	radius := 64.0

	comp.Hive.Get(hiveEntry).SpawnCooldown = 1
	comp.Hive.Get(hiveEntry).SpawnCountdown = 0
	comp.Hive.Get(hiveEntry).SpawnCount = 30

	x := rand.Float64() * float64(config.Game.Width)
	y := rand.Float64() * float64(config.Game.Height)

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
		dirtyMask := comp.Grid.Get(entry).DirtyMask

		// Outer circle, reduce by 0.1 each steps
		for i := 0; i < 15; i++ {
			r := radius * (1.5 + float64(i)*0.1)
			indices := util.GetCircleLatticeArea(x/config.Game.TileSize, y/config.Game.TileSize, r/config.Game.TileSize)
			delta := 1 / 15.0
			for _, index := range indices {
				xIndex := int(index[0])
				yIndex := int(index[1])

				grid[xIndex][yIndex] = float32(util.Clamp(float64(grid[xIndex][yIndex])-delta, 0.0, 1.0))
				dirtyMask[xIndex][yIndex] = true
			}
		}
	})

	return hive
}
