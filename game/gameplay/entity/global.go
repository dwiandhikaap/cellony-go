package ent

import (
	comp "autocell/game/gameplay/component"
	"autocell/game/scene"

	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/ecs"
)

func CreateGlobalStateEntity(ecs *ecs.ECS, sceneManager *scene.SceneManager) donburi.Entity {
	world := ecs.World
	globalState := world.Create(comp.GlobalState)

	globalStateEntry := world.Entry(globalState)
	globalStateData := comp.GlobalState.Get(globalStateEntry)

	globalStateData.BrushRadius = 12
	globalStateData.CurrentBrush = 0

	globalStateData.CellSpeed = 30

	return globalState
}
