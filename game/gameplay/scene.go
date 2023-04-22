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

	op := &ebiten.DrawImageOptions{}
	op = s.cam.GetScale(op, cam.Scale, cam.Scale)
	op = s.cam.GetTranslation(op, 0, 0)

	screen.DrawImage(s.cam.Surface, op)
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

	x, y := ebiten.CursorPosition()
	cam.SetPosition(float64(x/100)+1280/2, float64(y/100)+360)
}
