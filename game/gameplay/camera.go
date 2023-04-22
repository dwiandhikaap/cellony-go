package gameplay

import (
	"github.com/hajimehoshi/ebiten/v2"
	camera "github.com/melonfunction/ebiten-camera"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/ecs"
	"github.com/yohamta/donburi/filter"
)

// Separate file for camera ECS

type CameraData struct {
	cam *camera.Camera
}

var CameraComponent = donburi.NewComponentType[CameraData]()

var GlobalCamera *camera.Camera
var CameraCallStack = []func(*ecs.ECS, *camera.Camera){}

func createCameraEntity(world donburi.World) donburi.Entity {
	GlobalCamera = camera.NewCamera(1920, 1280, 0, 0, 0, 0.5)

	// Destroy any existing camera entities
	query := donburi.NewQuery(
		filter.Contains(CameraComponent),
	)

	query.Each(world, func(entry *donburi.Entry) {
		world.Remove(entry.Entity())
	})

	cam := world.Create(CameraComponent)
	cameraEntry := world.Entry(cam)

	CameraComponent.Get(cameraEntry).cam = GlobalCamera

	return cam
}

func addCameraRenderer(renderer func(*ecs.ECS, *camera.Camera)) {
	CameraCallStack = append(CameraCallStack, renderer)
}

func cameraRenderer(ecs *ecs.ECS, screen *ebiten.Image) {
	query := donburi.NewQuery(
		filter.Contains(CameraComponent),
	)

	query.Each(ecs.World, func(entry *donburi.Entry) {
		cam := CameraComponent.Get(entry).cam

		cam.Surface.Clear()

		for _, renderer := range CameraCallStack {
			renderer(ecs, cam)
		}

		cam.Blit(screen)
	})
}

func cameraSystem(ecs *ecs.ECS) {
	query := donburi.NewQuery(
		filter.Contains(CameraComponent),
	)

	query.Each(ecs.World, func(entry *donburi.Entry) {
		cam := CameraComponent.Get(entry).cam

		_, scrollAmount := ebiten.Wheel()
		if scrollAmount > 0 {
			cam.MovePosition(0, -10)
			cam.Zoom(1.1)
		} else if scrollAmount < 0 {
			cam.Zoom(0.9)
		}
	})
}
