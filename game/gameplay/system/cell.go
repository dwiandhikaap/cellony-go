package system

import (
	"cellony/game/config"
	comp "cellony/game/gameplay/component"
	"math"
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
	camera "github.com/melonfunction/ebiten-camera"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/ecs"
	"github.com/yohamta/donburi/filter"
)

func CellCollisionSystem(ecs *ecs.ECS) {
	cellQuery := donburi.NewQuery(
		filter.Contains(comp.Cell),
	)

	worldQuery := donburi.NewQuery(
		filter.Contains(comp.Grid),
	)

	worldQuery.Each(ecs.World, func(worldEntry *donburi.Entry) {
		grid := comp.Grid.Get(worldEntry)

		cellQuery.Each(ecs.World, func(cellEntry *donburi.Entry) {
			cellPosition := comp.Position.Get(cellEntry)

			x := int(cellPosition.X / config.Game.TileSize)
			y := int(cellPosition.Y / config.Game.TileSize)

			if cellPosition.X >= config.Game.Width ||
				cellPosition.X < 0 ||
				cellPosition.Y >= config.Game.Height ||
				cellPosition.Y < 0 {
				cellEntry.Remove()
			} else if grid.Grid[x][y] > 0 {
				cellEntry.Remove()

				grid.Grid[x][y] -= 0.1
				grid.DirtyMask[x][y] = true
			}
		})
	})
}

func CellMovementSystem(ecs *ecs.ECS) {
	query := donburi.NewQuery(
		filter.And(
			filter.Contains(comp.Position),
			filter.Contains(comp.Velocity),
			filter.Contains(comp.Speed),
		),
	)

	query.Each(ecs.World, func(entry *donburi.Entry) {
		// random chance to change direction
		if rand.Float32() < 0.01 {
			angle := rand.Float64() * 2 * 3.14159

			velocity := comp.Velocity.Get(entry)
			velocity.X = math.Cos(angle) * comp.Speed.Get(entry).Speed
			velocity.Y = math.Sin(angle) * comp.Speed.Get(entry).Speed
		}
	})
}

func CellSystem(ecs *ecs.ECS) {
	query := donburi.NewQuery(
		filter.And(
			filter.Contains(comp.Position),
			filter.Contains(comp.Velocity),
		),
	)

	query.Each(ecs.World, func(entry *donburi.Entry) {
		position := comp.Position.Get(entry)
		velocity := comp.Velocity.Get(entry)

		position.X += velocity.X * 1 / 60
		position.Y += velocity.Y * 1 / 60
	})
}

func CellRenderer(ecs *ecs.ECS, cam *camera.Camera) {
	query := donburi.NewQuery(
		filter.And(
			filter.Contains(comp.Position),
			filter.Contains(comp.Sprite),
		),
	)

	query.Each(ecs.World, func(entry *donburi.Entry) {
		sprite := comp.Sprite.Get(entry)
		position := comp.Position.Get(entry)
		screen := cam.Surface

		// Ass looking entity culling algorithm
		if !(position.X > (cam.X-4)-float64(cam.Width)/cam.Scale/2 &&
			position.X < (cam.X+4)+float64(cam.Width)/cam.Scale/2 &&
			position.Y > (cam.Y-4)-float64(cam.Height)/cam.Scale/2 &&
			position.Y < (cam.Y+4)+float64(cam.Height)/cam.Scale/2) {
			return
		}

		op := &ebiten.DrawImageOptions{}
		op = cam.GetTranslation(op, position.X, position.Y)
		screen.DrawImage(sprite.Sprite, op)
	})
}
