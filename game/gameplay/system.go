package gameplay

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/ecs"
	"github.com/yohamta/donburi/filter"
)

const (
	LayerBackground ecs.LayerID = iota
	LayerGame
)

func addSystem(ecs *ecs.ECS) {
	ecs.AddSystem(cellSystem)

	ecs.AddRenderer(LayerBackground, cellRenderer)
}

func cellSystem(ecs *ecs.ECS) {
	query := donburi.NewQuery(
		filter.And(
			filter.Contains(Position),
			filter.Contains(Velocity),
		),
	)

	query.Each(ecs.World, func(entry *donburi.Entry) {
		position := Position.Get(entry)
		velocity := Velocity.Get(entry)

		position.x += velocity.x * 1 / 60
		position.y += velocity.y * 1 / 60
	})
}

func cellRenderer(ecs *ecs.ECS, screen *ebiten.Image) {
	query := donburi.NewQuery(
		filter.And(
			filter.Contains(Position),
			filter.Contains(Sprite),
		),
	)

	query.Each(ecs.World, func(entry *donburi.Entry) {
		sprite := Sprite.Get(entry)
		position := Position.Get(entry)

		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(position.x, position.y)
		screen.DrawImage(sprite.sprite, op)
	})
}
