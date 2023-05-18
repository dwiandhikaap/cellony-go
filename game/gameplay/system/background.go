package system

import (
	"image/color"

	camera "github.com/melonfunction/ebiten-camera"
	"github.com/yohamta/donburi/ecs"
)

func BackgroundRenderer(ecs *ecs.ECS, cam *camera.Camera) {
	cam.Surface.Fill(color.RGBA{0x00, 0x01, 0x00, 255})
}
