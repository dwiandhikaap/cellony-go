package system

import (
	comp "cellony/game/gameplay/component"

	"github.com/hajimehoshi/ebiten/v2"
	camera "github.com/melonfunction/ebiten-camera"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/ecs"
	"github.com/yohamta/donburi/filter"
)

func BackgroundSpriteRenderer(ecs *ecs.ECS, cam *camera.Camera) {
	query := donburi.NewQuery(
		filter.And(
			filter.Contains(comp.Sprite),
			filter.Contains(comp.Position),
		),
	)

	// get unique z-indexes
	zIndexes := make(map[int32]bool)

	query.Each(ecs.World, func(entry *donburi.Entry) {
		sprite := comp.Sprite.Get(entry)
		zIndexes[sprite.Z] = true
	})

	for zIndex := int32(-1); zIndex < 0; zIndex++ {
		query.Each(ecs.World, func(entry *donburi.Entry) {
			sprite := comp.Sprite.Get(entry)

			if sprite.Z != zIndex {
				return
			}

			position := comp.Position.Get(entry)
			screen := cam.Surface

			w := screen.Bounds().Dx()
			h := screen.Bounds().Dy()

			// Ass looking entity culling algorithm
			if !(position.X > (cam.X-float64(w))-float64(cam.Width)/cam.Scale/2 &&
				position.X < (cam.X+float64(w))+float64(cam.Width)/cam.Scale/2 &&
				position.Y > (cam.Y-float64(h))-float64(cam.Height)/cam.Scale/2 &&
				position.Y < (cam.Y+float64(h))+float64(cam.Height)/cam.Scale/2) {
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

func ForegroundSpriteRenderer(ecs *ecs.ECS, cam *camera.Camera) {
	query := donburi.NewQuery(
		filter.And(
			filter.Contains(comp.Sprite),
			filter.Contains(comp.Position),
		),
	)

	for zIndex := int32(0); zIndex < 2; zIndex++ {
		query.Each(ecs.World, func(entry *donburi.Entry) {
			sprite := comp.Sprite.Get(entry)

			if sprite.Z != zIndex {
				return
			}

			position := comp.Position.Get(entry)
			screen := cam.Surface

			w := screen.Bounds().Dx()
			h := screen.Bounds().Dy()

			// Ass looking entity culling algorithm
			if !(position.X > (cam.X-float64(w))-float64(cam.Width)/cam.Scale/2 &&
				position.X < (cam.X+float64(w))+float64(cam.Width)/cam.Scale/2 &&
				position.Y > (cam.Y-float64(h))-float64(cam.Height)/cam.Scale/2 &&
				position.Y < (cam.Y+float64(h))+float64(cam.Height)/cam.Scale/2) {
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
