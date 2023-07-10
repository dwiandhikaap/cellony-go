package system

import (
	comp "autocell/game/gameplay/component"
	"math"

	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/ecs"
	"github.com/yohamta/donburi/filter"
)

func GlobalStateSystem(ecs *ecs.ECS) {
	query := donburi.NewQuery(
		filter.Contains(comp.GlobalState),
	)

	globalState, ok := query.First(ecs.World)
	if !ok {
		return
	}

	globalStateData := comp.GlobalState.Get(globalState)

	if !globalStateData.IsValueDirty {
		return
	}

	setCellSpeed(ecs, globalState, float64(globalStateData.CellSpeed))
	setHiveSpawnCount(ecs, globalState, globalStateData.SpawnCount)
	setHiveCooldown(ecs, globalState, globalStateData.SpawnCooldown)

	globalStateData.IsValueDirty = false
}

func setCellSpeed(ecs *ecs.ECS, cell *donburi.Entry, speed float64) {
	// get all cell
	query := donburi.NewQuery(
		filter.Contains(comp.Cell),
	)

	query.Each(ecs.World, func(entry *donburi.Entry) {
		vel := comp.Velocity.Get(entry)
		mag := math.Sqrt(vel.X*vel.X + vel.Y*vel.Y)
		vel.X = vel.X / mag * speed
		vel.Y = vel.Y / mag * speed

		cellData := comp.Speed.Get(entry)
		cellData.Speed = speed
	})
}

func setHiveSpawnCount(ecs *ecs.ECS, cell *donburi.Entry, count int) {
	// get all cell
	query := donburi.NewQuery(
		filter.Contains(comp.Hive),
	)

	query.Each(ecs.World, func(entry *donburi.Entry) {
		hiveData := comp.Hive.Get(entry)
		hiveData.SpawnCount = count
	})
}

func setHiveCooldown(ecs *ecs.ECS, cell *donburi.Entry, cooldown float64) {
	// get all cell
	query := donburi.NewQuery(
		filter.Contains(comp.Hive),
	)

	query.Each(ecs.World, func(entry *donburi.Entry) {
		hiveData := comp.Hive.Get(entry)
		hiveData.SpawnCooldown = cooldown
	})
}
