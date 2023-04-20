package gameplay

import "github.com/yohamta/donburi"

func createCellEntity(world donburi.World) donburi.Entity {
	return world.Create(Position, Velocity)
}
