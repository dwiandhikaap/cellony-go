package ent

import (
	"math"
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/filter"

	"cellony/game/assets"
	"cellony/game/config"
	"cellony/game/gameplay/comp"
	"cellony/game/graphics"
	"cellony/game/util"

	noise "github.com/ojrac/opensimplex-go"
)

type CreateCellOptions struct {
	X     float64
	Y     float64
	Speed float64
}

func CreateCellEntity(world donburi.World, options *CreateCellOptions) donburi.Entity {
	cell := world.Create(comp.Cell, comp.Position, comp.Velocity, comp.Speed, comp.Sprite)
	cellEntry := world.Entry(cell)

	comp.Position.Get(cellEntry).X = options.X
	comp.Position.Get(cellEntry).Y = options.Y

	comp.Speed.Get(cellEntry).Speed = options.Speed

	angle := rand.Float64() * 2 * 3.14159
	comp.Velocity.Get(cellEntry).X = math.Cos(angle) * comp.Speed.Get(cellEntry).Speed
	comp.Velocity.Get(cellEntry).Y = math.Sin(angle) * comp.Speed.Get(cellEntry).Speed

	comp.Sprite.Get(cellEntry).Sprite = assets.AssetsInstance.Sprites["circle64"]

	return cell
}

func CreateHiveEntity(world donburi.World) donburi.Entity {
	hive := world.Create(comp.Position, comp.Vertices, comp.Indices, comp.Color, comp.Hive)
	hiveEntry := world.Entry(hive)

	radius := 64.0

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
	vs, is := graphics.GeneratePolygonVertices(float32(x), float32(y), color, radius, 8, 0.0)

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

func CreateMapEntity(world donburi.World) {
	mapEntity := world.Create(comp.Grid, comp.Image)
	mapEntry := world.Entry(mapEntity)

	tileSize := config.Game.TileSize

	mapWidth := int(config.Game.Width / tileSize)
	mapHeight := int(config.Game.Height / tileSize)

	grid := make([][]float32, mapWidth)
	dirtyMask := make([][]bool, mapWidth)

	comp.Grid.Get(mapEntry).Grid = grid
	comp.Grid.Get(mapEntry).DirtyMask = dirtyMask
	comp.Image.Get(mapEntry).Img = ebiten.NewImage(int(config.Game.Width), int(config.Game.Height))

	n := noise.NewNormalized(1)

	for i := 0; i < mapWidth; i++ {
		grid[i] = make([]float32, mapHeight)
		dirtyMask[i] = make([]bool, mapHeight)
		for j := 0; j < mapHeight; j++ {
			val := float32(n.Eval2(float64(i)/tileSize, float64(j)/tileSize))
			if val > 0.45 {
				grid[i][j] = float32(util.RangeInterpolate(float64(val), 0.45, 1.0, 0.0, 1.0))
			} else {
				grid[i][j] = 0.0
			}
			dirtyMask[i][j] = true
		}
	}
}
