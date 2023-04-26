package gameplay

import (
	"cellony/game/assets"
	"cellony/game/config"
	util "cellony/game/util"
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
	ecs.AddSystem(cellCollisionSystem)
	ecs.AddSystem(hiveSystem)
	ecs.AddSystem(mapSystem)
	ecs.AddSystem(mapDestroySystem)
	ecs.AddSystem(cellMovementSystem)

	addCameraRenderer(mapRenderer)
	addCameraRenderer(cellRenderer)
	addCameraRenderer(hiveRenderer)

	ecs.AddRenderer(LayerBackground, cameraRenderer)
}

func cellCollisionSystem(ecs *ecs.ECS) {
	cellQuery := donburi.NewQuery(
		filter.Contains(Cell),
	)

	worldQuery := donburi.NewQuery(
		filter.Contains(Grid),
	)

	worldQuery.Each(ecs.World, func(worldEntry *donburi.Entry) {
		grid := Grid.Get(worldEntry)

		cellQuery.Each(ecs.World, func(cellEntry *donburi.Entry) {
			cellPosition := Position.Get(cellEntry)

			x := int(cellPosition.x / config.Game.TileSize)
			y := int(cellPosition.y / config.Game.TileSize)

			if cellPosition.x >= config.Game.Width ||
				cellPosition.x < 0 ||
				cellPosition.y >= config.Game.Height ||
				cellPosition.y < 0 {
				cellEntry.Remove()
			} else if grid.grid[x][y] > 0 {
				cellEntry.Remove()

				grid.grid[x][y] -= 0.1
				grid.dirtyMask[x][y] = true
			}
		})
	})
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

		tileSize := int(config.Game.TileSize)

		deadWall := ebiten.NewImage(tileSize, tileSize)
		deadWall.Fill(color.RGBA{0x00, 0x00, 0x00, 0xff})

		tileImg := []*ebiten.Image{
			util.ResizeImage(deadWall, tileSize, tileSize),
			util.ResizeImage(assets.AssetsInstance.Sprites["wall0"], tileSize, tileSize),
			util.ResizeImage(assets.AssetsInstance.Sprites["wall1"], tileSize, tileSize),
			util.ResizeImage(assets.AssetsInstance.Sprites["wall2"], tileSize, tileSize),
			util.ResizeImage(assets.AssetsInstance.Sprites["wall3"], tileSize, tileSize),
			util.ResizeImage(assets.AssetsInstance.Sprites["wall4"], tileSize, tileSize),
			util.ResizeImage(assets.AssetsInstance.Sprites["wall5"], tileSize, tileSize),
			util.ResizeImage(assets.AssetsInstance.Sprites["wall6"], tileSize, tileSize),
			util.ResizeImage(assets.AssetsInstance.Sprites["wall7"], tileSize, tileSize),
			util.ResizeImage(assets.AssetsInstance.Sprites["wall8"], tileSize, tileSize),
			util.ResizeImage(assets.AssetsInstance.Sprites["wall9"], tileSize, tileSize),
		}

		for i := 0; i < width; i++ {
			for j := 0; j < height; j++ {
				if !grid.dirtyMask[i][j] {
					continue
				}
				val := grid.grid[i][j]

				op := &ebiten.DrawImageOptions{}
				op.GeoM.Translate(float64(i*tileSize), float64(j*tileSize))

				index := int(val * float32(len(tileImg)-1))

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
		if !(position.x > (cam.X-4)-float64(cam.Width)/cam.Scale/2 &&
			position.x < (cam.X+4)+float64(cam.Width)/cam.Scale/2 &&
			position.y > (cam.Y-4)-float64(cam.Height)/cam.Scale/2 &&
			position.y < (cam.Y+4)+float64(cam.Height)/cam.Scale/2) {
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

func hiveSystem(ecs *ecs.ECS) {
	query := donburi.NewQuery(
		filter.Contains(Hive),
	)

	query.Each(ecs.World, func(entry *donburi.Entry) {
		cx := Position.Get(entry).x
		cy := Position.Get(entry).y

		op := &CreateCellOptions{
			X:     cx,
			Y:     cy,
			Speed: 50,
		}
		createCellEntity(ecs.World, op)
	})
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
