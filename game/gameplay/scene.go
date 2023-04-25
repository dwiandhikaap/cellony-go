package gameplay

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/yohamta/donburi"

	ecslib "github.com/yohamta/donburi/ecs"
)

type WorldScene struct {
	ecs *ecslib.ECS
}

func CreateWorldScene() *WorldScene {
	world := donburi.NewWorld()

	s := WorldScene{
		ecs: ecslib.NewECS(world),
	}

	createCameraEntity(s.ecs.World)
	createMapEntity(s.ecs.World)
	createHiveEntity(s.ecs.World)
	createHiveEntity(s.ecs.World)
	for i := 0; i < 50000; i++ {
		createCellEntity(s.ecs.World)
	}

	addSystem(s.ecs)

	return &s
}

func (s *WorldScene) Update() error {
	s.ecs.Update()
	return nil
}

func (s *WorldScene) Draw(screen *ebiten.Image) {
	s.ecs.Draw(screen)
}

func (s *WorldScene) Open() {
	s.ecs.Resume()
}

func (s *WorldScene) Close() {
	s.ecs.Pause()
}
