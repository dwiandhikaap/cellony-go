package system

import (
	"cellony/game/assets"
	"cellony/game/config"
	"cellony/game/gameplay/comp"
	"cellony/game/util"
	"image/color"
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
	camera "github.com/melonfunction/ebiten-camera"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/ecs"
	"github.com/yohamta/donburi/filter"
)

func MapSystem(ecs *ecs.ECS) {
	query := donburi.NewQuery(
		filter.And(
			filter.Contains(comp.Grid),
		),
	)

	query.Each(ecs.World, func(entry *donburi.Entry) {
		grid := comp.Grid.Get(entry)
		image := comp.Image.Get(entry)

		width := len(grid.Grid)
		height := len(grid.Grid[0])

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
				if !grid.DirtyMask[i][j] {
					continue
				}
				val := grid.Grid[i][j]

				op := &ebiten.DrawImageOptions{}
				op.GeoM.Translate(float64(i*tileSize), float64(j*tileSize))

				index := int(val * float32(len(tileImg)-1))

				image.Img.DrawImage(tileImg[index], op)

				grid.DirtyMask[i][j] = false
			}
		}
	})
}

func MapDestroySystem(ecs *ecs.ECS) {
	query := donburi.NewQuery(
		filter.And(
			filter.Contains(comp.Grid),
		),
	)

	query.Each(ecs.World, func(entry *donburi.Entry) {
		grid := comp.Grid.Get(entry)
		width := len(grid.Grid)
		height := len(grid.Grid[0])
		// Random tile got deleted
		if rand.Float32() < 0.1 {
			// random index
			i := rand.Intn(width)
			j := rand.Intn(height)
			grid.Grid[i][j] = 0
			grid.DirtyMask[i][j] = true
		}
	})
}

func MapRenderer(ecs *ecs.ECS, cam *camera.Camera) {
	query := donburi.NewQuery(
		filter.And(
			filter.Contains(comp.Grid),
		),
	)

	query.Each(ecs.World, func(entry *donburi.Entry) {
		image := comp.Image.Get(entry)

		op := &ebiten.DrawImageOptions{}
		op = cam.GetTranslation(op, 0, 0)
		cam.Surface.DrawImage(image.Img, op)
	})
}
