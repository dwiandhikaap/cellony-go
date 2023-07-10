package system

import (
	comp "autocell/game/gameplay/component"

	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/ecs"
	"github.com/yohamta/donburi/filter"
)

type CellAI func(ecs *ecs.ECS, grid *comp.GridData, cellEntry *donburi.Entry)

func BaseAI(ecs *ecs.ECS) {
	cellQuery := donburi.NewQuery(
		filter.Contains(comp.Cell),
	)

	worldQuery := donburi.NewQuery(
		filter.Contains(comp.Grid),
	)

	worldEntry, ok := worldQuery.FirstEntity(ecs.World)

	if !ok {
		return
	}

	grid := comp.Grid.Get(worldEntry)

	cellQuery.Each(ecs.World, func(entry *donburi.Entry) {
		activityAI(ecs, grid, entry)

		if comp.CellActivity.Get(entry).Activity == comp.Wandering {
			wanderingAI(ecs, grid, entry)
		} else if comp.CellActivity.Get(entry).Activity == comp.Fleeing {
			fleeingAI(ecs, grid, entry)
		} else if comp.CellActivity.Get(entry).Activity == comp.Attacking {
			attackingAI(ecs, grid, entry)
		} else if comp.CellActivity.Get(entry).Activity == comp.Mining {
			miningAI(ecs, grid, entry)
		} else if comp.CellActivity.Get(entry).Activity == comp.Delivering {
			deliveringAI(ecs, grid, entry)
		}
	})
}
