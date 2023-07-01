package system

import (
	"cellony/game/assets"
	"cellony/game/config"
	comp "cellony/game/gameplay/component"
	bitmask "cellony/lib/bit"
	"image/color"
	"math"
	"math/rand"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	camera "github.com/melonfunction/ebiten-camera"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/ecs"
	"github.com/yohamta/donburi/filter"
)

var wallImgs []*ebiten.Image
var foodImgs []*ebiten.Image
var deadWall *ebiten.Image
var highlight *ebiten.Image
var tileSize int

var initialized = false

func initailizeAssets() {
	if initialized {
		return
	}

	tileSize = int(config.Game.TileSize)
	deadWall = ebiten.NewImage(tileSize, tileSize)
	deadWall.Fill(color.RGBA{0x00, 0x00, 0x00, 0xff})

	highlight = ebiten.NewImage(tileSize, tileSize)
	highlight.Fill(color.RGBA{5, 20, 5, 0xff})

	wallImgs = []*ebiten.Image{
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

	foodImgs = []*ebiten.Image{
		assets.AssetsInstance.Sprites["food0"],
		assets.AssetsInstance.Sprites["food1"],
		assets.AssetsInstance.Sprites["food2"],
		assets.AssetsInstance.Sprites["food3"],
		assets.AssetsInstance.Sprites["food4"],
		assets.AssetsInstance.Sprites["food5"],
		assets.AssetsInstance.Sprites["food6"],
		assets.AssetsInstance.Sprites["food7"],
		assets.AssetsInstance.Sprites["food8"],
		assets.AssetsInstance.Sprites["food9"],
	}
}

var _lastMarkedDraw = time.Now().UnixMilli()

func MapSystem(ecs *ecs.ECS) {
	initailizeAssets()

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

		markedUpdate := (time.Now().UnixMilli() - _lastMarkedDraw) > 67
		if markedUpdate {
			_lastMarkedDraw = time.Now().UnixMilli()
		}

		for i := 0; i < width; i++ {
			for j := 0; j < height; j++ {
				val := grid.Grid[i][j]
				maskVal := grid.Mask[i][j]

				isMarked := bitmask.HasBit(maskVal, comp.MarkedMask)
				isDirty := bitmask.HasBit(maskVal, comp.DirtyMask)
				isWall := bitmask.HasBit(maskVal, comp.WallMask)

				index := int(val * float32(len(wallImgs)-1))
				if isDirty && !isMarked {
					op := &ebiten.DrawImageOptions{}
					op.GeoM.Translate(float64(i*tileSize), float64(j*tileSize))

					if val <= 0 && !isMarked {
						op.Blend = ebiten.BlendClear
						image.Img.DrawImage(deadWall, op)
					} else if isWall {
						image.Img.DrawImage(wallImgs[index], op)
					} else {
						image.Img.DrawImage(foodImgs[index], op)
					}

					grid.Mask[i][j] = bitmask.ClearBit(maskVal, comp.DirtyMask)
				} else if isMarked && markedUpdate {
					op := &ebiten.DrawImageOptions{}
					op.GeoM.Translate(float64(i*tileSize), float64(j*tileSize))
					op.ColorScale.SetA(0.4 + float32(1+(math.Sin(float64(time.Now().UnixMilli())/256)))/10)

					if val <= 0 {
						image.Img.DrawImage(highlight, op)
					} else if isWall {
						image.Img.DrawImage(wallImgs[index], op)
					} else {
						image.Img.DrawImage(foodImgs[index], op)
					}
				}
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
			//grid.DirtyMask[i][j] = true
			grid.Mask[i][j] = bitmask.SetBit(grid.Mask[i][j], comp.DirtyMask)
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
