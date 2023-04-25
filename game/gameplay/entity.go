package gameplay

import (
	"math"
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/yohamta/donburi"

	"cellony/game/assets"
	"cellony/game/config"
	"cellony/game/graphics"

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

	Position.Get(hiveEntry).x = config.Video.Width / 2
	Position.Get(hiveEntry).y = config.Video.Height / 2

	color := graphics.GenerateHiveColor()
	vs, is := graphics.GeneratePolygonVertices(float32(config.Video.Width/2), float32(config.Video.Height/2), color, 64.0, 8, 0.0)

	Vertices.Get(hiveEntry).vertices = vs
	Indices.Get(hiveEntry).indices = is

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
			if val > 0.5 {
				grid[i][j] = (val + 1) / 2
			} else {
				grid[i][j] = 0.0
			}
			dirtyMask[i][j] = true
		}
	}
}
