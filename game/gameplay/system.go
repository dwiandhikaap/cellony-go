package gameplay

import (
	"autocell/game/gameplay/camera"
	"autocell/game/gameplay/system"

	"github.com/yohamta/donburi/ecs"
)

const (
	LayerBackground ecs.LayerID = iota
	LayerHUD
)

func addSystem(ecs *ecs.ECS) {
	ecs.AddSystem(system.GlobalStateSystem)
	ecs.AddSystem(camera.CameraSystem)
	ecs.AddSystem(system.DebugPheroSystem)
	ecs.AddSystem(system.BrushSystem)
	//ecs.AddSystem(system.PathNodeSpawningSystem)
	ecs.AddSystem(system.BaseAI)
	ecs.AddSystem(system.CellHealthSystem)
	//ecs.AddSystem(system.CellCollisionSystem)
	ecs.AddSystem(system.HiveSystem)
	ecs.AddSystem(system.MapSystem)
	ecs.AddSystem(system.PheromoneSystem)
	ecs.AddSystem(system.HUDSystem)
	ecs.AddSystem(system.CellQTreeSystem)
	ecs.AddSystem(system.PheromoneQTreeSystem)
	//ecs.AddSystem(system.MapDestroySystem)

	camera.AddCameraRenderer(system.BackgroundRenderer)
	camera.AddCameraRenderer(system.BackgroundSpriteRenderer)
	camera.AddCameraRenderer(system.MapRenderer)
	camera.AddCameraRenderer(system.ForegroundSpriteRenderer)
	//camera.AddCameraRenderer(system.CellRenderer)
	camera.AddCameraRenderer(system.HiveRenderer)
	camera.AddCameraRenderer(system.BrushRenderer)

	ecs.AddRenderer(LayerBackground, camera.CameraRenderer)
	ecs.AddRenderer(LayerHUD, system.HUDRenderer)
}
