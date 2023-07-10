package system

import (
	"autocell/game/config"
	comp "autocell/game/gameplay/component"

	"github.com/s0rg/quadtree"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/ecs"
	"github.com/yohamta/donburi/filter"
)

var pheromoneQuadTree *quadtree.Tree[donburi.Entity]

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

func PheromoneQTreeSystem(ecs *ecs.ECS) {
	pheromoneQuadTree = quadtree.New[donburi.Entity](config.Game.Width, config.Game.Height, 4)
	//cellQuadTree.Del(config.Game.Width, config.Game.Height)

	query := donburi.NewQuery(
		filter.Contains(comp.Pheromone),
	)

	entCount := 0
	query.Each(ecs.World, func(entry *donburi.Entry) {
		position := comp.Position.Get(entry)
		ok := pheromoneQuadTree.Add(position.X, position.Y, 0, 0, entry.Entity())

		if !ok {
			println("failed to add to tree")
			entry.Remove()
		}

		entCount++
	})

	//println("Phereomone count: ", entCount)
}
