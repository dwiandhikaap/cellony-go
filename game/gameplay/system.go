package gameplay

import (
	"cellony/game/assets"
	"cellony/game/config"
	"image"
	"image/color"
	"math"
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
	camera "github.com/melonfunction/ebiten-camera"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/ecs"
	"github.com/yohamta/donburi/filter"
)

const (
	LayerBackground ecs.LayerID = iota
	LayerGame
)

func addSystem(ecs *ecs.ECS) {
	ecs.AddSystem(cameraSystem)
	ecs.AddSystem(cellSystem)
	ecs.AddSystem(mapSystem)
	ecs.AddSystem(mapDestroySystem)
	ecs.AddSystem(cellMovementSystem)

	addCameraRenderer(mapRenderer)
	addCameraRenderer(cellRenderer)
	addCameraRenderer(hiveRenderer)

	ecs.AddRenderer(LayerBackground, cameraRenderer)
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

func mapSystem(ecs *ecs.ECS) {
	query := donburi.NewQuery(
		filter.And(
			filter.Contains(Grid),
		),
	)

	query.Each(ecs.World, func(entry *donburi.Entry) {
		grid := Grid.Get(entry)
		image := Image.Get(entry)

		width := len(grid.grid)
		height := len(grid.grid[0])

		tileImg := []*ebiten.Image{
			assets.AssetsInstance.Sprites["wall0"],
			assets.AssetsInstance.Sprites["wall1"],
			assets.AssetsInstance.Sprites["wall2"],
			assets.AssetsInstance.Sprites["wall3"],
			assets.AssetsInstance.Sprites["wall4"],
			assets.AssetsInstance.Sprites["wall5"],
			assets.AssetsInstance.Sprites["wall6"],
			assets.AssetsInstance.Sprites["wall7"],
			assets.AssetsInstance.Sprites["wall8"],
			assets.AssetsInstance.Sprites["wall9"],
		}

		for i := 0; i < width; i++ {
			for j := 0; j < height; j++ {
				if !grid.dirtyMask[i][j] {
					continue
				}
				val := grid.grid[i][j]

				if val <= 0 {
					continue
				}

				op := &ebiten.DrawImageOptions{}
				op.GeoM.Translate(float64(i*10), float64(j*10))

				index := int(val * 10)

				image.img.DrawImage(tileImg[index], op)

				grid.dirtyMask[i][j] = false
			}
		}
	})
}

func mapDestroySystem(ecs *ecs.ECS) {
	query := donburi.NewQuery(
		filter.And(
			filter.Contains(Grid),
		),
	)

	query.Each(ecs.World, func(entry *donburi.Entry) {
		grid := Grid.Get(entry)
		width := len(grid.grid)
		height := len(grid.grid[0])
		// Random tile got deleted
		if rand.Float32() < 0.1 {
			// random index
			i := rand.Intn(width)
			j := rand.Intn(height)
			grid.grid[i][j] = 0
			grid.dirtyMask[i][j] = true
		}
	})
}

func cellRenderer(ecs *ecs.ECS, cam *camera.Camera) {
	query := donburi.NewQuery(
		filter.And(
			filter.Contains(Position),
			filter.Contains(Sprite),
		),
	)

	query.Each(ecs.World, func(entry *donburi.Entry) {
		sprite := Sprite.Get(entry)
		position := Position.Get(entry)
		screen := cam.Surface

		// Ass looking entity culling algorithm
		if !(position.x > cam.X-float64(cam.Width)/cam.Scale/2 && position.x < cam.X+float64(cam.Width)/cam.Scale/2 && position.y > cam.Y-float64(cam.Height)/cam.Scale/2 && position.y < cam.Y+float64(cam.Height)/cam.Scale/2) {
			return
		}

		op := &ebiten.DrawImageOptions{}
		op = cam.GetTranslation(op, position.x, position.y)
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

func hiveRenderer(ecs *ecs.ECS, cam *camera.Camera) {
	query := donburi.NewQuery(
		filter.And(
			filter.Contains(Vertices),
			filter.Contains(Indices),
		),
	)

	query.Each(ecs.World, func(entry *donburi.Entry) {
		vertices := Vertices.Get(entry).vertices
		indices := Indices.Get(entry).indices
		screen := cam.Surface

		op := &ebiten.DrawTrianglesOptions{}

		translatedVertices := make([]ebiten.Vertex, len(vertices))
		for i, v := range vertices {
			translatedVertices[i] = ebiten.Vertex{
				DstX:   float32(float64(v.DstX) - cam.X + float64(config.Video.Width/2)/cam.Scale),
				DstY:   float32(float64(v.DstY) - cam.Y + float64(config.Video.Height/2)/cam.Scale),
				SrcX:   v.SrcX,
				SrcY:   v.SrcY,
				ColorR: v.ColorR,
				ColorG: v.ColorG,
				ColorB: v.ColorB,
				ColorA: v.ColorA,
			}
		}
		screen.DrawTriangles(translatedVertices, indices, whiteSubImage, op)
	})
}

func mapRenderer(ecs *ecs.ECS, cam *camera.Camera) {
	query := donburi.NewQuery(
		filter.And(
			filter.Contains(Grid),
		),
	)

	query.Each(ecs.World, func(entry *donburi.Entry) {
		image := Image.Get(entry)

		op := &ebiten.DrawImageOptions{}
		op = cam.GetTranslation(op, 0, 0)
		cam.Surface.DrawImage(image.img, op)
	})
}
