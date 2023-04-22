package gameplay

import (
	"image"
	"image/color"
	"math"
	"math/rand"

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
	ecs.AddSystem(cellMovementSystem)

	ecs.AddRenderer(LayerBackground, cellRenderer)
	ecs.AddRenderer(LayerBackground, hiveRenderer)
}

func cellMovementSystem(ecs *ecs.ECS) {
	query := donburi.NewQuery(
		filter.And(
			filter.Contains(Position),
			filter.Contains(Velocity),
			filter.Contains(Speed),
		),
	)

	query.Each(ecs.World, func(entry *donburi.Entry) {
		// random chance to change direction
		if rand.Float32() < 0.01 {
			angle := rand.Float64() * 2 * 3.14159

			velocity := Velocity.Get(entry)
			velocity.x = math.Cos(angle) * Speed.Get(entry).speed
			velocity.y = math.Sin(angle) * Speed.Get(entry).speed
		}
	})
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

var (
	whiteImage = ebiten.NewImage(3, 3)

	// whiteSubImage is an internal sub image of whiteImage.
	// Use whiteSubImage at DrawTriangles instead of whiteImage in order to avoid bleeding edges.
	whiteSubImage = whiteImage.SubImage(image.Rect(1, 1, 2, 2)).(*ebiten.Image)
)

func init() {
	whiteImage.Fill(color.White)
}

func hiveRenderer(ecs *ecs.ECS, screen *ebiten.Image) {
	query := donburi.NewQuery(
		filter.And(
			filter.Contains(Vertices),
			filter.Contains(Indices),
		),
	)

	query.Each(ecs.World, func(entry *donburi.Entry) {
		vertices := Vertices.Get(entry).vertices
		indices := Indices.Get(entry).indices

		op := &ebiten.DrawTrianglesOptions{}
		screen.DrawTriangles(vertices, indices, whiteSubImage, op)
	})
}
