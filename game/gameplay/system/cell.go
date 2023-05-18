package system

import (
	"cellony/game/config"
	comp "cellony/game/gameplay/component"
	ent "cellony/game/gameplay/entity"
	"math"
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
	camera "github.com/melonfunction/ebiten-camera"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/ecs"
	"github.com/yohamta/donburi/filter"
)

func CellAISystem(ecs *ecs.ECS) {
	cellQuery := donburi.NewQuery(
		filter.Contains(comp.Cell),
	)

	worldQuery := donburi.NewQuery(
		filter.Contains(comp.Grid),
	)

	worldQuery.Each(ecs.World, func(worldEntry *donburi.Entry) {
		grid := comp.Grid.Get(worldEntry)

		cellQuery.Each(ecs.World, func(cellEntry *donburi.Entry) {
			pos := comp.Position.Get(cellEntry)
			vel := comp.Velocity.Get(cellEntry)

			activity := comp.CellActivity.Get(cellEntry).Activity

			// Randomly change direction
			if rand.Float32() < 0.01 {
				angle := rand.Float64() * 2 * 3.14159

				velocity := comp.Velocity.Get(cellEntry)
				velocity.X = math.Cos(angle) * comp.Speed.Get(cellEntry).Speed
				velocity.Y = math.Sin(angle) * comp.Speed.Get(cellEntry).Speed
			}

			nextX := int((pos.X + vel.X/30 + config.Game.TileSize/2) / config.Game.TileSize)
			nextY := int((pos.Y + vel.Y/30 + config.Game.TileSize/2) / config.Game.TileSize)

			// Check if cell is going towards a wall
			for nextX >= int(config.Game.Width/config.Game.TileSize) ||
				nextX < 0 ||
				nextY >= int(config.Game.Height/config.Game.TileSize) ||
				nextY < 0 ||
				grid.Grid[nextX][nextY] > 0 {
				// If cell is mining, stop
				if activity == comp.Mining {
					vel.X = 0
					vel.Y = 0
					break
				}
				angle := rand.Float64() * 2 * 3.14159

				vel.X = math.Cos(angle) * comp.Speed.Get(cellEntry).Speed
				vel.Y = math.Sin(angle) * comp.Speed.Get(cellEntry).Speed

				nextX = int((pos.X + vel.X/30 + config.Game.TileSize/2) / config.Game.TileSize)
				nextY = int((pos.Y + vel.Y/30 + config.Game.TileSize/2) / config.Game.TileSize)
			}

			// Set cell new position
			pos.X += vel.X * 1 / 60
			pos.Y += vel.Y * 1 / 60

			// Update internal cell clock
			cellData := comp.Cell.Get(cellEntry)
			if rand.Float64() < cellData.PheromoneChance {
				// check nearby pheromones
				pheromoneQuery := donburi.NewQuery(
					filter.Contains(comp.Pheromone),
				)

				nearbyIntensity := 0.0
				pheromoneQuery.Each(ecs.World, func(pheromoneEntry *donburi.Entry) {
					if nearbyIntensity > 2 {
						return
					}

					pheromone := comp.Pheromone.Get(pheromoneEntry)
					pheromonePos := comp.Position.Get(pheromoneEntry)

					if pheromone.HiveID == comp.Parent.Get(cellEntry).Id &&
						math.Abs(pheromonePos.X-pos.X) < 25 &&
						math.Abs(pheromonePos.Y-pos.Y) < 25 {
						nearbyIntensity += pheromone.Intensity
					}
				})

				if nearbyIntensity <= 2 {
					ent.CreatePheromoneEntity(ecs.World, &ent.CreatePheromoneOptions{
						X:         pos.X,
						Y:         pos.Y,
						Activity:  activity,
						HiveID:    comp.Parent.Get(cellEntry).Id,
						Intensity: 1,
					})
				}
			}
		})
	})
}

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

func CellRenderer(ecs *ecs.ECS, cam *camera.Camera) {
	query := donburi.NewQuery(
		filter.And(
			filter.Contains(comp.Cell),
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
