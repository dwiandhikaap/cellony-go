package gameplay

import (
	"cellony/game/gameplay/camera"
	"cellony/game/gameplay/system"

	"github.com/yohamta/donburi/ecs"
)

const (
	LayerBackground ecs.LayerID = iota
	LayerGame
)

func addSystem(ecs *ecs.ECS) {
	ecs.AddSystem(camera.CameraSystem)
	ecs.AddSystem(system.CellAISystem)
	//ecs.AddSystem(system.CellCollisionSystem)
	ecs.AddSystem(system.HiveSystem)
	ecs.AddSystem(system.MapSystem)
	ecs.AddSystem(system.PheromoneSystem)
	//ecs.AddSystem(system.MapDestroySystem)

	camera.AddCameraRenderer(system.MapRenderer)
	camera.AddCameraRenderer(system.CellRenderer)
	camera.AddCameraRenderer(system.HiveRenderer)

	ecs.AddRenderer(LayerBackground, camera.CameraRenderer)
}
