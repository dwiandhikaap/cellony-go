package gameplay

import (
	"cellony/game/config"

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
	GlobalCamera = camera.NewCamera(int(config.Video.Width), int(config.Video.Height), config.Video.Width/2, config.Video.Height/2, 0, 1)

	// Destroy any existing camera entities
	query := donburi.NewQuery(
		filter.Contains(CameraComponent),
	)

	query.Each(world, func(entry *donburi.Entry) {
		world.Remove(entry.Entity())
	})

	cam := world.Create(CameraComponent)
	cameraData := CameraComponent.Get(world.Entry(cam))

	cameraData.cam = GlobalCamera

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
		multiplier := config.Control.CamSpeed

		_, scrollAmount := ebiten.Wheel()
		if scrollAmount > 0 {
			cam.Zoom(1.1)
		} else if scrollAmount < 0 {
			cam.Zoom(0.9)
		}

		// cam panning like dota
		threshold := 40
		cx, cy := ebiten.CursorPosition()
		if cx < threshold || cx > int(config.Video.Width)-threshold || cy < threshold || cy > int(config.Video.Height)-threshold {
			if ebiten.IsKeyPressed(ebiten.KeyShiftLeft) {
				multiplier *= 3
			}

			dx := (float64(cx) - (config.Video.Width / 2)) / config.Video.Width * multiplier
			dy := (float64(cy) - (config.Video.Height / 2)) / config.Video.Height * multiplier

			// move camera
			cam.MovePosition(dx, dy)
		}
	})
}
