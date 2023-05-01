package gameplay

import (
	"cellony/game/gameplay/camera"
	comp "cellony/game/gameplay/component"
	ent "cellony/game/gameplay/entity"

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

	playerCam := camera.CreateCameraEntity(s.ecs.World)
	ent.CreateMapEntity(s.ecs.World)
	playerHive := ent.CreateHiveEntity(s.ecs.World)
	ent.CreateHiveEntity(s.ecs.World)

	playerHiveEntry := world.Entry(playerHive)
	playerCamEntry := world.Entry(playerCam)

	playerHivePos := comp.Position.Get(playerHiveEntry)
	playerCamData := camera.CameraComponent.Get(playerCamEntry)
	playerCamData.Cam.X = playerHivePos.X
	playerCamData.Cam.Y = playerHivePos.Y

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
