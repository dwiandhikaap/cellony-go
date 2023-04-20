package gameplay

import (
	"math/rand"

	"github.com/yohamta/donburi"

	"cellony/game/assets"
)

func createCellEntity(world donburi.World) donburi.Entity {
	cell := world.Create(Position, Velocity, Sprite)
	cellEntry := world.Entry(cell)

	Position.Get(cellEntry).x = rand.Float64() * 640
	Position.Get(cellEntry).y = rand.Float64() * 480

	Velocity.Get(cellEntry).x = (rand.Float64() - 0.5) * 2 * 100
	Velocity.Get(cellEntry).y = (rand.Float64() - 0.5) * 2 * 100

	Sprite.Get(cellEntry).sprite = assets.AssetsInstance.Sprites["circle64"]

	return cell
}
