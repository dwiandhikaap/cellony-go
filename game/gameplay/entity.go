package gameplay

import (
	"math"
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/filter"

	"cellony/game/assets"
	"cellony/game/config"
	"cellony/game/graphics"
	"cellony/game/util"

	noise "github.com/ojrac/opensimplex-go"
)

func createCellEntity(world donburi.World) donburi.Entity {
	cell := world.Create(Position, Velocity, Speed, Sprite)
	cellEntry := world.Entry(cell)

	Position.Get(cellEntry).x = rand.Float64() * float64(config.Game.Width)
	Position.Get(cellEntry).y = rand.Float64() * float64(config.Game.Height)

	Speed.Get(cellEntry).speed = (rand.Float64() + 1) / 2 * 100

	angle := rand.Float64() * 2 * 3.14159
	Velocity.Get(cellEntry).x = math.Cos(angle) * Speed.Get(cellEntry).speed
	Velocity.Get(cellEntry).y = math.Sin(angle) * Speed.Get(cellEntry).speed

	Sprite.Get(cellEntry).sprite = assets.AssetsInstance.Sprites["circle64"]

	return cell
}

func createHiveEntity(world donburi.World) donburi.Entity {
	hive := world.Create(Position, Vertices, Indices)
	hiveEntry := world.Entry(hive)

	radius := 64.0

	x := rand.Float64() * float64(config.Game.Width)
	y := rand.Float64() * float64(config.Game.Height)

	// padding
	xPadding := 0.2 * config.Game.Width
	yPadding := 0.2 * config.Game.Height

	x = util.RangeInterpolate(x, 0.0, config.Game.Width, xPadding, float64(config.Game.Width)-xPadding)
	y = util.RangeInterpolate(y, 0.0, config.Game.Height, yPadding, float64(config.Game.Height)-yPadding)

	println(int(x), int(y))

	Position.Get(hiveEntry).x = x
	Position.Get(hiveEntry).y = y

	color := graphics.GenerateHiveColor()
	vs, is := graphics.GeneratePolygonVertices(float32(x), float32(y), color, radius, 8, 0.0)

	Vertices.Get(hiveEntry).vertices = vs
	Indices.Get(hiveEntry).indices = is

	// adjust map near hive
	mapQuery := donburi.NewQuery(
		filter.Contains(Grid),
	)

	mapQuery.Each(world, func(entry *donburi.Entry) {
		grid := Grid.Get(entry).grid
		dirtyMask := Grid.Get(entry).dirtyMask

		// Outer circle, reduce by 0.1 each steps
		for i := 0; i < 10; i++ {
			indices := util.GetCircleLatticeArea(x, y, radius*(2+float64(i)*0.1))
			for _, index := range indices {
				xIndex := int(index[0] / 10)
				yIndex := int(index[1] / 10)

				asd := float64(grid[xIndex][yIndex])
				grid[xIndex][yIndex] = float32(util.Clamp(asd-0.001, 0.0, 1.0))
				dirtyMask[xIndex][yIndex] = true
			}
		}
	})

	return hive
}

func createMapEntity(world donburi.World) {
	mapEntity := world.Create(Grid, Image)
	mapEntry := world.Entry(mapEntity)

	mapDownscale := 10.0

	mapWidth := int(config.Game.Width / mapDownscale)
	mapHeight := int(config.Game.Height / mapDownscale)

	grid := make([][]float32, mapWidth)
	dirtyMask := make([][]bool, mapWidth)

	Grid.Get(mapEntry).grid = grid
	Grid.Get(mapEntry).dirtyMask = dirtyMask
	Image.Get(mapEntry).img = ebiten.NewImage(int(config.Game.Width), int(config.Game.Height))

	n := noise.NewNormalized(1)

	for i := 0; i < mapWidth; i++ {
		grid[i] = make([]float32, mapHeight)
		dirtyMask[i] = make([]bool, mapHeight)
		for j := 0; j < mapHeight; j++ {
			val := float32(n.Eval2(float64(i)/mapDownscale, float64(j)/mapDownscale))
			if val > 0.45 {
				grid[i][j] = float32(util.RangeInterpolate(float64(val), 0.45, 1.0, 0.0, 1.0))
			} else {
				grid[i][j] = 0.0
			}
			dirtyMask[i][j] = true
		}
	}
}
