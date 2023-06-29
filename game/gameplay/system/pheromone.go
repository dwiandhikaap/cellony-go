package system

import (
	comp "cellony/game/gameplay/component"

	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/ecs"
	"github.com/yohamta/donburi/filter"
)

func PheromoneSystem(ecs *ecs.ECS) {
	query := donburi.NewQuery(
		filter.Contains(comp.Pheromone),
	)

	query.Each(ecs.World, func(entry *donburi.Entry) {
		pheromone := comp.Pheromone.Get(entry)

		if pheromone.Intensity < 0 {
			ecs.World.Remove(entry.Entity())
			return
		}

		pheromone.Intensity -= 0.001

		sprite := comp.Sprite.Get(entry)
		sprite.Opacity = 0.3 * (pheromone.Intensity / pheromone.MaxIntensity)
	})
}
