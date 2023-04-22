package gameplay

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/yohamta/donburi"

	camera "github.com/melonfunction/ebiten-camera"
	ecslib "github.com/yohamta/donburi/ecs"
)

var cam *camera.Camera

type WorldScene struct {
	ecs *ecslib.ECS
	cam *camera.Camera
}

func CreateWorldScene() *WorldScene {
	world := donburi.NewWorld()
	cam = camera.NewCamera(1280, 720, 0, 0, 0, 1)

	s := WorldScene{
		ecs: ecslib.NewECS(world),
		cam: cam,
	}

	createHiveEntity(s.ecs.World)
	for i := 0; i < 10; i++ {
		createCellEntity(s.ecs.World)
	}

	addSystem(s.ecs)

	return &s
}

func (s *WorldScene) Update() error {
	cameraUpdate()
	s.ecs.Update()
	return nil
}

func (s *WorldScene) Draw(screen *ebiten.Image) {
	s.cam.Surface.Clear()
	s.ecs.Draw(s.cam.Surface)
	s.cam.Blit(screen)
}

func (s *WorldScene) Open() {
	s.ecs.Resume()
}

func (s *WorldScene) Close() {
	s.ecs.Pause()
}

func cameraUpdate() {
	_, scrollAmount := ebiten.Wheel()
	if scrollAmount > 0 {
		cam.Zoom(1.1)
	} else if scrollAmount < 0 {
		cam.Zoom(0.9)
	}

	cam.MovePosition(1, 0)
	cam.Rotate(0.01)
}
