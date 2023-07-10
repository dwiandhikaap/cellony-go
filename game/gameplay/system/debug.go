package system

import (
	camera "autocell/game/gameplay/camera"
	comp "autocell/game/gameplay/component"
	ent "autocell/game/gameplay/entity"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/ecs"
	"github.com/yohamta/donburi/filter"
)

var (
	justPressed = false
)

func DebugPheroSystem(ecs *ecs.ECS) {
	if ebiten.IsKeyPressed(ebiten.KeySpace) {
		if justPressed && !ebiten.IsKeyPressed(ebiten.KeyShift) {
			return
		}

		activity := comp.Delivering

		cursorX, cursorY := camera.CursorWorldPosition()

		cx := (cursorX)
		cy := (cursorY)

		cellQuery := donburi.NewQuery(
			filter.Contains(comp.Cell),
		)

		cellEntry, ok := cellQuery.First(ecs.World)
		if !ok {
			print("no cell")
			return
		}

		ent.CreatePheromoneEntity(ecs.World, &ent.CreatePheromoneOptions{
			X:         cx,
			Y:         cy,
			Activity:  activity,
			HiveID:    comp.Parent.Get(cellEntry).Id,
			Intensity: 1,
		})

		justPressed = true
	} else if ebiten.IsKeyPressed(ebiten.KeyZ) {
		if justPressed && !ebiten.IsKeyPressed(ebiten.KeyShift) {
			return
		}

		activity := comp.Wandering

		cursorX, cursorY := camera.CursorWorldPosition()

		cx := (cursorX)
		cy := (cursorY)

		cellQuery := donburi.NewQuery(
			filter.Contains(comp.Cell),
		)

		cellEntry, ok := cellQuery.First(ecs.World)
		if !ok {
			print("no cell")
			return
		}

		ent.CreatePheromoneEntity(ecs.World, &ent.CreatePheromoneOptions{
			X:         cx,
			Y:         cy,
			Activity:  activity,
			HiveID:    comp.Parent.Get(cellEntry).Id,
			Intensity: 1,
		})

		justPressed = true

	} else {
		justPressed = false
	}
}
