package gameplay

import (
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/ecs"
	"github.com/yohamta/donburi/filter"
)

func addSystem(ecs *ecs.ECS) {
	ecs.AddSystem(cellSystem)
}

func cellSystem(ecs *ecs.ECS) {
	query := donburi.NewQuery(
		filter.And(
			filter.Contains(Position),
			filter.Contains(Velocity),
		),
	)

	query.Each(ecs.World, func(entry *donburi.Entry) {
		// print pos
		println(Position.Get(entry).x, Position.Get(entry).y)

		position := Position.Get(entry)
		velocity := Velocity.Get(entry)

		position.x += velocity.x
		position.y += velocity.y
	})
}
