package ent

import (
	"github.com/hajimehoshi/ebiten/v2"
	noise "github.com/ojrac/opensimplex-go"
	"github.com/yohamta/donburi"

	"cellony/game/config"
	comp "cellony/game/gameplay/component"
	"cellony/game/util"
	bitmask "cellony/lib/bit"
)

func CreateMapEntity(world donburi.World) {
	mapEntity := world.Create(comp.Grid, comp.Image)
	mapEntry := world.Entry(mapEntity)

	tileSize := config.Game.TileSize

	mapWidth := int(config.Game.Width / tileSize)
	mapHeight := int(config.Game.Height / tileSize)

	grid := make([][]float32, mapWidth)
	mask := make([][]uint8, mapWidth)

	comp.Grid.Get(mapEntry).Grid = grid
	comp.Grid.Get(mapEntry).Mask = mask
	comp.Image.Get(mapEntry).Img = ebiten.NewImage(int(config.Game.Width), int(config.Game.Height))

	terrainNoise := noise.NewNormalized(1)
	foodNoise := noise.NewNormalized(2)

	for i := 0; i < mapWidth; i++ {
		grid[i] = make([]float32, mapHeight)
		mask[i] = make([]uint8, mapHeight)
		for j := 0; j < mapHeight; j++ {
			val := float32(terrainNoise.Eval2(float64(i)/tileSize, float64(j)/tileSize))
			if val > 0.45 {
				grid[i][j] = float32(util.RangeInterpolate(float64(val), 0.45, 1.0, 0.0, 1.0))
			} else {
				grid[i][j] = 0.0
			}

			mask[i][j] = bitmask.SetBit(mask[i][j], (comp.DirtyMask))
			mask[i][j] = bitmask.SetBit(mask[i][j], (comp.WallMask))
		}
	}

	for i := 0; i < mapWidth; i++ {
		for j := 0; j < mapHeight; j++ {
			val := float32(foodNoise.Eval2(float64(i)/tileSize, float64(j)/tileSize))
			if val > 0.8 {
				gridValue := float32(util.RangeInterpolate(float64(val), 0.5, 1.0, 0.0, 1.0))
				grid[i][j] = (grid[i][j]*0.25 + gridValue*1.75) / 2.0

				mask[i][j] = bitmask.SetBit(mask[i][j], comp.DirtyMask)
				mask[i][j] = bitmask.ClearBit(mask[i][j], comp.WallMask)
			}
		}
	}
}

func IsTileDirty(bitmask uint8) bool {
	return (bitmask & uint8(comp.DirtyMask)) != 0
}

func IsTileWall(bitmask uint8) bool {
	return (bitmask & uint8(comp.WallMask)) != 0
}

func IsTileFood(bitmask uint8) bool {
	return (bitmask & uint8(comp.WallMask)) == 0
}

func IsTileMarked(bitmask uint8) bool {
	return (bitmask & uint8(comp.MarkedMask)) != 0
}
