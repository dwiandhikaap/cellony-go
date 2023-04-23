package gameplay

import (
	"math"
	"math/rand"

	"github.com/yohamta/donburi"

	"cellony/game/assets"
	"cellony/game/config"
	"cellony/game/graphics"
)

func createCellEntity(world donburi.World) donburi.Entity {
	cell := world.Create(Position, Velocity, Speed, Sprite)
	cellEntry := world.Entry(cell)

	Position.Get(cellEntry).x = rand.Float64() * float64(config.Video.Width)
	Position.Get(cellEntry).y = rand.Float64() * float64(config.Video.Height)

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
