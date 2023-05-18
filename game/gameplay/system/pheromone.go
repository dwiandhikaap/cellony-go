package system

import (
	comp "cellony/game/gameplay/component"

	"github.com/hajimehoshi/ebiten/v2"
	camera "github.com/melonfunction/ebiten-camera"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/ecs"
	"github.com/yohamta/donburi/filter"
)

func PheromoneSystem(ecs *ecs.ECS) {
	query := donburi.NewQuery(
		filter.Contains(comp.Pheromone),
	)

	query.Each(ecs.World, func(entry *donburi.Entry) {
		pheromone := comp.Pheromone.Get(entry)

		if pheromone.Intensity < 0 {
			ecs.World.Remove(entry.Entity())
			return
		}

		pheromone.Intensity -= 0.001

		sprite := comp.Sprite.Get(entry)
		sprite.Opacity = 0.3 * (pheromone.Intensity / pheromone.MaxIntensity)
	})
}

func PheromoneRenderer(ecs *ecs.ECS, cam *camera.Camera) {
	query := donburi.NewQuery(
		filter.And(
			filter.Contains(comp.Pheromone),
		),
	)

	for zIndex := uint8(0); zIndex < 8; zIndex++ {
		query.Each(ecs.World, func(entry *donburi.Entry) {
			sprite := comp.Sprite.Get(entry)

			if sprite.Z != zIndex {
				return
			}

			position := comp.Position.Get(entry)
			screen := cam.Surface

			// Ass looking entity culling algorithm
			if !(position.X > (cam.X-4)-float64(cam.Width)/cam.Scale/2 &&
				position.X < (cam.X+4)+float64(cam.Width)/cam.Scale/2 &&
				position.Y > (cam.Y-4)-float64(cam.Height)/cam.Scale/2 &&
				position.Y < (cam.Y+4)+float64(cam.Height)/cam.Scale/2) {
				return
			}

			scale := sprite.Scale
			opacity := sprite.Opacity
			spriteWidth := float64(sprite.Sprite.Bounds().Dx()) * scale
			spriteHeight := float64(sprite.Sprite.Bounds().Dy()) * scale

			op := &ebiten.DrawImageOptions{}
			op.GeoM.Scale(scale, scale)
			op.ColorScale.ScaleAlpha(float32(opacity))

			op = cam.GetTranslation(op, position.X-spriteWidth/2, position.Y-spriteHeight/2)
			screen.DrawImage(sprite.Sprite, op)
		})
	}
}
